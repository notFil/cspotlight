import { apiClient } from './api';
import type { User } from '@/types';

class UserService {
  public async getUsers(): Promise<User[]> {
    const response = await apiClient.get<User[]>('/api/v1/users');
    return response.data;
  }
}

export const userService = new UserService();
