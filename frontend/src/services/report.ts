import { apiClient } from './api';
import type { CSPReport } from '@/types';

class ReportService {
  public async getReports(projectId: string): Promise<CSPReport[]> {
    const response = await apiClient.get<CSPReport[]>(`/api/v1/reports/${projectId}/csp`);
    return response.data;
  }
}

export const reportService = new ReportService();
