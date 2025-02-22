    -- +goose Up
    -- +goose StatementBegin
    SELECT 'up SQL query';
    CREATE OR REPLACE FUNCTION update_modified_at_column()
    RETURNS TRIGGER AS $$
    BEGIN
        NEW.modified_at = CURRENT_TIMESTAMP;
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    CREATE TRIGGER set_modified_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_at_column();

    CREATE TRIGGER set_modified_at
    BEFORE UPDATE ON repetitive_task_templates
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_at_column();

    CREATE TRIGGER set_modified_at
    BEFORE UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_at_column();

    CREATE TRIGGER set_modified_at
    BEFORE UPDATE ON tags
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_at_column();

    CREATE TRIGGER set_modified_at
    BEFORE UPDATE ON spaces
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_at_column();

    -- +goose StatementEnd

    -- +goose Down
    -- +goose StatementBegin
    SELECT 'down SQL query';
    DROP TRIGGER IF EXISTS set_modified_at ON users;
    DROP TRIGGER IF EXISTS set_modified_at ON repetitive_task_templates;
    DROP TRIGGER IF EXISTS set_modified_at ON tasks;
    DROP TRIGGER IF EXISTS set_modified_at ON tags;
    DROP TRIGGER IF EXISTS set_modified_at ON spaces;

    -- +goose StatementEnd
