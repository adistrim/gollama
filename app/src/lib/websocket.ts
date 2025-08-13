import { WS_URL } from '@/config';
import { useEffect, useState } from 'react';

export type ChatMessage = {
  content?: string;
  response?: string;
  sessionId?: string;
  error?: string;
  isProcessing?: boolean;
};

export function useWebSocket() {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [sessionId, setSessionId] = useState<string>('');
  const [connected, setConnected] = useState(false);
  const [processingIndex, setProcessingIndex] = useState<number | null>(null);

  useEffect(() => {
    const ws = new WebSocket(WS_URL);
    
    ws.onopen = () => {
      console.log('Connected to chat server');
      setConnected(true);
      setSocket(ws);
    };
    
    ws.onclose = () => {
      console.log('Disconnected from chat server');
      setConnected(false);
      setSocket(null);
    };
    
    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      setConnected(false);
    };
    
    return () => {
      ws.close();
    };
  }, []);

  useEffect(() => {
    if (!socket) return;
    
    const handleMessage = (event: MessageEvent) => {
      try {
        const messageRaw = JSON.parse(event.data);
        const message: ChatMessage = {
          ...messageRaw,
          sessionId: messageRaw.session_id,
          isProcessing: messageRaw.is_processing,
        };
        
        if (message.sessionId) setSessionId(message.sessionId);
        
        if (message.isProcessing) {
          setMessages((prev) => {
            if (processingIndex !== null) {
              const newMessages = [...prev];
              newMessages[processingIndex] = message;
              return newMessages;
            } else {
              setProcessingIndex(prev.length);
              return [...prev, message];
            }
          });
        } else {
          setMessages((prev) => {
            const filteredMessages = prev.filter(msg => !msg.isProcessing);
            return [...filteredMessages, message];
          });
          setProcessingIndex(null);
        }
      } catch (error) {
        console.error('Error handling message:', error);
      }
    };
    
    socket.addEventListener('message', handleMessage);
    
    return () => {
      socket.removeEventListener('message', handleMessage);
    };
  }, [socket, processingIndex]);

  const sendMessage = (content: string) => {
    if (socket && connected) {
      const message = { content, session_id: sessionId };
      socket.send(JSON.stringify(message));
      setMessages((prev) => [...prev, { content }]);
      setMessages((prev) => [...prev, { isProcessing: true, response: "Thinking..." }]);
      setProcessingIndex(messages.length + 1);
    }
  };

  return { messages, sendMessage, connected };
}
