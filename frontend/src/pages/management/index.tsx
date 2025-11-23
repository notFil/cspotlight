import { useState } from 'react';
import { Button } from "@/components/ui/button";
import { UserTable } from './components/user-table';
import { ProjectTable } from './components/project-table';
import { TeamTable } from './components/team-table';

const Management = () => {
  const [activeTab, setActiveTab] = useState<'users' | 'projects' | 'teams'>('users');

  return (
    <div className="p-6 bg-background">
      <div className="main-content max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-2" style={{ color: 'var(--color-text)' }}>Management</h1>
          <p className="text-sm" style={{ color: 'var(--color-text-secondary)' }}>Manage your accounts and projects</p>
        </div>

        {/* Tabs */}
        <div className="flex space-x-2 mb-6 border-b pb-2">
          <Button
            variant={activeTab === 'users' ? 'default' : 'ghost'}
            onClick={() => setActiveTab('users')}
          >
            Manage Users
          </Button>
          <Button
            variant={activeTab === 'projects' ? 'default' : 'ghost'}
            onClick={() => setActiveTab('projects')}
          >
            Manage Projects
          </Button>
          <Button
            variant={activeTab === 'teams' ? 'default' : 'ghost'}
            onClick={() => setActiveTab('teams')}
          >
            Manage Teams
          </Button>
        </div>

        {/* Content */}
        <div>
          {activeTab === 'users' && <UserTable />}
          {activeTab === 'projects' && <ProjectTable />}
          {activeTab === 'teams' && <TeamTable />}
        </div>

      </div>
    </div>
  );
};

export default Management;