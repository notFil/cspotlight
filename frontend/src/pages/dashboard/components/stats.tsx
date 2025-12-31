import { Card, CardContent } from "@/components/ui/card";
import { useReportSummaryStats } from "@/hooks/use-reports";
import { ArrowUp, ArrowDown } from "lucide-react";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { LoadingPage } from "@/components/common/loading-page";
import { ErrorPage } from "@/components/common/error-page";

const StatsChange = ({ change, betterPositive = false }: { change: number, betterPositive?: boolean }) => {
    const isNeutral = change === 0;
    const isPositive = change > 0;

    let colorClass = "";
    if (isNeutral) {
        colorClass = "text-gray-600 bg-gray-50 dark:bg-gray-800 dark:text-gray-400";
    } else {
        const isGood = betterPositive ? isPositive : !isPositive;
        colorClass = isGood
            ? "text-emerald-600 bg-emerald-50 dark:bg-emerald-900/20 dark:text-emerald-400"
            : "text-red-600 bg-red-50 dark:bg-red-900/20 dark:text-red-400";
    }

    const Icon = change >= 0 ? ArrowUp : ArrowDown;

    return (
        <Tooltip>
            <TooltipTrigger asChild>
                <div className={`mt-2 flex items-center text-xs font-medium w-fit px-2 py-0.5 rounded-full ${colorClass}`}>
                    {!isNeutral && <Icon className="w-3 h-3 mr-1" />}
                    <span>{Math.abs(change).toFixed(2)}%</span>
                </div>
            </TooltipTrigger>
            <TooltipContent>
                <p>Changes over the last 24 hours</p>
            </TooltipContent>
        </Tooltip>
    );
};

const Stats = ({ projectId }: { projectId: string }) => {
    const { data: statsData, isLoading, error } = useReportSummaryStats(projectId);

    if (isLoading) {
        return <LoadingPage className="h-48" />;
    }

    if (error || !statsData?.data) {
        return <ErrorPage error={error as Error} className="h-48" />;
    }

    const stats = statsData.data;

    return (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-6">
            <Card>
                <CardContent className="p-6">
                    <div className="text-2xl font-bold">{stats.totalViolations.value}</div>
                    <div className="text-sm font-medium text-muted-foreground mt-1">Total Violations</div>
                    <StatsChange change={stats.totalViolations.change} />
                </CardContent>
            </Card>
            <Card>
                <CardContent className="p-6">
                    <div className="text-2xl font-bold">{stats.totalCriticalViolations.value}</div>
                    <div className="text-sm font-medium text-muted-foreground mt-1">Critical Violations</div>
                    <StatsChange change={stats.totalCriticalViolations.change} />
                </CardContent>
            </Card>
            <Card>
                <CardContent className="p-6">
                    <div className="text-2xl font-bold">{stats.affectedDomains.value}</div>
                    <div className="text-sm font-medium text-muted-foreground mt-1">Affected Paths</div>
                    <StatsChange change={stats.affectedDomains.change} />
                </CardContent>
            </Card>
            <Card>
                <CardContent className="p-6">
                    <div className="text-2xl font-bold">{stats.policyEnforcement.value}%</div>
                    <div className="text-sm font-medium text-muted-foreground mt-1">Policy Enforcement Status</div>
                    <StatsChange change={stats.policyEnforcement.change} betterPositive={true} />
                </CardContent>
            </Card>
        </div>
    );
};

export default Stats;
