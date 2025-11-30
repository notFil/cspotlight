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
import type { User } from "@/types";
import { UserModal } from "./user-modal";
import { DeleteConfirmModal } from "./delete-confirmation-modal";
import { useUsers, useCreateUser, useUpdateUser, useDeleteUser } from "@/hooks/use-users";
import { useCreateTeam } from "@/hooks/use-teams";

export const UserTable = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | undefined>(undefined);

  // React Query hooks
  const { data: users = [], isLoading, isError, error } = useUsers();
  const { mutateAsync: createUser } = useCreateUser();
  const { mutateAsync: updateUser } = useUpdateUser();
  const { mutateAsync: deleteUser } = useDeleteUser();

  const handleAdd = () => {
    setSelectedUser(undefined);
    setIsModalOpen(true);
  };

  const handleEdit = (user: User) => {
    setSelectedUser(user);
    setIsModalOpen(true);
  };

  const handleDeleteClick = (user: User) => {
    setSelectedUser(user);
    setIsDeleteModalOpen(true);
  };

  const handleSave = async (user: User) => {
    try {
      if (selectedUser) {
        // Edit - send PUT request
        await updateUser({ id: user.id, data: user });
      } else {
        // Add - send POST request
        await createUser(user);
      }
      setIsModalOpen(false);
    } catch (error) {
      console.error("Failed to save user:", error);
    }
  };

  const handleDeleteConfirm = async () => {
    try {
      if (selectedUser) {
        // Send DELETE request
        await deleteUser(selectedUser.id);
        setIsDeleteModalOpen(false);
      }
    } catch (error) {
      console.error("Failed to delete user:", error);
    }
  };

  if (isLoading) {
    return <div>Loading users...</div>;
  }

  if (isError) {
    return <div>Error loading users: {error instanceof Error ? error.message : 'Unknown error'}</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-semibold">Manage Users</h2>
        <Button onClick={handleAdd}>
          <Plus className="mr-2 h-4 w-4" /> Add User
        </Button>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>First Name</TableHead>
            <TableHead>Last Name</TableHead>
            <TableHead>Role</TableHead>
            <TableHead>Team</TableHead>
            <TableHead>Last Updated</TableHead>
            <TableHead className="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {users.length === 0 ? (
            <TableRow>
              <TableCell colSpan={5} className="text-center h-24 text-muted-foreground">
                No user available in this instance
              </TableCell>
            </TableRow>
          ) : users.map((user) => (
            <TableRow key={user.id}>
              <TableCell>{user.firstName}</TableCell>
              <TableCell>{user.lastName}</TableCell>
              <TableCell className="capitalize">{user.role}</TableCell>
              <TableCell>{user.teamName || '-'}</TableCell>
              <TableCell>{user.updatedAt || '-'}</TableCell>
              <TableCell className="text-right">
                <Button variant="ghost" size="icon" onClick={() => handleEdit(user)}>
                  <Edit className="h-4 w-4" />
                </Button>
                <Button variant="ghost" size="icon" onClick={() => handleDeleteClick(user)}>
                  <Trash2 className="h-4 w-4 text-destructive" />
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>

      <UserModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSave={handleSave}
        user={selectedUser}
      />

      <DeleteConfirmModal
        isOpen={isDeleteModalOpen}
        onClose={() => setIsDeleteModalOpen(false)}
        onConfirm={handleDeleteConfirm}
        title="Delete User"
        description="Are you sure you want to delete this user? This action cannot be undone."
      />
    </div>
  );
};
