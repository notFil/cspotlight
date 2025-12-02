import { apiClient } from './api';
import type { User, UserRegister, ChangePasswordRequest } from '@/types';

class UserService {
  public async getUsers(): Promise<User[]> {
    const response = await apiClient.get<User[]>('/api/v1/users');
    return response.data;
  }

  public async getCurrentUser(): Promise<User> {
    const response = await apiClient.get<User>('/api/v1/users/me');
    return response.data;
  }

  public async createUser(user: Omit<User, 'id'>): Promise<User> {
    const response = await apiClient.post<User>('/api/v1/users', user);
    return response.data;
  }

  public async registerUser(user: UserRegister): Promise<UserRegister> {
    const response = await apiClient.post<UserRegister>('/api/v1/users/register', user);
    return response.data;
  }

  public async updateUser(id: string, user: Partial<User>): Promise<User> {
    const response = await apiClient.put<User>(`/api/v1/users/${id}`, user);
    return response.data;
  }

  public async deleteUser(id: string): Promise<void> {
    await apiClient.delete(`/api/v1/users/${id}`);
  }

  public async setDefaultProject(id: string, projectID: string): Promise<User> {
    const response = await apiClient.patch<User>(`/api/v1/users/${id}/project`, { projectID });
    return response.data;
  }

  public async changePassword(id: string, request: ChangePasswordRequest): Promise<void> {
    await apiClient.patch(`/api/v1/users/${id}/password`, request);
  }

  public async changeImage(id: string, image: File): Promise<void> {
    const formData = new FormData();
    formData.append('image', image);
    await apiClient.patch(`/api/v1/users/${id}/image`, formData, {
      headers: {
        'Content-Type': undefined,
      },
    });
  }
}

export const userService = new UserService();
