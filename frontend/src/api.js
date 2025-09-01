const USER_API = "http://localhost:8081"; 
const CHAT_API = "http://localhost:8082"; 

export async function register(name, password) {
  const res = await fetch(`${USER_API}/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, password }),
  });
  if (!res.ok) throw new Error("Registration failed");
  return res.json();
}

export async function login(name, password) {
  const res = await fetch(`${USER_API}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, password }),
  });
  if (!res.ok) throw new Error("Invalid credentials");
  return res.json();
}

export async function loginAnon() {
  const res = await fetch(`${USER_API}/login/anon`, { method: "POST" });
  if (!res.ok) throw new Error("Anon login failed");
  return res.json();
}

export async function getMessages(token) {
  const res = await fetch(`${CHAT_API}/messages`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw new Error("Failed to fetch messages");
  return res.json();
}

export async function sendMessage(token, content) {
  const res = await fetch(`${CHAT_API}/messages`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ content }),
  });
  if (!res.ok) throw new Error("Failed to send message");
  return res.json();
}

export function connectWS(token, onMessage) {
  const ws = new WebSocket(`ws://localhost:8082/ws?token=${token}`);
  ws.onmessage = (e) => {
    try {
      onMessage(JSON.parse(e.data));
    } catch (err) {
      console.error("WS parse error", err);
    }
  };
  ws.onclose = () => console.log("WS closed");
  ws.onerror = (e) => console.error("WS error", e);
  return ws;
}
