import { useEffect, useState } from "react";
import { auth } from "@/lib/firebase";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { fetchNotes, syncUser } from "@/services/api";
import { Link } from "react-router-dom";
import { Settings, Users } from "lucide-react";

interface Note {
  id: string;
  title: string;
  updated_at: string;
}

export default function Dashboard() {
  const [notes, setNotes] = useState<Note[]>([]);
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    syncUser();
  }, []);

  useEffect(() => {
      setLoading(true);
      fetchNotes(page, 20, search).then(newNotes => {
          if (page === 1) {
              setNotes(newNotes);
          } else {
              setNotes(prev => [...prev, ...newNotes]);
          }
          if (newNotes.length < 20) {
              setHasMore(false);
          } else {
              setHasMore(true);
          }
      }).catch(console.error).finally(() => setLoading(false));
  }, [page, search]);

  const handleSearch = (val: string) => {
      setSearch(val);
      setPage(1);
      setHasMore(true);
  };

  return (
    <div className="p-4 md:p-8 max-w-7xl mx-auto">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-2xl font-bold">My Journal</h1>
        <div className="flex gap-2">
          <Link to="/connections">
            <Button variant="ghost" size="icon" title="Connections">
              <Users className="w-5 h-5" />
            </Button>
          </Link>
          <Link to="/settings">
            <Button variant="ghost" size="icon" title="Settings">
              <Settings className="w-5 h-5" />
            </Button>
          </Link>
          <Button variant="outline" onClick={() => auth.signOut()}>
            Sign Out
          </Button>
        </div>
      </div>

      <div className="mb-6">
        <Input
          placeholder="Search notes..."
          value={search}
          onChange={(e) => handleSearch(e.target.value)}
          className="max-w-md"
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {notes.length === 0 && !loading && <p className="text-gray-500 col-span-full">No notes found.</p>}
        {notes.map((note) => (
          <Link key={note.id} to={`/notes/${note.id}`} className="block">
            <div className="p-6 bg-white hover:bg-gray-50 transition rounded-lg shadow border h-40 flex flex-col">
                <h3 className="font-semibold mb-2 line-clamp-2">{note.title || "Untitled Note"}</h3>
                <p className="text-gray-400 text-xs mt-auto">{new Date(note.updated_at).toLocaleDateString()}</p>
            </div>
          </Link>
        ))}
      </div>

      {hasMore && notes.length > 0 && (
          <div className="mt-8 text-center">
              <Button onClick={() => setPage(p => p + 1)} disabled={loading} variant="outline">
                  {loading ? "Loading..." : "Load More"}
              </Button>
          </div>
      )}

      <div className="fixed bottom-8 right-8 flex flex-col gap-4">
        <Link to="/chat">
            <Button className="w-14 h-14 rounded-full shadow-lg" size="icon" variant="secondary" title="Ask AI">
                🤖
            </Button>
        </Link>
        <Link to="/notes/new">
            <Button className="w-14 h-14 rounded-full shadow-lg text-2xl" size="icon" title="New Note">
                +
            </Button>
        </Link>
      </div>
    </div>
  );
}
