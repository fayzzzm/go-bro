import { useState, type FormEvent } from 'react';

interface AddTodoFormProps {
  onAdd: (title: string) => Promise<void>;
}

export function AddTodoForm({ onAdd }: AddTodoFormProps) {
  const [title, setTitle] = useState('');
  const [isAdding, setIsAdding] = useState(false);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!title.trim() || isAdding) return;

    setIsAdding(true);
    try {
      await onAdd(title.trim());
      setTitle('');
    } catch (err) {
      console.error(err);
    } finally {
      setIsAdding(false);
    }
  };

  return (
    <form className="add-todo-form" onSubmit={handleSubmit}>
      <input
        type="text"
        className="add-todo-input"
        placeholder="What needs to be done?"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        disabled={isAdding}
      />
      <button type="submit" className="btn-add" disabled={isAdding || !title.trim()}>
        {isAdding ? 'Adding...' : 'Add Task'}
      </button>
    </form>
  );
}
