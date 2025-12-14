import { render, screen } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import Dashboard from './index'
import { mockUser } from '@/test/fixtures'

// Mock child components to isolate Dashboard test
vi.mock('@/pages/dashboard/components/dashboard-fallback', () => ({
  default: () => <div data-testid="dashboard-fallback">Fallback</div>,
}))
vi.mock('@/pages/dashboard/components/stats', () => ({
  default: () => <div data-testid="stats">Stats</div>,
}))
vi.mock('@/pages/dashboard/components/violation-trend', () => ({
  default: () => <div data-testid="violation-trend">ViolationTrend</div>,
}))
vi.mock('@/pages/dashboard/components/top-violated-directives', () => ({
  default: () => <div data-testid="top-violated-directives">TopViolatedDirectives</div>,
}))
vi.mock('@/pages/dashboard/components/top-violated-document-urls', () => ({
  default: () => <div data-testid="top-violated-document-urls">TopViolatedDocumentUrls</div>,
}))
vi.mock('@/pages/dashboard/components/violation-by-browser-os', () => ({
  default: () => <div data-testid="violation-by-browser-os">ViolationByBrowserOS</div>,
}))
vi.mock('@/pages/dashboard/components/violation-sources', () => ({
  default: () => <div data-testid="violation-sources">ViolationSources</div>,
}))
vi.mock('@/pages/dashboard/components/reports-graph', () => ({
  default: () => <div data-testid="reports-graph">ReportsGraph</div>,
}))

// Mock useAuth
const useAuthMock = vi.fn()
vi.mock('@/hooks/use-auth', () => ({
  useAuth: () => useAuthMock(),
}))

describe('Dashboard Page', () => {
  it('renders fallback when no defaultProjectId', () => {
    useAuthMock.mockReturnValue({
      user: { ...mockUser, defaultProjectId: null },
    })
    render(<Dashboard />)
    expect(screen.getByTestId('dashboard-fallback')).toBeInTheDocument()
  })

  it('renders dashboard content when defaultProjectId exists', () => {
    useAuthMock.mockReturnValue({
      user: { ...mockUser, defaultProjectId: 'project-123' },
    })
    render(<Dashboard />)

    expect(screen.getByText('Dashboard')).toBeInTheDocument()
    expect(screen.getByText('Real-time Content Security Policy monitoring and analytics')).toBeInTheDocument()

    expect(screen.getByTestId('stats')).toBeInTheDocument()
    expect(screen.getByTestId('reports-graph')).toBeInTheDocument()
    expect(screen.getByTestId('violation-trend')).toBeInTheDocument()
    expect(screen.getByTestId('top-violated-directives')).toBeInTheDocument()
    expect(screen.getByTestId('top-violated-document-urls')).toBeInTheDocument()
    expect(screen.getByTestId('violation-by-browser-os')).toBeInTheDocument()
    expect(screen.getByTestId('violation-sources')).toBeInTheDocument()
  })
})
