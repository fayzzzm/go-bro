import { useState, useEffect, type FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { todosApi } from '../api';
import type { Todo } from '../api';

export default function DashboardPage() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [newTodoTitle, setNewTodoTitle] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [isAdding, setIsAdding] = useState(false);
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    loadTodos();
  }, []);

  const loadTodos = async () => {
    try {
      const data = await todosApi.list();
      setTodos(data);
    } catch (error) {
      console.error('Failed to load todos:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleAddTodo = async (e: FormEvent) => {
    e.preventDefault();
    if (!newTodoTitle.trim()) return;

    setIsAdding(true);
    try {
      const todo = await todosApi.create(newTodoTitle.trim());
      setTodos([todo, ...todos]);
      setNewTodoTitle('');
    } catch (error) {
      console.error('Failed to add todo:', error);
    } finally {
      setIsAdding(false);
    }
  };

  const handleToggle = async (id: number) => {
    try {
      const updated = await todosApi.toggle(id);
      setTodos(todos.map((t) => (t.id === id ? updated : t)));
    } catch (error) {
      console.error('Failed to toggle todo:', error);
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await todosApi.delete(id);
      setTodos(todos.filter((t) => t.id !== id));
    } catch (error) {
      console.error('Failed to delete todo:', error);
    }
  };

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  return (
    <div className="dashboard">
      <nav className="navbar">
        <h1>📝 Todos</h1>
        <div className="navbar-right">
          <span className="user-email">{user?.email}</span>
          <button className="btn-logout" onClick={handleLogout}>
            Logout
          </button>
        </div>
      </nav>

      <main className="main-content">
        <section className="todo-section">
          <h2>My Tasks</h2>

          <form className="add-todo-form" onSubmit={handleAddTodo}>
            <input
              type="text"
              className="add-todo-input"
              placeholder="What needs to be done?"
              value={newTodoTitle}
              onChange={(e) => setNewTodoTitle(e.target.value)}
              disabled={isAdding}
            />
            <button type="submit" className="btn-add" disabled={isAdding || !newTodoTitle.trim()}>
              {isAdding ? 'Adding...' : 'Add Task'}
            </button>
          </form>

          {isLoading ? (
            <div className="loading">Loading todos...</div>
          ) : todos.length === 0 ? (
            <div className="empty-state">
              <p>No tasks yet</p>
              <span>Add your first task above!</span>
            </div>
          ) : (
            <div className="todo-list">
              {todos.map((todo) => (
                <div key={todo.id} className={`todo-item ${todo.completed ? 'completed' : ''}`}>
                  <div
                    className={`todo-checkbox ${todo.completed ? 'checked' : ''}`}
                    onClick={() => handleToggle(todo.id)}
                  />
                  <div className="todo-content">
                    <div className="todo-title">{todo.title}</div>
                    {todo.description && <div className="todo-description">{todo.description}</div>}
                  </div>
                  <button className="todo-delete" onClick={() => handleDelete(todo.id)}>
                    ×
                  </button>
                </div>
              ))}
            </div>
          )}
        </section>
      </main>
    </div>
  );
}
