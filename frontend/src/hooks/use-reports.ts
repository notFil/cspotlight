import { useQuery } from '@tanstack/react-query';
import { reportService } from '@/services/report';

export function useReports(projectId: string) {
  return useQuery({
    queryKey: ['reports', projectId],
    queryFn: () => reportService.getReports(projectId),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

