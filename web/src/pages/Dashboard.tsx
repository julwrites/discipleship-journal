import { useEffect, useState, useRef } from "react";
import { auth } from "@/lib/firebase";
import { useDebounce } from "@/hooks/useDebounce";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { fetchNotes, syncUser, NoteFilter, deleteNote, getGroups, shareNote, getNote, askAI } from "@/services/api";
import { Link } from "react-router-dom";
import { Settings, Users, BookOpen, Filter, CalendarIcon, User as UserIcon, Book, LogOut } from "lucide-react";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { format } from "date-fns";
import { cn } from "@/lib/utils";
import { NoteCard, Note } from "@/components/NoteCard";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
    DialogFooter,
} from "@/components/ui/dialog";
import { Textarea } from "@/components/ui/textarea";
import { toast } from "sonner";

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

  // --- Actions State ---
  const [selectedNote, setSelectedNote] = useState<Note | null>(null);

  // Delete State
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [deleting, setDeleting] = useState(false);

  // Share State
  const [shareDialogOpen, setShareDialogOpen] = useState(false);
  const [sharing, setSharing] = useState(false);
  const [myGroups, setMyGroups] = useState<{ id: string, name: string }[]>([]);
  const [selectedGroupId, setSelectedGroupId] = useState("");
  const [shareComment, setShareComment] = useState("");

  // AI State
  const [aiDialogOpen, setAiDialogOpen] = useState(false);
  const [askingAI, setAskingAI] = useState(false);
  const [aiPrompt, setAiPrompt] = useState("");
  const [aiResponse, setAiResponse] = useState("");
  const [aiNoteContent, setAiNoteContent] = useState("");
  const [loadingNoteContent, setLoadingNoteContent] = useState(false);

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
                  const newNotes = response.data || response;

                  if (page === 1) {
                      setNotes(newNotes);
                  } else {
                      setNotes(prev => [...prev, ...newNotes]);
                  }

                  if (response.meta) {
                      setHasMore(page < response.meta.total_pages);
                  } else {
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

  // --- Handlers ---

  const handleDeleteClick = (note: Note) => {
      setSelectedNote(note);
      setDeleteDialogOpen(true);
  };

  const confirmDelete = async () => {
      if (!selectedNote) return;
      setDeleting(true);
      try {
          await deleteNote(selectedNote.id);
          setNotes(prev => prev.filter(n => n.id !== selectedNote.id));
          toast.success("Note deleted");
          setDeleteDialogOpen(false);
      } catch (e) {
          console.error(e);
          toast.error("Failed to delete note");
      } finally {
          setDeleting(false);
      }
  };

  const handleShareClick = async (note: Note) => {
      setSelectedNote(note);
      setShareDialogOpen(true);
      setShareComment("");
      setSelectedGroupId("");

      // Fetch groups if not already fetched (or refetch to be safe)
      try {
          const groups = await getGroups();
          setMyGroups(groups || []);
      } catch (e) {
          console.error(e);
          toast.error("Failed to load groups");
      }
  };

  const confirmShare = async () => {
      if (!selectedNote || !selectedGroupId) return;
      setSharing(true);
      try {
          await shareNote(selectedGroupId, selectedNote.id, shareComment);
          toast.success("Note shared successfully");
          setShareDialogOpen(false);
      } catch (e) {
          console.error(e);
          toast.error("Failed to share note");
      } finally {
          setSharing(false);
      }
  };

  const handleAskAIClick = (note: Note) => {
      setSelectedNote(note);
      setAiDialogOpen(true);
      setAiPrompt("");
      setAiResponse("");
      setAiNoteContent("");
      setLoadingNoteContent(true);

      // Fetch full content
      getNote(note.id).then(fullNote => {
          let content = fullNote.content || "";
          // Handle legacy format if needed, though getNote usually returns normalized object if we adjusted it,
          // but based on NoteEditor it might return object with markdown
          if (typeof content === 'object' && content.markdown) {
              content = content.markdown;
          } else if (typeof content === 'object') {
              content = JSON.stringify(content); // Fallback
          }
          setAiNoteContent(content);
      }).catch(e => {
          console.error(e);
          toast.error("Failed to load note content for AI");
          setAiDialogOpen(false);
      }).finally(() => {
          setLoadingNoteContent(false);
      });
  };

  const confirmAskAI = async () => {
      if (!aiNoteContent || !aiPrompt) return;
      setAskingAI(true);
      try {
          const res = await askAI(aiNoteContent, aiPrompt);
          setAiResponse(res.response);
      } catch (e) {
          console.error(e);
          toast.error("Failed to get AI response");
      } finally {
          setAskingAI(false);
      }
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

          {/* Resources Dropdown */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" title="Resources">
                <BookOpen className="w-5 h-5" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Resources</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <Link to="/reading-plans">
                <DropdownMenuItem className="cursor-pointer">
                  <CalendarIcon className="mr-2 h-4 w-4" /> Reading Plans
                </DropdownMenuItem>
              </Link>
              <Link to="/memory-verses">
                <DropdownMenuItem className="cursor-pointer">
                  <Book className="mr-2 h-4 w-4" /> Memory Verses
                </DropdownMenuItem>
              </Link>
              <Link to="/templates">
                <DropdownMenuItem className="cursor-pointer">
                  <Filter className="mr-2 h-4 w-4" /> Study Templates
                </DropdownMenuItem>
              </Link>
            </DropdownMenuContent>
          </DropdownMenu>

          {/* User Dropdown */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" title="User Menu">
                <UserIcon className="w-5 h-5" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>My Account</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <Link to="/settings">
                <DropdownMenuItem className="cursor-pointer">
                  <Settings className="mr-2 h-4 w-4" /> Settings
                </DropdownMenuItem>
              </Link>
              <DropdownMenuItem className="cursor-pointer text-destructive focus:text-destructive" onClick={() => auth.signOut()}>
                <LogOut className="mr-2 h-4 w-4" /> Sign Out
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
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
        {notes.length === 0 && !loading && <p className="text-muted-foreground col-span-full">No notes found.</p>}
        {notes.map((note) => (
          <NoteCard
            key={note.id}
            note={note}
            onDelete={handleDeleteClick}
            onShare={handleShareClick}
            onAskAI={handleAskAIClick}
          />
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

      {/* Delete Dialog */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
          <DialogContent>
              <DialogHeader>
                  <DialogTitle>Delete Note</DialogTitle>
                  <DialogDescription>
                      Are you sure you want to delete "{selectedNote?.title}"? This action cannot be undone.
                  </DialogDescription>
              </DialogHeader>
              <DialogFooter>
                  <Button variant="outline" onClick={() => setDeleteDialogOpen(false)}>Cancel</Button>
                  <Button variant="destructive" onClick={confirmDelete} disabled={deleting}>
                      {deleting ? "Deleting..." : "Delete"}
                  </Button>
              </DialogFooter>
          </DialogContent>
      </Dialog>

      {/* Share Dialog */}
      <Dialog open={shareDialogOpen} onOpenChange={setShareDialogOpen}>
          <DialogContent>
              <DialogHeader>
                  <DialogTitle>Share to Group</DialogTitle>
                  <DialogDescription>Share "{selectedNote?.title}" with your study group.</DialogDescription>
              </DialogHeader>
              <div className="space-y-4">
                  <Select value={selectedGroupId} onValueChange={setSelectedGroupId}>
                      <SelectTrigger>
                          <SelectValue placeholder="Select a Group..." />
                      </SelectTrigger>
                      <SelectContent>
                          {myGroups.map(g => (
                              <SelectItem key={g.id} value={g.id}>{g.name}</SelectItem>
                          ))}
                      </SelectContent>
                  </Select>
                  <Input
                      placeholder="Add a comment (optional)..."
                      value={shareComment}
                      onChange={(e) => setShareComment(e.target.value)}
                  />
                  <Button onClick={confirmShare} disabled={sharing || !selectedGroupId} className="w-full">
                      {sharing ? "Sharing..." : "Share Note"}
                  </Button>
              </div>
          </DialogContent>
      </Dialog>

      {/* AI Dialog */}
      <Dialog open={aiDialogOpen} onOpenChange={setAiDialogOpen}>
          <DialogContent className="sm:max-w-[500px]">
              <DialogHeader>
                  <DialogTitle>Ask AI</DialogTitle>
                  <DialogDescription>
                      Ask a question about "{selectedNote?.title}".
                  </DialogDescription>
              </DialogHeader>
              <div className="space-y-4">
                  {loadingNoteContent ? (
                      <div className="text-center py-4 text-muted-foreground">Loading note content...</div>
                  ) : (
                      <>
                        <div className="flex flex-col gap-2">
                            <Textarea
                                placeholder="Ask a question..."
                                value={aiPrompt}
                                onChange={(e) => setAiPrompt(e.target.value)}
                            />
                            <Button onClick={confirmAskAI} disabled={askingAI || !aiNoteContent}>
                                {askingAI ? "Thinking..." : "Ask"}
                            </Button>
                        </div>
                        {aiResponse && (
                            <div className="p-4 bg-muted border rounded max-h-60 overflow-auto text-sm">
                                <p className="font-semibold mb-2">Answer:</p>
                                <div className="prose prose-sm dark:prose-invert max-w-none break-words">
                                    <div dangerouslySetInnerHTML={{ __html: aiResponse }} />
                                </div>
                            </div>
                        )}
                      </>
                  )}
              </div>
          </DialogContent>
      </Dialog>

    </div>
  );
}
