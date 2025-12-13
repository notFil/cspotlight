// src/types/index.ts
export interface AppState {
  isLoading: boolean;
  user?: User | null;
  theme: 'light' | 'dark';
  error?: string | null;
}

export interface User {
  id?: string;
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  role: 'user' | 'admin' | 'superadmin';
  image?: string;
  disabled: boolean;
  teamId: string;
  teamName: string;
  updatedAt?: string;
  defaultProjectId?: string;
}

export interface UserRegister {
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
}

export interface UserLogin {
  username: string;
  password: string;
}

export interface Team {
  id: string;
  name: string;
  description: string;
  updatedAt?: string;
}

export interface Project {
  id?: string;
  name: string;
  description: string;
  teamId: string;
  teamName: string;
  reportingUrl?: string;
  lastActive?: string;
  disabled: boolean;
}

export interface CSPReport {
  directive: string;
  documentURL: string;
  disposition: string;
  blockedURL: string;
  body: string;
  sourceIP: string;
  userAgent: string;
  count: number;
  lastSeen: string;
}

export interface AuthTokenData {
  accessToken: string;
  refreshToken: string;
  createdAt: string;
  expiresIn: string;
  refreshExpiresIn: string;
}

export interface AuthRequest {
  username: string;
  password: string;
}

export interface ChangePasswordRequest {
  currentPassword: string;
  newPassword: string;
  confirmNewPassword: string;
}

export interface Pagination {
  page: number;
  pageSize: number;
  totalRows: number;
  totalPages: number;
}

export interface APIResponse<T> {
  data: T;
  message?: string;
  pagination?: Pagination;
  error?: string;
}

export interface ReportGraphData {
  date: string;
  violations: Violation[];
}

export interface Violation {
  directive: string;
  count: number;
}

export interface ReportStatsMetrics {
  totalViolations: MetricSummary;
  totalCriticalViolations: MetricSummary;
  affectedDomains: MetricSummary;
  policyEnforcement: MetricSummary;
}

export interface MetricSummary {
  value: number;
  change: number;
}

export interface ReportViolationTrends extends Array<ViolationTrend> { }

export interface ViolationTrend {
  day: string;
  critical: number;
  high: number;
  medium: number;
}

export interface ReportTopViolatedDirectives {
  totalViolations: number;
  violations: ViolatedDirectives[];
}

export interface ViolatedDirectives {
  directive: string;
  count: number;
  percentage: number;
}

export interface ReportTopViolatedDocumentURLs extends Array<ViolatedDocumentURLs> { }

export interface ViolatedDocumentURLs {
  url: string;
  count: number;
}

export interface ReportBrowserOSViolations {
  browser: StatsItem[];
  os: StatsItem[];
}

export interface StatsItem {
  name: string;
  value: number;
}

export interface ReportTopViolationSources extends Array<ViolationSource> { }

export interface ViolationSource {
  blockedURL: string;
  count: number;
  severity: string;
  lastSeen: string;
}

