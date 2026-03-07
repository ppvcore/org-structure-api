-- +goose Up
CREATE TABLE departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    parent_id INT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT departments_name_length
        CHECK (char_length(trim(name)) BETWEEN 1 AND 200),

    CONSTRAINT fk_departments_parent
        FOREIGN KEY (parent_id)
        REFERENCES departments(id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX ux_departments_parent_name
ON departments (
    COALESCE(parent_id, 0),
    lower(trim(name))
);

CREATE INDEX idx_departments_parent_id
ON departments(parent_id);

-- +goose Down
DROP TABLE IF EXISTS departments CASCADE;