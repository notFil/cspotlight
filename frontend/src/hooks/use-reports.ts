import { useQuery } from '@tanstack/react-query';
import { reportService } from '@/services/report';

export function useReports(projectId: string, page: number = 1, pageSize: number = 10) {
  return useQuery({
    queryKey: ['reports', projectId, page, pageSize],
    queryFn: () => reportService.getReports(projectId, page, pageSize),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

