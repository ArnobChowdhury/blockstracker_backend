-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

-- Add user_id column to spaces table
ALTER TABLE spaces ADD COLUMN user_id UUID REFERENCES users(id) ON DELETE CASCADE;
CREATE INDEX idx_spaces_user_id ON spaces(user_id);

-- Add user_id column to tags table
ALTER TABLE tags ADD COLUMN user_id UUID REFERENCES users(id) ON DELETE CASCADE;
CREATE INDEX idx_tags_user_id ON tags(user_id);

-- Add user_id column to repetitive_task_templates table
ALTER TABLE repetitive_task_templates ADD COLUMN user_id UUID REFERENCES users(id) ON DELETE CASCADE;
CREATE INDEX idx_repetitive_task_templates_user_id ON repetitive_task_templates(user_id);

-- Add user_id column to tasks table
ALTER TABLE tasks ADD COLUMN user_id UUID REFERENCES users(id) ON DELETE CASCADE;
CREATE INDEX idx_tasks_user_id ON tasks(user_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP INDEX IF EXISTS idx_tasks_user_id;
ALTER TABLE tasks DROP COLUMN user_id;

DROP INDEX IF EXISTS idx_repetitive_task_templates_user_id;
ALTER TABLE repetitive_task_templates DROP COLUMN user_id;

DROP INDEX IF EXISTS idx_tags_user_id;
ALTER TABLE tags DROP COLUMN user_id;

DROP INDEX IF EXISTS idx_spaces_user_id;
ALTER TABLE spaces DROP COLUMN user_id;

-- +goose StatementEnd
