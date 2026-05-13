ALTER TABLE students 
ADD COLUMN IF NOT EXISTS class_id UUID;

DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_student_class'
    ) THEN
        ALTER TABLE students
        ADD CONSTRAINT fk_student_class
            FOREIGN KEY(class_id)
            REFERENCES classes(id)
            ON DELETE SET NULL;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_students_class_id ON students(class_id);
