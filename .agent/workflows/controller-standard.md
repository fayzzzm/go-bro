---
description: Standards for writing Go Gin Controllers (Standardized Reply Pattern)
---

# Controller Standards (Standardized Reply Pattern)

Follow these rules to ensure all controllers are lean, consistent, and boilerplate-free.

## 1. Naming Conventions
- **Receiver**: Always use `ctrl` for the controller struct receiver.
- **Context**: Always use `c` for the `*gin.Context` variable.
- **Example**: `func (ctrl *TodoController) Create(c *gin.Context)`

## 2. Shared Reply Package
Never use `c.JSON()` or `c.AbortWithStatus()` directly. Always use the `pkg/reply` package for symmetry and centralization.

- **Success**:
    - `reply.OK(c, data)` -> 200 OK
    - `reply.Created(c, data)` -> 201 Created
- **Errors**:
    - `reply.Error(c, code, message, err)` -> Custom code
    - `reply.NotFound(c, err)` -> 404 Not Found
    - `reply.InternalError(c, err)` -> 500 Internal Error

## 3. Request Binding
Use the generic `middleware.GetBody[T](c)` helper to retrieve JSON payloads. It ensures the type is correct and reduces manual binding noise.

```go
req := middleware.GetBody[CreateRequest](c)
```

## 4. Logical Flow
Handlers should be high-level orchestrators:
1.  Extract parameters (ID from path, Query from query).
2.  Get body (if applicable).
3.  Call Service/UseCase.
4.  Check for errors using a single `if reply.Error(...) { return }` block.
5.  Send back the successful result using `reply.OK(...)`.

## 5. Lean Handlers
If a handler exceeds 20-30 lines, move logic to the Service or UseCase layer. Controllers should only handle HTTP concerns.
