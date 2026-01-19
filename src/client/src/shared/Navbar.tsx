import { useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

export function Navbar() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  return (
    <nav className="navbar">
      <h1>📝 Todos</h1>
      <div className="navbar-right">
        <span className="user-email">{user?.email}</span>
        <button className="btn-logout" onClick={handleLogout}>
          Logout
        </button>
      </div>
    </nav>
  );
}
