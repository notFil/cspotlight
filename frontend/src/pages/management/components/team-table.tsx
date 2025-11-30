import { useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Edit, Trash2, Plus } from "lucide-react";
import type { Team } from "@/types";
import { TeamModal } from "./team-modal";
import { DeleteConfirmModal } from "./delete-confirmation-modal";
import { useTeams, useCreateTeam, useUpdateTeam, useDeleteTeam } from "@/hooks/use-teams";

export const TeamTable = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [selectedTeam, setSelectedTeam] = useState<Team | undefined>(undefined);

  const { data: teams = [], isLoading, isError, error } = useTeams();
  const { mutateAsync: createTeam } = useCreateTeam();
  const { mutateAsync: updateTeam } = useUpdateTeam();
  const { mutateAsync: deleteTeam } = useDeleteTeam();

  const handleAdd = () => {
    setSelectedTeam(undefined);
    setIsModalOpen(true);
  };

  const handleEdit = (team: Team) => {
    setSelectedTeam(team);
    setIsModalOpen(true);
  };

  const handleDeleteClick = (team: Team) => {
    setSelectedTeam(team);
    setIsDeleteModalOpen(true);
  };

  const handleSave = async (team: Team) => {
    if (selectedTeam) {
      // Edit
      await updateTeam({ id: selectedTeam.id, data: team });
    } else {
      // Add
      await createTeam(team);
    }
    setIsModalOpen(false);
  };

  const handleDeleteConfirm = async () => {
    if (selectedTeam) {
      await deleteTeam(selectedTeam.id);
      setIsDeleteModalOpen(false);
    }
  };

  if (isLoading) {
    return <div>Loading teams...</div>;
  }

  if (isError) {
    return <div>Error loading teams: {error instanceof Error ? error.message : 'Unknown error'}</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-semibold">Manage Teams</h2>
        <Button onClick={handleAdd}>
          <Plus className="mr-2 h-4 w-4" /> Add Team
        </Button>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>Description</TableHead>
            <TableHead>Last Updated</TableHead>
            <TableHead className="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {teams.length === 0 ? (
            <TableRow>
              <TableCell colSpan={5} className="text-center h-24 text-muted-foreground">
                No team available in this instance
              </TableCell>
            </TableRow>
          ) : teams.map((team) => (
            <TableRow key={team.id}>
              <TableCell className="font-medium">{team.name}</TableCell>
              <TableCell>{team.description}</TableCell>
              <TableCell>{team.updatedAt}</TableCell>
              <TableCell className="text-right">
                <Button variant="ghost" size="icon" onClick={() => handleEdit(team)}>
                  <Edit className="h-4 w-4" />
                </Button>
                <Button variant="ghost" size="icon" onClick={() => handleDeleteClick(team)}>
                  <Trash2 className="h-4 w-4 text-destructive" />
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>

      <TeamModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSave={handleSave}
        team={selectedTeam}
      />

      <DeleteConfirmModal
        isOpen={isDeleteModalOpen}
        onClose={() => setIsDeleteModalOpen(false)}
        onConfirm={handleDeleteConfirm}
        title="Delete Team"
        description="Are you sure you want to delete this team? This action cannot be undone."
      />
    </div>
  );
};
