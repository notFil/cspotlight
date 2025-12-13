import { useNavigate } from "react-router-dom";
import { Activity, Star, Copy } from "lucide-react";
import { LoadingPage } from "@/components/common/loading-page";
import { ErrorPage } from "@/components/common/error-page";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { useProjects } from "@/hooks/use-projects";
import { useSetDefaultProject } from "@/hooks/use-users";
import { useAuth } from "@/hooks/use-auth";
import { timeAgo } from "@/utils/date";
import { toast } from "sonner";

const Projects = () => {
  const navigate = useNavigate();
  const { data: projects, isLoading, error } = useProjects();
  const { user } = useAuth();
  const setDefaultProject = useSetDefaultProject();

  if (isLoading) {
    return <LoadingPage />;
  }

  if (error) {
    return <ErrorPage error={error as Error} />;
  }

  const handleSetDefault = (projectId: string) => {
    if (user?.id) {
      setDefaultProject.mutate({ id: user.id, data: projectId });
    }
  };

  const handleCopyUrl = (url: string) => {
    navigator.clipboard.writeText(url);
    toast.success("Reporting URL copied to clipboard");
  };

  return (
    <div className="p-6 min-h-screen" style={{ backgroundColor: 'var(--color-background)' }}>
      <div className="main-content max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold mb-2" style={{ color: 'var(--color-text)' }}>Projects</h1>
            <p className="text-sm" style={{ color: 'var(--color-text-secondary)' }}>
              Manage your projects and view CSP violation reports
            </p>
          </div>
        </div>

        {/* Projects Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {projects?.map((project) => {
            const isDefault = user?.defaultProjectId === project.id;
            return (
              <Card key={project.id} className="hover:shadow-lg transition-shadow duration-200 flex flex-col">
                <CardHeader>
                  <div className="flex justify-between items-start">
                    <CardTitle className="text-xl flex items-center gap-2">
                      {project.name}
                    </CardTitle>
                    <span
                      className={`px-2 py-1 rounded-full text-xs font-medium ${!project.disabled
                        ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-100"
                        : "bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-100"
                        }`}
                    >
                      {!project.disabled ? "active" : "inactive"}
                    </span>
                  </div>
                  <CardDescription>{project.description}</CardDescription>
                </CardHeader>
                <CardContent className="flex-grow">
                  <div className="space-y-4">
                    <div className="text-sm">
                      <p className="text-muted-foreground mb-1">Reporting URL</p>
                      <div className="flex items-center gap-2">
                        <code className="bg-muted px-2 py-1 rounded text-xs block truncate flex-1">
                          {project.reportingUrl || "N/A"}
                        </code>
                        {project.reportingUrl && (
                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-6 w-6"
                            onClick={() => project.reportingUrl && handleCopyUrl(project.reportingUrl)}
                          >
                            <Copy className="h-3 w-3" />
                          </Button>
                        )}
                      </div>
                    </div>
                    <div className="flex items-center text-xs text-muted-foreground">
                      <Activity className="mr-1 h-3 w-3" />
                      Last active: {timeAgo(project.lastActive?.toString() || "") || "Never"}
                    </div>
                  </div>
                </CardContent>
                <CardFooter className="flex gap-2">
                  <Button
                    className="flex-1"
                    variant="outline"
                    onClick={() => navigate(`/reports/${project.id}`)}
                  >
                    View Details
                  </Button>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="icon"
                        disabled={isDefault || setDefaultProject.isPending}
                        onClick={() => project.id && handleSetDefault(project.id)}
                        className={isDefault ? "opacity-100" : ""}
                      >
                        <Star
                          className={`h-4 w-4 ${isDefault ? "fill-yellow-500 text-yellow-500" : ""
                            }`}
                        />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>
                      <p>{isDefault ? "Default Project" : "Set as default"}</p>
                    </TooltipContent>
                  </Tooltip>
                </CardFooter>
              </Card>
            );
          })}
          {projects?.length === 0 && (
            <div className="col-span-full text-center py-12 text-muted-foreground">
              No projects found.
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Projects;