import { useQuery } from '@tanstack/react-query';
import { reportService } from '@/services/report';

export function useReports(projectId: string, page: number = 1, pageSize: number = 10) {
  return useQuery({
    queryKey: ['reports', projectId, page, pageSize],
    queryFn: () => reportService.getReports(projectId, page, pageSize),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

export function useReportSummaryStats(projectId: string) {
  return useQuery({
    queryKey: ['report-summary-stats', projectId],
    queryFn: () => reportService.getReportSummaryStats(projectId),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

export function useReportGraphData(projectId: string) {
  return useQuery({
    queryKey: ['report-graph-data', projectId],
    queryFn: () => reportService.getReportGraphData(projectId),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

export function useReportViolationTrend(projectId: string) {
  return useQuery({
    queryKey: ['report-violation-trend', projectId],
    queryFn: () => reportService.getReportViolationTrend(projectId),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

export function useReportBrowserOSViolation(projectId: string) {
  return useQuery({
    queryKey: ['report-browser-os-violation', projectId],
    queryFn: () => reportService.getReportBrowserOSViolation(projectId),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}



