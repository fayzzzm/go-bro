-- SQL Functions for User Operations
-- These functions encapsulate all business logic for security and maintainability
-- Uses function composition for DRY and maintainability

-- =============================================================================
-- HELPER FUNCTIONS (Private/Internal)
-- =============================================================================

-- -----------------------------------------------------------------------------
-- fn_validate_name: Validates a name string
-- Returns: TRUE if valid, FALSE otherwise
-- -----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION fn_validate_name(p_name VARCHAR(255))
RETURNS BOOLEAN
LANGUAGE plpgsql
IMMUTABLE  -- Same input always gives same output
AS $$
BEGIN
    -- Name cannot be empty or whitespace-only
    RETURN TRIM(COALESCE(p_name, '')) != '';
END;
$$;

-- -----------------------------------------------------------------------------
-- fn_validate_email: Validates an email string
-- Returns: TRUE if valid, FALSE otherwise
-- -----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION fn_validate_email(p_email VARCHAR(255))
RETURNS BOOLEAN
LANGUAGE plpgsql
IMMUTABLE
AS $$
BEGIN
    -- Email must not be null and must contain @
    IF p_email IS NULL THEN
        RETURN FALSE;
    END IF;
    
    -- Must contain @ symbol (basic validation)
    IF POSITION('@' IN p_email) = 0 THEN
        RETURN FALSE;
    END IF;
    
    -- Must have something before and after @
    IF p_email LIKE '@%' OR p_email LIKE '%@' THEN
        RETURN FALSE;
    END IF;
    
    RETURN TRUE;
END;
$$;

-- -----------------------------------------------------------------------------
-- fn_email_exists: Checks if an email already exists in the database
-- Returns: TRUE if email exists, FALSE otherwise
-- -----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION fn_email_exists(p_email VARCHAR(255))
RETURNS BOOLEAN
LANGUAGE plpgsql
STABLE  -- Result depends on database state but stable within transaction
SECURITY DEFINER
AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 FROM users u 
        WHERE u.email = LOWER(TRIM(p_email))
    );
END;
$$;

-- -----------------------------------------------------------------------------
-- fn_normalize_email: Normalizes an email (lowercase, trimmed)
-- Returns: Normalized email string
-- -----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION fn_normalize_email(p_email VARCHAR(255))
RETURNS VARCHAR(255)
LANGUAGE plpgsql
IMMUTABLE
AS $$
BEGIN
    RETURN LOWER(TRIM(p_email));
END;
$$;

-- -----------------------------------------------------------------------------
-- fn_normalize_name: Normalizes a name (trimmed)
-- Returns: Normalized name string
-- -----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION fn_normalize_name(p_name VARCHAR(255))
RETURNS VARCHAR(255)
LANGUAGE plpgsql
IMMUTABLE
AS $$
BEGIN
    RETURN TRIM(p_name);
END;
$$;


-- =============================================================================
-- PUBLIC API FUNCTIONS
-- =============================================================================

-- =============================================================================
-- fn_get_user_by_id: Retrieves a user by their ID
-- Returns: The user record or raises exception if not found
-- This is a foundational function that other functions can call
-- =============================================================================
CREATE OR REPLACE FUNCTION fn_get_user_by_id(
    p_id INTEGER
)
RETURNS TABLE (
    id INTEGER,
    name VARCHAR(255),
    email VARCHAR(255),
    created_at TIMESTAMPTZ
)
LANGUAGE plpgsql
SECURITY DEFINER
STABLE
AS $$
BEGIN
    -- Validate input
    IF p_id IS NULL OR p_id <= 0 THEN
        RAISE EXCEPTION 'invalid user id'
            USING ERRCODE = 'invalid_parameter_value';
    END IF;
    
    RETURN QUERY
    SELECT u.id, u.name, u.email, u.created_at
    FROM users u
    WHERE u.id = p_id;
    
    -- Check if user was found
    IF NOT FOUND THEN
        RAISE EXCEPTION 'user not found'
            USING ERRCODE = 'no_data_found';
    END IF;
END;
$$;

