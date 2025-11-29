import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Dashboard from './pages/dashboard';
import Auth from './pages/auth';
import Layout from './layout';
import Management from './pages/management';
import Projects from './pages/projects';
import Profile from './pages/profile';
import ProtectedRoute from './components/common/protected-route';
import '@/services/auth'; // Import to ensure AuthService is initialized and callback registered

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Auth />} />

        <Route element={<ProtectedRoute />}>
          <Route element={<Layout><Dashboard /></Layout>} path="/" />
          <Route element={<Layout><Dashboard /></Layout>} path="/project/:id" />
          <Route element={<Layout><Management /></Layout>} path="/management" />
          <Route element={<Layout><Projects /></Layout>} path="/projects" />
          <Route element={<Layout><Profile /></Layout>} path="/profile" />
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Router>
  );
}

export default App;
