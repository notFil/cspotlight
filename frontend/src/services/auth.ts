import { apiClient } from "./api";
import type { AuthTokenData } from "@/types";

class AuthService {
  constructor() {
    apiClient.setRefreshAuthTokenMethod(this.refreshAuthToken.bind(this));
  }

  public async login(credentials: any): Promise<AuthTokenData> {
    const response = await apiClient.post<AuthTokenData>('/api/v1/auth/login', credentials);
    return response.data;
  }

  public async refreshAuthToken(refreshToken: string): Promise<AuthTokenData> {
    const response = await apiClient.post<AuthTokenData>(
      '/api/v1/auth/refresh',
      { refreshToken }
    );
    return response.data;
  }
}

export const authService = new AuthService();
