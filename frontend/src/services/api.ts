import axios from 'axios';
import type { AxiosInstance, AxiosResponse } from 'axios';
import type { APIResponse, AuthTokenData } from '@/types';

class APIClient {
    private client: AxiosInstance;

    private isRefreshing = false;
    private failedQueue: any[] = [];
    private refreshAuthToken: ((refreshToken: string) => Promise<AuthTokenData>) | null = null;

    constructor() {
        this.client = axios.create({
            baseURL: import.meta.env.VITE_API_URL,
            timeout: 5000,
            headers: {
                'Content-Type': 'application/json',
            },
        });

        this.setupInterceptors();
    }

    public setRefreshAuthTokenMethod(method: (refreshToken: string) => Promise<AuthTokenData>) {
        this.refreshAuthToken = method;
    }

    private processQueue(error: any, token: string | null = null) {
        this.failedQueue.forEach(prom => {
            if (error) {
                prom.reject(error);
            } else {
                prom.resolve(token);
            }
        });

        this.failedQueue = [];
    }

    private setupInterceptors() {
        this.client.interceptors.request.use(
            (config) => {
                const token = localStorage.getItem('authToken');
                if (token) {
                    config.headers['Authorization'] = `Bearer ${token}`;
                }
                return config;
            },
            (error) => Promise.reject(error)
        );

        this.client.interceptors.response.use(
            (response: AxiosResponse) => response,
            async (error) => {
                const originalRequest = error.config;

                // Prevent infinite loops if the refresh endpoint itself returns 401
                if (originalRequest.url.includes('/auth/refresh')) {
                    return Promise.reject(error);
                }

                if (error.response && error.response.status === 401 && !originalRequest._retry) {
                    if (this.isRefreshing) {
                        return new Promise<string>((resolve, reject) => {
                            this.failedQueue.push({ resolve, reject });
                        }).then(token => {
                            originalRequest.headers['Authorization'] = 'Bearer ' + token;
                            return this.client(originalRequest);
                        }).catch(err => {
                            return Promise.reject(err);
                        });
                    }

                    originalRequest._retry = true;
                    this.isRefreshing = true;

                    const refreshToken = localStorage.getItem('refreshToken');

                    if (!refreshToken || !this.refreshAuthToken) {
                        // No refresh token or no refresh method registered, logout
                        this.isRefreshing = false;
                        localStorage.removeItem('authToken');
                        localStorage.removeItem('refreshToken');
                        window.location.href = '/login';
                        return Promise.reject(error);
                    }

                    try {
                        // Use the injected auth service to refresh the token
                        const { accessToken, refreshToken: newRefreshToken } = await this.refreshAuthToken(refreshToken);

                        localStorage.setItem('authToken', accessToken);
                        localStorage.setItem('refreshToken', newRefreshToken);

                        this.client.defaults.headers.common['Authorization'] = 'Bearer ' + accessToken;
                        originalRequest.headers['Authorization'] = 'Bearer ' + accessToken;

                        this.processQueue(null, accessToken);
                        this.isRefreshing = false;

                        return this.client(originalRequest);
                    } catch (err) {
                        this.processQueue(err, null);
                        this.isRefreshing = false;
                        localStorage.removeItem('authToken');
                        localStorage.removeItem('refreshToken');
                        window.location.href = '/login';
                        return Promise.reject(err);
                    }
                }

                return Promise.reject(error);
            }
        );
    }

    public async get<T>(url: string): Promise<APIResponse<T>> {
        const response = await this.client.get<APIResponse<T>>(url);
        return response.data;
    }

    public async post<T>(url: string, data: any): Promise<APIResponse<T>> {
        const response = await this.client.post<APIResponse<T>>(url, data);
        return response.data;
    }

    public async put<T>(url: string, data: any): Promise<APIResponse<T>> {
        const response = await this.client.put<APIResponse<T>>(url, data);
        return response.data;
    }

    public async delete<T>(url: string): Promise<APIResponse<T>> {
        const response = await this.client.delete<APIResponse<T>>(url);
        return response.data;
    }
}

export const apiClient = new APIClient();