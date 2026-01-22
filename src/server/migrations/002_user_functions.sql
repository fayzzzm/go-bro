-- User & Auth SQL API (Standardized)
-- Schema: users
-- Pattern: Request/Response Composite Types

-- =============================================================================
-- API CONTRACT TYPES
-- =============================================================================

-- INPUT: One object for all user-related parameters
CREATE TYPE users.user_request AS (
    id            INTEGER,
    name          TEXT,
    email         TEXT,
    password_hash TEXT,
    limit_val     INTEGER,
    offset_val    INTEGER
);

-- OUTPUT: Standard user data (Public/General)
CREATE TYPE users.user_response AS (
    id         INTEGER,
    name       TEXT,
    email      TEXT,
    created_at TIMESTAMPTZ
);

-- OUTPUT: Sensitive user data (Auth only)
CREATE TYPE users.user_auth_response AS (
    id            INTEGER,
    name          TEXT,
    email         TEXT,
    password_hash TEXT,
    created_at    TIMESTAMPTZ
);

-- =============================================================================
-- INTERNAL HELPERS (Private to schema)
-- =============================================================================

CREATE OR REPLACE FUNCTION users.normalize_email(p_email TEXT) 
RETURNS TEXT AS $$ BEGIN RETURN LOWER(TRIM(p_email)); END; $$ LANGUAGE plpgsql IMMUTABLE;

CREATE OR REPLACE FUNCTION users.normalize_name(p_name TEXT) 
RETURNS TEXT AS $$ BEGIN RETURN TRIM(p_name); END; $$ LANGUAGE plpgsql IMMUTABLE;

-- =============================================================================
-- PUBLIC API FUNCTIONS
-- =============================================================================

-- CREATE
CREATE OR REPLACE FUNCTION users.create(r users.user_request)
RETURNS SETOF users.user_response AS $$
DECLARE
    v_id INTEGER;
    v_email TEXT := users.normalize_email(r.email);
    v_name TEXT := users.normalize_name(r.name);
BEGIN
    -- Validations
    IF v_name = '' THEN RAISE EXCEPTION 'name required' USING ERRCODE = 'check_violation'; END IF;
    IF v_email = '' OR POSITION('@' IN v_email) = 0 THEN RAISE EXCEPTION 'invalid email' USING ERRCODE = 'check_violation'; END IF;
    
    IF EXISTS (SELECT 1 FROM public.users WHERE email = v_email) THEN
        RAISE EXCEPTION 'email already exists' USING ERRCODE = 'unique_violation';
    END IF;

    INSERT INTO public.users (name, email, password_hash)
    VALUES (v_name, v_email, r.password_hash)
    RETURNING id INTO v_id;
    
    RETURN QUERY SELECT id, name, email, created_at FROM public.users WHERE id = v_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- GET BY EMAIL (Auth)
CREATE OR REPLACE FUNCTION users.get_by_email(r users.user_request)
RETURNS SETOF users.user_auth_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, name, email, password_hash, created_at
    FROM public.users
    WHERE email = users.normalize_email(r.email);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

-- GET BY ID
CREATE OR REPLACE FUNCTION users.get(r users.user_request)
RETURNS SETOF users.user_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, name, email, created_at
    FROM public.users
    WHERE id = r.id;
    
    IF NOT FOUND THEN RAISE EXCEPTION 'user not found' USING ERRCODE = 'no_data_found'; END IF;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

-- LIST
CREATE OR REPLACE FUNCTION users.list(r users.user_request)
RETURNS SETOF users.user_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, name, email, created_at
    FROM public.users
    ORDER BY created_at DESC
    LIMIT COALESCE(r.limit_val, 100)
    OFFSET COALESCE(r.offset_val, 0);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;
