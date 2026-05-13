-- Seed data untuk audit_logs
-- Menggunakan user yang sudah ada di database

-- Contoh audit log untuk CREATE student (by Admin User)
INSERT INTO audit_logs (
    id, user_id, full_name, role_id, role_name, method, endpoint, 
    status_code, ip_address, user_agent, duration_ms, action, 
    entity_id, entity_type, new_data, created_at
) VALUES (
    'a4c5e7c1-3b7b-4d2d-9c44-1f9c5d87c101',
    '80cba4ea-f7c5-413c-b74d-4b66b26ed422',
    'Admin User',
    '9e4b5838-fc29-4cee-9b00-f8c089841a42',
    'admin',
    'POST',
    '/api/v1/students',
    201,
    '[::1]',
    'PostmanRuntime/7.54.0',
    25,
    'create',
    '59e99021-0aa1-4849-8c2f-9bf99c97f679',
    'student',
    '{"id": "59e99021-0aa1-4849-8c2f-9bf99c97f679", "full_name": "John Doe", "nis": "20260001", "status": "active"}'::jsonb,
    '2026-05-13 23:20:10'
);

-- Contoh audit log untuk UPDATE class (by Teacher)
INSERT INTO audit_logs (
    id, user_id, full_name, role_id, role_name, method, endpoint, 
    status_code, ip_address, user_agent, duration_ms, action, 
    entity_id, entity_type, changes, created_at
) VALUES (
    'c8f91d1a-2b2e-4f82-b0a8-5b5b5c0ef202',
    'a4036981-fa4f-4082-9fad-efd882ddafba',
    'Alvinda',
    '01a3ce30-b248-49b3-ade2-147076a2d8e3',
    'teacher',
    'PUT',
    '/api/v1/classes/38bd672e-2706-4eca-96a1-548171ef4fc4',
    200,
    '[::1]',
    'Mozilla/5.0',
    15,
    'update',
    '38bd672e-2706-4eca-96a1-548171ef4fc4',
    'class',
    '{"class_name": {"old": "X IPA 1", "new": "X IPA 2"}, "capacity": {"old": 30, "new": 35}}'::jsonb,
    '2026-05-13 22:57:19'
);

-- Contoh audit log untuk PARTIAL UPDATE student (by Admin)
INSERT INTO audit_logs (
    id, user_id, full_name, role_id, role_name, method, endpoint, 
    status_code, ip_address, user_agent, duration_ms, action, 
    entity_id, entity_type, changes, created_at
) VALUES (
    'd2a4b71f-9f11-4f93-a5a3-3b0fd4bb3303',
    '80cba4ea-f7c5-413c-b74d-4b66b26ed422',
    'Admin User',
    '9e4b5838-fc29-4cee-9b00-f8c089841a42',
    'admin',
    'PATCH',
    '/api/v1/students/59e99021-0aa1-4849-8c2f-9bf99c97f679',
    200,
    '[::1]',
    'PostmanRuntime/7.54.0',
    12,
    'partial_update',
    '59e99021-0aa1-4849-8c2f-9bf99c97f679',
    'student',
    '{"phone_number": {"old": "08123456789", "new": "08987654321"}}'::jsonb,
    '2026-05-13 23:15:01'
);

-- Contoh audit log untuk DELETE student (by Admin)
INSERT INTO audit_logs (
    id, user_id, full_name, role_id, role_name, method, endpoint, 
    status_code, ip_address, user_agent, duration_ms, action, 
    entity_id, entity_type, deleted_data, created_at
) VALUES (
    'f1d93b8e-7a22-4c7a-bfa7-9f8d2b9a4404',
    '80cba4ea-f7c5-413c-b74d-4b66b26ed422',
    'Admin User',
    '9e4b5838-fc29-4cee-9b00-f8c089841a42',
    'admin',
    'DELETE',
    '/api/v1/students/59e99021-0aa1-4849-8c2f-9bf99c97f679',
    200,
    '[::1]',
    'PostmanRuntime/7.54.0',
    9,
    'delete',
    '59e99021-0aa1-4849-8c2f-9bf99c97f679',
    'student',
    '{"id": "59e99021-0aa1-4849-8c2f-9bf99c97f679", "full_name": "John Doe", "nis": "20260001"}'::jsonb,
    '2026-05-13 23:25:40'
);
