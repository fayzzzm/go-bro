-- 1. Public User info (for listing, profile, etc)
CREATE TYPE user_row AS (
    id INTEGER,
    name VARCHAR(255),
    email VARCHAR(255),
    created_at TIMESTAMPTZ
);

-- 2. Auth User info (includes password for verification)
CREATE TYPE user_auth_row AS (
    id INTEGER,
    name VARCHAR(255),
    email VARCHAR(255),
    password_hash VARCHAR(255),
    created_at TIMESTAMPTZ
);

-- 3. Todo row
CREATE TYPE todo_row AS (
    id INTEGER,
    user_id INTEGER,
    title VARCHAR(500),
    description TEXT,
    completed BOOLEAN,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
