import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { authService } from '@/services/auth';
import { toast } from 'sonner';
import type { UserRegister } from '@/types';
import { useAppContext } from '@/context/AppContext';

export function useAuth() {
  const { login, logout, state } = useAppContext();
  const [localLoading, setLocalLoading] = useState(false);
  const [localError, setLocalError] = useState<string | null>(null);
  const navigate = useNavigate();

  const register = async (data: UserRegister) => {
    setLocalLoading(true);
    setLocalError(null);
    try {
      await authService.register(data);
      toast.success('Account created successfully. Please sign in.');
      navigate('/login');
    } catch (err: any) {
      console.error('Registration failed', err);
      const message = err.response?.data?.message || 'Failed to create account. Please try again.';
      setLocalError(message);
      toast.error(message);
    } finally {
      setLocalLoading(false);
    }
  };

  return {
    login,
    register,
    logout,
    loading: state.isLoading || localLoading,
    error: state.error || localError,
    user: state.user,
  };
}

