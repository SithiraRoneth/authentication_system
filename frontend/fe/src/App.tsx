import type { ReactNode } from 'react';
import './App.css'
import Header from './components/Header';
import { useAuth } from './context/AuthContext';
import { Navigate, Route, Routes } from 'react-router-dom';
import Register from './pages/RegisterForm';
import Login from './pages/LoginForm';
import Dashboard from './pages/DashBoard';

const ProtectedRoute = ({ children }: { children: ReactNode }) => {
  const { user } = useAuth();
  return user ? children : <Navigate to="/login" />;
};

function App() {
  return (
    <div className="min-h-screen bg-gray-100">
      <Header />
      <Routes>
        <Route path="/register" element={<Register />} />
        <Route path="/login" element={<Login />} />
        <Route path="/dashboard" element={<ProtectedRoute><Dashboard /></ProtectedRoute>} />
        <Route path="*" element={<Navigate to="/login" />} />
      </Routes>
    </div>
  );
}

export default App;

