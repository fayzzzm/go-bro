-- Todo SQL Functions
-- CRUD operations for todos

-- =============================================================================
-- fn_create_todo: Creates a new todo for a user
-- =============================================================================
CREATE OR REPLACE FUNCTION fn_create_todo(
    p_user_id INTEGER,
    p_title VARCHAR(500),
    p_description TEXT DEFAULT NULL
)
RETURNS TABLE (
    id INTEGER,
    user_id INTEGER,
    title VARCHAR(500),
    description TEXT,
    completed BOOLEAN,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
DECLARE
    v_todo_id INTEGER;
BEGIN
    -- Validate user exists
    IF NOT EXISTS (SELECT 1 FROM users WHERE users.id = p_user_id) THEN
        RAISE EXCEPTION 'user not found'
            USING ERRCODE = 'foreign_key_violation';
    END IF;
    
    -- Validate title
    IF TRIM(COALESCE(p_title, '')) = '' THEN
        RAISE EXCEPTION 'title cannot be empty'
            USING ERRCODE = 'check_violation';
    END IF;
    
    -- Insert todo
    INSERT INTO todos (user_id, title, description)
    VALUES (p_user_id, TRIM(p_title), p_description)
    RETURNING todos.id INTO v_todo_id;
    
    -- Return created todo
    RETURN QUERY
    SELECT t.id, t.user_id, t.title, t.description, t.completed, t.created_at, t.updated_at
    FROM todos t
    WHERE t.id = v_todo_id;
END;
$$;

-- =============================================================================
-- fn_get_todos_by_user: Gets all todos for a user
-- =============================================================================
CREATE OR REPLACE FUNCTION fn_get_todos_by_user(
    p_user_id INTEGER,
    p_limit INTEGER DEFAULT 100,
    p_offset INTEGER DEFAULT 0
)
RETURNS TABLE (
    id INTEGER,
    user_id INTEGER,
    title VARCHAR(500),
    description TEXT,
    completed BOOLEAN,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
LANGUAGE plpgsql
SECURITY DEFINER
STABLE
AS $$
BEGIN
    -- Apply constraints
    IF p_limit IS NULL OR p_limit <= 0 THEN
        p_limit := 100;
    ELSIF p_limit > 1000 THEN
        p_limit := 1000;
    END IF;
    
    IF p_offset IS NULL OR p_offset < 0 THEN
        p_offset := 0;
    END IF;
    
    RETURN QUERY
    SELECT t.id, t.user_id, t.title, t.description, t.completed, t.created_at, t.updated_at
    FROM todos t
    WHERE t.user_id = p_user_id
    ORDER BY t.created_at DESC
    LIMIT p_limit
    OFFSET p_offset;
END;
$$;

-- =============================================================================
-- fn_get_todo_by_id: Gets a specific todo (with ownership check)
-- =============================================================================
CREATE OR REPLACE FUNCTION fn_get_todo_by_id(
    p_todo_id INTEGER,
    p_user_id INTEGER
)
RETURNS TABLE (
    id INTEGER,
    user_id INTEGER,
    title VARCHAR(500),
    description TEXT,
    completed BOOLEAN,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
LANGUAGE plpgsql
SECURITY DEFINER
STABLE
AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, t.user_id, t.title, t.description, t.completed, t.created_at, t.updated_at
    FROM todos t
    WHERE t.id = p_todo_id AND t.user_id = p_user_id;
    
    IF NOT FOUND THEN
        RAISE EXCEPTION 'todo not found'
            USING ERRCODE = 'no_data_found';
    END IF;
END;
$$;

-- =============================================================================
-- fn_update_todo: Updates a todo (with ownership check)
-- =============================================================================
CREATE OR REPLACE FUNCTION fn_update_todo(
    p_todo_id INTEGER,
    p_user_id INTEGER,
    p_title VARCHAR(500) DEFAULT NULL,
    p_description TEXT DEFAULT NULL,
    p_completed BOOLEAN DEFAULT NULL
)
RETURNS TABLE (
    id INTEGER,
    user_id INTEGER,
    title VARCHAR(500),
    description TEXT,
    completed BOOLEAN,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
BEGIN
    -- Check ownership
    IF NOT EXISTS (SELECT 1 FROM todos t WHERE t.id = p_todo_id AND t.user_id = p_user_id) THEN
        RAISE EXCEPTION 'todo not found'
            USING ERRCODE = 'no_data_found';
    END IF;
    
    -- Update only provided fields
    UPDATE todos t
    SET 
        title = COALESCE(NULLIF(TRIM(p_title), ''), t.title),
        description = COALESCE(p_description, t.description),
        completed = COALESCE(p_completed, t.completed)
    WHERE t.id = p_todo_id AND t.user_id = p_user_id;
    
    -- Return updated todo
    RETURN QUERY
    SELECT t.id, t.user_id, t.title, t.description, t.completed, t.created_at, t.updated_at
    FROM todos t
    WHERE t.id = p_todo_id;
END;
$$;

-- =============================================================================
-- fn_delete_todo: Deletes a todo (with ownership check)
-- =============================================================================
CREATE OR REPLACE FUNCTION fn_delete_todo(
    p_todo_id INTEGER,
    p_user_id INTEGER
)
RETURNS BOOLEAN
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
DECLARE
    v_deleted BOOLEAN;
BEGIN
    DELETE FROM todos t
    WHERE t.id = p_todo_id AND t.user_id = p_user_id;
    
    GET DIAGNOSTICS v_deleted = ROW_COUNT;
    
    IF NOT v_deleted THEN
        RAISE EXCEPTION 'todo not found'
            USING ERRCODE = 'no_data_found';
    END IF;
    
    RETURN TRUE;
END;
$$;

-- =============================================================================
-- fn_toggle_todo: Toggles the completed status of a todo
-- =============================================================================
CREATE OR REPLACE FUNCTION fn_toggle_todo(
    p_todo_id INTEGER,
    p_user_id INTEGER
)
RETURNS TABLE (
    id INTEGER,
    user_id INTEGER,
    title VARCHAR(500),
    description TEXT,
    completed BOOLEAN,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
BEGIN
    -- Check ownership
    IF NOT EXISTS (SELECT 1 FROM todos t WHERE t.id = p_todo_id AND t.user_id = p_user_id) THEN
        RAISE EXCEPTION 'todo not found'
            USING ERRCODE = 'no_data_found';
    END IF;
    
    -- Toggle completed
    UPDATE todos t
    SET completed = NOT t.completed
    WHERE t.id = p_todo_id AND t.user_id = p_user_id;
    
    -- Return updated todo
    RETURN QUERY
    SELECT t.id, t.user_id, t.title, t.description, t.completed, t.created_at, t.updated_at
    FROM todos t
    WHERE t.id = p_todo_id;
END;
$$;

-- =============================================================================
-- COMMENTS
-- =============================================================================
COMMENT ON FUNCTION fn_create_todo IS 'Creates a new todo for a user';
COMMENT ON FUNCTION fn_get_todos_by_user IS 'Gets all todos for a user with pagination';
COMMENT ON FUNCTION fn_get_todo_by_id IS 'Gets a specific todo with ownership check';
COMMENT ON FUNCTION fn_update_todo IS 'Updates a todo with ownership check';
COMMENT ON FUNCTION fn_delete_todo IS 'Deletes a todo with ownership check';
COMMENT ON FUNCTION fn_toggle_todo IS 'Toggles the completed status of a todo';
