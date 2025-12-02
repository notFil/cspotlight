import { apiClient } from "./api";
import type { AuthTokenData, UserRegister, UserLogin } from "@/types";

class AuthService {
  constructor() {
    apiClient.setRefreshAuthTokenMethod(this.refreshAuthToken.bind(this));
  }

  public async login(credentials: UserLogin): Promise<AuthTokenData> {
    const response = await apiClient.post<AuthTokenData>('/api/v1/auth/login', credentials);
    return response.data;
  }

  public async register(user: UserRegister): Promise<UserRegister> {
    const response = await apiClient.post<UserRegister>('/api/v1/auth/register', user);
    return response.data;
  }

  public async refreshAuthToken(refreshToken: string): Promise<AuthTokenData> {
    const response = await apiClient.post<AuthTokenData>(
      '/api/v1/auth/refresh',
      { refreshToken }
    );
    return response.data;
  }

  public async logout(): Promise<void> {
    await apiClient.post('/api/v1/signout', {});
  }
}

export const authService = new AuthService();
