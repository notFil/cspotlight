import { apiClient } from './api';
import type { CSPReport, APIResponse, ReportGraphData, ReportStatsMetrics, ReportViolationTrends, ReportBrowserOSViolations, ReportTopViolatedDirectives, ReportTopViolatedDocumentURLs, ReportTopViolationSources, ReportFilters } from '@/types';

class ReportService {
  public async getReports(projectId: string, page: number = 1, pageSize: number = 10, filters?: ReportFilters): Promise<APIResponse<CSPReport[]>> {
    const params = new URLSearchParams({
      page: page.toString(),
      page_size: pageSize.toString(),
    });

    if (filters) {
      if (filters.directive) params.append('directive', filters.directive);
      if (filters.disposition) params.append('disposition', filters.disposition);
      if (filters.blockedURL) params.append('blocked_url', filters.blockedURL);
      if (filters.userAgent) params.append('user_agent', filters.userAgent);
      if (filters.documentURL) params.append('document_url', filters.documentURL);
    }

    return await apiClient.get<CSPReport[]>(`/api/v1/reports/${projectId}?${params.toString()}`);
  }

  public async getReportSummaryStats(projectId: string): Promise<APIResponse<ReportStatsMetrics>> {
    return await apiClient.get<ReportStatsMetrics>(`/api/v1/analytics/${projectId}/summary-stats`);
  }

  public async getReportGraphData(projectId: string): Promise<APIResponse<ReportGraphData[]>> {
    return await apiClient.get<ReportGraphData[]>(`/api/v1/analytics/${projectId}/graph-data`);
  }

  public async getReportViolationTrend(projectId: string): Promise<APIResponse<ReportViolationTrends>> {
    return await apiClient.get<ReportViolationTrends>(`/api/v1/analytics/${projectId}/violation-trend`);
  }

  public async getReportTopViolatedDirectives(projectId: string): Promise<APIResponse<ReportTopViolatedDirectives>> {
    return await apiClient.get<ReportTopViolatedDirectives>(`/api/v1/analytics/${projectId}/violated-directives`);
  }

  public async getReportTopViolatedDocumentURLs(projectId: string): Promise<APIResponse<ReportTopViolatedDocumentURLs>> {
    return await apiClient.get<ReportTopViolatedDocumentURLs>(`/api/v1/analytics/${projectId}/violated-document-urls`);
  }

  public async getReportBrowserOSViolation(projectId: string): Promise<APIResponse<ReportBrowserOSViolations>> {
    return await apiClient.get<ReportBrowserOSViolations>(`/api/v1/analytics/${projectId}/software-stats`);
  }

  public async getReportTopViolationSources(projectId: string): Promise<APIResponse<ReportTopViolationSources>> {
    return await apiClient.get<ReportTopViolationSources>(`/api/v1/analytics/${projectId}/violation-sources`);
  }
}

export const reportService = new ReportService();
