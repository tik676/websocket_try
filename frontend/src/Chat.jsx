import { useEffect, useRef, useState } from "react";

export default function Chat({ token, wsBase }) {
  const [messages, setMessages] = useState([]);
  const [content, setContent] = useState("");
  const wsRef = useRef(null);
  const endRef = useRef(null);

  useEffect(() => {
    fetch("http://127.0.0.1:8082/messages?limit=50", {
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    })
      .then((r) => r.json())
      .then((data) => setMessages(Array.isArray(data) ? data : []))
      .catch(console.error);

    const url = token
      ? `${wsBase}/ws?token=${encodeURIComponent(token)}`
      : `${wsBase}/ws`;

    const ws = new WebSocket(url);
    wsRef.current = ws;

    ws.onopen = () => {
    };

    ws.onmessage = (ev) => {
      try {
        const msg = JSON.parse(ev.data);
        setMessages((prev) => [...prev, msg]);
      } catch {
      }
    };

    ws.onclose = () => {
    };

    ws.onerror = (e) => {
    };

    return () => {
      ws.close();
    };
  }, [token, wsBase]);

  const handleSend = () => {
    const text = content.trim();
    if (!text) return;

    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(text);
      setContent("");
    } else {
      fetch("http://127.0.0.1:8082/messages", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
        body: JSON.stringify({ content: text }),
      })
        .then((r) => r.json())
        .then((msg) => {
          setMessages((prev) => [...prev, msg]);
          setContent("");
        })
        .catch(console.error);
    }
  };

  useEffect(() => {
    endRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  return (
    <div style={{ maxWidth: 600, margin: "auto", padding: 20 }}>
      <h2>Chat</h2>
      <div style={{ height: 400, overflowY: "auto", border: "1px solid #ccc", padding: 10 }}>
        {messages.map((m, i) => (
          <div key={i}>
            <b>{m.user_name || "Anon"}:</b> {m.content}
          </div>
        ))}
        <div ref={endRef} />
      </div>
      <input
        value={content}
        onChange={(e) => setContent(e.target.value)}
        onKeyDown={(e) => e.key === "Enter" && handleSend()}
        placeholder="Message..."
        style={{ width: "80%", padding: 5, marginTop: 5 }}
      />
      <button onClick={handleSend} style={{ marginLeft: 5, padding: 5 }}>
        Send
      </button>
    </div>
  );
}
