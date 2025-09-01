import { useState } from "react";
import { login, loginAnon } from "./api";

export default function Login({ onLogin, switchToRegister }) {
  const [name, setName] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = async () => {
    try {
      const data = await login(name, password);
      onLogin(data.access_token);
    } catch {
      alert("Login failed");
    }
  };

  const handleAnon = async () => {
    try {
      const data = await loginAnon();
      onLogin(data.access_token);
    } catch {
      alert("Anon login failed");
    }
  };

  return (
    <div className="auth-container">
      <h2>Login</h2>
      <input placeholder="Name" value={name} onChange={e => setName(e.target.value)} />
      <input placeholder="Password" type="password" value={password} onChange={e => setPassword(e.target.value)} />
      <button onClick={handleLogin}>Login</button>
      <button onClick={handleAnon}>Login Anon</button>
      <p onClick={switchToRegister} style={{cursor: "pointer"}}>Go to Register</p>
    </div>
  );
}
