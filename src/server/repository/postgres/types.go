package postgres

// UserRequest matches the PostgreSQL type users.user_request
type UserRequest struct {
	ID           *int    `db:"id"`
	Name         *string `db:"name"`
	Email        *string `db:"email"`
	PasswordHash *string `db:"password_hash"`
	LimitVal     *int    `db:"limit_val"`
	OffsetVal    *int    `db:"offset_val"`
}

// TodoRequest matches the PostgreSQL type todos.todo_request
type TodoRequest struct {
	ID          *int    `db:"id"`
	UserID      *int    `db:"user_id"`
	Title       *string `db:"title"`
	Description *string `db:"description"`
	Completed   *bool   `db:"completed"`
	LimitVal    *int    `db:"limit_val"`
	OffsetVal   *int    `db:"offset_val"`
}
