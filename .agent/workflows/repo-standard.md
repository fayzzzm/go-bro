---
description: Standards for writing PostgreSQL Repositories (Generic Pattern)
---

# Repository Standards (Generic & Object Pattern)

Follow these rules for clean, type-safe database interactions using `pgx` and generics.

## 1. Generic Helpers
Always inherit from or utilize the generic query helpers in `repository/postgres/base.go`:
- `queryOne[T](ctx, pool, query, args...)`
- `queryRows[T](ctx, pool, query, args...)`

## 2. Object Payload Pattern
When calling SQL functions, package parameters into a `map[string]any` to satisfy the "Single Object" SQL signature.

```go
payload := map[string]any{
    "id":      todoID,
    "user_id": userID,
}
return queryOne[models.Todo](ctx, r.pool, "SELECT * FROM todos.get($1)", payload)
```

## 3. Error Handling
- Use `errors.Is(err, pgx.ErrNoRows)` for specific "Not Found" logic if the SQL function doesn't raise its own exception.
- Return raw errors to the Service layer; let the Controller decide the HTTP status code.

## 4. Context usage
Always pass `context.Context` through to the database calls to allow for cancellation and timeouts.
