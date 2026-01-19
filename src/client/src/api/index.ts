const API_BASE = '/api/v1';

export interface User {
  id: number;
  name: string;
  email: string;
  created_at: string;
}

export interface Todo {
  id: number;
  user_id: number;
  title: string;
  description?: string;
  completed: boolean;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  user: User;
  message?: string;
}

export interface TodosResponse {
  todos: Todo[];
  count: number;
}

// Auth API
export const authApi = {
  signup: async (name: string, email: string, password: string): Promise<AuthResponse> => {
    const res = await fetch(`${API_BASE}/auth/signup`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include', // Important for cookies
      body: JSON.stringify({ name, email, password }),
    });
    if (!res.ok) {
      const error = await res.json();
      throw new Error(error.error || 'Signup failed');
    }
    return res.json();
  },

  login: async (email: string, password: string): Promise<AuthResponse> => {
    const res = await fetch(`${API_BASE}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ email, password }),
    });
    if (!res.ok) {
      const error = await res.json();
      throw new Error(error.error || 'Login failed');
    }
    return res.json();
  },

  logout: async (): Promise<void> => {
    await fetch(`${API_BASE}/auth/logout`, {
      method: 'POST',
      credentials: 'include',
    });
  },

  me: async (): Promise<{ user_id: number; email: string } | null> => {
    try {
      const res = await fetch(`${API_BASE}/me`, {
        credentials: 'include',
      });
      if (!res.ok) return null;
      return res.json();
    } catch {
      return null;
    }
  },
};

// Todos API
export const todosApi = {
  list: async (): Promise<Todo[]> => {
    const res = await fetch(`${API_BASE}/todos`, {
      credentials: 'include',
    });
    if (!res.ok) {
      throw new Error('Failed to fetch todos');
    }
    const data: TodosResponse = await res.json();
    return data.todos || [];
  },

  create: async (title: string, description?: string): Promise<Todo> => {
    const res = await fetch(`${API_BASE}/todos`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ title, description }),
    });
    if (!res.ok) {
      const error = await res.json();
      throw new Error(error.error || 'Failed to create todo');
    }
    const data = await res.json();
    return data.todo;
  },

  toggle: async (id: number): Promise<Todo> => {
    const res = await fetch(`${API_BASE}/todos/${id}/toggle`, {
      method: 'PATCH',
      credentials: 'include',
    });
    if (!res.ok) {
      throw new Error('Failed to toggle todo');
    }
    const data = await res.json();
    return data.todo;
  },

  delete: async (id: number): Promise<void> => {
    const res = await fetch(`${API_BASE}/todos/${id}`, {
      method: 'DELETE',
      credentials: 'include',
    });
    if (!res.ok) {
      throw new Error('Failed to delete todo');
    }
  },
};
