CREATE TABLE courses (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE groups (
    id UUID PRIMARY KEY,
    course_id UUID NOT NULL REFERENCES courses(id),
    name TEXT NOT NULL,
    expire_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(course_id, name)
);

CREATE TYPE member_role AS ENUM ('student', 'teacher');
CREATE TABLE group_members (
    id UUID PRIMARY KEY,
    group_id UUID NOT NULL REFERENCES groups(id),
    user_id UUID NOT NULL,
    role member_role NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(group_id, user_id, role)
);

CREATE TABLE course_instructors (
    id UUID PRIMARY KEY,
    course_id UUID NOT NULL REFERENCES courses(id),
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(course_id, user_id)
);
