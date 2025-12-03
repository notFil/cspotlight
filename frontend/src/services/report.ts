import { apiClient } from './api';
import type { CSPReport, APIResponse, ReportGraphDataDTO } from '@/types';

class ReportService {
  public async getReports(projectId: string, page: number = 1, pageSize: number = 10): Promise<APIResponse<CSPReport[]>> {
    return await apiClient.get<CSPReport[]>(`/api/v1/reports/${projectId}/csp?page=${page}&page_size=${pageSize}`);
  }

  public async getReportGraphData(projectId: string): Promise<APIResponse<ReportGraphDataDTO[]>> {
    return await apiClient.get<ReportGraphDataDTO[]>(`/api/v1/analytics/${projectId}/graph-data`);
  }
}

export const reportService = new ReportService();
