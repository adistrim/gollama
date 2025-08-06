import { Button } from '@/components/ui/button'
import { Link } from '@tanstack/react-router'
import { PaperPlaneIcon, CodeIcon, GitHubLogoIcon } from '@radix-ui/react-icons'

export default function HomePage() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen p-6 bg-background">
      <div className="text-center space-y-6 max-w-sm">
        <div className="flex justify-center gap-2">
          <div className="p-3 rounded-full bg-primary/10 text-primary">
            <CodeIcon className="h-6 w-6" />
          </div>
          <div className="p-3 rounded-full bg-primary/10 text-primary">
            <GitHubLogoIcon className="h-6 w-6" />
          </div>
        </div>
        
        <h1 className="text-4xl font-bold">Gollama</h1>
        <p className="text-muted-foreground">
          AI software engineer powered by MCP server and gpt-oss with tools for GitHub
        </p>
        
        <Link to="/chat" className="block pt-4">
          <Button size="lg" className="gap-2 w-full">
            Start Coding <PaperPlaneIcon className="h-4 w-4" />
          </Button>
        </Link>
      </div>
    </div>
  )
}
