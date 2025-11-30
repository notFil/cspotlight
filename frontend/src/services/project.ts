import { apiClient } from './api';
import type { Project } from '@/types';

class ProjectService {
  public async getProjects(): Promise<Project[]> {
    const response = await apiClient.get<Project[]>('/api/v1/projects');
    return response.data;
  }

  public async createProject(project: Omit<Project, 'id'>): Promise<Project> {
    const response = await apiClient.post<Project>('/api/v1/projects', project);
    return response.data;
  }

  public async updateProject(id: string, project: Partial<Project>): Promise<Project> {
    const response = await apiClient.put<Project>(`/api/v1/projects/${id}`, project);
    return response.data;
  }

  public async deleteProject(id: string): Promise<void> {
    await apiClient.delete(`/api/v1/projects/${id}`);
  }
}

export const projectService = new ProjectService();
