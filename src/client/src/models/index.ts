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
