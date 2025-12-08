
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

const TopViolatedDocumentUris = () => {
  const data = [
    { uri: "https://example.com/login", count: 1250 },
    { uri: "https://example.com/checkout", count: 980 },
    { uri: "https://example.com/dashboard", count: 750 },
    { uri: "https://example.com/profile", count: 450 },
    { uri: "https://example.com/settings", count: 200 },
  ];

  const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
      return (
        <div className="recharts-custom-tooltip">
          <p className="font-semibold text-foreground mb-2">{label}</p>
          <div className="space-y-1 text-muted-foreground text-xs">
            <p>
              Violations: <span className="font-medium text-foreground">{payload[0].value}</span>
            </p>
          </div>
        </div>
      );
    }
    return null;
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Top Violated Document URIs</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="h-[300px] w-full mt-4">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart
              layout="vertical"
              data={data}
              margin={{ top: 5, right: 30, left: 0, bottom: 5 }}
            >
              <CartesianGrid strokeDasharray="3 3" horizontal={false} className="stroke-border" opacity={0.3} />
              <XAxis type="number" hide />
              <YAxis
                dataKey="uri"
                type="category"
                width={80}
                tick={{ fontSize: 12 }}
                className="text-muted-foreground"
                tickFormatter={(value) => {
                  try {
                    const url = new URL(value);
                    return url.pathname;
                  } catch {
                    return value;
                  }
                }}
              />
              <Tooltip content={<CustomTooltip />} cursor={{ fill: 'hsl(var(--muted))', opacity: 0.3 }} />
              <Bar
                dataKey="count"
                fill="hsl(var(--primary))"
                radius={[0, 4, 4, 0]}
                barSize={20}
                background={{ fill: 'hsl(var(--muted))', opacity: 0.1 }}
              />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </CardContent>
    </Card>
  );
};

export default TopViolatedDocumentUris;
