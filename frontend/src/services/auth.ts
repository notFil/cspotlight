import { apiClient } from "./api";
import type { UserRegister, UserLogin } from "@/types";

class AuthService {
  constructor() {
  }

  public async login(credentials: UserLogin): Promise<void> {
    await apiClient.post('/api/v1/auth/login', credentials);
  }

  public async register(user: UserRegister): Promise<any> {
    const response = await apiClient.post<any>('/api/v1/auth/register', user);
    return response.data;
  }

  public async logout(): Promise<void> {
    await apiClient.post('/api/v1/signout', {});
  }
}

export const authService = new AuthService();
