-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS spaces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- Indexes
CREATE INDEX idx_spaces_name ON spaces(name);

CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
-- Indexes
CREATE INDEX idx_tags_name ON tags(name);

CREATE TYPE task_time_of_day AS ENUM ('morning', 'afternoon', 'evening', 'night');

CREATE TABLE IF NOT EXISTS repetitive_task_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    is_active BOOLEAN DEFAULT TRUE,
    title VARCHAR NOT NULL,
    description TEXT,
    schedule VARCHAR NOT NULL,
    priority INT DEFAULT 3,
    should_be_scored BOOLEAN DEFAULT FALSE,
    monday BOOLEAN DEFAULT FALSE,
    tuesday BOOLEAN DEFAULT FALSE,
    wednesday BOOLEAN DEFAULT FALSE,
    thursday BOOLEAN DEFAULT FALSE,
    friday BOOLEAN DEFAULT FALSE,
    saturday BOOLEAN DEFAULT FALSE,
    sunday BOOLEAN DEFAULT FALSE,
    time_of_day task_time_of_day,
    last_date_of_task_generation TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    space_id UUID,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE SET NULL
);

-- Indexes
CREATE INDEX idx_repetitive_task_templates_space_id ON repetitive_task_templates(space_id);
CREATE INDEX idx_repetitive_task_templates_title ON repetitive_task_templates(title);

CREATE TYPE task_status AS ENUM ('INCOMPLETE', 'FAILED', 'COMPLETED');

CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    is_active BOOLEAN DEFAULT TRUE,
    title VARCHAR NOT NULL,
    description TEXT,
    schedule VARCHAR,
    priority INT DEFAULT 3,
    completion_status task_status DEFAULT 'INCOMPLETE',
    due_date TIMESTAMPTZ,
    should_be_scored BOOLEAN,
    score INT,
    time_of_day task_time_of_day,
    repetitive_task_template_id UUID NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    space_id UUID,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (repetitive_task_template_id) REFERENCES repetitive_task_templates(id) ON DELETE SET NULL,
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE
);


-- Indexes
CREATE INDEX idx_tasks_repetitive_task_template_id ON tasks(repetitive_task_template_id);
CREATE INDEX idx_tasks_space_id ON tasks(space_id);
CREATE INDEX idx_tasks_completion_status ON tasks(completion_status);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);


CREATE TABLE IF NOT EXISTS task_tags (
    task_id UUID,
    tag_id UUID,
    PRIMARY KEY (task_id, tag_id),
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_task_tags_task_id ON task_tags(task_id);
CREATE INDEX idx_task_tags_tag_id ON task_tags(tag_id);

CREATE TABLE IF NOT EXISTS repetitive_task_template_tags (
    repetitive_task_template_id UUID,
    tag_id UUID,
    PRIMARY KEY (repetitive_task_template_id, tag_id),
    FOREIGN KEY (repetitive_task_template_id) REFERENCES repetitive_task_templates(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_repetitive_task_template_tags_template_id ON repetitive_task_template_tags(repetitive_task_template_id);
CREATE INDEX idx_repetitive_task_template_tags_tag_id ON repetitive_task_template_tags(tag_id);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR NOT NULL UNIQUE,
    password VARCHAR,
    provider VARCHAR,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT password_or_provider CHECK (
        password IS NOT NULL OR provider IS NOT NULL
    )
);

-- Indexes
CREATE INDEX idx_users_email ON users(email);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_spaces_name;
DROP INDEX IF EXISTS idx_tags_name;
DROP INDEX IF EXISTS idx_repetitive_task_templates_space_id;
DROP INDEX IF EXISTS idx_repetitive_task_templates_title;
DROP INDEX IF EXISTS idx_tasks_repetitive_task_template_id;
DROP INDEX IF EXISTS idx_tasks_space_id;
DROP INDEX IF EXISTS idx_tasks_completion_status;
DROP INDEX IF EXISTS idx_tasks_due_date;
DROP INDEX IF EXISTS idx_task_tags_task_id;
DROP INDEX IF EXISTS idx_task_tags_tag_id;
DROP INDEX IF EXISTS idx_repetitive_task_template_tags_template_id;
DROP INDEX IF EXISTS idx_repetitive_task_template_tags_tag_id;
DROP INDEX IF EXISTS idx_users_email;

DROP TABLE IF EXISTS repetitive_task_template_tags CASCADE;
DROP TABLE IF EXISTS task_tags CASCADE;
DROP TABLE IF EXISTS tasks CASCADE;
DROP TABLE IF EXISTS repetitive_task_templates CASCADE;
DROP TABLE IF EXISTS tags CASCADE;
DROP TABLE IF EXISTS spaces CASCADE;
DROP TABLE IF EXISTS users CASCADE;

DROP TYPE IF EXISTS task_status;
-- +goose StatementEnd
