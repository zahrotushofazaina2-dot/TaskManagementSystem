
-- Tambahkan SUPER_ADMIN ke daftar role yang diizinkan.
ALTER TABLE users
    DROP CONSTRAINT chk_users_role;

ALTER TABLE users
    ADD CONSTRAINT chk_users_role
    CHECK (role IN ('USER', 'ADMIN', 'SUPER_ADMIN'));

-- Tetapkan akun utama sebagai SUPER_ADMIN.
UPDATE users
SET role = 'SUPER_ADMIN'
WHERE email = 'chicken@gmail.com';

-- Verifikasi hasil migration.
SELECT id, name, email, role, is_active
FROM users
WHERE email = 'chicken@gmail.com';
