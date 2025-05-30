package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/core/domain"
)

type PostgresMembersRepository struct {
	db *sql.DB
}

func NewMembersRepository(db *sql.DB) *PostgresMembersRepository {
	return &PostgresMembersRepository{db: db}
}

func (r *PostgresMembersRepository) Create(ctx context.Context, gm *domain.GroupMember) error {
	const op = "PostgresMembersRepository.Create"
	return runInTx(ctx, r.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO group_members (id, group_id, user_id, role, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			gm.ID, gm.GroupID, gm.UserID, gm.Role, gm.CreatedAt, gm.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
}

func (r *PostgresMembersRepository) Delete(ctx context.Context, groupID, userID uuid.UUID, role domain.MemberRole) error {
	const op = "PostgresMembersRepository.Delete"
	return runInTx(ctx, r.db, func(tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx,
			`DELETE FROM group_members WHERE group_id = $1 AND user_id = $2 AND role = $3`,
			groupID, userID, role,
		)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if affected == 0 {
			return domain.ErrGroupMemberNotFound
		}
		return nil
	})
}

func (r *PostgresMembersRepository) ListByGroup(ctx context.Context, groupID uuid.UUID) ([]*domain.GroupMember, error) {
	const op = "PostgresMembersRepository.ListByGroup"
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, group_id, user_id, role, created_at, updated_at
		 FROM group_members WHERE group_id = $1`,
		groupID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var members []*domain.GroupMember
	for rows.Next() {
		var gm domain.GroupMember
		if err := rows.Scan(&gm.ID, &gm.GroupID, &gm.UserID, &gm.Role, &gm.CreatedAt, &gm.UpdatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		members = append(members, &gm)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return members, nil
}

func (r *PostgresMembersRepository) Exists(ctx context.Context, groupID, userID uuid.UUID, role domain.MemberRole) (bool, error) {
	const op = "PostgresMembersRepository.Exists"
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS (
			SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2 AND role = $3
		)`,
		groupID, userID, role,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}
