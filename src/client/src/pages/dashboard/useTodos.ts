import { useState, useEffect, useCallback } from 'react';
import { todoService } from '../../api/services';
import type { Todo } from '../../models';

export function useTodos() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadTodos = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const data = await todoService.list();
      setTodos(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load todos');
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    loadTodos();
  }, [loadTodos]);

  const addTodo = async (title: string, description?: string) => {
    try {
      const todo = await todoService.create(title, description);
      setTodos((prev) => [todo, ...prev]);
      return todo;
    } catch (err) {
      const msg = err instanceof Error ? err.message : 'Failed to add todo';
      setError(msg);
      throw err;
    }
  };

  const toggleTodo = async (id: number) => {
    try {
      const updated = await todoService.toggle(id);
      setTodos((prev) => prev.map((t) => (t.id === id ? updated : t)));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to toggle todo');
    }
  };

  const deleteTodo = async (id: number) => {
    try {
      await todoService.delete(id);
      setTodos((prev) => prev.filter((t) => t.id !== id));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete todo');
    }
  };

  return {
    todos,
    isLoading,
    error,
    addTodo,
    toggleTodo,
    deleteTodo,
    refreshTodos: loadTodos,
  };
}
