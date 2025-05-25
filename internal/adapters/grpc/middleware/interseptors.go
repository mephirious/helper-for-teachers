package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/mephirious/helper-for-teachers/services/auth-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/auth-svc/pkg/logger"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const UserCtxKey string = "userClaims"

type AuthInterceptor struct {
	publicMethods map[string]struct{}
	permissions   map[string][]domain.Role
	jwtSvc        domain.JWTService
	log           *logger.Logger
}

func NewAuthInterceptor(
	publicMethods []string,
	permissions map[string][]domain.Role,
	jwt domain.JWTService,
	log *logger.Logger,
) *AuthInterceptor {
	// build a set for quick public-check
	publicSet := make(map[string]struct{}, len(publicMethods))
	for _, m := range publicMethods {
		publicSet[m] = struct{}{}
	}

	return &AuthInterceptor{
		publicMethods: publicSet, // store the map
		permissions:   permissions,
		jwtSvc:        jwt,
		log:           log,
	}
}

// Validate the token, then auth with role and pass tokenPayload into ctx
func (i *AuthInterceptor) UnaryAuthentificate() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// Skip public methods
		if _, ok := i.publicMethods[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		// Extract token
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeaders := md["authorization"]
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization header not supplied")
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeaders[0], "Bearer "))

		// Validate token
		claims, err := i.jwtSvc.Validate(ctx, token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Role check
		if allowedRoles, ok := i.permissions[info.FullMethod]; ok {
			var found bool
			for _, role := range allowedRoles {
				if role == claims.Role {
					found = true
					break
				}
			}
			if !found {
				i.log.Warn("role is not authorized to pass", "role", claims.Role)
				return nil, status.Error(codes.PermissionDenied, "role not authorized to pass")
			}
		}

		// Inject token payload into ctx
		ctx = context.WithValue(ctx, UserCtxKey, claims)
		return handler(ctx, req)
	}
}

// Logging incoming req
func (i *AuthInterceptor) UnaryLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// Extract or generate request ID
		md, _ := metadata.FromIncomingContext(ctx)
		var reqID string
		if vals := md.Get("x-request-id"); len(vals) > 0 {
			reqID = vals[0]
		} else {
			// Append reqID to ctx
			reqID = uuid.New().String()
			ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", reqID)
		}

		// Start timer for req
		start := time.Now()

		// Log incoming req
		i.log.Info("incoming gRPC request", "method", info.FullMethod, "request_id", reqID)

		// Call handler
		resp, err = handler(ctx, req)

		// Log duration, and returning value or err
		duration := time.Since(start)
		if err != nil {
			i.log.Error("gRPC request FAILED",
				"method", info.FullMethod,
				"request_id", reqID,
				"duration", duration,
				"error", err,
			)
		} else {
			i.log.Info("gRPC request SUCCEDED",
				"method", info.FullMethod,
				"request_id", reqID,
				"duration", duration,
			)

		}

		return resp, err
	}
}

// TODO: rate limiting
