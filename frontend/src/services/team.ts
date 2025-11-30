import { apiClient } from './api';
import type { Team } from '@/types';

class TeamService {
  public async getTeams(): Promise<Team[]> {
    const response = await apiClient.get<Team[]>('/api/v1/teams');
    return response.data;
  }

  public async createTeam(team: Omit<Team, 'id'>): Promise<Team> {
    const response = await apiClient.post<Team>('/api/v1/teams', team);
    return response.data;
  }

  public async updateTeam(id: string, team: Partial<Team>): Promise<Team> {
    const response = await apiClient.put<Team>(`/api/v1/teams/${id}`, team);
    return response.data;
  }

  public async deleteTeam(id: string): Promise<void> {
    await apiClient.delete(`/api/v1/teams/${id}`);
  }
}

export const teamService = new TeamService();
