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
  
    const register = (email: string, password: string) => {
      setUser({ email });
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