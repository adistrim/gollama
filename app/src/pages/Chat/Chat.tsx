import { useState, useRef, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { useWebSocket } from '@/lib/websocket';
import { MarkdownRenderer } from '@/components/markdown-renderer';
import { PaperPlaneIcon } from '@radix-ui/react-icons';
import { Textarea } from '@/components/ui/textarea';
import { ScrollArea } from '@/components/ui/scroll-area';

export default function ChatPage() {
  const [inputValue, setInputValue] = useState('');
  const { messages, sendMessage, connected } = useWebSocket();
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  useEffect(() => {
    textareaRef.current?.focus();
  }, []);

  const handleSendMessage = () => {
    if (inputValue.trim() && connected) {
      sendMessage(inputValue);
      setInputValue('');
      if (textareaRef.current) {
        textareaRef.current.style.height = 'auto';
      }
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  const handleTextareaChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInputValue(e.target.value);
    e.target.style.height = 'auto';
    e.target.style.height = `${Math.min(e.target.scrollHeight, 200)}px`;
  };

  const isBotThinking = () => {
    if (messages.length === 0) return false;
    const lastMessage = messages[messages.length - 1];
    return !!lastMessage.content && !lastMessage.response && !lastMessage.error;
  };

  return (
    <div className="flex flex-col h-screen bg-background">
      {/* Main content area */}
      <main className="flex-1 overflow-hidden">
        {/* Welcome screen for empty chat */}
        {messages.length === 0 && (
          <div className="absolute inset-0 flex flex-col items-center justify-center text-center p-8">
            <div className="max-w-md space-y-4">
              <h2 className="text-2xl font-semibold tracking-tight">Gollama Chat</h2>
              <p className="text-muted-foreground">
                Start a conversation with Gollama. Ask a question or describe what you need help with.
              </p>
            </div>
          </div>
        )}

        {/* Messages area */}
        {messages.length > 0 && (
          <ScrollArea className="h-full pb-32">
            <div className="max-w-2xl mx-auto py-8 px-4">
              {messages.map((message, index) => (
                <div 
                  key={index} 
                  className="mb-6"
                >
                  {message.content ? (
                    /* User message - right aligned with background */
                    <div className="flex justify-end">
                      <div className="bg-primary text-primary-foreground px-4 py-3 rounded-xl max-w-[80%]">
                        <div className="whitespace-pre-wrap">{message.content}</div>
                      </div>
                    </div>
                  ) : (
                    /* AI message - left aligned */
                    <div className="flex">
                      <div className="max-w-[80%]">
                        {message.response ? (
                          <MarkdownRenderer content={message.response} />
                        ) : (
                          <div className="text-destructive">{message.error || "An error occurred"}</div>
                        )}
                      </div>
                    </div>
                  )}
                </div>
              ))}
              
              {/* Bot typing indicator */}
              {isBotThinking() && (
                <div className="flex mb-6">
                  <div className="max-w-[80%]">
                    <div className="text-muted-foreground flex items-center gap-2 px-4 py-3">
                      <div className="size-4 rounded-full border-2 border-muted-foreground/30 border-t-muted-foreground animate-spin"></div>
                      <span className="text-sm">Thinking...</span>
                    </div>
                  </div>
                </div>
              )}
                            
              <div ref={messagesEndRef} />
            </div>
          </ScrollArea>
        )}
      </main>

      {/* Input area */}
      <div className="fixed bottom-0 left-0 right-0 bg-background border-t border-border">
        <div className="max-w-2xl mx-auto p-4">
          <div className="relative">
            <Textarea
              ref={textareaRef}
              value={inputValue}
              onChange={handleTextareaChange}
              onKeyDown={handleKeyDown}
              placeholder="Message Gollama..."
              className="min-h-[60px] max-h-[200px] pr-12 resize-none rounded-lg focus-visible:ring-1 focus-visible:ring-primary/50"
              disabled={!connected}
              rows={1}
            />
            <Button 
              onClick={handleSendMessage}
              disabled={!connected || !inputValue.trim()}
              className="absolute right-2 bottom-2 rounded-md p-2 size-9"
              size="icon"
              variant={inputValue.trim() ? "default" : "ghost"}
            >
              <PaperPlaneIcon className="h-4 w-4" />
            </Button>
          </div>
          <div className="text-xs text-center text-muted-foreground mt-2">
            Gollama can make mistakes. <span>Powered by gpt-oss</span>
          </div>
        </div>
      </div>
    </div>
  );
}
