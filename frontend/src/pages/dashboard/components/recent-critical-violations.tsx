import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const RecentCriticalViolations = () => {
    return (
        <Card>
            <CardHeader>
                <CardTitle>Recent Critical Violations</CardTitle>
            </CardHeader>
            <CardContent>
                <div className="space-y-3">
                    <div className="p-3 rounded-lg bg-red-50 dark:bg-red-900/10 border-l-4 border-red-500">
                        <div className="flex items-start justify-between mb-1">
                            <span className="text-sm font-medium text-foreground">Inline script execution blocked</span>
                            <span className="text-xs text-muted-foreground">2m ago</span>
                        </div>
                        <div className="text-xs text-muted-foreground">
                            Domain: app.example.com | Directive: script-src
                        </div>
                    </div>
                    <div className="p-3 rounded-lg bg-red-50 dark:bg-red-900/10 border-l-4 border-red-500">
                        <div className="flex items-start justify-between mb-1">
                            <span className="text-sm font-medium text-foreground">Unauthorized external resource</span>
                            <span className="text-xs text-muted-foreground">8m ago</span>
                        </div>
                        <div className="text-xs text-muted-foreground">
                            Domain: checkout.example.com | Directive: img-src
                        </div>
                    </div>
                    <div className="p-3 rounded-lg bg-orange-50 dark:bg-orange-900/10 border-l-4 border-orange-500">
                        <div className="flex items-start justify-between mb-1">
                            <span className="text-sm font-medium text-foreground">Third-party script violation</span>
                            <span className="text-xs text-muted-foreground">15m ago</span>
                        </div>
                        <div className="text-xs text-muted-foreground">
                            Domain: dashboard.example.com | Directive: script-src
                        </div>
                    </div>
                    <div className="p-3 rounded-lg bg-red-50 dark:bg-red-900/10 border-l-4 border-red-500">
                        <div className="flex items-start justify-between mb-1">
                            <span className="text-sm font-medium text-foreground">Unsafe eval() detected</span>
                            <span className="text-xs text-muted-foreground">22m ago</span>
                        </div>
                        <div className="text-xs text-muted-foreground">
                            Domain: api.example.com | Directive: script-src
                        </div>
                    </div>
                </div>
            </CardContent>
        </Card>
    );
};

export default RecentCriticalViolations;
