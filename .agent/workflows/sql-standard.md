---
description: Standards for writing Premium SQL Functions (Object Pattern)
---

# SQL Function Standards (Premium Object Pattern)

Follow these rules when creating or refactoring PostgreSQL functions to maintain a "Premium" API layer.

## 1. Schema Namespacing
- Always place functions within a dedicated schema (e.g., `users`, `todos`, `billing`).
- Use RESTful method names: `create`, `get`, `list`, `update`, `delete`, `toggle`.

## 2. The Contract Pattern (Composite Types)
Every schema MUST define two primary composite types as its API contract:
- **`schema.request`**: A single "Object" containing all possible input fields (id, filters, data, pagination).
- **`schema.response`**: A clean output "Object" for consistent returns.

```sql
CREATE TYPE todos.todo_request AS (
    id          INTEGER,
    user_id     INTEGER,
    title       TEXT,
    -- ...
    limit_val   INTEGER,
    offset_val  INTEGER
);
```

## 3. Strict Signature
Every public API function MUST accept exactly one parameter of the request type.
- **Good**: `CREATE FUNCTION todos.create(r todos.todo_request) ...`
- **Bad**: `CREATE FUNCTION todos.create(p_title TEXT, p_user_id INT) ...`

## 4. Implementation Rules
- **Security**: Use `SECURITY DEFINER` to encapsulate logic and control access.
- **Normalization**: Use internal helpers for trimming and lowercasing (e.g., `users.normalize_email`).
- **Clean Logic**: Use `COALESCE(r.field, table.field)` in updates to support partial updates.
- **Errors**: Use descriptive `RAISE EXCEPTION` with proper SQLSTATEs (e.g., `no_data_found`, `check_violation`).

## 5. Return Types
- Use `RETURNS SETOF schema.response` for consistency.
- Return the full record after `INSERT` or `UPDATE` using `RETURNING`.
