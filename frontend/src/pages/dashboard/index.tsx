import { useAuth } from '@/hooks/use-auth';
import DashboardFallback from '@/pages/dashboard/components/dashboard-fallback';
import Stats from '@/pages/dashboard/components/stats';
import ViolationTrend from '@/pages/dashboard/components/violation-trend';
import TopViolatedDirectives from '@/pages/dashboard/components/top-violated-directives';
import TopViolatedDocumentUrls from '@/pages/dashboard/components/top-violated-document-urls';
import ViolationByBrowserOS from '@/pages/dashboard/components/violation-by-browser-os';
import ViolationSources from '@/pages/dashboard/components/violation-sources';
import ReportsGraph from '@/pages/dashboard/components/reports-graph';

const Dashboard = () => {
  const { user } = useAuth();

  if (!user?.defaultProjectId) {
    return <DashboardFallback />;
  }

  return (
    <div className="p-6" style={{ backgroundColor: 'var(--color-background)' }}>
      <div className="main-content max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-2" style={{ color: 'var(--color-text)' }}>Dashboard</h1>
          <p className="text-sm" style={{ color: 'var(--color-text-secondary)' }}>Real-time Content Security Policy monitoring and analytics</p>
        </div>

        {/* Top Stats Row */}
        <Stats projectId={user.defaultProjectId} />

        {/* Reports Graph */}
        <ReportsGraph projectId={user.defaultProjectId} />

        {/* Main Content Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
          <ViolationTrend projectId={user.defaultProjectId} />
          <TopViolatedDirectives projectId={user.defaultProjectId} />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
          <TopViolatedDocumentUrls projectId={user.defaultProjectId} />
          <ViolationByBrowserOS projectId={user.defaultProjectId} />
        </div>

        {/* Bottom Row */}
        <ViolationSources projectId={user.defaultProjectId} />
      </div>
    </div>
  );
};

export default Dashboard;