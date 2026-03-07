-- +goose Up
CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    department_id INT NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    position VARCHAR(200) NOT NULL,
    hired_at DATE NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT employees_full_name_length
        CHECK (char_length(trim(full_name)) BETWEEN 1 AND 200),

    CONSTRAINT employees_position_length
        CHECK (char_length(trim(position)) BETWEEN 1 AND 200),

    CONSTRAINT fk_employees_department
        FOREIGN KEY (department_id)
        REFERENCES departments(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_employees_department_id
ON employees(department_id);

-- +goose Down
DROP TABLE IF EXISTS employees;