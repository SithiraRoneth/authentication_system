import {
  createContext,
  useContext,
  useState,
  useEffect,
  type ReactNode,
} from "react";

type AuthContextType = {
  user: { email: string } | null;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string) => Promise<void>;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<{ email: string } | null>(null);

  // ✅ Token check every time app starts or reloads
  const verifyToken = async () => {
    const token = localStorage.getItem("token");
    if (!token) {
      setUser(null);
      return;
    }

    try {
      const res = await fetch("http://localhost:8080/api/user/me", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      if (!res.ok) throw new Error("Unauthorized");
      const data = await res.json();
      setUser({ email: data.email });
    } catch (err) {
      console.warn("Invalid token, logging out");
      localStorage.removeItem("token");
      setUser(null);
    }
  };

  useEffect(() => {
    verifyToken();
  }, []);

  const login = async (email: string, password: string) => {
    const res = await fetch("http://localhost:8080/api/user/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
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

  // Do not set the user here
  // Instead, inform the user that registration was successful and navigate to the login page
  alert("Registration successful");
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
