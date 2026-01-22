-- Authentication SQL Functions
-- RESTful naming: users.create, users.get_by_email

-- =============================================================================
-- users.create: Creates a new user with hashed password
-- =============================================================================
CREATE OR REPLACE FUNCTION users.create(
    p_name VARCHAR(255),
    p_email VARCHAR(255),
    p_password_hash VARCHAR(255)
)
RETURNS SETOF user_row
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
DECLARE
    v_user_id INTEGER;
BEGIN
    -- Validate name
    IF TRIM(COALESCE(p_name, '')) = '' THEN
        RAISE EXCEPTION 'name cannot be empty'
            USING ERRCODE = 'check_violation';
    END IF;
    
    -- Validate email
    IF p_email IS NULL OR POSITION('@' IN p_email) = 0 THEN
        RAISE EXCEPTION 'invalid email address'
            USING ERRCODE = 'check_violation';
    END IF;
    
    -- Check email uniqueness
    IF EXISTS (SELECT 1 FROM public.users u WHERE u.email = LOWER(TRIM(p_email))) THEN
        RAISE EXCEPTION 'email already exists'
            USING ERRCODE = 'unique_violation';
    END IF;
    
    -- Validate password hash
    IF TRIM(COALESCE(p_password_hash, '')) = '' THEN
        RAISE EXCEPTION 'password is required'
            USING ERRCODE = 'check_violation';
    END IF;
    
    -- Insert user
    INSERT INTO public.users (name, email, password_hash)
    VALUES (TRIM(p_name), LOWER(TRIM(p_email)), p_password_hash)
    RETURNING id INTO v_user_id;
    
    -- Return user data (without password_hash)
    RETURN QUERY
    SELECT u.id, u.name, u.email, u.created_at
    FROM public.users u
    WHERE u.id = v_user_id;
END;
$$;

-- =============================================================================
-- users.get_by_email: Gets user by email for login verification
-- =============================================================================
CREATE OR REPLACE FUNCTION users.get_by_email(
    p_email VARCHAR(255)
)
RETURNS SETOF user_auth_row
LANGUAGE plpgsql
SECURITY DEFINER
STABLE
AS $$
BEGIN
    IF p_email IS NULL OR POSITION('@' IN p_email) = 0 THEN
        RAISE EXCEPTION 'invalid email address'
            USING ERRCODE = 'check_violation';
    END IF;
    
    RETURN QUERY
    SELECT u.id, u.name, u.email, u.password_hash, u.created_at
    FROM public.users u
    WHERE u.email = LOWER(TRIM(p_email));
    
    IF NOT FOUND THEN
        RAISE EXCEPTION 'user not found'
            USING ERRCODE = 'no_data_found';
    END IF;
END;
$$;

-- =============================================================================
-- COMMENTS
-- =============================================================================
COMMENT ON FUNCTION users.create IS 'Creates a new user with password hash, validates inputs';
COMMENT ON FUNCTION users.get_by_email IS 'Gets user by email including password_hash for login verification';
