-- Todo SQL Functions
-- RESTful naming: todos.create, todos.list, todos.get, todos.update, todos.delete, todos.toggle

-- =============================================================================
-- todos.create: Creates a new todo for a user
-- =============================================================================
CREATE TYPE todo.list_response AS (
    id INTEGER,
    user_id INTEGER,
    title VARCHAR(500),
    description TEXT,
    completed BOOLEAN,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE OR REPLACE FUNCTION todos.create(
    p_user_id INTEGER,
    p_title VARCHAR(500),
    p_description TEXT DEFAULT NULL
)
RETURNS SETOF todos.list_response
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
DECLARE
    v_todo_id INTEGER;
BEGIN
    -- Validate user exists
    IF NOT EXISTS (SELECT 1 FROM public.users WHERE id = p_user_id) THEN
        RAISE EXCEPTION 'user not found'
            USING ERRCODE = 'foreign_key_violation';
    END IF;
    
    -- Validate title
    IF TRIM(COALESCE(p_title, '')) = '' THEN
        RAISE EXCEPTION 'title cannot be empty'
            USING ERRCODE = 'check_violation';
    END IF;
    
    -- Insert todo
    INSERT INTO public.todos (user_id, title, description)
    VALUES (p_user_id, TRIM(p_title), p_description)
    RETURNING id INTO v_todo_id;
    
    -- Return created todo
    RETURN QUERY
    SELECT t.id, t.user_id, t.title, t.description, t.completed, t.created_at, t.updated_at
    FROM public.todos t
    WHERE t.id = v_todo_id;
END;
$$;

-- =============================================================================
-- todos.list: Gets all todos for a user with pagination
-- =============================================================================
CREATE OR REPLACE FUNCTION todos.list(
    p_user_id INTEGER,
    p_limit INTEGER DEFAULT 100,
    p_offset INTEGER DEFAULT 0
)
RETURNS SETOF todos.list_response
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
    FROM public.todos t
    WHERE t.user_id = p_user_id
    ORDER BY t.created_at DESC
    LIMIT p_limit
    OFFSET p_offset;
END;
$$;

-- =============================================================================
-- todos.get: Gets a specific todo (with ownership check)
-- =============================================================================
CREATE OR REPLACE FUNCTION todos.get(
    p_todo_id INTEGER,
    p_user_id INTEGER
)
RETURNS SETOF todo_row
LANGUAGE plpgsql
SECURITY DEFINER
STABLE
AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, t.user_id, t.title, t.description, t.completed, t.created_at, t.updated_at
    FROM public.todos t
    WHERE t.id = p_todo_id AND t.user_id = p_user_id;
    
    IF NOT FOUND THEN
        RAISE EXCEPTION 'todo not found'
            USING ERRCODE = 'no_data_found';
    END IF;
END;
$$;

-- =============================================================================
-- todos.update: Updates a todo (with ownership check)
-- =============================================================================
CREATE OR REPLACE FUNCTION todos.update(
    p_todo_id INTEGER,
    p_user_id INTEGER,
    p_title VARCHAR(500) DEFAULT NULL,
    p_description TEXT DEFAULT NULL,
    p_completed BOOLEAN DEFAULT NULL
)
RETURNS SETOF todo_row
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
BEGIN
    -- Check ownership
    IF NOT EXISTS (SELECT 1 FROM public.todos t WHERE t.id = p_todo_id AND t.user_id = p_user_id) THEN
        RAISE EXCEPTION 'todo not found'
            USING ERRCODE = 'no_data_found';
    END IF;
    
    -- Update only provided fields
    UPDATE public.todos t
    SET 
        title = COALESCE(NULLIF(TRIM(p_title), ''), t.title),
        description = COALESCE(p_description, t.description),
        completed = COALESCE(p_completed, t.completed)
    WHERE t.id = p_todo_id AND t.user_id = p_user_id;
    
    -- Return updated todo
    RETURN QUERY
    SELECT t.id, t.user_id, t.title, t.description, t.completed, t.created_at, t.updated_at
    FROM public.todos t
    WHERE t.id = p_todo_id;
END;
$$;

-- =============================================================================
-- todos.delete: Deletes a todo (with ownership check)
-- =============================================================================
CREATE OR REPLACE FUNCTION todos.delete(
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
    DELETE FROM public.todos t
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
-- todos.toggle: Toggles the completed status of a todo
-- =============================================================================
CREATE OR REPLACE FUNCTION todos.toggle(
    p_todo_id INTEGER,
    p_user_id INTEGER
)
RETURNS SETOF todo_row
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
BEGIN
    -- Check ownership
    IF NOT EXISTS (SELECT 1 FROM public.todos t WHERE t.id = p_todo_id AND t.user_id = p_user_id) THEN
        RAISE EXCEPTION 'todo not found'
            USING ERRCODE = 'no_data_found';
    END IF;
    
    -- Toggle completed
    UPDATE public.todos t
    SET completed = NOT t.completed
    WHERE t.id = p_todo_id AND t.user_id = p_user_id;
    
    -- Return updated todo
    RETURN QUERY
    SELECT t.id, t.user_id, t.title, t.description, t.completed, t.created_at, t.updated_at
    FROM public.todos t
    WHERE t.id = p_todo_id;
END;
$$;

-- =============================================================================
-- COMMENTS
-- =============================================================================
COMMENT ON FUNCTION todos.create IS 'Creates a new todo for a user';
COMMENT ON FUNCTION todos.list IS 'Gets all todos for a user with pagination';
COMMENT ON FUNCTION todos.get IS 'Gets a specific todo with ownership check';
COMMENT ON FUNCTION todos.update IS 'Updates a todo with ownership check';
COMMENT ON FUNCTION todos.delete IS 'Deletes a todo with ownership check';
COMMENT ON FUNCTION todos.toggle IS 'Toggles the completed status of a todo';
