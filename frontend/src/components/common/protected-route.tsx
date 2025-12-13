import { Navigate, Outlet } from 'react-router-dom';
import { useAppContext } from '@/context/AppContext';
import { LoadingPage } from './loading-page';

const ProtectedRoute = () => {
  const { state } = useAppContext();
  const { user, isLoading } = state;

  if (isLoading) {
    return <LoadingPage />;
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return <Outlet />;
};

export default ProtectedRoute;
