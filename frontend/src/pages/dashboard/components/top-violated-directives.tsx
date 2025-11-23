import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useState } from 'react';

const TopViolatedDirectives = () => {
  const [showPercentages, setShowPercentages] = useState(true);

  const directives = [
    { name: 'script-src', violations: 523, percentage: '42%' },
    { name: 'img-src', violations: 312, percentage: '25%' },
    { name: 'style-src', violations: 187, percentage: '15%' },
    { name: 'connect-src', violations: 145, percentage: '12%' },
    { name: 'frame-src', violations: 80, percentage: '6%' },
  ];

  const getWidth = (violations: number) => {
    const max = 523;
    return `${(violations / max) * 100}%`;
  };

  const getTooltipContent = (name: string) => {
    switch (name) {
      case 'script-src':
        return '<strong>script-src violations:</strong><br>Inline scripts: 312<br>External scripts: 189<br>Unsafe eval: 22';
      case 'img-src':
        return '<strong>img-src violations:</strong><br>Third-party images: 245<br>Data URIs: 48<br>Unknown sources: 19';
      case 'style-src':
        return '<strong>style-src violations:</strong><br>Inline styles: 134<br>External stylesheets: 42<br>Unsafe inline: 11';
      case 'connect-src':
        return '<strong>connect-src violations:</strong><br>API calls: 98<br>WebSocket: 32<br>XHR requests: 15';
      case 'frame-src':
        return '<strong>frame-src violations:</strong><br>Embedded iframes: 56<br>Third-party widgets: 19<br>Other: 5';
      default:
        return '';
    }
  };

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-lg font-semibold">Top Violated Directives</CardTitle>
        <label className="flex items-center gap-2 cursor-pointer">
          <span className="text-xs text-muted-foreground">Show percentages</span>
          <input
            type="checkbox"
            checked={showPercentages}
            onChange={(e) => setShowPercentages(e.target.checked)}
            className="accent-primary h-4 w-4"
          />
        </label>
      </CardHeader>
      <CardContent>
        <div className="space-y-4 mt-4">
          {directives.map((directive) => (
            <div key={directive.name} className="group cursor-pointer relative" data-violations={directive.violations} data-percentage={directive.percentage}>
              <div className="flex justify-between items-center mb-1">
                <span className="text-sm font-medium text-foreground">{directive.name}</span>
                <span className="text-sm directive-value text-muted-foreground">
                  {showPercentages ? `${directive.percentage} of total` : `${directive.violations} violations`}
                </span>
              </div>
              <div className="w-full bg-secondary rounded-full h-2 overflow-hidden">
                <div className="h-2 bg-primary rounded-full transition-all duration-500 ease-out" style={{ width: getWidth(directive.violations) }} title={`${directive.name}: ${directive.violations} violations`}></div>
              </div>
              <div
                className="tooltip absolute left-0 top-full mt-2 hidden group-hover:block z-10 p-2 rounded text-xs bg-popover border border-border shadow-lg text-popover-foreground"
                style={{ minWidth: '200px' }}
                dangerouslySetInnerHTML={{ __html: getTooltipContent(directive.name) }}
              />
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
};

export default TopViolatedDirectives;
