import { useState } from "react";
import { register } from "./api";

export default function Register({ switchToLogin }) {
  const [name, setName] = useState("");
  const [password, setPassword] = useState("");

  const handleRegister = async () => {
    try {
      await register(name, password);
      alert("Registered successfully");
      switchToLogin();
    } catch {
      alert("Registration failed");
    }
  };

  return (
    <div className="auth-container">
      <h2>Register</h2>
      <input placeholder="Name" value={name} onChange={e => setName(e.target.value)} />
      <input placeholder="Password" type="password" value={password} onChange={e => setPassword(e.target.value)} />
      <button onClick={handleRegister}>Register</button>
      <p onClick={switchToLogin} style={{cursor: "pointer"}}>Go to Login</p>
    </div>
  );
}
