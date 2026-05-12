-- 000003_create_subjects.sql

-- Create subjects table
CREATE TABLE IF NOT EXISTS subjects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    code VARCHAR(50) UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Insert default subjects
INSERT INTO subjects (name, code, description) VALUES
    ('Matematika', 'MTK', 'Mata pelajaran Matematika'),
    ('Ilmu Pengetahuan Alam', 'IPA', 'Mata pelajaran IPA'),
    ('Ilmu Pengetahuan Sosial', 'IPS', 'Mata pelajaran IPS'),
    ('Bahasa Indonesia', 'BIND', 'Mata pelajaran Bahasa Indonesia'),
    ('Bahasa Inggris', 'BING', 'Mata pelajaran Bahasa Inggris')
ON CONFLICT (name) DO NOTHING;
