-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE changes (
    user_id UUID NOT NULL,
    change_id BIGINT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id UUID NOT NULL,
    operation TEXT NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, change_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_changes_user_id_changed_at ON changes(user_id, changed_at DESC);

ALTER TABLE spaces ADD COLUMN last_change_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE tags ADD COLUMN last_change_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE repetitive_task_templates ADD COLUMN last_change_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE tasks ADD COLUMN last_change_id BIGINT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE tasks DROP COLUMN last_change_id;
ALTER TABLE repetitive_task_templates DROP COLUMN last_change_id;
ALTER TABLE tags DROP COLUMN last_change_id;
ALTER TABLE spaces DROP COLUMN last_change_id;

DROP INDEX IF EXISTS idx_changes_user_id_changed_at;
DROP TABLE IF EXISTS changes;
-- +goose StatementEnd