-- =============================================================================
-- fn_register_user: Validates and creates a new user
-- Returns: The created user record (via fn_get_user_by_id)
-- Raises Exception: If validation fails
-- Uses: fn_validate_name, fn_validate_email, fn_email_exists, fn_get_user_by_id
-- =============================================================================
CREATE OR REPLACE FUNCTION fn_register_user(
    p_name VARCHAR(255),
    p_email VARCHAR(255)
)
RETURNS TABLE (
    id INTEGER,
    name VARCHAR(255),
    email VARCHAR(255),
    created_at TIMESTAMPTZ
)
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
DECLARE
    v_user_id INTEGER;
    v_normalized_name VARCHAR(255);
    v_normalized_email VARCHAR(255);
BEGIN
    -- =========================================================================
    -- VALIDATION RULES (Using helper functions!)
    -- =========================================================================
    
    -- Rule 1: Validate name using helper function
    IF NOT fn_validate_name(p_name) THEN
        RAISE EXCEPTION 'name cannot be empty'
            USING ERRCODE = 'check_violation';
    END IF;
    
    -- Rule 2: Validate email format using helper function
    IF NOT fn_validate_email(p_email) THEN
        RAISE EXCEPTION 'invalid email address'
            USING ERRCODE = 'check_violation';
    END IF;
    
    -- Rule 3: Check email uniqueness using helper function
    IF fn_email_exists(p_email) THEN
        RAISE EXCEPTION 'email already exists'
            USING ERRCODE = 'unique_violation';
    END IF;
    
    -- =========================================================================
    -- NORMALIZATION (Using helper functions!)
    -- =========================================================================
    v_normalized_name := fn_normalize_name(p_name);
    v_normalized_email := fn_normalize_email(p_email);
    
    -- =========================================================================
    -- INSERTION
    -- =========================================================================
    INSERT INTO users (name, email, password_hash)
    VALUES (v_normalized_name, v_normalized_email, 'temporary_hash')
    RETURNING users.id INTO v_user_id;
    
    -- =========================================================================
    -- RETURN (Using fn_get_user_by_id for consistency!)
    -- =========================================================================
    -- This ensures the return format is always consistent with fn_get_user_by_id
    -- If we ever add more fields or change formatting, it's automatic!
    RETURN QUERY SELECT * FROM fn_get_user_by_id(v_user_id);
END;
$$;

-- =============================================================================
-- fn_list_users: Lists users with pagination
-- Returns: List of user records
-- =============================================================================
CREATE OR REPLACE FUNCTION fn_list_users(
    p_limit INTEGER DEFAULT 100,
    p_offset INTEGER DEFAULT 0
)
RETURNS TABLE (
    id INTEGER,
    name VARCHAR(255),
    email VARCHAR(255),
    created_at TIMESTAMPTZ
)
LANGUAGE plpgsql
SECURITY DEFINER
STABLE
AS $$
DECLARE
    v_limit INTEGER;
    v_offset INTEGER;
BEGIN
    -- Apply default constraints (using local vars to avoid modifying params)
    v_limit := COALESCE(p_limit, 100);
    v_offset := COALESCE(p_offset, 0);
    
    -- Validate and constrain limit
    IF v_limit <= 0 THEN
        v_limit := 100;
    ELSIF v_limit > 1000 THEN
        v_limit := 1000;  -- Cap at 1000 for performance
    END IF;
    
    -- Validate offset
    IF v_offset < 0 THEN
        v_offset := 0;
    END IF;
    
    RETURN QUERY
    SELECT u.id, u.name, u.email, u.created_at
    FROM users u
    ORDER BY u.created_at DESC, u.id DESC
    LIMIT v_limit
    OFFSET v_offset;
END;
$$;


-- =============================================================================
-- COMMENTS (Documentation)
-- =============================================================================

-- Helper functions
COMMENT ON FUNCTION fn_validate_name IS 'Validates a name string - returns true if non-empty after trimming';
COMMENT ON FUNCTION fn_validate_email IS 'Validates an email string - checks for @ symbol and basic format';
COMMENT ON FUNCTION fn_email_exists IS 'Checks if an email already exists in the users table';
COMMENT ON FUNCTION fn_normalize_email IS 'Normalizes email to lowercase and trimmed';
COMMENT ON FUNCTION fn_normalize_name IS 'Normalizes name by trimming whitespace';

-- Public API functions
COMMENT ON FUNCTION fn_register_user IS 'Validates and creates a new user with email uniqueness check. Uses helper functions for validation and fn_get_user_by_id for return.';
COMMENT ON FUNCTION fn_get_user_by_id IS 'Retrieves a user by ID with validation - foundational function used by other functions';
COMMENT ON FUNCTION fn_list_users IS 'Lists users with pagination, ordered by creation date desc';
