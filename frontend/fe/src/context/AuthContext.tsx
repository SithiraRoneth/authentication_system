import { createContext, useContext, useState, type ReactNode} from "react";

type AuthContextType = {
  user: { email: string } | null;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string) => Promise<void>;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<{ email: string } | null>(null);

  const login = async (email: string, password: string) => {
    const res = await fetch("http://localhost:8080/api/user/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({ username: email, password }),
    });
    if (!res.ok) throw new Error("Login failed");

    const data = await res.json();
    localStorage.setItem("token", data.token);
    setUser({ email });
  };

  const register = async (email: string, password: string) => {
    const res = await fetch("http://localhost:8080/api/user/", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
    if (!res.ok) throw new Error("Registration failed");
    setUser({ email });
  };

  const logout = () => {
    localStorage.removeItem("token");
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) throw new Error("useAuth must be used within AuthProvider");
  return context;
};
