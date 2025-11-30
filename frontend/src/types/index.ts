// src/types/index.ts
export interface Team {
  id: string;
  name: string;
  description: string;
  last_updated: string;
}

export interface User {
  id?: string;
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

export interface UserCreate {
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  role: 'user' | 'admin' | 'superadmin';
  disabled: boolean;
  teamId: string;
}

export interface Project {
  id?: string;
  name: string;
  description: string;
  teamId: string;
  teamName: string;
  reportingUrl?: string;
  lastActive?: string;
  disabled: boolean;
}

export interface AuthTokenData {
  accessToken: string;
  refreshToken: string;
  createdAt: string;
  expiresIn: string;
  refreshExpiresIn: string;
}

export interface Pagination {
  page: number;
  pageSize: number;
  totalRows: number;
  totalPages: number;
}

export interface APIResponse<T> {
  data: T;
  message?: string;
  pagination?: Pagination;
  error?: string;
}

export interface AppState {
  isLoading: boolean;
  user?: User | null;
  theme: 'light' | 'dark';
}

export interface CSPReport {
  url: string;
  directive: string;
  ipAddress: string;
  raw: string;
  userAgent: string;
  count: number;
  lastSeen: string;
}

export interface Team {
  id: string;
  name: string;
  description: string;
  updatedAt: string;
}
