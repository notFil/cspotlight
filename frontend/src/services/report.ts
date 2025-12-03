import { apiClient } from './api';
import type { CSPReport, APIResponse } from '@/types';

class ReportService {
  public async getReports(projectId: string, page: number = 1, pageSize: number = 10): Promise<APIResponse<CSPReport[]>> {
    return await apiClient.get<CSPReport[]>(`/api/v1/reports/${projectId}/csp?page=${page}&page_size=${pageSize}`);
  }
}

export const reportService = new ReportService();
