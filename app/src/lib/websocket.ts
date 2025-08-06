import { WS_URL } from '@/config';
import { useEffect, useState } from 'react';

export type ChatMessage = {
  content?: string;
  response?: string;
  sessionId?: string;
  error?: string;
};

export function useWebSocket() {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [sessionId, setSessionId] = useState<string>('');
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    const ws = new WebSocket(WS_URL);
    
    ws.onopen = () => {
      console.log('Connected to chat server');
      setConnected(true);
      setSocket(ws);
    };
    
    ws.onmessage = (event) => {
      const message = JSON.parse(event.data) as ChatMessage;
      if (message.sessionId) setSessionId(message.sessionId);
      setMessages((prev) => [...prev, message]);
    };
    
    ws.onclose = () => {
      console.log('Disconnected from chat server');
      setConnected(false);
    };
    
    return () => {
      ws.close();
    };
  }, []);

  const sendMessage = (content: string) => {
    if (socket && connected) {
      const message = { content, session_id: sessionId };
      socket.send(JSON.stringify(message));
      setMessages((prev) => [...prev, { content }]);
    }
  };

  return { messages, sendMessage, connected };
}
