import { render, screen } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import Profile from './index'
import { mockUser } from '@/test/fixtures'

// Mock child components
vi.mock('./components/user-info-card', () => ({
  UserInfoCard: ({ name }: { name: string }) => <div data-testid="user-info-card">{name}</div>,
}))
vi.mock('./components/password-update-form', () => ({
  PasswordUpdateForm: () => <div data-testid="password-update-form">PasswordUpdateForm</div>,
}))
vi.mock('./components/avatar-upload', () => ({
  AvatarUpload: () => <div data-testid="avatar-upload">AvatarUpload</div>,
}))

// Mock useAuth
const useAuthMock = vi.fn()
vi.mock('@/hooks/use-auth', () => ({
  useAuth: () => useAuthMock(),
}))

describe('Profile Page', () => {
  it('renders loading when no user', () => {
    useAuthMock.mockReturnValue({ user: null })
    render(<Profile />)
    expect(screen.getByText('Loading...')).toBeInTheDocument()
  })

  it('renders profile content when user exists', () => {
    useAuthMock.mockReturnValue({
      user: mockUser,
    })
    render(<Profile />)

    expect(screen.getByText('Profile')).toBeInTheDocument()
    expect(screen.getByTestId('avatar-upload')).toBeInTheDocument()
    expect(screen.getByTestId('user-info-card')).toHaveTextContent('John Doe')
    expect(screen.getByTestId('password-update-form')).toBeInTheDocument()
  })
})
