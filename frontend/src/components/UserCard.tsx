import React, { memo } from 'react';
import { User } from '@/types';

interface UserCardProps {
  user: User;
  onEdit: (userId: string) => void;
  onDelete: (userId: string) => void;
}

export const UserCard = memo(({ user, onEdit, onDelete }: UserCardProps) => {
  return (
    <div className="p-4 border rounded-lg">
      <h3 className="font-semibold">{user.first_name} {user.last_name}</h3>
      <p className="text-gray-600">{user.email}</p>
      <div className="mt-2 space-x-2">
        <button
          onClick={() => onEdit(user.id)}
          className="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Edit
        </button>
        <button
          onClick={() => onDelete(user.id)}
          className="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600"
        >
          Delete
        </button>
      </div>
    </div>
  );
});

UserCard.displayName = 'UserCard';