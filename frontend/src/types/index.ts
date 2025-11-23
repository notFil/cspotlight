// src/types/index.ts
export interface Team {
  id: string;
  name: string;
  description: string;
  last_updated: string;
}

export interface User {
  id: string;
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  role: 'user' | 'admin' | 'superadmin';
  disabled: boolean;
  teamId: string;
  teamName: string;
  updatedAt?: string; // Making optional as it's not in the sample response, but keeping for compatibility if needed
}

export interface Project {
  id: string;
  name: string;
  description: string;
  team: string; // Team ID or Name
  reportingUrl: string;
  updatedAt: string;
}

export interface APIResponse<T> {
  data: T;
  message?: string;
  error?: string;
}

export interface AppState {
  isLoading: boolean;
  user?: User | null;
  theme: 'light' | 'dark';
}