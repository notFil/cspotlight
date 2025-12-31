import { useState, useEffect } from 'react';
import { Button } from "@/components/ui/button";
import { UserTable } from './components/user-table';
import { ProjectTable } from './components/project-table';
import { TeamTable } from './components/team-table';
import { CSPGuide } from './components/csp-guide';
import { useAuth } from '@/hooks/use-auth';

const Management = () => {
  const { user, loading } = useAuth();
  const isSuperAdmin = user?.role === 'superadmin';

  // Default to projects, but if superadmin, we can let them switch to users/teams.
  // We'll use an effect to ensure non-superadmins don't get stuck on a hidden tab if logic changes.
  const [activeTab, setActiveTab] = useState<'users' | 'projects' | 'teams'>('projects');

  useEffect(() => {
    if (!loading && !isSuperAdmin && (activeTab === 'users' || activeTab === 'teams')) {
      setActiveTab('projects');
    }
  }, [loading, isSuperAdmin, activeTab]);

  if (loading) {
    return <div className="p-6">Loading...</div>;
  }

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
          {isSuperAdmin && (
            <Button
              variant={activeTab === 'users' ? 'default' : 'ghost'}
              onClick={() => setActiveTab('users')}
            >
              Manage Users
            </Button>
          )}
          <Button
            variant={activeTab === 'projects' ? 'default' : 'ghost'}
            onClick={() => setActiveTab('projects')}
          >
            Manage Projects
          </Button>
          {isSuperAdmin && (
            <Button
              variant={activeTab === 'teams' ? 'default' : 'ghost'}
              onClick={() => setActiveTab('teams')}
            >
              Manage Teams
            </Button>
          )}
        </div>

        {/* Content */}
        <div>
          {isSuperAdmin && activeTab === 'users' && <UserTable />}
          {activeTab === 'projects' && (
            <>
              <ProjectTable />
              <CSPGuide />
            </>
          )}
          {isSuperAdmin && activeTab === 'teams' && <TeamTable />}
        </div>

      </div>
    </div>
  );
};

export default Management;