import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { useReportTopViolationSources } from "@/hooks/use-reports";
import { LoadingPage } from "@/components/common/loading-page";
import { ErrorPage } from "@/components/common/error-page";
import { timeAgo } from "@/utils/date";

const getSeverityColor = (severity: string) => {
    switch (severity.toLowerCase()) {
        case "critical":
            return "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400";
        case "high":
            return "bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400";
        case "medium":
            return "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400";
        case "low":
            return "bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400";
        default:
            return "bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400";
    }
};

const ViolationSources = ({ projectId }: { projectId: string }) => {
    const { data, isLoading, error } = useReportTopViolationSources(projectId);

    if (isLoading) {
        return <LoadingPage className="h-48" />;
    }

    if (error || !data) {
        return <ErrorPage error={error as Error} className="h-48" />;
    }

    const topViolationSources = data?.data;
    return (
        <Card className="mb-6">
            <CardHeader>
                <CardTitle>Top Violation Sources</CardTitle>
            </CardHeader>
            <CardContent>
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>Source</TableHead>
                            <TableHead>Violations</TableHead>
                            <TableHead>Severity</TableHead>
                            <TableHead>Last Seen</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {topViolationSources?.map((source) => (
                            <TableRow key={source.blockedURL}>
                                <TableCell className="font-medium">
                                    {source.blockedURL}
                                </TableCell>
                                <TableCell>{source.count}</TableCell>
                                <TableCell>
                                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getSeverityColor(source.severity)}`}>
                                        {source.severity}
                                    </span>
                                </TableCell>
                                <TableCell className="text-muted-foreground">{timeAgo(source.lastSeen || "") || "N/A"}</TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </CardContent>
        </Card>
    );
};

export default ViolationSources;
