import { apiClient } from "./api";

export const login = async (credentials: any): Promise<string> => {
  const response = await apiClient.post<string>('/api/v1/auth/login', credentials);
  return response.data;
};
