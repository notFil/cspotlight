import axios from 'axios';
import type { AxiosInstance, AxiosRequestConfig } from 'axios';
import type { APIResponse } from '@/types';

class APIClient {
    private client: AxiosInstance;

    constructor() {
        this.client = axios.create({
            baseURL: import.meta.env.DEV ? '' : import.meta.env.VITE_API_URL,
            timeout: 5000,
            withCredentials: true,
            headers: {
                'Content-Type': 'application/json',
            },
        });
    }

    public async get<T>(url: string, config?: AxiosRequestConfig): Promise<APIResponse<T>> {
        const response = await this.client.get<APIResponse<T>>(url, config);
        return response.data;
    }

    public async post<T>(url: string, data: any, config?: AxiosRequestConfig): Promise<APIResponse<T>> {
        const response = await this.client.post<APIResponse<T>>(url, data, config);
        return response.data;
    }

    public async put<T>(url: string, data: any, config?: AxiosRequestConfig): Promise<APIResponse<T>> {
        const response = await this.client.put<APIResponse<T>>(url, data, config);
        return response.data;
    }

    public async patch<T>(url: string, data: any, config?: AxiosRequestConfig): Promise<APIResponse<T>> {
        const response = await this.client.patch<APIResponse<T>>(url, data, config);
        return response.data;
    }

    public async delete<T>(url: string, config?: AxiosRequestConfig): Promise<APIResponse<T>> {
        const response = await this.client.delete<APIResponse<T>>(url, config);
        return response.data;
    }
}

export const apiClient = new APIClient();