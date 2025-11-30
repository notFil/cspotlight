import type Dashboard from "../dashboard";

import { useNavigate } from "react-router-dom";
import { Activity, Star } from "lucide-react";
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

interface Project {
  id: string;
  name: string;
  description: string;
  reportingUrl: string;
  status: "active" | "inactive";
  lastActive: string;
  isDefault?: boolean;
}

const MOCK_PROJECTS: Project[] = [
  {
    id: "1",
    name: "Production Environment",
    description: "Main production environment for the e-commerce platform",
    reportingUrl: "https://api.example.com/csp-report",
    status: "active",
    lastActive: "2 mins ago",
    isDefault: true,
  },
  {
    id: "2",
    name: "Staging Environment",
    description: "Staging environment for testing new features",
    reportingUrl: "https://api-staging.example.com/csp-report",
    status: "active",
    lastActive: "1 hour ago",
  },
  {
    id: "3",
    name: "Dev Environment",
    description: "Development environment for internal use",
    reportingUrl: "https://api-dev.example.com/csp-report",
    status: "inactive",
    lastActive: "2 days ago",
  },
];

const Projects = () => {
  const navigate = useNavigate();

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
          {MOCK_PROJECTS.map((project) => (
            <Card key={project.id} className="hover:shadow-lg transition-shadow duration-200 flex flex-col">
              <CardHeader>
                <div className="flex justify-between items-start">
                  <CardTitle className="text-xl flex items-center gap-2">
                    {project.name}

                  </CardTitle>
                  <span
                    className={`px-2 py-1 rounded-full text-xs font-medium ${project.status === "active"
                      ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-100"
                      : "bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-100"
                      }`}
                  >
                    {project.status}
                  </span>
                </div>
                <CardDescription>{project.description}</CardDescription>
              </CardHeader>
              <CardContent className="flex-grow">
                <div className="space-y-4">
                  <div className="text-sm">
                    <p className="text-muted-foreground mb-1">Reporting URL</p>
                    <code className="bg-muted px-2 py-1 rounded text-xs block truncate">
                      {project.reportingUrl}
                    </code>
                  </div>
                  <div className="flex items-center text-xs text-muted-foreground">
                    <Activity className="mr-1 h-3 w-3" />
                    Last active: {project.lastActive}
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
                  <TooltipTrigger>
                    <span tabIndex={0} className="cursor-default">
                      <Button
                        variant="ghost"
                        size="icon"
                        disabled={project.isDefault}
                        onClick={() => console.log(`Set project ${project.id} as default`)}
                        className={project.isDefault ? "opacity-100" : ""}
                      >
                        <Star
                          className={`h-4 w-4 ${project.isDefault ? "fill-yellow-500 text-yellow-500" : ""
                            }`}
                        />
                      </Button>
                    </span>
                  </TooltipTrigger>
                  <TooltipContent>
                    <p>{project.isDefault ? "Default Project" : "Set as default"}</p>
                  </TooltipContent>
                </Tooltip>
              </CardFooter>
            </Card>
          ))}
        </div>
      </div>
    </div>
  );
};

export default Projects;