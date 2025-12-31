import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export const CSPGuide = () => {
  return (
    <Card className="mt-8">
      <CardHeader>
        <CardTitle>Setting up Content Security Policy (CSP)</CardTitle>
        <CardDescription>
          Protect your application by configuring the Content-Security-Policy header.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <h3 className="text-sm font-medium">1. Configure the Report-To Header</h3>
          <p className="text-sm text-muted-foreground">
            Add the <code>Report-To</code> header to your server's response. Replace <code>YOUR_REPORTING_URL</code> with the URL from the table above.
          </p>
          <pre className="bg-muted p-4 rounded-md overflow-x-auto text-xs">
            <code>
              Report-To: &#123;"group":"csp-endpoint","max_age":10886400,"endpoints":[&#123;"url":"YOUR_REPORTING_URL"&#125;]&#125;
            </code>
          </pre>
        </div>

        <div className="space-y-2">
          <h3 className="text-sm font-medium">2. Configure the Content-Security-Policy Header</h3>
          <p className="text-sm text-muted-foreground">
            Add the <code>Content-Security-Policy</code> header to start enforcing policies and reporting violations.
          </p>
          <pre className="bg-muted p-4 rounded-md overflow-x-auto text-xs">
            <code>
              Content-Security-Policy: default-src 'self'; report-to csp-endpoint;
            </code>
          </pre>
        </div>
      </CardContent>
    </Card>
  );
};
