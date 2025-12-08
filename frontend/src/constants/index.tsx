export const CSP_DIRECTIVE_COLORS: Record<string, string> = {
  'default-src': '#e11d48',
  'script-src': '#21808d',
  'script-src-elem': '#0d9488',
  'script-src-attr': '#14b8a6',
  'style-src': '#518a91',
  'style-src-elem': '#0891b2',
  'style-src-attr': '#06b6d4',
  'img-src': '#a84b2f',
  'connect-src': '#e68161',
  'font-src': '#7c3aed',
  'object-src': '#db2777',
  'media-src': '#9333ea',
  'frame-src': '#c0152f',
  'sandbox': '#ea580c',
  'report-uri': '#475569',
  'child-src': '#d97706',
  'form-action': '#65a30d',
  'frame-ancestors': '#16a34a',
  'plugin-types': '#059669',
  'base-uri': '#2563eb',
  'worker-src': '#ca8a04',
  'manifest-src': '#be123c',
  'prefetch-src': '#8b5cf6',
  'navigate-to': '#0284c7',
};

export const directives = Object.keys(CSP_DIRECTIVE_COLORS);

export const BROWSER_COLORS: Record<string, string> = {
  Chrome: '#3b82f6', // blue-500
  Firefox: '#f97316', // orange-500
  Safari: '#10b981', // emerald-500
  Edge: '#8b5cf6',   // violet-500
  Other: '#94a3b8',  // slate-400
};

export const OS_COLORS: Record<string, string> = {
  Windows: 'bg-blue-500',
  MacOS: 'bg-zinc-500',
  Linux: 'bg-orange-500',
  Android: 'bg-green-500',
  Other: 'bg-slate-400',
};