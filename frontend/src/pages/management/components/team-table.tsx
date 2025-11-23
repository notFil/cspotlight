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

// Mock data
const MOCK_TEAMS: Team[] = [
  {
    id: '1',
    name: 'Team Alpha',
    description: 'Primary development team',
    last_updated: '2023-10-26',
  },
  {
    id: '2',
    name: 'Team Beta',
    description: 'Support and maintenance',
    last_updated: '2023-10-25',
  },
];

export const TeamTable = () => {
  const [teams, setTeams] = useState<Team[]>(MOCK_TEAMS);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [selectedTeam, setSelectedTeam] = useState<Team | undefined>(undefined);

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

  const handleSave = (team: Team) => {
    if (selectedTeam) {
      // Edit
      setTeams(teams.map(t => t.id === team.id ? team : t));
    } else {
      // Add
      setTeams([...teams, { ...team, id: Math.random().toString(36).substr(2, 9), last_updated: new Date().toISOString().split('T')[0] }]);
    }
    setIsModalOpen(false);
  };

  const handleDeleteConfirm = () => {
    if (selectedTeam) {
      setTeams(teams.filter(t => t.id !== selectedTeam.id));
      setIsDeleteModalOpen(false);
    }
  };

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
          {teams.map((team) => (
            <TableRow key={team.id}>
              <TableCell className="font-medium">{team.name}</TableCell>
              <TableCell>{team.description}</TableCell>
              <TableCell>{team.last_updated}</TableCell>
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
