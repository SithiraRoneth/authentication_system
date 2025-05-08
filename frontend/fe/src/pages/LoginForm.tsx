import { useState } from "react";
import { useAuth } from "../context/AuthContext";
import { useNavigate } from "react-router-dom";

export default function Login() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const { login } = useAuth();
    const navigate = useNavigate();
  
    const handleSubmit = (e: React.FormEvent) => {
      e.preventDefault();
      if (!email || !password) return alert("All fields required");
      login(email, password);
      navigate('/dashboard');
    };
  
    return (
      <form onSubmit={handleSubmit} className="max-w-md mx-auto mt-10 bg-white p-6 rounded shadow">
        <h2 className="text-2xl mb-4">Login</h2>
        <input
          className="w-full p-2 border mb-4"
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <input
          className="w-full p-2 border mb-4"
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <button className="w-full bg-blue-500 text-white p-2 rounded">Login</button>
      </form>
    );
  }