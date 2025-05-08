import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function Header() {
  const { user } = useAuth();

  return (
    <nav className="bg-white shadow p-4 flex justify-between items-center">
      <h1 className="text-xl font-bold">MyApp</h1>
      <div className="space-x-4">
        {user ? (
          <>
            <Link to="/dashboard" className="text-blue-500">Dashboard</Link>
            <Link to="/profile" className="text-blue-500">Profile</Link>
          </>
        ) : (
          <>
            <Link to="/login" className="text-blue-500">Login</Link>
            <Link to="/register" className="text-blue-500">Register</Link>
          </>
        )}
      </div>
    </nav>
  );
}
