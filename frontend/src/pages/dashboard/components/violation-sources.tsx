import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";

const ViolationSources = () => {
    return (
        <Card className="mb-6">
            <CardHeader>
                <CardTitle>Top Violation Sources</CardTitle>
            </CardHeader>
            <CardContent>
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>Domain</TableHead>
                            <TableHead>Violations</TableHead>
                            <TableHead>Severity</TableHead>
                            <TableHead>Last Seen</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        <TableRow className="group cursor-pointer" title="Click for details">
                            <TableCell className="font-medium">
                                cdn.analytics-tracker.com
                                <div className="hidden group-hover:block text-xs mt-1 p-2 rounded bg-red-50 border-l-2 border-red-500 dark:bg-red-900/10">
                                    <strong>Risk:</strong> Unauthorized tracking script<br />
                                    <strong>Action:</strong> Block and investigate source
                                </div>
                            </TableCell>
                            <TableCell>387</TableCell>
                            <TableCell>
                                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400">
                                    Critical
                                </span>
                            </TableCell>
                            <TableCell className="text-muted-foreground">2 minutes ago</TableCell>
                        </TableRow>
                        <TableRow className="group cursor-pointer" title="Click for details">
                            <TableCell className="font-medium">
                                third-party-ads.net
                                <div className="hidden group-hover:block text-xs mt-1 p-2 rounded bg-orange-50 border-l-2 border-orange-500 dark:bg-orange-900/10">
                                    <strong>Risk:</strong> Ad injection attempts<br />
                                    <strong>Action:</strong> Review ad network policies
                                </div>
                            </TableCell>
                            <TableCell>264</TableCell>
                            <TableCell>
                                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400">
                                    High
                                </span>
                            </TableCell>
                            <TableCell className="text-muted-foreground">15 minutes ago</TableCell>
                        </TableRow>
                        <TableRow className="group cursor-pointer" title="Click for details">
                            <TableCell className="font-medium">
                                external-fonts.googleapis.com
                                <div className="hidden group-hover:block text-xs mt-1 p-2 rounded bg-emerald-50 border-l-2 border-emerald-500 dark:bg-emerald-900/10">
                                    <strong>Risk:</strong> Font loading from CDN<br />
                                    <strong>Action:</strong> Consider whitelisting trusted source
                                </div>
                            </TableCell>
                            <TableCell>156</TableCell>
                            <TableCell>
                                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400">
                                    Medium
                                </span>
                            </TableCell>
                            <TableCell className="text-muted-foreground">1 hour ago</TableCell>
                        </TableRow>
                        <TableRow className="group cursor-pointer" title="Click for details">
                            <TableCell className="font-medium">
                                widgets.social-media.io
                                <div className="hidden group-hover:block text-xs mt-1 p-2 rounded bg-orange-50 border-l-2 border-orange-500 dark:bg-orange-900/10">
                                    <strong>Risk:</strong> Social media widget loading<br />
                                    <strong>Action:</strong> Verify widget necessity
                                </div>
                            </TableCell>
                            <TableCell>142</TableCell>
                            <TableCell>
                                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400">
                                    High
                                </span>
                            </TableCell>
                            <TableCell className="text-muted-foreground">32 minutes ago</TableCell>
                        </TableRow>
                        <TableRow className="group cursor-pointer" title="Click for details">
                            <TableCell className="font-medium">
                                cdn.customer-chat.com
                                <div className="hidden group-hover:block text-xs mt-1 p-2 rounded bg-emerald-50 border-l-2 border-emerald-500 dark:bg-emerald-900/10">
                                    <strong>Risk:</strong> Chat widget integration<br />
                                    <strong>Action:</strong> Monitor for suspicious activity
                                </div>
                            </TableCell>
                            <TableCell>98</TableCell>
                            <TableCell>
                                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400">
                                    Medium
                                </span>
                            </TableCell>
                            <TableCell className="text-muted-foreground">3 hours ago</TableCell>
                        </TableRow>
                    </TableBody>
                </Table>
            </CardContent>
        </Card>
    );
};

export default ViolationSources;
