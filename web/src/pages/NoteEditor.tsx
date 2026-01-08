import { useEffect, useState, useRef } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import {
    createNote,
    getNote,
    updateNote,
    deleteNote,
    getBiblePassage,
    askAI,
    getGroups,
    shareNote,
    searchMemoryVerses,
    MemoryVerse
} from "@/services/api";
import RichTextEditor from "@/components/RichTextEditor";
import { Editor } from "@tiptap/react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import {
    DropdownMenu,
    DropdownMenuTrigger,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator
} from "@/components/ui/dropdown-menu";
import { MoreVertical, Book, Sparkles, Share2, Trash2, Quote } from "lucide-react";
import { useDebounce } from "@/hooks/useDebounce";

export default function NoteEditor() {
    const { id } = useParams();
    const navigate = useNavigate();
    const [title, setTitle] = useState("");
    const [content, setContent] = useState("");
    const [mode, setMode] = useState<"edit" | "preview">("edit");
    const [saving, setSaving] = useState(false);
    const [lastSaved, setLastSaved] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);

    const editorRef = useRef<Editor | null>(null);

    // Bible Passage State
    const [passageRef, setPassageRef] = useState("");
    const [bibleText, setBibleText] = useState("");
    const [loadingPassage, setLoadingPassage] = useState(false);
    const [passageDialogOpen, setPassageDialogOpen] = useState(false);

    // Memory Verse State
    const [verseDialogOpen, setVerseDialogOpen] = useState(false);
    const [verseSearch, setVerseSearch] = useState("");
    const [verses, setVerses] = useState<MemoryVerse[]>([]);
    const [loadingVerses, setLoadingVerses] = useState(false);
    const debouncedVerseSearch = useDebounce(verseSearch, 300);

    // AI State
    const [aiPrompt, setAiPrompt] = useState("");
    const [aiResponse, setAiResponse] = useState("");
    const [askingAI, setAskingAI] = useState(false);
    const [aiDialogOpen, setAiDialogOpen] = useState(false);

    // Sharing State
    const [myGroups, setMyGroups] = useState<{ id: string, name: string }[]>([]);
    const [selectedGroupId, setSelectedGroupId] = useState("");
    const [shareComment, setShareComment] = useState("");
    const [sharing, setSharing] = useState(false);
    const [deleting, setDeleting] = useState(false);
    const [saveError, setSaveError] = useState(false);
    const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false);
    const [shareDialogOpen, setShareDialogOpen] = useState(false);

    useEffect(() => {
        if (id && id !== "new") {
            setLoading(true);
            getNote(id).then(note => {
                setTitle(note.title);
                // Handle legacy content format (object with markdown)
                let noteContent = note.content || "";
                if (typeof noteContent === 'object' && noteContent.markdown) {
                    noteContent = noteContent.markdown;
                }
                setContent(noteContent);
                setLastSaved("Loaded");
            }).catch(e => {
                console.error(e);
                toast.error("Failed to load note");
            }).finally(() => {
                setLoading(false);
            });
        } else {
            setLoading(false);
        }
    }, [id]);

    useEffect(() => {
        if (verseDialogOpen) {
            handleSearchVerses(debouncedVerseSearch);
        }
    }, [debouncedVerseSearch, verseDialogOpen]);

    const handleSave = async (manual = true) => {
        setSaving(true);
        setSaveError(false);
        try {
            if (id === "new") {
                if (!manual) return; // Don't auto-save new notes until title/content exists or manual save
                const res = await createNote(title, content);
                navigate(`/notes/${res.id}`, { replace: true });
                setLastSaved(new Date().toLocaleTimeString());
            } else if (id) {
                await updateNote(id, title, content);
                setLastSaved(new Date().toLocaleTimeString());
            }
        } catch (e) {
            console.error(e);
            setSaveError(true);
            if (manual) {
                const message = e instanceof Error ? e.message : "Failed to save";
                // If it's a TypeError (usually network/CORS), it often lacks details, but we can hint at it.
                if (message === "Failed to fetch" || (e instanceof TypeError && message.includes("fetch"))) {
                    toast.error("Network error: Cannot reach server. Please check your connection.");
                } else {
                    toast.error(`Failed to save: ${message}`);
                }
            }
        } finally {
            setSaving(false);
        }
    };

    // Auto-save effect
    useEffect(() => {
        if (!id || id === "new") return;

        const timer = setTimeout(() => {
            if (title || content) {
                handleSave(false);
            }
        }, 2000); // 2 second debounce

        return () => clearTimeout(timer);
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [title, content, id]);

    const handleDelete = () => {
        if (!id || id === "new") return;
        setDeleteConfirmOpen(true);
    };

    const confirmDelete = async () => {
        setDeleting(true);
        try {
            await deleteNote(id!);
            navigate("/", { replace: true });
        } catch (e) {
            console.error(e);
            toast.error("Failed to delete note.");
            setDeleting(false);
            setDeleteConfirmOpen(false);
        }
    };

    const fetchMyGroups = async () => {
        try {
            const data = await getGroups();
            setMyGroups(data || []);
        } catch (error) {
            console.error("Failed to fetch groups", error);
        }
    };

    const handleShare = async () => {
        if (!id || id === "new") {
            toast.error("Please save the note first.");
            return;
        }
        if (!selectedGroupId) return;
        setSharing(true);
        try {
            await shareNote(selectedGroupId, id, shareComment);
            toast.success("Note shared!");
            setSelectedGroupId("");
            setShareComment("");
            setShareDialogOpen(false);
        } catch (error) {
            console.error("Share failed", error);
            toast.error("Failed to share.");
        } finally {
            setSharing(false);
        }
    };

    const handleFetchPassage = async () => {
        if (!passageRef) return;
        setLoadingPassage(true);
        try {
            const res = await getBiblePassage(passageRef);
            // API returns "verse" (real) or "text" (mock). Support both.
            const text = res.verse || res.text || res.content || "Passage found but no text returned.";
            setBibleText(text);

            if (res.reference) {
                setPassageRef(res.reference);
            }
        } catch (e) {
            console.error(e);
            setBibleText("Error fetching passage.");
        } finally {
            setLoadingPassage(false);
        }
    };

    const handleAddPassage = () => {
        if (!editorRef.current) return;

        // Build HTML for the passage
        // Backend returns HTML, so we don't need to manually wrap paragraphs unless it's a legacy response
        const html = `<blockquote><p><strong>${passageRef}</strong></p>${bibleText}</blockquote><p></p>`;

        editorRef.current.chain().focus().insertContent(html).run();

        setPassageRef("");
        setBibleText("");
        setPassageDialogOpen(false);
    };

    const handleSearchVerses = async (q: string) => {
        setLoadingVerses(true);
        try {
            const res = await searchMemoryVerses(q);
            setVerses(res.data || []);
        } catch (error) {
            console.error(error);
        } finally {
            setLoadingVerses(false);
        }
    };

    const handleInsertVerse = (verse: MemoryVerse) => {
        if (!editorRef.current) return;
        const html = `<blockquote><p><strong>${verse.reference} (${verse.version})</strong></p><p>${verse.text}</p></blockquote><p></p>`;
        editorRef.current.chain().focus().insertContent(html).run();
        setVerseDialogOpen(false);
    };

    const handleAskAI = async () => {
        setAskingAI(true);
        setAiResponse("");
        try {
            const res = await askAI(content, aiPrompt);
            setAiResponse(res.response);
        } catch (e) {
            console.error(e);
            setAiResponse("Error asking AI.");
        } finally {
            setAskingAI(false);
        }
    };

    const handleAddAIResponse = () => {
        if (!aiResponse || !editorRef.current) return;

        // Backend returns HTML now (via system prompt to LLM), so we insert directly
        const html = `<blockquote><p><em>Question: ${aiPrompt}</em></p>${aiResponse}</blockquote><p></p>`;

        editorRef.current.chain().focus().insertContent(html).run();

        setAiPrompt("");
        setAiResponse("");
        setAiDialogOpen(false);
    };

    return (
        <div className="flex flex-col h-screen max-w-4xl mx-auto p-4">
            <div className="flex justify-between items-center mb-4">
                <Button variant="ghost" onClick={() => navigate("/")}>&larr; Back</Button>

                <div className="flex items-center gap-2">
                    {saveError && <span className="text-sm text-destructive hidden sm:inline">Error saving</span>}
                    {!saveError && lastSaved && <span className="text-sm text-muted-foreground hidden sm:inline">{saving ? "Saving..." : `Saved at ${lastSaved}`}</span>}

                    <Button variant={mode === "edit" ? "default" : "outline"} onClick={() => setMode("edit")} className="hidden md:inline-flex">Edit</Button>
                    <Button variant={mode === "preview" ? "default" : "outline"} onClick={() => setMode("preview")} className="hidden md:inline-flex">Preview</Button>
                    <Button onClick={() => handleSave(true)} disabled={saving}>{saving ? "Saving..." : "Save"}</Button>

                    {/* Desktop Toolbar Buttons */}
                    <div className="hidden md:flex items-center gap-2">
                         <Button variant="outline" onClick={() => setPassageDialogOpen(true)} title="Lookup Bible Passage" aria-label="Add Scripture">
                            <Book className="h-4 w-4" />
                         </Button>
                         <Button variant="outline" onClick={() => setVerseDialogOpen(true)} title="Insert Memory Verse" aria-label="Insert Memory Verse">
                            <Quote className="h-4 w-4" />
                         </Button>
                         <Button variant="outline" onClick={() => setAiDialogOpen(true)} title="Ask AI" aria-label="Ask AI">
                            <Sparkles className="h-4 w-4" />
                         </Button>

                        {id && id !== "new" && (
                            <>
                                <Button variant="outline" onClick={() => {
                                    setShareDialogOpen(true);
                                    fetchMyGroups();
                                }}>Share</Button>
                                <Button variant="destructive" onClick={handleDelete} disabled={deleting}>
                                    {deleting ? "..." : "Delete"}
                                </Button>
                            </>
                        )}
                    </div>

                    {/* Mobile Dropdown Menu */}
                    <div className="md:hidden">
                        <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                                <Button variant="ghost" size="icon">
                                    <MoreVertical className="h-4 w-4" />
                                </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                                <DropdownMenuItem onClick={() => setMode(mode === "edit" ? "preview" : "edit")}>
                                    {mode === "edit" ? "Preview Mode" : "Edit Mode"}
                                </DropdownMenuItem>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem onClick={() => setPassageDialogOpen(true)}>
                                    <Book className="mr-2 h-4 w-4" /> Add Scripture
                                </DropdownMenuItem>
                                <DropdownMenuItem onClick={() => setVerseDialogOpen(true)}>
                                    <Quote className="mr-2 h-4 w-4" /> Add Memory Verse
                                </DropdownMenuItem>
                                <DropdownMenuItem onClick={() => setAiDialogOpen(true)}>
                                    <Sparkles className="mr-2 h-4 w-4" /> Ask AI
                                </DropdownMenuItem>
                                {id && id !== "new" && (
                                    <>
                                        <DropdownMenuItem onClick={() => {
                                            setShareDialogOpen(true);
                                            fetchMyGroups();
                                        }}>
                                            <Share2 className="mr-2 h-4 w-4" /> Share
                                        </DropdownMenuItem>
                                        <DropdownMenuSeparator />
                                        <DropdownMenuItem className="text-destructive" onClick={handleDelete}>
                                            <Trash2 className="mr-2 h-4 w-4" /> Delete
                                        </DropdownMenuItem>
                                    </>
                                )}
                            </DropdownMenuContent>
                        </DropdownMenu>
                    </div>
                </div>
            </div>

            <Dialog open={deleteConfirmOpen} onOpenChange={setDeleteConfirmOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Delete Note</DialogTitle>
                        <DialogDescription>
                            Are you sure you want to delete this note? This action cannot be undone.
                        </DialogDescription>
                    </DialogHeader>
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setDeleteConfirmOpen(false)}>Cancel</Button>
                        <Button variant="destructive" onClick={confirmDelete} disabled={deleting}>
                            {deleting ? "Deleting..." : "Delete"}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

             <Dialog open={passageDialogOpen} onOpenChange={setPassageDialogOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Add Bible Passage</DialogTitle>
                    </DialogHeader>
                    <div className="space-y-4">
                        <div className="flex flex-col gap-2">
                            <Textarea
                                placeholder="e.g. John 3:16"
                                value={passageRef}
                                onChange={(e) => setPassageRef(e.target.value)}
                                className="min-h-[100px]"
                            />
                            <Button onClick={handleFetchPassage} disabled={loadingPassage} className="w-full">
                                {loadingPassage ? "..." : "Search"}
                            </Button>
                        </div>
                        {bibleText && (
                            <div className="p-2 bg-muted border rounded max-h-40 overflow-auto text-sm italic">
                                <div className="prose prose-sm dark:prose-invert max-w-none break-words">
                                    <div dangerouslySetInnerHTML={{ __html: bibleText }} />
                                </div>
                            </div>
                        )}
                        {bibleText && (
                            <Button onClick={handleAddPassage} className="w-full">Insert into Note</Button>
                        )}
                    </div>
                </DialogContent>
            </Dialog>

            <Dialog open={verseDialogOpen} onOpenChange={setVerseDialogOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Insert Memory Verse</DialogTitle>
                    </DialogHeader>
                    <div className="space-y-4">
                        <Input
                            placeholder="Search verses..."
                            value={verseSearch}
                            onChange={(e) => setVerseSearch(e.target.value)}
                        />
                        <div className="max-h-60 overflow-auto space-y-2">
                            {loadingVerses ? (
                                <div className="text-center text-sm text-muted-foreground">Loading...</div>
                            ) : verses.length === 0 ? (
                                <div className="text-center text-sm text-muted-foreground">No verses found</div>
                            ) : (
                                verses.map(v => (
                                    <div
                                        key={v.id}
                                        className="p-2 border rounded hover:bg-muted cursor-pointer"
                                        onClick={() => handleInsertVerse(v)}
                                    >
                                        <div className="font-semibold text-sm">{v.reference}</div>
                                        <div className="text-xs text-muted-foreground line-clamp-2">{v.text}</div>
                                    </div>
                                ))
                            )}
                        </div>
                    </div>
                </DialogContent>
            </Dialog>

            <Dialog open={aiDialogOpen} onOpenChange={setAiDialogOpen}>
                <DialogContent className="sm:max-w-[500px]">
                    <DialogHeader>
                        <DialogTitle>Ask AI about this note</DialogTitle>
                    </DialogHeader>
                    <div className="space-y-4">
                        <div className="flex flex-col gap-2">
                            <Textarea
                                placeholder="Ask a question..."
                                value={aiPrompt}
                                onChange={(e) => setAiPrompt(e.target.value)}
                            />
                            <Button onClick={handleAskAI} disabled={askingAI || !content}>
                                {askingAI ? "Thinking..." : "Ask"}
                            </Button>
                        </div>
                        {aiResponse && (
                            <>
                                <div className="p-4 bg-muted border rounded max-h-60 overflow-auto text-sm">
                                    <p className="font-semibold mb-2">Answer:</p>
                                    <div className="prose prose-sm dark:prose-invert max-w-none break-words">
                                        <div dangerouslySetInnerHTML={{ __html: aiResponse }} />
                                    </div>
                                </div>
                                <Button onClick={handleAddAIResponse} className="w-full">
                                    Insert into Note
                                </Button>
                            </>
                        )}
                    </div>
                </DialogContent>
            </Dialog>
             <Dialog open={shareDialogOpen} onOpenChange={(open) => {
                setShareDialogOpen(open);
                if(open && myGroups.length === 0) fetchMyGroups();
            }}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Share to Group</DialogTitle>
                    </DialogHeader>
                    <div className="space-y-4">
                        <select
                            className="w-full p-2 border rounded"
                            value={selectedGroupId}
                            onChange={(e) => setSelectedGroupId(e.target.value)}
                        >
                            <option value="">Select a Group...</option>
                            {myGroups.map(g => (
                                <option key={g.id} value={g.id}>{g.name}</option>
                            ))}
                        </select>
                        <Input
                            placeholder="Add a comment (optional)..."
                            value={shareComment}
                            onChange={(e) => setShareComment(e.target.value)}
                        />
                        <Button onClick={handleShare} disabled={sharing || !selectedGroupId} className="w-full">
                            {sharing ? "Sharing..." : "Share Note"}
                        </Button>
                    </div>
                </DialogContent>
            </Dialog>

            <input
                className="text-4xl font-bold w-full mb-4 p-2 border-b outline-none bg-transparent text-foreground placeholder:text-muted-foreground"
                placeholder="Title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
            />

            <div className="flex-1 overflow-hidden flex flex-col">
                {mode === "edit" ? (
                    !loading && (
                        <RichTextEditor
                            initialContent={content}
                            onChange={setContent}
                            editable={true}
                            onEditorReady={(editor) => { editorRef.current = editor; }}
                        />
                    )
                ) : (
                    <div className="flex-1 border rounded-lg overflow-auto p-4 prose dark:prose-invert max-w-none bg-muted">
                        <div dangerouslySetInnerHTML={{ __html: content }} />
                    </div>
                )}
            </div>
        </div>
    );
}
