-- +migrate Up
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample roles
INSERT INTO roles (name, description) VALUES 
    ('super_admin', 'Super Administrator with full system access'),
    ('admin', 'Administrator with management access'),
    ('teacher', 'Teacher with teaching and grading access'),
    ('student', 'Student with learning access')
ON CONFLICT (name) DO NOTHING;

-- +migrate Down
DROP TABLE IF EXISTS roles;
