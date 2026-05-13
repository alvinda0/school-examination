-- Seed data untuk classes
-- Pastikan sudah ada data teachers terlebih dahulu

-- Insert sample classes untuk tahun ajaran 2024/2025
INSERT INTO classes (id, name, grade_level, academic_year, homeroom_teacher_id, max_students, status, created_at, updated_at) VALUES
-- Kelas 10 (X)
(gen_random_uuid(), 'X IPA 1', 10, '2024/2025', NULL, 40, 'ACTIVE', NOW(), NOW()),
(gen_random_uuid(), 'X IPA 2', 10, '2024/2025', NULL, 40, 'ACTIVE', NOW(), NOW()),
(gen_random_uuid(), 'X IPS 1', 10, '2024/2025', NULL, 40, 'ACTIVE', NOW(), NOW()),
(gen_random_uuid(), 'X IPS 2', 10, '2024/2025', NULL, 40, 'ACTIVE', NOW(), NOW()),

-- Kelas 11 (XI)
(gen_random_uuid(), 'XI IPA 1', 11, '2024/2025', NULL, 40, 'ACTIVE', NOW(), NOW()),
(gen_random_uuid(), 'XI IPA 2', 11, '2024/2025', NULL, 40, 'ACTIVE', NOW(), NOW()),
(gen_random_uuid(), 'XI IPS 1', 11, '2024/2025', NULL, 40, 'ACTIVE', NOW(), NOW()),
(gen_random_uuid(), 'XI IPS 2', 11, '2024/2025', NULL, 40, 'ACTIVE', NOW(), NOW()),

-- Kelas 12 (XII)
(gen_random_uuid(), 'XII IPA 1', 12, '2024/2025', NULL, 40, 'ACTIVE', NOW(), NOW()),
(gen_random_uuid(), 'XII IPA 2', 12, '2024/2025', NULL, 40, 'ACTIVE', NOW(), NOW()),
(gen_random_uuid(), 'XII IPS 1', 12, '2024/2025', NULL, 40, 'ACTIVE', NOW(), NOW()),
(gen_random_uuid(), 'XII IPS 2', 12, '2024/2025', NULL, 40, 'ACTIVE', NOW(), NOW())
ON CONFLICT (name, academic_year) DO NOTHING;

-- Kelas untuk tahun ajaran 2023/2024 (archived)
INSERT INTO classes (id, name, grade_level, academic_year, homeroom_teacher_id, max_students, status, created_at, updated_at) VALUES
(gen_random_uuid(), 'X IPA 1', 10, '2023/2024', NULL, 40, 'INACTIVE', NOW(), NOW()),
(gen_random_uuid(), 'X IPA 2', 10, '2023/2024', NULL, 40, 'INACTIVE', NOW(), NOW()),
(gen_random_uuid(), 'XI IPA 1', 11, '2023/2024', NULL, 40, 'INACTIVE', NOW(), NOW()),
(gen_random_uuid(), 'XII IPA 1', 12, '2023/2024', NULL, 40, 'INACTIVE', NOW(), NOW())
ON CONFLICT (name, academic_year) DO NOTHING;

-- Display inserted classes
SELECT 
    name, 
    grade_level, 
    academic_year, 
    max_students, 
    status,
    created_at
FROM classes
ORDER BY academic_year DESC, grade_level ASC, name ASC;

-- Count classes by academic year
SELECT 
    academic_year,
    status,
    COUNT(*) as total_classes
FROM classes
GROUP BY academic_year, status
ORDER BY academic_year DESC, status;
