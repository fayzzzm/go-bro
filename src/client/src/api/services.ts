import { request } from './client';
import type { User, Todo } from '../models';

export interface AuthResponse {
  user: User;
  message?: string;
}

export const authService = {
  signup: (name: string, email: string, password: string) =>
    request<AuthResponse>('/auth/signup', {
      method: 'POST',
      body: JSON.stringify({ name, email, password }),
    }),

  login: (email: string, password: string) =>
    request<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),

  logout: () =>
    request<void>('/auth/logout', {
      method: 'POST',
    }),

  me: () =>
    request<{ user_id: number; email: string }>('/me'),
};

export const todoService = {
  list: async () => {
    const data = await request<{ todos: Todo[]; count: number }>('/todos');
    return data.todos || [];
  },

  create: (title: string, description?: string) =>
    request<{ todo: Todo }>('/todos', {
      method: 'POST',
      body: JSON.stringify({ title, description }),
    }).then(res => res.todo),

  toggle: (id: number) =>
    request<{ todo: Todo }>(`/todos/${id}/toggle`, {
      method: 'PATCH',
    }).then(res => res.todo),

  delete: (id: number) =>
    request<void>(`/todos/${id}`, {
      method: 'DELETE',
    }),
};
