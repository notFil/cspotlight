import { useNavigate } from 'react-router-dom';
import { LayoutDashboard, ArrowRight } from 'lucide-react';
import { Button } from '@/components/ui/button';

const DashboardFallback = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-[80vh] flex items-center justify-center p-6" style={{ backgroundColor: 'var(--color-background)' }}>
      <div className="max-w-2xl w-full flex flex-col items-center text-center p-12">
        <div className="w-24 h-24 bg-primary/5 rounded-full flex items-center justify-center mb-8">
          <LayoutDashboard className="w-12 h-12 text-primary/80" />
        </div>

        <h1 className="text-4xl font-bold mb-4 tracking-tight" style={{ color: 'var(--color-text)' }}>
          Welcome to cspotlight
        </h1>

        <p className="text-xl text-muted-foreground mb-10 max-w-lg font-light leading-relaxed">
          To start viewing analytics and monitoring your Content Security Policy, please select a default project.
        </p>

        <div className="flex flex-col sm:flex-row gap-4 w-full justify-center">
          <Button
            className="group px-6 rounded-full shadow-sm hover:shadow-md transition-all duration-300"
            onClick={() => navigate('/projects')}
          >
            Go to Projects
            <ArrowRight className="ml-2 w-4 h-4 transition-transform group-hover:translate-x-1" />
          </Button>
        </div>

        <p className="mt-12 text-sm text-muted-foreground/50">
          Select the star icon on any project to set it as your default.
        </p>
      </div>
    </div>
  );
};

export default DashboardFallback;
