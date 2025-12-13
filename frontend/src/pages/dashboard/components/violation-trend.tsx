import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
    BarChart,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip as RechartsTooltip,
    ResponsiveContainer,
} from "recharts";
import { LoadingPage } from "@/components/common/loading-page";
import { ErrorPage } from "@/components/common/error-page";
import { useReportViolationTrend } from "@/hooks/use-reports";

const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
        const data = payload[0].payload;
        return (
            <div className="recharts-custom-tooltip">
                <p className="font-semibold text-foreground mb-2">
                    {data.day} Details
                </p>
                <div className="space-y-1 text-muted-foreground text-xs">
                    <p>
                        <span className="font-medium text-destructive">Critical:</span> {data.critical}
                    </p>
                    <p>
                        <span className="font-medium text-orange-500">High:</span> {data.high}
                    </p>
                    <p>
                        <span className="font-medium text-yellow-500">Medium:</span> {data.medium}
                    </p>
                    <p className="mt-2 pt-2 border-t border-border font-semibold text-foreground">
                        Total: {data.total}
                    </p>
                </div>
            </div>
        );
    }
    return null;
};

const ViolationTrend = ({ projectId }: { projectId: string }) => {
    const { data: violationTrendData, isLoading, error } = useReportViolationTrend(projectId);

    if (isLoading) {
        return <LoadingPage className="h-48" />;
    }

    if (error || !violationTrendData?.data) {
        return <ErrorPage error={error as Error} className="h-48" />;
    }

    const violationTrend = violationTrendData.data;

    return (
        <Card>
            <CardHeader>
                <CardTitle>Violation Trend (7 Days)</CardTitle>
            </CardHeader>
            <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                    <BarChart data={violationTrend} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
                        <CartesianGrid strokeDasharray="3 3" className="stroke-border" opacity={0.3} />
                        <XAxis
                            dataKey="day"
                            tick={{ fontSize: 12 }}
                            className="text-muted-foreground"
                        />
                        <YAxis
                            tick={{ fontSize: 12 }}
                            className="text-muted-foreground"
                        />
                        <RechartsTooltip content={<CustomTooltip />} cursor={{ fill: 'hsl(var(--muted))', opacity: 0.3 }} />
                        <Bar
                            dataKey="total"
                            fill="hsl(var(--primary))"
                            radius={[4, 4, 0, 0]}
                            animationDuration={500}
                        />
                    </BarChart>
                </ResponsiveContainer>
            </CardContent>
        </Card>
    );
};

export default ViolationTrend;
