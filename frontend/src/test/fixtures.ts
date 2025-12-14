
import type {
  User,
  Team,
  Project,
  CSPReport,
  Pagination,
  ReportStatsMetrics,
  ReportViolationTrends,
  ReportTopViolatedDirectives,
  ReportBrowserOSViolations,
  ReportTopViolationSources,
} from '@/types'

export const mockUser: User = {
  id: 'u123',
  firstName: 'Rick',
  lastName: 'James',
  username: 'rjames',
  email: 'rjames@example.com',
  role: 'user',
  image: '',
  disabled: false,
  teamId: 't123',
  teamName: 'Demo Team',
  updatedAt: new Date().toISOString(),
  defaultProjectId: 'p123',
}

export const mockAdminUser: User = {
  ...mockUser,
  id: 'a123',
  username: 'admin',
  role: 'admin',
}

export const mockSuperAdminUser: User = {
  ...mockUser,
  id: 'sa123',
  username: 'superadmin',
  role: 'superadmin',
}

export const mockTeam: Team = {
  id: 't123',
  name: 'Demo Team',
  description: 'A demo team for testing',
  updatedAt: new Date().toISOString(),
}

export const mockProject: Project = {
  id: 'p123',
  name: 'Demo Project',
  description: 'A demo project',
  teamId: 't123',
  teamName: 'Demo Team',
  reportingUrl: 'https://api.example.com/report/p123',
  lastActive: new Date().toISOString(),
  disabled: false,
}

export const mockCSPReport: CSPReport = {
  directive: 'script-src',
  documentURL: 'https://example.com/page',
  disposition: 'enforce',
  blockedURL: 'https://evil.com/script.js',
  body: '{}',
  sourceIP: '127.0.0.1',
  userAgent: 'Mozilla/5.0',
  count: 5,
  lastSeen: new Date().toISOString(),
}

export const mockPagination: Pagination = {
  page: 1,
  pageSize: 10,
  totalRows: 50,
  totalPages: 5,
}

export const mockReportStatsMetrics: ReportStatsMetrics = {
  totalViolations: { value: 1200, change: 10 },
  totalCriticalViolations: { value: 45, change: -5 },
  affectedDomains: { value: 3, change: 0 },
  policyEnforcement: { value: 98, change: 2 },
}

export const mockViolationTrends: ReportViolationTrends = [
  { day: 'Mon', critical: 5, high: 10, medium: 20 },
  { day: 'Tue', critical: 3, high: 15, medium: 25 },
  { day: 'Wed', critical: 6, high: 8, medium: 30 },
]

export const mockTopViolatedDirectives: ReportTopViolatedDirectives = {
  totalViolations: 1000,
  violations: [
    { directive: 'script-src', count: 500, percentage: 50 },
    { directive: 'img-src', count: 300, percentage: 30 },
    { directive: 'style-src', count: 200, percentage: 20 },
  ],
}

export const mockBrowserOSViolations: ReportBrowserOSViolations = {
  browser: [
    { name: 'Chrome', value: 60 },
    { name: 'Firefox', value: 30 },
    { name: 'Safari', value: 10 },
  ],
  os: [
    { name: 'Windows', value: 50 },
    { name: 'MacOS', value: 40 },
    { name: 'Linux', value: 10 },
  ],
}

export const mockTopViolationSources: ReportTopViolationSources = [
  {
    blockedURL: 'https://evil.com/bad.js',
    count: 150,
    severity: 'critical',
    lastSeen: new Date().toISOString(),
  },
  {
    blockedURL: 'http://insecure.com/img.png',
    count: 80,
    severity: 'medium',
    lastSeen: new Date().toISOString(),
  },
]
