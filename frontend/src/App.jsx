import { useEffect, useMemo, useState } from "react";
import Login from "./Login";
import Register from "./Register";
import Chat from "./Chat";

export default function App() {
  const [token, setToken] = useState(() => localStorage.getItem("jwt"));
  const [view, setView] = useState(() => (localStorage.getItem("jwt") ? "chat" : "login"));

  const wsBase = useMemo(() => "ws://127.0.0.1:8082", []);

  const handleLogin = (t) => {
    localStorage.setItem("jwt", t);
    setToken(t);
    setView("chat");
  };

  const handleLogout = () => {
    localStorage.removeItem("jwt");
    setToken(null);
    setView("login");
  };

  useEffect(() => {
    const onStorage = (e) => {
      if (e.key === "jwt") {
        setToken(e.newValue);
        setView(e.newValue ? "chat" : "login");
      }
    };
    window.addEventListener("storage", onStorage);
    return () => window.removeEventListener("storage", onStorage);
  }, []);

  if (view === "login")
    return <Login onLogin={handleLogin} switchToRegister={() => setView("register")} />;
  if (view === "register")
    return <Register switchToLogin={() => setView("login")} />;

  return (
    <div>
      <button onClick={handleLogout} style={{ float: "right", margin: 10 }}>
        Logout
      </button>
      <Chat token={token} wsBase={wsBase} />
    </div>
  );
}
