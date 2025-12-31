export const CSP_DIRECTIVE_COLORS: Record<string, string> = {
  'default-src': '#FACC15',
  'script-src': '#EAB308',
  'script-src-elem': '#CA8A04',
  'script-src-attr': '#A3E635',
  'style-src': '#84CC16',
  'style-src-elem': '#65A30D',
  'style-src-attr': '#4ADE80',
  'img-src': '#22C55E',
  'connect-src': '#16A34A',
  'font-src': '#34D399',
  'object-src': '#10B981',
  'media-src': '#059669',
  'frame-src': '#2DD4BF',
  'sandbox': '#14B8A6',
  'child-src': '#22D3EE',
  'form-action': '#06B6D4',
  'frame-ancestors': '#0891B2',
  'plugin-types': '#38BDF8',
  'base-uri': '#0EA5E9',
  'worker-src': '#0284C7',
  'manifest-src': '#60A5FA',
  'prefetch-src': '#3B82F6',
  'navigate-to': '#2563EB',
};

export const DIRECTIVES = Object.keys(CSP_DIRECTIVE_COLORS);

export const BROWSER_COLORS: Record<string, string> = {
  Chrome: '#3b82f6',
  Firefox: '#f97316',
  Safari: '#10b981',
  Edge: '#8b5cf6',
  Opera: '#e41511ff',
  IE: '#6bb4e8ff',
  Samsung: '#7430e2ff',
  UCBrowser: '#d1702aff',
  Others: '#94a3b8',
};

export const OS_COLORS: Record<string, string> = {
  Windows: 'bg-red-500',
  macOS: 'bg-zinc-500',
  Linux: 'bg-orange-500',
  iOS: 'bg-blue-500',
  Android: 'bg-green-500',
  ChromeOS: 'bg-yellow-500',
  Others: 'bg-slate-400',
};