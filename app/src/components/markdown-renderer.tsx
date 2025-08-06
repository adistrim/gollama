import * as React from 'react';
import ReactMarkdown from 'react-markdown';

interface MarkdownRendererProps {
  content: string;
  className?: string;
}

export const MarkdownRenderer: React.FC<MarkdownRendererProps> = ({ content }) => {
  return (
    <ReactMarkdown
      components={{
        pre: ({ children, ...props }) => (
          <pre className="bg-muted p-2 rounded-md overflow-auto" {...props}>
            {children}
          </pre>
        ),
        code: ({ children, ...props }) => (
          <code className="bg-muted px-1 py-0.5 rounded-sm font-mono text-sm" {...props}>
            {children}
          </code>
        )
      }}
    >
      {content}
    </ReactMarkdown>
  );
};
