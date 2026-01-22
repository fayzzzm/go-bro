-- Todo SQL API: Standardized Request/Response Pattern
-- RESTful naming: todos.create, todos.list, todos.get, todos.update, todos.delete, todos.toggle

-- =============================================================================
-- API CONTRACT TYPES
-- =============================================================================

-- INPUT: One object to handle all potential todo parameters
CREATE TYPE todo.todo_request AS (
    id          INTEGER,
    user_id     INTEGER,
    title       TEXT,
    description TEXT,
    completed   BOOLEAN,
    limit_val   INTEGER,
    offset_val  INTEGER
);

-- OUTPUT: One object for all todo responses
CREATE TYPE todo.todo_response AS (
    id          INTEGER,
    user_id     INTEGER,
    title       VARCHAR(500),
    description TEXT,
    completed   BOOLEAN,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ
);

-- =============================================================================
-- API FUNCTIONS
-- =============================================================================

-- CREATE
CREATE OR REPLACE FUNCTION todos.create(r todo.todo_request)
RETURNS SETOF todo.todo_response AS $$
BEGIN
    RETURN QUERY
    INSERT INTO public.todos (user_id, title, description)
    VALUES (r.user_id, r.title, r.description)
    RETURNING id, user_id, title, description, completed, created_at, updated_at;
END;
$$ LANGUAGE plpgsql;

-- LIST
CREATE OR REPLACE FUNCTION todos.list(r todo.todo_request)
RETURNS SETOF todo.todo_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, user_id, title, description, completed, created_at, updated_at
    FROM public.todos
    WHERE user_id = r.user_id
    ORDER BY created_at DESC
    LIMIT COALESCE(r.limit_val, 100)
    OFFSET COALESCE(r.offset_val, 0);
END;
$$ LANGUAGE plpgsql;

-- GET
CREATE OR REPLACE FUNCTION todos.get(r todo.todo_request)
RETURNS SETOF todo.todo_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, user_id, title, description, completed, created_at, updated_at
    FROM public.todos
    WHERE id = r.id AND user_id = r.user_id;
END;
$$ LANGUAGE plpgsql;

-- UPDATE
CREATE OR REPLACE FUNCTION todos.update(r todo.todo_request)
RETURNS SETOF todo.todo_response AS $$
BEGIN
    RETURN QUERY
    UPDATE public.todos
    SET 
        title = COALESCE(r.title, title),
        description = COALESCE(r.description, description),
        completed = COALESCE(r.completed, completed),
        updated_at = NOW()
    WHERE id = r.id AND user_id = r.user_id
    RETURNING id, user_id, title, description, completed, created_at, updated_at;
END;
$$ LANGUAGE plpgsql;

-- DELETE
CREATE OR REPLACE FUNCTION todos.delete(r todo.todo_request)
RETURNS BOOLEAN AS $$
BEGIN
    DELETE FROM public.todos 
    WHERE id = r.id AND user_id = r.user_id;
    RETURN FOUND;
END;
$$ LANGUAGE plpgsql;

-- TOGGLE
CREATE OR REPLACE FUNCTION todos.toggle(r todo.todo_request)
RETURNS SETOF todo.todo_response AS $$
BEGIN
    RETURN QUERY
    UPDATE public.todos
    SET completed = NOT completed, updated_at = NOW()
    WHERE id = r.id AND user_id = r.user_id
    RETURNING id, user_id, title, description, completed, created_at, updated_at;
END;
$$ LANGUAGE plpgsql;
