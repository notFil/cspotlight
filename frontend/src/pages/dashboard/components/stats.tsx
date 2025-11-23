import { Card, CardContent } from "@/components/ui/card";
import { ArrowUp, ArrowDown } from "lucide-react";

const Stats = () => {
    return (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-6">
            <Card>
                <CardContent className="p-6">
                    <div className="text-2xl font-bold">1,247</div>
                    <div className="text-sm font-medium text-muted-foreground mt-1">Total Violations (24h)</div>
                    <div className="mt-2 flex items-center text-xs font-medium text-red-600 bg-red-50 w-fit px-2 py-0.5 rounded-full dark:bg-red-900/20 dark:text-red-400">
                        <ArrowUp className="w-3 h-3 mr-1" />
                        <span>12%</span>
                    </div>
                </CardContent>
            </Card>
            <Card>
                <CardContent className="p-6">
                    <div className="text-2xl font-bold">23</div>
                    <div className="text-sm font-medium text-muted-foreground mt-1">Critical Violations</div>
                    <div className="mt-2 flex items-center text-xs font-medium text-red-600 bg-red-50 w-fit px-2 py-0.5 rounded-full dark:bg-red-900/20 dark:text-red-400">
                        <ArrowUp className="w-3 h-3 mr-1" />
                        <span>8%</span>
                    </div>
                </CardContent>
            </Card>
            <Card>
                <CardContent className="p-6">
                    <div className="text-2xl font-bold">12</div>
                    <div className="text-sm font-medium text-muted-foreground mt-1">Affected Domains</div>
                    <div className="mt-2 flex items-center text-xs font-medium text-emerald-600 bg-emerald-50 w-fit px-2 py-0.5 rounded-full dark:bg-emerald-900/20 dark:text-emerald-400">
                        <ArrowDown className="w-3 h-3 mr-1" />
                        <span>3%</span>
                    </div>
                </CardContent>
            </Card>
            <Card>
                <CardContent className="p-6">
                    <div className="text-2xl font-bold">94.3%</div>
                    <div className="text-sm font-medium text-muted-foreground mt-1">Policy Compliance</div>
                    <div className="mt-2 flex items-center text-xs font-medium text-emerald-600 bg-emerald-50 w-fit px-2 py-0.5 rounded-full dark:bg-emerald-900/20 dark:text-emerald-400">
                        <ArrowUp className="w-3 h-3 mr-1" />
                        <span>2.1%</span>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
};

export default Stats;
