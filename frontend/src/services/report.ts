import { apiClient } from './api';
import type { CSPReport, APIResponse, ReportGraphDataDTO, ReportMetricsDTO, ReportViolationTrendDTO } from '@/types';

class ReportService {
  public async getReports(projectId: string, page: number = 1, pageSize: number = 10): Promise<APIResponse<CSPReport[]>> {
    return await apiClient.get<CSPReport[]>(`/api/v1/reports/${projectId}/csp?page=${page}&page_size=${pageSize}`);
  }

  public async getReportSummaryStats(projectId: string): Promise<APIResponse<ReportMetricsDTO>> {
    return await apiClient.get<ReportMetricsDTO>(`/api/v1/analytics/${projectId}/summary-stats`);
  }

  public async getReportGraphData(projectId: string): Promise<APIResponse<ReportGraphDataDTO[]>> {
    return await apiClient.get<ReportGraphDataDTO[]>(`/api/v1/analytics/${projectId}/graph-data`);
  }

  public async getReportViolationTrend(projectId: string): Promise<APIResponse<ReportViolationTrendDTO>> {
    return await apiClient.get<ReportViolationTrendDTO>(`/api/v1/analytics/${projectId}/violation-trend`);
  }
}

export const reportService = new ReportService();
