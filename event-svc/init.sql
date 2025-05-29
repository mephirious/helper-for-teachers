-- Create extension for UUID generation if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create lesson_schedules table
CREATE TABLE IF NOT EXISTS lesson_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    group_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    valid_from TIMESTAMP WITH TIME ZONE NOT NULL,
    valid_to TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    course_id UUID NOT NULL,
    lesson_ids UUID[] NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CHECK (valid_from < valid_to)
);

-- Create indexes for lesson_schedules
CREATE INDEX IF NOT EXISTS idx_lesson_schedules_group_id ON lesson_schedules(group_id);
CREATE INDEX IF NOT EXISTS idx_lesson_schedules_course_id ON lesson_schedules(course_id);
CREATE INDEX IF NOT EXISTS idx_lesson_schedules_valid_from ON lesson_schedules(valid_from);
CREATE INDEX IF NOT EXISTS idx_lesson_schedules_valid_to ON lesson_schedules(valid_to);
CREATE INDEX IF NOT EXISTS idx_lesson_schedules_is_active ON lesson_schedules(is_active);

-- Create task_schedules table
CREATE TABLE IF NOT EXISTS task_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    group_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    valid_from TIMESTAMP WITH TIME ZONE NOT NULL,
    valid_to TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    course_id UUID NOT NULL,
    task_ids UUID[] NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CHECK (valid_from < valid_to)
);

-- Create indexes for task_schedules
CREATE INDEX IF NOT EXISTS idx_task_schedules_group_id ON task_schedules(group_id);
CREATE INDEX IF NOT EXISTS idx_task_schedules_course_id ON task_schedules(course_id);
CREATE INDEX IF NOT EXISTS idx_task_schedules_valid_from ON task_schedules(valid_from);
CREATE INDEX IF NOT EXISTS idx_task_schedules_valid_to ON task_schedules(valid_to);
CREATE INDEX IF NOT EXISTS idx_task_schedules_is_active ON task_schedules(is_active);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for automatic updated_at updates
CREATE TRIGGER update_lesson_schedules_updated_at
BEFORE UPDATE ON lesson_schedules
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_task_schedules_updated_at
BEFORE UPDATE ON task_schedules
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();