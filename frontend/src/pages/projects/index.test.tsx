import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import Projects from './index'
import { mockProject, mockUser } from '@/test/fixtures'

// Mock hooks
const useProjectsMock = vi.fn()
vi.mock('@/hooks/use-projects', () => ({
  useProjects: () => useProjectsMock(),
}))

const useAuthMock = vi.fn()
vi.mock('@/hooks/use-auth', () => ({
  useAuth: () => useAuthMock(),
}))

const useSetDefaultProjectMock = { mutate: vi.fn(), isPending: false }
vi.mock('@/hooks/use-users', () => ({
  useSetDefaultProject: () => useSetDefaultProjectMock,
}))

// Mock router
const navigateMock = vi.fn()
vi.mock('react-router-dom', () => ({
  useNavigate: () => navigateMock,
}))

// Mock Date utils
vi.mock('@/utils/date', () => ({
  timeAgo: () => '2 days ago',
}))

// Mock LoadingPage
vi.mock('@/components/common/loading-page', () => ({
  LoadingPage: () => <div>Loading...</div>,
}))

describe('Projects Page', () => {
  it('renders loading state', () => {
    useProjectsMock.mockReturnValue({ isLoading: true })
    useAuthMock.mockReturnValue({})
    render(<Projects />)
    expect(screen.getByText('Loading...')).toBeInTheDocument()
  })

  it('renders error state', () => {
    useProjectsMock.mockReturnValue({
      isLoading: false,
      error: new Error('Failed to fetch'),
    })
    useAuthMock.mockReturnValue({})
    render(<Projects />)
    expect(screen.getByText('Failed to fetch')).toBeInTheDocument()
  })

  it('renders empty state', () => {
    useProjectsMock.mockReturnValue({
      isLoading: false,
      data: [],
    })
    useAuthMock.mockReturnValue({ user: { id: 'user-1' } })
    render(<Projects />)
    expect(screen.getByText('No projects found. Reach to your admin to create a project.')).toBeInTheDocument()
  })

  it('renders projects list', () => {
    const projects = [
      { ...mockProject, id: 'p1', name: 'Project 1', description: 'Desc 1', disabled: false },
      { ...mockProject, id: 'p2', name: 'Project 2', description: 'Desc 2', disabled: true },
    ]
    useProjectsMock.mockReturnValue({
      isLoading: false,
      data: projects,
    })
    useAuthMock.mockReturnValue({ user: { ...mockUser, defaultProjectId: 'p2' } })

    render(<Projects />)

    expect(screen.getByText('Project 1')).toBeInTheDocument()
    expect(screen.getByText('Desc 1')).toBeInTheDocument()
    expect(screen.getByText('active')).toBeInTheDocument()

    expect(screen.getByText('Project 2')).toBeInTheDocument()
    expect(screen.getByText('Desc 2')).toBeInTheDocument()
    expect(screen.getByText('inactive')).toBeInTheDocument()
  })

  it('navigates to details', () => {
    const projects = [
      { ...mockProject, id: 'p1', name: 'Project 1' },
    ]
    useProjectsMock.mockReturnValue({
      isLoading: false,
      data: projects,
    })
    useAuthMock.mockReturnValue({ user: { ...mockUser } })

    render(<Projects />)

    fireEvent.click(screen.getByText('View Details'))
    expect(navigateMock).toHaveBeenCalledWith('/reports/p1')
  })
})
