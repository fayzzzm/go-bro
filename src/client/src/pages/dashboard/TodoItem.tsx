import type { Todo } from '../../models';

interface TodoItemProps {
  todo: Todo;
  onToggle: (id: number) => void;
  onDelete: (id: number) => void;
}

export function TodoItem({ todo, onToggle, onDelete }: TodoItemProps) {
  return (
    <div className={`todo-item ${todo.completed ? 'completed' : ''}`}>
      <div
        className={`todo-checkbox ${todo.completed ? 'checked' : ''}`}
        onClick={() => onToggle(todo.id)}
      />
      <div className="todo-content">
        <div className="todo-title">{todo.title}</div>
        {todo.description && <div className="todo-description">{todo.description}</div>}
      </div>
      <button className="todo-delete" onClick={() => onDelete(todo.id)}>
        ×
      </button>
    </div>
  );
}
