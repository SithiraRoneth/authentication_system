import axios from 'axios';
import { createContext, useContext, useState, type ReactNode } from 'react';

type AuthContextType = {
    user: { email: string } | null;
    login: (email: string, password: string) => void;
    register: (email: string, password: string) => void;
    logout: () => void;
  };
  
  const AuthContext = createContext<AuthContextType | undefined>(undefined);
  
  export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [user, setUser] = useState<{ email: string } | null>(null);
  
    const login = (email: string, password: string) => {
      setUser({ email });
    };
  
    const register = async (email: string, password: string) => {
  try {
    const res = await axios.post("http://localhost:8080/api/user/", {
      email,
      password,
    });

    // Optional: adjust this if backend returns different data
    if (res.status === 200 || res.status === 201) {
      setUser({ email });
    } else {
      alert("Registration failed");
    }
  } catch (error: any) {
    console.error("Registration error:", error);
    alert(error.response?.data?.message || "Registration failed");
  }
};
  
    const logout = () => setUser(null);
  
    return (
      <AuthContext.Provider value={{ user, login, register, logout }}>
        {children}
      </AuthContext.Provider>
    );
  };
  export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within an AuthProvider');
  return context;
};