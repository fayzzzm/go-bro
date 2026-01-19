export const API_BASE = '/api/v1';

export const request = async <T>(url: string, options: RequestInit = {}): Promise<T> => {
  const res = await fetch(`${API_BASE}${url}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    credentials: 'include',
  });

  if (!res.ok) {
    const errorBody = await res.json().catch(() => ({}));
    throw new Error(errorBody.error || `Request failed with status ${res.status}`);
  }

  // Handle 204 No Content
  if (res.status === 204) {
    return {} as T;
  }

  return res.json();
};
