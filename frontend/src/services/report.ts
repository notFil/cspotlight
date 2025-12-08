import { apiClient } from './api';
import type { CSPReport, APIResponse, ReportGraphDataDTO, ReportMetricsDTO, ReportViolationTrendDTO, ReportBrowserOSViolationDTO } from '@/types';

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

  public async getReportBrowserOSViolation(projectId: string): Promise<APIResponse<ReportBrowserOSViolationDTO>> {
    return await apiClient.get<ReportBrowserOSViolationDTO>(`/api/v1/analytics/${projectId}/software-stats`);
  }
}

export const reportService = new ReportService();
