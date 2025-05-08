import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';

export default function Profile() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="max-w-md mx-auto mt-10 bg-white p-6 rounded shadow">
      <h2 className="text-2xl mb-4">Your Profile</h2>
      <p className="mb-6">Logged in as: <strong>{user?.email}</strong></p>
      <button
        onClick={handleLogout}
        className="w-full bg-red-500 text-white p-2 rounded"
      >
        Logout
      </button>
    </div>
  );
}
