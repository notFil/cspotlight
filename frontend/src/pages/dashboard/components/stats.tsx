import { Card, CardContent } from "@/components/ui/card";
import { useReportSummaryStats } from "@/hooks/use-reports";
import { ArrowUp, ArrowDown } from "lucide-react";
import { LoadingPage } from "@/components/loading-page";
import { ErrorPage } from "@/components/error-page";

const Stats = ({ projectId }: { projectId: string }) => {
    const { data: statsData, isLoading, error } = useReportSummaryStats(projectId);

    if (isLoading) {
        return <LoadingPage className="h-48" />;
    }

    if (error || !statsData?.data) {
        return <ErrorPage error={error as Error} className="h-48" />;
    }

    const stats = statsData.data;

    const renderChange = (change: number) => {
        const isPositive = change > 0;
        const isNeutral = change === 0;
        const colorClass = isPositive ? "text-red-600 bg-red-50 dark:bg-red-900/20 dark:text-red-400" :
            (isNeutral ? "text-gray-600 bg-gray-50 dark:bg-gray-800 dark:text-gray-400" : "text-emerald-600 bg-emerald-50 dark:bg-emerald-900/20 dark:text-emerald-400");
        const Icon = isPositive ? ArrowUp : ArrowDown;

        return (
            <div className={`mt-2 flex items-center text-xs font-medium w-fit px-2 py-0.5 rounded-full ${colorClass}`}>
                {!isNeutral && <Icon className="w-3 h-3 mr-1" />}
                <span>{Math.abs(change).toFixed(2)}%</span>
            </div>
        );
    };

    return (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-6">
            <Card>
                <CardContent className="p-6">
                    <div className="text-2xl font-bold">{stats.totalViolations.value}</div>
                    <div className="text-sm font-medium text-muted-foreground mt-1">Total Violations (24h)</div>
                    {renderChange(stats.totalViolations.change)}
                </CardContent>
            </Card>
            <Card>
                <CardContent className="p-6">
                    <div className="text-2xl font-bold">{stats.totalCriticalViolations.value}</div>
                    <div className="text-sm font-medium text-muted-foreground mt-1">Critical Violations</div>
                    {renderChange(stats.totalCriticalViolations.change)}
                </CardContent>
            </Card>
            <Card>
                <CardContent className="p-6">
                    <div className="text-2xl font-bold">{stats.affectedDomains.value}</div>
                    <div className="text-sm font-medium text-muted-foreground mt-1">Affected Domains</div>
                    {renderChange(stats.affectedDomains.change)}
                </CardContent>
            </Card>
            <Card>
                <CardContent className="p-6">
                    <div className="text-2xl font-bold">{stats.policyEnforcement.value}%</div>
                    <div className="text-sm font-medium text-muted-foreground mt-1">Policy Enforcement Status</div>
                    {renderChange(stats.policyEnforcement.change)}
                </CardContent>
            </Card>
        </div>
    );
};

export default Stats;
