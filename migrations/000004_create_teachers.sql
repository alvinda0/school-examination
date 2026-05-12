-- 000004_create_teachers.sql

-- Create teachers table
CREATE TABLE IF NOT EXISTS teachers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL,
    nip VARCHAR(50) UNIQUE,
    gender VARCHAR(20),
    birth_place VARCHAR(100),
    birth_date DATE,
    religion VARCHAR(50),
    phone_number VARCHAR(20),
    address TEXT,
    photo_url TEXT,
    status VARCHAR(20) DEFAULT 'ACTIVE',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_teacher_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- Insert default teacher (using existing user with role teacher)
INSERT INTO teachers (user_id, nip, gender, status) 
SELECT id, 'NIP-001', 'Laki-laki', 'ACTIVE'
FROM users 
WHERE email = 'teacher@example.com'
ON CONFLICT (user_id) DO NOTHING;
