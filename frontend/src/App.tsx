import { Routes, Route, Navigate } from 'react-router-dom';
import Dashboard from './pages/dashboard';
import Auth from './pages/auth';
import Signup from './pages/auth/signup';
import Layout from './layout';
import Management from './pages/management';
import Projects from './pages/projects';
import Profile from './pages/profile';
import Reports from "./pages/reports";
import ProtectedRoute from './components/common/protected-route';
import { Toaster } from "@/components/ui/sonner";
import '@/services/auth';

function App() {
  return (
    <>
      <Routes>
        <Route path="/login" element={<Auth />} />
        <Route path="/signup" element={<Signup />} />

        <Route element={<ProtectedRoute />}>
          <Route element={<Layout><Dashboard /></Layout>} path="/" />
          <Route element={<Layout><Dashboard /></Layout>} path="/project/:id" />
          <Route element={<Layout><Management /></Layout>} path="/management" />
          <Route element={<Layout><Projects /></Layout>} path="/projects" />
          <Route element={<Layout><Profile /></Layout>} path="/profile" />
          <Route element={<Layout><Reports /></Layout>} path="/reports/:projectId" />
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
      <Toaster />
    </>
  );
}

export default App;
