import { render, screen, fireEvent } from '@testing-library/react';
import { UserCard } from '../UserCard';
import type { User } from '@/types';

const mockUser: User = {
  id: '26997316-ccb6-46e1-b224-712541788a53',
  first_name: 'John',
  last_name: 'Doe',
  username: 'johnd',
  email: 'john@example.com'
};

describe('UserCard', () => {
  it('renders user information', () => {
    const mockOnEdit = jest.fn();
    const mockOnDelete = jest.fn();

    render(
      <UserCard
        user={mockUser}
        onEdit={mockOnEdit}
        onDelete={mockOnDelete}
      />
    );

    expect(screen.getByText('John Doe')).toBeInTheDocument();
    expect(screen.getByText('john@example.com')).toBeInTheDocument();
  });

  it('calls onEdit when edit button is clicked', () => {
    const mockOnEdit = jest.fn();
    const mockOnDelete = jest.fn();

    render(
      <UserCard
        user={mockUser}
        onEdit={mockOnEdit}
        onDelete={mockOnDelete}
      />
    );

    fireEvent.click(screen.getByText('Edit'));
    expect(mockOnEdit).toHaveBeenCalledWith('1');
  });
});