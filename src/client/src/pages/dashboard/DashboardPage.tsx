import { useTodos } from './useTodos';
import { Navbar } from '../../shared/Navbar';
import { TodoItem } from './TodoItem';
import { AddTodoForm } from './AddTodoForm';

export default function DashboardPage() {
  const { todos, isLoading, addTodo, toggleTodo, deleteTodo } = useTodos();

  return (
    <div className="dashboard">
      <Navbar />

      <main className="main-content">
        <section className="todo-section">
          <h2>My Tasks</h2>

          <AddTodoForm onAdd={async (title) => { await addTodo(title); }} />

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
                <TodoItem 
                  key={todo.id} 
                  todo={todo} 
                  onToggle={toggleTodo} 
                  onDelete={deleteTodo} 
                />
              ))}
            </div>
          )}
        </section>
      </main>
    </div>
  );
}
