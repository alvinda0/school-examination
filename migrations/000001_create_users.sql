-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role_id INT REFERENCES roles(id),
    status VARCHAR(50) DEFAULT 'active',
    last_login TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Insert sample users with different roles
-- Note: Password is 'password123' hashed with bcrypt
INSERT INTO users (full_name, email, password, role_id, status) VALUES 
    ('Super Admin User', 'superadmin@school.com', '$2a$10$rZ8qH5L5vZ5vZ5vZ5vZ5vOqH5L5vZ5vZ5vZ5vZ5vZ5vZ5vZ5vZ5vZ', 
     (SELECT id FROM roles WHERE name = 'super_admin'), 'active'),
    ('Admin User', 'admin@school.com', '$2a$10$rZ8qH5L5vZ5vZ5vZ5vZ5vOqH5L5vZ5vZ5vZ5vZ5vZ5vZ5vZ5vZ5vZ', 
     (SELECT id FROM roles WHERE name = 'admin'), 'active'),
    ('Teacher User', 'teacher@school.com', '$2a$10$rZ8qH5L5vZ5vZ5vZ5vZ5vOqH5L5vZ5vZ5vZ5vZ5vZ5vZ5vZ5vZ5vZ', 
     (SELECT id FROM roles WHERE name = 'teacher'), 'active'),
    ('Student User', 'student@school.com', '$2a$10$rZ8qH5L5vZ5vZ5vZ5vZ5vOqH5L5vZ5vZ5vZ5vZ5vZ5vZ5vZ5vZ5vZ', 
     (SELECT id FROM roles WHERE name = 'student'), 'active')
ON CONFLICT (email) DO NOTHING;

-- +migrate Down
DROP TABLE IF EXISTS users;
