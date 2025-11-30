INSERT INTO admins (
    email,
    password_hash,
    first_name,
    last_name,
    phone,
    role,
    status,
    email_verified,
    mfa_enabled,
    created_at,
    updated_at
) VALUES (
    'superadmin@example.com',
    '$2a$14$3lcx5a8xFi5Zs8z3UEs0nuZq1lo71MYP3VinRiJ.WQjwZ99KzHBsa', -- bcrypt hash for password 'secret' '
    'Alexander',
    'John',
    '+00000000000',
    'superadmin',
    'active',
    TRUE,
    FALSE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);
