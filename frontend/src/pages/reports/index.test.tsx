import { render, screen } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import Reports from './index'
import { mockCSPReport, mockPagination } from '@/test/fixtures'

// Mock hooks
const useReportsMock = vi.fn()
vi.mock('@/hooks/use-reports', () => ({
  useReports: () => useReportsMock(),
}))

// Mock router
vi.mock('react-router-dom', () => ({
  useParams: () => ({ projectId: 'p123' }),
  useSearchParams: () => [new URLSearchParams(), vi.fn()],
}))

// Mock DataTable
vi.mock('./data-table', () => ({
  DataTable: ({ data }: { data: any[] }) => (
    <div data-testid="data-table">
      {data.map((item, i) => (
        <div key={i}>{item.directive}</div>
      ))}
    </div>
  ),
}))

// Mock LoadingPage
vi.mock('@/components/common/loading-page', () => ({
  LoadingPage: () => <div>Loading...</div>,
}))

describe('Reports Page', () => {
  it('renders loading state', () => {
    useReportsMock.mockReturnValue({ isLoading: true })
    render(<Reports />)
    expect(screen.getByText('Loading...')).toBeInTheDocument()
  })

  it('renders error state', () => {
    useReportsMock.mockReturnValue({
      isLoading: false,
      error: new Error('Failed to fetch'),
    })
    render(<Reports />)
    expect(screen.getByText('Failed to fetch')).toBeInTheDocument()
  })

  it('renders reports list', () => {
    useReportsMock.mockReturnValue({
      isLoading: false,
      data: { data: [mockCSPReport], pagination: mockPagination },
    })

    render(<Reports />)

    expect(screen.getByText(/Viewing reports for Report ID: p123/i)).toBeInTheDocument()
    expect(screen.getByTestId('data-table')).toBeInTheDocument()
    expect(screen.getByText(mockCSPReport.directive)).toBeInTheDocument()
  })
})
