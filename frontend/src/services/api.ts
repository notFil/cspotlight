import axios from 'axios';
import type { AxiosInstance, AxiosResponse } from 'axios';
import type { APIResponse } from '@/types';

class APIClient {
    private client: AxiosInstance;

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
            (error) => {
                if (error.response && error.response.status === 401) {
                    localStorage.removeItem('authToken');
                    window.location.href = '/login';
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