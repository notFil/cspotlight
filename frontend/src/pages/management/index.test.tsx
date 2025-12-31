import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import Management from './index'
import { mockUser, mockSuperAdminUser } from '@/test/fixtures'

// Mock child components
vi.mock('./components/user-table', () => ({
  UserTable: () => <div data-testid="user-table">UserTable</div>,
}))
vi.mock('./components/project-table', () => ({
  ProjectTable: () => <div data-testid="project-table">ProjectTable</div>,
}))
vi.mock('./components/team-table', () => ({
  TeamTable: () => <div data-testid="team-table">TeamTable</div>,
}))
vi.mock('./components/csp-guide', () => ({
  CSPGuide: () => <div data-testid="csp-guide">CSPGuide</div>,
}))

// Mock useAuth
const useAuthMock = vi.fn()
vi.mock('@/hooks/use-auth', () => ({
  useAuth: () => useAuthMock(),
}))

describe('Management Page', () => {
  it('loading state', () => {
    useAuthMock.mockReturnValue({ loading: true })
    render(<Management />)
    expect(screen.getByText('Loading...')).toBeInTheDocument()
  })

  it('renders correctly for superadmin', () => {
    useAuthMock.mockReturnValue({
      user: mockSuperAdminUser,
      loading: false,
    })
    render(<Management />)

    expect(screen.getByText('Management')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /manage users/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /manage projects/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /manage teams/i })).toBeInTheDocument()
  })

  it('renders correctly for regular user', () => {
    useAuthMock.mockReturnValue({
      user: mockUser,
      loading: false,
    })
    render(<Management />)

    expect(screen.getByText('Management')).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: /manage users/i })).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: /manage projects/i })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: /manage teams/i })).not.toBeInTheDocument()
  })

  it('switches tabs for superadmin', () => {
    useAuthMock.mockReturnValue({
      user: mockSuperAdminUser,
      loading: false,
    })
    render(<Management />)

    // Check default tab 
    expect(screen.getByTestId('project-table')).toBeInTheDocument()

    // Switch to Users
    fireEvent.click(screen.getByRole('button', { name: /manage users/i }))
    expect(screen.getByTestId('user-table')).toBeInTheDocument()
    expect(screen.queryByTestId('project-table')).not.toBeInTheDocument()

    // Switch to Teams
    fireEvent.click(screen.getByRole('button', { name: /manage teams/i }))
    expect(screen.getByTestId('team-table')).toBeInTheDocument()
  })
})
