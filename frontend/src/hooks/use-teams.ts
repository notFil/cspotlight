import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { teamService } from '@/services/team';
import type { Team } from '@/types';

export function useTeams() {
  return useQuery({
    queryKey: ['teams'],
    queryFn: () => teamService.getTeams(),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

export function useCreateTeam() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (teamData: Omit<Team, 'id'>) => teamService.createTeam(teamData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['teams'] });
    },
  });
}

export function useUpdateTeam() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Team> }) => teamService.updateTeam(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['teams'] });
    },
  });
}

export function useDeleteTeam() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => teamService.deleteTeam(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['teams'] });
    },
  });
}