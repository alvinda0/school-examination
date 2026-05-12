-- 000005_create_teacher_subjects.sql

-- Create teacher_subjects table (junction table)
CREATE TABLE IF NOT EXISTS teacher_subjects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    teacher_id UUID NOT NULL,
    subject_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_ts_teacher
        FOREIGN KEY (teacher_id)
        REFERENCES teachers(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_ts_subject
        FOREIGN KEY (subject_id)
        REFERENCES subjects(id)
        ON DELETE CASCADE,
    CONSTRAINT uq_teacher_subject
        UNIQUE (teacher_id, subject_id)
);

-- Insert default teacher-subject relationships
INSERT INTO teacher_subjects (teacher_id, subject_id)
SELECT t.id, s.id
FROM teachers t
CROSS JOIN subjects s
WHERE t.nip = 'NIP-001'
AND s.code IN ('MTK', 'IPA')
ON CONFLICT (teacher_id, subject_id) DO NOTHING;
