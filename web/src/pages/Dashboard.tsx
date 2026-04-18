import { useEffect, useState, useRef } from "react";
import { auth } from "@/lib/firebase";
import { useAuth } from "@/hooks/useAuth";
import { useDebounce } from "@/hooks/useDebounce";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { fetchNotes, syncUser, NoteFilter, deleteNote, getGroups, shareNote, getNote, askAIStream, getConnections, getOrCreateDirectGroup, Connection, getTags, Tag } from "@/services/api";
import { getCachedNotes, setCachedNotes } from "@/services/cache";
import { Link } from "react-router-dom";
import { Settings, Users, BookOpen, Filter, CalendarIcon, User as UserIcon, Book, LogOut, Tag as TagIcon } from "lucide-react";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
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
import { Skeleton } from "@/components/ui/skeleton";

export default function Dashboard() {
  const { user } = useAuth();
  const [notes, setNotes] = useState<Note[]>([]);
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search, 500);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(true);
  const isInitialLoad = useRef(true);
  const prevSearchRef = useRef(debouncedSearch);

  // Filter state
  const [startDate, setStartDate] = useState<Date | undefined>();
  const [endDate, setEndDate] = useState<Date | undefined>();
  const [sortBy, setSortBy] = useState<"updated_at" | "created_at" | "title">("updated_at");
  const [sortOrder, setSortOrder] = useState<"asc" | "desc">("desc");
  const [tagFilter, setTagFilter] = useState<string>("");
  const [availableTags, setAvailableTags] = useState<string[]>([]);

  // Ref to track if filters changed to reset page
  const prevFilterRef = useRef({ startDate, endDate, sortBy, sortOrder, tag: tagFilter });

  useEffect(() => {
    getTags().then(tags => setAvailableTags(tags.map((t: Tag) => t.name))).catch(console.error);
  }, []);

  // --- Actions State ---
  const [selectedNote, setSelectedNote] = useState<Note | null>(null);

  // Delete State
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [deleting, setDeleting] = useState(false);

  // Share State
  const [shareDialogOpen, setShareDialogOpen] = useState(false);
  const [sharing, setSharing] = useState(false);
  const [myGroups, setMyGroups] = useState<{ id: string; name: string; type?: string }[]>([]);
  const [connections, setConnections] = useState<Connection[]>([]);
  const [shareType, setShareType] = useState<"group" | "person">("group");
  const [selectedGroupId, setSelectedGroupId] = useState("");
  const [selectedPersonId, setSelectedPersonId] = useState("");
  const [shareComment, setShareComment] = useState("");

  // AI State
  const [aiDialogOpen, setAiDialogOpen] = useState(false);
  const [askingAI, setAskingAI] = useState(false);
  const [aiPrompt, setAiPrompt] = useState("");
  const [aiResponse, setAiResponse] = useState("");
  const [aiNoteContent, setAiNoteContent] = useState("");
  const [loadingNoteContent, setLoadingNoteContent] = useState(false);
  const [currentUserId, setCurrentUserId] = useState<string>("");

  useEffect(() => {
    syncUser().then(u => {
        if (u.id) setCurrentUserId(u.id);
    });
  }, []);

  useEffect(() => {
      let ignore = false;

      const currentFilters = { startDate, endDate, sortBy, sortOrder, tag: tagFilter };
      const filtersChanged =
          prevFilterRef.current.startDate !== startDate ||
          prevFilterRef.current.endDate !== endDate ||
          prevFilterRef.current.sortBy !== sortBy ||
          prevFilterRef.current.sortOrder !== sortOrder ||
          prevFilterRef.current.tag !== tagFilter;

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
          const isDefaultView = page === 1 && !debouncedSearch && !startDate && !endDate && sortBy === "updated_at" && (!tagFilter || tagFilter === "_all");

          if (isInitialLoad.current && isDefaultView && user?.uid) {
              const cached = getCachedNotes(user.uid);
              if (cached) {
                  setNotes(cached.notes);
                  setHasMore(cached.hasMore);
                  setLoading(false); // Show cached content immediately
              } else {
                  setLoading(true);
              }
          } else {
              setLoading(true);
          }

          try {
              const filter: NoteFilter = {
                  search: debouncedSearch,
                  startDate,
                  endDate,
                  sortBy,
                  sortOrder,
                  tag: tagFilter === "_all" ? undefined : tagFilter
              };
              const response = await fetchNotes(page, 20, filter);
              if (!ignore) {
                  const newNotes = response.data || response;
                  let newHasMore = false;

                  if (response.meta) {
                      newHasMore = page < response.meta.total_pages;
                  } else {
                      newHasMore = newNotes.length >= 20;
                  }

                  if (page === 1) {
                      setNotes(newNotes);
                      if (isDefaultView && user?.uid) {
                          setCachedNotes(user.uid, newNotes, newHasMore);
                      }
                  } else {
                      setNotes(prev => [...prev, ...newNotes]);
                  }

                  setHasMore(newHasMore);
              }
          } catch (error) {
              if (!ignore) console.error(error);
          } finally {
              if (!ignore) {
                  setLoading(false);
                  isInitialLoad.current = false;
              }
          }
      };
      load();
      return () => { ignore = true; };
  }, [page, debouncedSearch, startDate, endDate, sortBy, sortOrder, tagFilter, user?.uid]);

  const handleSearch = (val: string) => {
      setSearch(val);
  };

  const clearFilters = () => {
      setStartDate(undefined);
      setEndDate(undefined);
      setSortBy("updated_at");
      setSortOrder("desc");
      setTagFilter("");
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
      setSelectedPersonId("");
      setShareType("group");

      // Fetch groups and connections
      try {
          const [groups, conns] = await Promise.all([
              getGroups(),
              getConnections()
          ]);
          setMyGroups(groups || []);
          setConnections(conns || []);
      } catch (e) {
          console.error(e);
          toast.error("Failed to load sharing options");
      }
  };

  const confirmShare = async () => {
      if (!selectedNote) return;

      let targetGroupId = selectedGroupId;

      setSharing(true);
      try {
          if (shareType === "person") {
              if (!selectedPersonId) return;
              // Get or create direct group
              const group = await getOrCreateDirectGroup(selectedPersonId);
              targetGroupId = group.id;
          }

          if (!targetGroupId) return;

          await shareNote(targetGroupId, selectedNote.id, shareComment);
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
      setAiResponse("");
      try {
          await askAIStream(aiNoteContent, aiPrompt, undefined, {
              onStart: () => setAiResponse(""),
              onChunk: (chunk) => setAiResponse(prev => prev + chunk), // Currently receives full text, but supports appending
              onDone: () => setAskingAI(false),
              onError: (e) => {
                  console.error(e);
                  toast.error("Failed to get AI response");
                  setAskingAI(false);
              }
          });
      } catch (e) {
          console.error(e);
          toast.error("Failed to start AI request");
          setAskingAI(false);
      }
  };

  return (
    <div className="p-4 md:p-8 max-w-7xl mx-auto">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-2xl font-bold">My Journal</h1>
        <div className="flex gap-2">
          <Link to="/connections">
            <Button variant="ghost" size="icon" title="Connections" aria-label="Connections">
              <Users className="w-5 h-5" />
            </Button>
          </Link>

          {/* Resources Dropdown */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" title="Resources" aria-label="Resources">
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
                  <Filter className="mr-2 h-4 w-4" /> Templates
                </DropdownMenuItem>
              </Link>
              <Link to="/tags">
                <DropdownMenuItem className="cursor-pointer">
                  <TagIcon className="mr-2 h-4 w-4" /> Tags
                </DropdownMenuItem>
              </Link>
            </DropdownMenuContent>
          </DropdownMenu>

          {/* User Dropdown */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" title="User Menu" aria-label="User Menu">
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
                <Button variant="outline" size="icon" title="Filter & Sort" aria-label="Filter and sort options">
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
                        <h4 className="font-medium leading-none">Tag</h4>
                        <Select value={tagFilter} onValueChange={setTagFilter}>
                            <SelectTrigger>
                                <SelectValue placeholder="All Tags" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="_all">All Tags</SelectItem>
                                {availableTags.map(t => (
                                    <SelectItem key={t} value={t}>{t}</SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
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
        {loading && notes.length === 0 && Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="flex flex-col space-y-3 p-6 border rounded-xl shadow-sm bg-card text-card-foreground">
                <div className="space-y-2">
                    <Skeleton className="h-5 w-3/4" />
                    <Skeleton className="h-3 w-1/2" />
                </div>
                <div className="space-y-2 pt-4">
                    <Skeleton className="h-4 w-full" />
                    <Skeleton className="h-4 w-full" />
                    <Skeleton className="h-4 w-2/3" />
                </div>
                <div className="flex gap-2 pt-4 mt-auto">
                    <Skeleton className="h-6 w-16 rounded-full" />
                    <Skeleton className="h-6 w-16 rounded-full" />
                </div>
            </div>
        ))}
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
            <Button className="w-14 h-14 rounded-full shadow-lg" size="icon" variant="secondary" title="Ask AI" aria-label="Ask AI">
                🤖
            </Button>
        </Link>
        <Link to="/notes/new">
            <Button className="w-14 h-14 rounded-full shadow-lg text-2xl" size="icon" title="New Note" aria-label="New Note">
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
                  <DialogTitle>Share Note</DialogTitle>
                  <DialogDescription>Share "{selectedNote?.title}" with others.</DialogDescription>
              </DialogHeader>

              <Tabs value={shareType} onValueChange={(v) => setShareType(v as "group" | "person")} className="w-full">
                  <TabsList className="grid w-full grid-cols-2 mb-4">
                      <TabsTrigger value="group">Share to Group</TabsTrigger>
                      <TabsTrigger value="person">Share to Person</TabsTrigger>
                  </TabsList>

                  <TabsContent value="group" className="space-y-4">
                      <Select value={selectedGroupId} onValueChange={setSelectedGroupId}>
                          <SelectTrigger>
                              <SelectValue placeholder="Select a Group..." />
                          </SelectTrigger>
                          <SelectContent>
                              {myGroups.filter(g => g.type !== 'direct').length === 0 && (
                                  <div className="p-2 text-sm text-muted-foreground text-center">No groups found</div>
                              )}
                              {myGroups.filter(g => g.type !== 'direct').map(g => (
                                  <SelectItem key={g.id} value={g.id}>{g.name}</SelectItem>
                              ))}
                          </SelectContent>
                      </Select>
                  </TabsContent>

                  <TabsContent value="person" className="space-y-4">
                      <Select value={selectedPersonId} onValueChange={setSelectedPersonId}>
                          <SelectTrigger>
                              <SelectValue placeholder="Select a Person..." />
                          </SelectTrigger>
                          <SelectContent>
                              {connections.filter(c => c.status === 'accepted').length === 0 && (
                                  <div className="p-2 text-sm text-muted-foreground text-center">No connections found</div>
                              )}
                              {connections.filter(c => c.status === 'accepted').map(c => {
                                  // Determine other user using Database ID (currentUserId) not Firebase UID (user.uid)
                                  // Fallback to email check if ID not yet loaded to prevent race conditions
                                  const isReq = currentUserId ? c.requester_id === currentUserId : (user?.email === c.requester_email);
                                  const otherId = isReq ? c.receiver_id : c.requester_id;
                                  const otherEmail = isReq ? c.receiver_email : c.requester_email;
                                  const otherName = isReq ? c.receiver_username : c.requester_username;

                                  return (
                                      <SelectItem key={c.id} value={otherId}>{otherName || otherEmail}</SelectItem>
                                  );
                              })}
                          </SelectContent>
                      </Select>
                  </TabsContent>
              </Tabs>

              <div className="space-y-4 mt-4">
                  <Input
                      placeholder="Add a comment (optional)..."
                      value={shareComment}
                      onChange={(e) => setShareComment(e.target.value)}
                  />
                  <Button
                      onClick={confirmShare}
                      disabled={sharing || (shareType === "group" ? !selectedGroupId : !selectedPersonId)}
                      className="w-full"
                  >
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
