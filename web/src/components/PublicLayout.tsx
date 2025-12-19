import { Link, Outlet } from "react-router-dom";
import { Book } from "lucide-react";

export default function PublicLayout() {
  return (
    <div className="min-h-screen bg-background font-sans antialiased">
      <header className="border-b border-border bg-card">
        <div className="container mx-auto max-w-4xl px-4 h-16 flex items-center justify-between">
          <Link to="/" className="flex items-center space-x-2 text-foreground hover:opacity-80 transition-opacity">
            <Book className="h-6 w-6" />
            <span className="text-xl font-bold">Discipleship Journal</span>
          </Link>
          <nav>
            <Link to="/login" className="text-sm font-medium hover:underline underline-offset-4">
              Sign In
            </Link>
          </nav>
        </div>
      </header>
      <main className="container mx-auto max-w-4xl px-4 py-8">
        <Outlet />
      </main>
      <footer className="border-t border-border py-6 bg-muted/30">
        <div className="container mx-auto max-w-4xl px-4 flex flex-col md:flex-row justify-between items-center gap-4 text-sm text-muted-foreground">
          <p>&copy; {new Date().getFullYear()} Discipleship Journal. All rights reserved.</p>
          <div className="flex gap-4">
            <Link to="/privacy-policy" className="hover:underline underline-offset-4">Privacy Policy</Link>
            <Link to="/terms-of-service" className="hover:underline underline-offset-4">Terms of Service</Link>
          </div>
        </div>
      </footer>
    </div>
  );
}
