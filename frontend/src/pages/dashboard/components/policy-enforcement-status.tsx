import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const PolicyEnforcementStatus = () => {
    return (
        <Card>
            <CardHeader>
                <CardTitle>Policy Enforcement Status</CardTitle>
            </CardHeader>
            <CardContent>
                <div className="space-y-4">
                    <div>
                        <div className="flex justify-between items-center mb-2">
                            <span className="text-sm font-medium text-foreground">Enforcing Mode</span>
                            <span className="text-sm font-semibold text-foreground">73%</span>
                        </div>
                        <div className="w-full bg-secondary rounded-full h-3 overflow-hidden">
                            <div className="h-3 rounded-full bg-emerald-500" style={{ width: '73%' }}></div>
                        </div>
                        <div className="mt-1 text-xs text-muted-foreground">18 domains actively blocking violations</div>
                    </div>
                    <div>
                        <div className="flex justify-between items-center mb-2">
                            <span className="text-sm font-medium text-foreground">Report-Only Mode</span>
                            <span className="text-sm font-semibold text-foreground">27%</span>
                        </div>
                        <div className="w-full bg-secondary rounded-full h-3 overflow-hidden">
                            <div className="h-3 rounded-full bg-orange-500" style={{ width: '27%' }}></div>
                        </div>
                        <div className="mt-1 text-xs text-muted-foreground">7 domains in monitoring mode</div>
                    </div>
                    <div className="mt-6 p-4 rounded-lg bg-primary/10 border-l-4 border-primary">
                        <p className="text-sm text-foreground">
                            <strong>Recommendation:</strong> Consider moving staging.app.example.com and api-v2.example.com to enforcing mode after 7 days of monitoring.
                        </p>
                    </div>
                </div>
            </CardContent>
        </Card>
    );
};

export default PolicyEnforcementStatus;
