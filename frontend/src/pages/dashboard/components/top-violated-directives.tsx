import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useState } from 'react';
import { useReportTopViolatedDirectives } from "@/hooks/use-reports";
import { LoadingPage } from "@/components/common/loading-page";
import { ErrorPage } from "@/components/common/error-page";

const TopViolatedDirectives = ({ projectId }: { projectId: string }) => {
  const [showPercentages, setShowPercentages] = useState(true);

  const { data, isLoading, error } = useReportTopViolatedDirectives(projectId);

  if (isLoading) {
    return <LoadingPage className="h-48" />;
  }

  if (error || !data) {
    return <ErrorPage error={error as Error} className="h-48" />;
  }

  const topViolatedDirectives = data?.data;

  const getWidth = (violations: number) => {
    const max = topViolatedDirectives?.totalViolations;
    return `${(violations / max) * 100}%`;
  };

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-lg font-semibold">Top Violated Directives</CardTitle>
        <label className="flex items-center gap-2 cursor-pointer">
          <span className="text-xs text-muted-foreground">Show percentages</span>
          <input
            type="checkbox"
            checked={showPercentages}
            onChange={(e) => setShowPercentages(e.target.checked)}
            className="accent-primary h-4 w-4"
          />
        </label>
      </CardHeader>
      <CardContent>
        <div className="space-y-4 mt-4">
          {topViolatedDirectives?.violations.map((directive) => (
            <div key={directive.directive} className="group cursor-pointer relative" data-violations={directive.count} data-percentage={directive.percentage}>
              <div className="flex justify-between items-center mb-1">
                <span className="text-sm font-medium text-foreground">{directive.directive}</span>
                <span className="text-sm directive-value text-muted-foreground">
                  {showPercentages ? `${directive.percentage.toFixed(2)}% of total` : `${directive.count} violations`}
                </span>
              </div>
              <div className="w-full bg-secondary rounded-full h-2 overflow-hidden">
                <div className="h-2 bg-primary rounded-full transition-all duration-500 ease-out" style={{ width: getWidth(directive.count) }} title={`${directive.directive}: ${directive.count} violations`}></div>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
};

export default TopViolatedDirectives;
