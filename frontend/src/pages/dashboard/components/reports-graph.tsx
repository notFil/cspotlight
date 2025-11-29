import React, { useState, useMemo } from 'react';
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

// Color palette for directives
const DIRECTIVE_COLORS = {
    'script-src': '#21808d',
    'img-src': '#a84b2f',
    'style-src': '#518a91ff',
    'connect-src': '#e68161',
    'frame-src': '#c0152f'
} as const;

type DirectiveKey = keyof typeof DIRECTIVE_COLORS;

const ReportsGraph = () => {
    const [duration, setDuration] = useState(30);
    const [enabledDirectives, setEnabledDirectives] = useState<Set<DirectiveKey>>(
        new Set(['script-src', 'img-src', 'style-src', 'connect-src', 'frame-src'])
    );

    // Generate sample data
    const chartData = useMemo(() => {
        const data = [];
        for (let i = 0; i < duration; i++) {
            data.push({
                day: i + 1,
                'script-src': Math.floor(Math.random() * 50) + 30,
                'img-src': Math.floor(Math.random() * 40) + 20,
                'style-src': Math.floor(Math.random() * 30) + 10,
                'connect-src': Math.floor(Math.random() * 25) + 8,
                'frame-src': Math.floor(Math.random() * 15) + 5,
            });
        }
        return data;
    }, [duration]);

    return (
        <Card className="mb-6">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-lg font-semibold">Reports Graph</CardTitle>
            </CardHeader>
            <CardContent>
                {/* Select and Directive toggles on same line */}
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
                        {(Object.keys(DIRECTIVE_COLORS) as DirectiveKey[]).map((directive) => (
                            <ToggleGroupItem
                                key={directive}
                                value={directive}
                                aria-label={`Toggle ${directive}`}
                                className="gap-1.5"
                            >
                                <span
                                    className="inline-block w-3 h-3 rounded-full"
                                    style={{ backgroundColor: DIRECTIVE_COLORS[directive] }}
                                />
                                <span className="text-xs font-medium">{directive}</span>
                            </ToggleGroupItem>
                        ))}
                    </ToggleGroup>
                </div>

                {/* Recharts Line Chart */}
                <div className="w-full" style={{ height: '350px' }}>
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
                            <Legend />
                            {(Object.keys(DIRECTIVE_COLORS) as DirectiveKey[]).map((directive) => (
                                enabledDirectives.has(directive) && (
                                    <Line
                                        key={directive}
                                        type="monotone"
                                        dataKey={directive}
                                        stroke={DIRECTIVE_COLORS[directive]}
                                        strokeWidth={2.5}
                                        dot={{ r: 3 }}
                                        activeDot={{ r: 5 }}
                                    />
                                )
                            ))}
                        </LineChart>
                    </ResponsiveContainer>
                </div>
            </CardContent>
        </Card>
    );
};

export default ReportsGraph;
