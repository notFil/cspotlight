import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import Auth from './index'

// Mock useAuth
const useAuthMock = vi.fn()
const loginMock = vi.fn()

vi.mock('@/hooks/use-auth', () => ({
  useAuth: () => useAuthMock(),
}))

// Mock ThemeToggle to avoid theme context issues if any
vi.mock('@/components/common/theme-toggle', () => ({
  default: () => <div data-testid="theme-toggle">ThemeToggle</div>,
}))

describe('Auth Page', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    useAuthMock.mockReturnValue({
      login: loginMock,
      loading: false,
      error: null,
    })
  })

  it('renders login form correctly', () => {
    render(
      <MemoryRouter>
        <Auth />
      </MemoryRouter>
    )
    expect(screen.getByText(/Welcome to cspotlight!/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/Username/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/Password/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /Sign In/i })).toBeInTheDocument()
  })

  it('handles input changes', () => {
    render(
      <MemoryRouter>
        <Auth />
      </MemoryRouter>
    )
    const usernameInput = screen.getByLabelText(/Username/i)
    const passwordInput = screen.getByLabelText(/Password/i)

    fireEvent.change(usernameInput, { target: { value: 'testuser' } })
    fireEvent.change(passwordInput, { target: { value: 'password123' } })

    expect((usernameInput as HTMLInputElement).value).toBe('testuser')
    expect((passwordInput as HTMLInputElement).value).toBe('password123')
  })

  it('toggles password visibility', () => {
    render(
      <MemoryRouter>
        <Auth />
      </MemoryRouter>
    )
    const passwordInput = screen.getByLabelText(/Password/i)
    const toggleButton = screen.getByRole('button', { name: /Show password/i })

    expect(passwordInput).toHaveAttribute('type', 'password')

    fireEvent.click(toggleButton)
    expect(passwordInput).toHaveAttribute('type', 'text')
    expect(screen.getByRole('button', { name: /Hide password/i })).toBeInTheDocument()

    fireEvent.click(toggleButton)
    expect(passwordInput).toHaveAttribute('type', 'password')
  })

  it('submits the form with correct data', async () => {
    render(
      <MemoryRouter>
        <Auth />
      </MemoryRouter>
    )
    const usernameInput = screen.getByLabelText(/Username/i)
    const passwordInput = screen.getByLabelText(/Password/i)
    const submitButton = screen.getByRole('button', { name: /Sign In/i })

    fireEvent.change(usernameInput, { target: { value: 'testuser' } })
    fireEvent.change(passwordInput, { target: { value: 'password123' } })
    fireEvent.click(submitButton)

    await waitFor(() => {
      expect(loginMock).toHaveBeenCalledTimes(1)
      expect(loginMock).toHaveBeenCalledWith({
        username: 'testuser',
        password: 'password123',
      })
    })
  })

  it('displays loading state', () => {
    useAuthMock.mockReturnValue({
      login: loginMock,
      loading: true,
      error: null,
    })

    render(
      <MemoryRouter>
        <Auth />
      </MemoryRouter>
    )

    expect(screen.getByRole('button', { name: /Signing in/i })).toBeDisabled()
  })

  it('displays error message', () => {
    useAuthMock.mockReturnValue({
      login: loginMock,
      loading: false,
      error: 'Invalid credentials',
    })

    render(
      <MemoryRouter>
        <Auth />
      </MemoryRouter>
    )

    expect(screen.getByText(/Invalid credentials/i)).toBeInTheDocument()
  })
})
