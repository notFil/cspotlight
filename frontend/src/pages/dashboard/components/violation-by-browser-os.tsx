import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from "recharts";
import { BROWSER_COLORS, OS_COLORS } from "@/constants";

const ViolationByBrowserOS = () => {
  // Mock Data (Values will come from API)
  const browserData = [
    { name: 'Chrome', value: 45 },
    { name: 'Firefox', value: 30 },
    { name: 'Safari', value: 15 },
    { name: 'Edge', value: 10 },
  ];

  const osData = [
    { name: 'Windows', value: 60 },
    { name: 'MacOS', value: 25 },
    { name: 'Linux', value: 10 },
    { name: 'Android', value: 5 },
  ];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Violation by Browser/OS</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Browser Section - Pie Chart */}
          <div className="flex flex-col items-center">
            <h3 className="text-sm font-medium mb-4 text-muted-foreground">Top Browsers</h3>
            <div className="h-[200px] w-full">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={browserData}
                    cx="50%"
                    cy="50%"
                    innerRadius={60}
                    outerRadius={80}
                    paddingAngle={5}
                    dataKey="value"
                  >
                    {browserData.map((entry: { name: string; value: number }, index: number) => (
                      <Cell key={`cell-${index}`} fill={BROWSER_COLORS[entry.name] || BROWSER_COLORS.Other} />
                    ))}
                  </Pie>
                  <Tooltip
                    contentStyle={{ backgroundColor: 'var(--color-background)', borderRadius: '8px', border: '1px solid var(--color-border)' }}
                    itemStyle={{ color: 'var(--color-text)' }}
                  />
                </PieChart>
              </ResponsiveContainer>
            </div>
            <div className="flex flex-wrap justify-center gap-2 mt-4">
              {browserData.map((entry: { name: string; value: number }) => (
                <div key={entry.name} className="flex items-center text-xs">
                  <span className="w-2 h-2 rounded-full mr-1" style={{ backgroundColor: BROWSER_COLORS[entry.name] || BROWSER_COLORS.Other }}></span>
                  <span className="text-muted-foreground">{entry.name} ({entry.value}%)</span>
                </div>
              ))}
            </div>
          </div>

          {/* OS Section - Progress Bars */}
          <div>
            <h3 className="text-sm font-medium mb-4 text-muted-foreground">Top Operating Systems</h3>
            <div className="space-y-4">
              {osData.map((os: { name: string; value: number }) => (
                <div key={os.name}>
                  <div className="flex justify-between items-center mb-1">
                    <span className="text-xs font-medium">{os.name}</span>
                    <span className="text-xs text-muted-foreground">{os.value}%</span>
                  </div>
                  <div className="w-full bg-secondary rounded-full h-2 overflow-hidden">
                    <div
                      className={`h-2 rounded-full ${OS_COLORS[os.name] || OS_COLORS.Other}`}
                      style={{ width: `${os.value}%` }}
                    ></div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

export default ViolationByBrowserOS;
