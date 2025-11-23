import { apiClient } from './api';
import type { User } from '@/types';

class UserService {
  public async getUsers(): Promise<User[]> {
    const response = await apiClient.get<User[]>('/api/v1/users');
    return response.data;
  }

  public async createUser(user: Omit<User, 'id'>): Promise<User> {
    const response = await apiClient.post<User>('/api/v1/users', user);
    return response.data;
  }

  public async updateUser(id: string, user: Partial<User>): Promise<User> {
    const response = await apiClient.put<User>(`/api/v1/users/${id}`, user);
    return response.data;
  }

  public async deleteUser(id: string): Promise<void> {
    await apiClient.delete(`/api/v1/users/${id}`);
  }
}

export const userService = new UserService();
