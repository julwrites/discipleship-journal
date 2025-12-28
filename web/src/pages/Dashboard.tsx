import { useEffect, useState, useRef } from "react";
import { auth } from "@/lib/firebase";
import { useDebounce } from "@/hooks/useDebounce";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { fetchNotes, syncUser, NoteFilter } from "@/services/api";
import { Link } from "react-router-dom";
import { Settings, Users, BookOpen, Filter, CalendarIcon } from "lucide-react";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { format } from "date-fns";
import { cn } from "@/lib/utils";

interface Note {
  id: string;
  title: string;
  updated_at: string;
}

export default function Dashboard() {
  const [notes, setNotes] = useState<Note[]>([]);
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search, 500);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);
  const prevSearchRef = useRef(debouncedSearch);

  // Filter state
  const [startDate, setStartDate] = useState<Date | undefined>();
  const [endDate, setEndDate] = useState<Date | undefined>();
  const [sortBy, setSortBy] = useState<"updated_at" | "created_at" | "title">("updated_at");
  const [sortOrder, setSortOrder] = useState<"asc" | "desc">("desc");

  // Ref to track if filters changed to reset page
  const prevFilterRef = useRef({ startDate, endDate, sortBy, sortOrder });

  useEffect(() => {
    syncUser();
  }, []);

  useEffect(() => {
      let ignore = false;

      const currentFilters = { startDate, endDate, sortBy, sortOrder };
      const filtersChanged =
          prevFilterRef.current.startDate !== startDate ||
          prevFilterRef.current.endDate !== endDate ||
          prevFilterRef.current.sortBy !== sortBy ||
          prevFilterRef.current.sortOrder !== sortOrder;

      // Handle search or filter changes
      if (prevSearchRef.current !== debouncedSearch || filtersChanged) {
          prevSearchRef.current = debouncedSearch;
          prevFilterRef.current = currentFilters;

          if (page !== 1) {
              setPage(1);
              setHasMore(true);
              // Skip fetch, wait for re-render with page=1
              return;
          } else {
              // Already at page 1, reset hasMore just in case
              setHasMore(true);
          }
      }

      const load = async () => {
          setLoading(true);
          try {
              const filter: NoteFilter = {
                  search: debouncedSearch,
                  startDate,
                  endDate,
                  sortBy,
                  sortOrder
              };
              const response = await fetchNotes(page, 20, filter);
              if (!ignore) {
                  // Check if response has data/meta structure or is just array (for backward compat if needed, though we updated API)
                  const newNotes = response.data || response;

                  if (page === 1) {
                      setNotes(newNotes);
                  } else {
                      setNotes(prev => [...prev, ...newNotes]);
                  }

                  if (response.meta) {
                      setHasMore(page < response.meta.total_pages);
                  } else {
                      // Fallback
                      if (newNotes.length < 20) {
                          setHasMore(false);
                      } else {
                          setHasMore(true);
                      }
                  }
              }
          } catch (error) {
              if (!ignore) console.error(error);
          } finally {
              if (!ignore) setLoading(false);
          }
      };
      load();
      return () => { ignore = true; };
  }, [page, debouncedSearch, startDate, endDate, sortBy, sortOrder]);

  const handleSearch = (val: string) => {
      setSearch(val);
  };

  const clearFilters = () => {
      setStartDate(undefined);
      setEndDate(undefined);
      setSortBy("updated_at");
      setSortOrder("desc");
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
          <Link to="/reading-plans">
            <Button variant="ghost" size="icon" title="Reading Plans">
              <BookOpen className="w-5 h-5" />
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

      <div className="mb-6 flex gap-2 items-center">
        <Input
          placeholder="Search notes..."
          value={search}
          onChange={(e) => handleSearch(e.target.value)}
          className="max-w-md"
        />

        <Popover>
            <PopoverTrigger asChild>
                <Button variant="outline" size="icon" title="Filter & Sort">
                    <Filter className="w-4 h-4" />
                </Button>
            </PopoverTrigger>
            <PopoverContent className="w-80 p-4">
                <div className="space-y-4">
                    <div className="space-y-2">
                        <h4 className="font-medium leading-none">Sort By</h4>
                        <div className="flex gap-2">
                            <Select value={sortBy} onValueChange={(val: "updated_at" | "created_at" | "title") => setSortBy(val)}>
                                <SelectTrigger className="w-[140px]">
                                    <SelectValue placeholder="Sort By" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="updated_at">Updated</SelectItem>
                                    <SelectItem value="created_at">Created</SelectItem>
                                    <SelectItem value="title">Title</SelectItem>
                                </SelectContent>
                            </Select>
                            <Select value={sortOrder} onValueChange={(val: "asc" | "desc") => setSortOrder(val)}>
                                <SelectTrigger className="w-[110px]">
                                    <SelectValue placeholder="Order" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="desc">Desc</SelectItem>
                                    <SelectItem value="asc">Asc</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>
                    </div>

                    <div className="space-y-2">
                        <h4 className="font-medium leading-none">Date Range</h4>
                        <div className="grid gap-2">
                            <Popover>
                                <PopoverTrigger asChild>
                                    <Button
                                        variant={"outline"}
                                        className={cn(
                                            "w-full justify-start text-left font-normal",
                                            !startDate && "text-muted-foreground"
                                        )}
                                    >
                                        <CalendarIcon className="mr-2 h-4 w-4" />
                                        {startDate ? format(startDate, "PPP") : <span>Start Date</span>}
                                    </Button>
                                </PopoverTrigger>
                                <PopoverContent className="w-auto p-0" align="start">
                                    <Calendar
                                        mode="single"
                                        selected={startDate}
                                        onSelect={setStartDate}
                                        initialFocus
                                    />
                                </PopoverContent>
                            </Popover>
                            <Popover>
                                <PopoverTrigger asChild>
                                    <Button
                                        variant={"outline"}
                                        className={cn(
                                            "w-full justify-start text-left font-normal",
                                            !endDate && "text-muted-foreground"
                                        )}
                                    >
                                        <CalendarIcon className="mr-2 h-4 w-4" />
                                        {endDate ? format(endDate, "PPP") : <span>End Date</span>}
                                    </Button>
                                </PopoverTrigger>
                                <PopoverContent className="w-auto p-0" align="start">
                                    <Calendar
                                        mode="single"
                                        selected={endDate}
                                        onSelect={setEndDate}
                                        initialFocus
                                    />
                                </PopoverContent>
                            </Popover>
                        </div>
                    </div>

                    <Button variant="ghost" className="w-full" onClick={clearFilters}>
                        Clear Filters
                    </Button>
                </div>
            </PopoverContent>
        </Popover>
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
