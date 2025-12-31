import { useState, useMemo } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import {
    LineChart,
    Line,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    Legend,
    ResponsiveContainer
} from 'recharts';
import { LoadingPage } from '@/components/common/loading-page';
import { ErrorPage } from '@/components/common/error-page';

import { CSP_DIRECTIVE_COLORS, DIRECTIVES } from '@/constants';

type DirectiveKey = keyof typeof CSP_DIRECTIVE_COLORS;

import { useReportGraphData } from '@/hooks/use-reports';

const ReportsGraph = ({ projectId }: { projectId: string }) => {
    const [duration, setDuration] = useState(30);
    const [enabledDirectives, setEnabledDirectives] = useState<Set<DirectiveKey>>(
        new Set(DIRECTIVES)
    );

    const { data: graphData, isLoading, error } = useReportGraphData(projectId);

    const chartData = useMemo(() => {
        if (isLoading) return [];
        if (error) {
            console.error('ReportsGraph: error', error);
            return [];
        }

        const rawData = graphData?.data;

        if (!rawData || !Array.isArray(rawData)) {
            console.warn('ReportsGraph: data is not an array', rawData);
            return [];
        }

        const today = new Date();
        today.setHours(0, 0, 0, 0);

        const cutoffDate = new Date(today);
        cutoffDate.setDate(today.getDate() - duration);


        const filteredData = rawData
            .filter((point: any) => {
                if (!point.date) return false;

                const [year, month, day] = point.date.split('-').map(Number);
                const pointDate = new Date(year, month - 1, day);

                return pointDate >= cutoffDate;
            });

        return filteredData.map((point: any) => {
            const [year, month, day] = point.date.split('-').map(Number);
            const dateObj = new Date(year, month - 1, day);
            const displayDay = dateObj.toLocaleDateString('en-US', { month: 'long', day: 'numeric' });

            const dataPoint: any = {
                day: displayDay,
            };

            if (point.violations) {
                point.violations.forEach((v: any) => {
                    dataPoint[v.directive] = v.count;
                });
            }
            return dataPoint;
        });
    }, [graphData, duration, isLoading, error]);

    return (
        <Card className="mb-6">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-lg font-semibold">Reports Graph</CardTitle>
            </CardHeader>
            <CardContent>
                <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4 mb-4 pb-4 border-b border-border">
                    <div className="shrink-0">
                        <Select value={duration.toString()} onValueChange={(value) => setDuration(parseInt(value))}>
                            <SelectTrigger className="w-[140px]">
                                <SelectValue placeholder="Select duration" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="7">1 Week</SelectItem>
                                <SelectItem value="14">2 Weeks</SelectItem>
                                <SelectItem value="21">3 Weeks</SelectItem>
                                <SelectItem value="30">4 Weeks</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>

                    <ToggleGroup
                        type="multiple"
                        value={Array.from(enabledDirectives)}
                        onValueChange={(values) => setEnabledDirectives(new Set(values as DirectiveKey[]))}
                        className="flex flex-wrap gap-2 flex-1"
                        variant="outline"
                    >
                        {DIRECTIVES.map((directive) => (
                            <ToggleGroupItem
                                key={directive}
                                value={directive}
                                aria-label={`Toggle ${directive}`}
                                className="gap-1.5"
                            >
                                <span
                                    className="inline-block w-3 h-3 rounded-full"
                                    style={{ backgroundColor: CSP_DIRECTIVE_COLORS[directive] }}
                                />
                                <span className="text-xs font-medium">{directive}</span>
                            </ToggleGroupItem>
                        ))}
                    </ToggleGroup>
                </div>

                {/* Recharts Line Chart */}
                <div className="w-full" style={{ height: '350px' }}>
                    {isLoading ? (
                        <LoadingPage className="h-full" />
                    ) : error ? (
                        <ErrorPage error={error} />
                    ) : (
                        <ResponsiveContainer width="100%" height="100%">
                            <LineChart
                                data={chartData}
                                margin={{ top: 5, right: 20, left: 0, bottom: 5 }}
                            >
                                <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
                                <XAxis
                                    dataKey="day"
                                    label={{ value: 'Day', position: 'insideBottom', offset: -5 }}
                                    className="text-xs"
                                    tick={{ fill: 'hsl(var(--foreground))' }}
                                />
                                <YAxis
                                    label={{ value: 'Violations', angle: -90, position: 'insideLeft' }}
                                    className="text-xs"
                                    tick={{ fill: 'hsl(var(--foreground))' }}
                                />
                                <Tooltip
                                    contentStyle={{
                                        backgroundColor: 'hsl(var(--background))',
                                        border: '1px solid hsl(var(--border))',
                                        borderRadius: '6px',
                                        color: 'hsl(var(--foreground))'
                                    }}
                                />
                                <Legend wrapperStyle={{ paddingTop: '20px' }} />
                                {DIRECTIVES.map((directive) => (
                                    enabledDirectives.has(directive) && (
                                        <Line
                                            key={directive}
                                            type="monotone"
                                            dataKey={directive}
                                            stroke={CSP_DIRECTIVE_COLORS[directive]}
                                            strokeWidth={3.0}
                                            dot={{ r: 3 }}
                                            activeDot={{ r: 5 }}
                                        />
                                    )
                                ))}
                            </LineChart>
                        </ResponsiveContainer>
                    )}
                </div>
            </CardContent>
        </Card >
    );
};

export default ReportsGraph;
