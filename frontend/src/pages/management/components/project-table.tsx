import { useState } from 'react';
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
import type { Project } from "@/types";
import { ProjectModal } from "./project-modal";
import { DeleteConfirmModal } from "./delete-confirmation-modal";
import { useProjects, useCreateProject, useUpdateProject, useDeleteProject } from "@/hooks/use-projects";
import { LoadingPage } from "@/components/loading-page";
import { ErrorPage } from "@/components/error-page";


export const ProjectTable = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [selectedProject, setSelectedProject] = useState<Project | undefined>(undefined);

  const { data: projects = [], isLoading, isError, error } = useProjects();
  const { mutateAsync: createProject } = useCreateProject();
  const { mutateAsync: updateProject } = useUpdateProject();
  const { mutateAsync: deleteProject } = useDeleteProject();

  const handleAdd = () => {
    setSelectedProject(undefined);
    setIsModalOpen(true);
  };

  const handleEdit = (project: Project) => {
    setSelectedProject(project);
    setIsModalOpen(true);
  };

  const handleDeleteClick = (project: Project) => {
    setSelectedProject(project);
    setIsDeleteModalOpen(true);
  };

  const handleSave = async (project: Project) => {
    if (selectedProject?.id) {
      // Edit
      await updateProject({ id: selectedProject.id, data: project });
    } else {
      // Add
      await createProject(project);
    }
    setIsModalOpen(false);
  };

  const handleDeleteConfirm = async () => {
    if (selectedProject && selectedProject.id) {
      await deleteProject(selectedProject.id);
      setIsDeleteModalOpen(false);
    }
  };

  if (isLoading) {
    return <LoadingPage />;
  }

  if (isError) {
    return <ErrorPage error={error as Error} />;
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-semibold">Manage Projects</h2>
        <Button onClick={handleAdd}>
          <Plus className="mr-2 h-4 w-4" /> Add Project
        </Button>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>Team</TableHead>
            <TableHead>Reporting URL</TableHead>
            <TableHead>Last Updated</TableHead>
            <TableHead className="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {projects.length === 0 ? (
            <TableRow>
              <TableCell colSpan={5} className="text-center h-24 text-muted-foreground">
                No project available in this instance
              </TableCell>
            </TableRow>
          ) : (
            projects.map((project) => (
              <TableRow key={project.id}>
                <TableCell className="font-medium">{project.name}</TableCell>
                <TableCell>{project.teamName}</TableCell>
                <TableCell className="max-w-xs truncate">{project.reportingUrl}</TableCell>
                <TableCell>{project.lastActive || 'No recent activity'}</TableCell>
                <TableCell className="text-right">
                  <Button variant="ghost" size="icon" onClick={() => handleEdit(project)}>
                    <Edit className="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="icon" onClick={() => handleDeleteClick(project)}>
                    <Trash2 className="h-4 w-4 text-destructive" />
                  </Button>
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>

      <ProjectModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSave={handleSave}
        project={selectedProject}
      />

      <DeleteConfirmModal
        isOpen={isDeleteModalOpen}
        onClose={() => setIsDeleteModalOpen(false)}
        onConfirm={handleDeleteConfirm}
        title="Delete Project"
        description="Are you sure you want to delete this project? This action cannot be undone."
      />
    </div>
  );
};
