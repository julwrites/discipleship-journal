import { useEffect, useState, useRef, useCallback } from "react";
import { useParams, useNavigate, useBlocker, useLocation } from "react-router-dom";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import {
    createNote,
    getNote,
    updateNote,
    deleteNote,
    getBiblePassage,
    askAIStream,
    getGroups,
    shareItem,
    searchMemoryVerses,
    MemoryVerse,
    syncUser,
    getConnections,
    getOrCreateDirectGroup,
    getTags,
    Tag
} from "@/services/api";
import RichTextEditor from "@/components/RichTextEditor";
import { TagInput } from "@/components/TagInput";
import { Editor } from "@tiptap/react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog"
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
import { BibleVersionSelector } from "@/components/BibleVersionSelector";
import { BibleReferenceInput } from "@/components/BibleReferenceInput";

export default function NoteEditor() {
    const { id } = useParams();
    const navigate = useNavigate();
    const location = useLocation();
    const [title, setTitle] = useState("");
    const [content, setContent] = useState("");
    const [tags, setTags] = useState<string[]>([]);
    const [suggestions, setSuggestions] = useState<string[]>([]);
    const [initialTitle, setInitialTitle] = useState("");
    const [initialContent, setInitialContent] = useState("");
    const [initialTags, setInitialTags] = useState<string[]>([]);
    const [mode, setMode] = useState<"edit" | "preview">("edit");

    // Derived state for dirty check
    const isDirty = (title !== initialTitle) || (content !== initialContent) || (JSON.stringify(tags) !== JSON.stringify(initialTags));

    const [saving, setSaving] = useState(false);
    const [lastSaved, setLastSaved] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);

    const editorRef = useRef<Editor | null>(null);

    // Default Version from User Settings
    const [userVersion, setUserVersion] = useState("ESV");
    const [currentUserEmail, setCurrentUserEmail] = useState("");

    // Bible Passage State
    const [passageRef, setPassageRef] = useState("");
    const [bibleText, setBibleText] = useState("");
    const [passageVersion, setPassageVersion] = useState("ESV");
    const [loadingPassage, setLoadingPassage] = useState(false);
    const [passageDialogOpen, setPassageDialogOpen] = useState(false);

    // Memory Verse State
    const [verseDialogOpen, setVerseDialogOpen] = useState(false);
    const [verseSearch, setVerseSearch] = useState("");
    const [verses, setVerses] = useState<MemoryVerse[]>([]);
    const [loadingVerses, setLoadingVerses] = useState(false);
    const [insertingVerse, setInsertingVerse] = useState(false);
    const debouncedVerseSearch = useDebounce(verseSearch, 300);

    // AI State
    const [aiPrompt, setAiPrompt] = useState("");
    const [aiResponse, setAiResponse] = useState("");
    const [aiVersion, setAiVersion] = useState("ESV");
    const [askingAI, setAskingAI] = useState(false);
    const [aiDialogOpen, setAiDialogOpen] = useState(false);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const fetchingPromiseRef = useRef<Promise<any> | null>(null);

    // Sharing State
    const [myGroups, setMyGroups] = useState<{ id: string, name: string }[]>([]);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const [myConnections, setMyConnections] = useState<any[]>([]);
    const [selectedGroupId, setSelectedGroupId] = useState("");
    const [selectedConnectionId, setSelectedConnectionId] = useState("");
    const [shareComment, setShareComment] = useState("");
    const [sharing, setSharing] = useState(false);
    const [deleting, setDeleting] = useState(false);
    const [saveError, setSaveError] = useState(false);
    const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false);
    const [shareDialogOpen, setShareDialogOpen] = useState(false);

    useEffect(() => {
        // Load user settings for default version
        syncUser().then(u => {
            if (u.email) setCurrentUserEmail(u.email);
            if (u.settings?.bible_version) {
                setUserVersion(u.settings.bible_version);
                setPassageVersion(u.settings.bible_version);
                setAiVersion(u.settings.bible_version);
            }
        }).catch(console.error);

        // Load tags for autocomplete
        getTags().then(tags => setSuggestions(tags.map((t: Tag) => t.name))).catch(console.error);

        if (id && id !== "new") {
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const state = location.state as any;
            // Optimization: If we just created the note (auto-save), don't set loading to true
            // to avoid unmounting the editor and losing cursor focus/scroll.
            if (!state?.fromCreate) {
                setLoading(true);
            }

            getNote(id).then(note => {
                setTitle(note.title);
                setInitialTitle(note.title);

                // Handle legacy content format (object with markdown)
                let noteContent = note.content || "";
                if (typeof noteContent === 'object') {
                    // eslint-disable-next-line @typescript-eslint/no-explicit-any
                    const anyContent = noteContent as any;
                    if (anyContent.markdown) {
                         noteContent = anyContent.markdown;
                    }
                }
                setContent(noteContent as string);
                setInitialContent(noteContent as string);

                // Set tags
                const noteTags = note.tags ? note.tags.map((t: Tag) => t.name) : [];
                setTags(noteTags);
                setInitialTags(noteTags);

                setLastSaved("Loaded");
            }).catch(e => {
                console.error(e);
                toast.error("Failed to load note");
            }).finally(() => {
                setLoading(false);
            });
        } else {
            // Handle pre-population from location state
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const state = location.state as any;
            let promise = Promise.resolve();

            if (state) {
                if (state.title) setTitle(state.title);
                if (state.tags) setTags(state.tags);
                if (state.content) setContent(state.content);

                // Fetch passage if provided
                if (state.passageRef) {
                    if (!fetchingPromiseRef.current) {
                        const version = state.passageVersion || "ESV";
                        const refs = state.passageRef.split(';').map((r: string) => r.trim()).filter((r: string) => r.length > 0);

                        fetchingPromiseRef.current = Promise.all(refs.map((ref: string) =>
                            getBiblePassage(ref, version).then(res => ({ ref, res })).catch(err => ({ ref, err }))
                        )).then((results) => {
                            let combinedHtml = "";
                            // eslint-disable-next-line @typescript-eslint/no-explicit-any
                            results.forEach(({ ref, res, err }: any) => {
                                if (err) {
                                    console.error(`Failed to fetch ${ref}`, err);
                                    combinedHtml += `<blockquote><p><strong>${ref} (${version})</strong></p>Failed to load text.</blockquote><p></p>`;
                                } else {
                                    const text = res.verse || res.text || res.content || "";
                                    combinedHtml += `<blockquote><p><strong>${ref} (${version})</strong></p>${text}</blockquote><p></p>`;
                                }
                            });
                            setContent(prev => prev + combinedHtml);
                        });
                    }
                    promise = fetchingPromiseRef.current;
                }
            }

            promise.finally(() => setLoading(false));
        }
    }, [id, location.state]);

    useEffect(() => {
        if (verseDialogOpen) {
            handleSearchVerses(debouncedVerseSearch);
        }
    }, [debouncedVerseSearch, verseDialogOpen]);

    const saveNote = useCallback(async (shouldNavigate = true) => {
        if (!title.trim()) {
            toast.error("Please provide a title for your note.");
            return false;
        }

        setSaving(true);
        setSaveError(false);
        try {
            if (id === "new") {
                const res = await createNote(title, content, tags);
                if (shouldNavigate) {
                    // Pass fromCreate state to prevent reloading/flashing
                    navigate(`/notes/${res.id}`, { replace: true, state: { ...location.state, fromCreate: true } });
                }
                setLastSaved(new Date().toLocaleTimeString());
                setInitialTitle(title);
                setInitialContent(content);
                setInitialTags(tags);
                return true;
            } else if (id) {
                await updateNote(id, title, content, tags);
                setLastSaved(new Date().toLocaleTimeString());
                setInitialTitle(title);
                setInitialContent(content);
                setInitialTags(tags);
                return true;
            }
        } catch (e) {
            console.error(e);
            setSaveError(true);
            const message = e instanceof Error ? e.message : "Failed to save";
            if (message === "Failed to fetch" || (e instanceof TypeError && message.includes("fetch"))) {
                toast.error("Network error: Cannot reach server. Please check your connection.");
            } else {
                toast.error(`Failed to save: ${message}`);
            }
            return false;
        } finally {
            setSaving(false);
        }
        return false;
    }, [id, title, content, tags, navigate]);

    const handleSave = useCallback(() => {
        saveNote(true);
    }, [saveNote]);

    // Auto-save effect
    useEffect(() => {
        const timeoutId = setTimeout(() => {
            // Auto-save if dirty, not currently saving, and we have a title
            if (isDirty && !saving && title.trim()) {
                saveNote(true);
            }
        }, 2000);

        return () => clearTimeout(timeoutId);
    }, [isDirty, saving, title, saveNote]);

    // Handle Ctrl+S / Cmd+S
    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if ((e.ctrlKey || e.metaKey) && e.key === 's') {
                e.preventDefault();
                handleSave();
            }
        };

        window.addEventListener("keydown", handleKeyDown);
        return () => window.removeEventListener("keydown", handleKeyDown);
    }, [handleSave]);

    // Use useBlocker to warn about unsaved changes
    const blocker = useBlocker(
        ({ currentLocation, nextLocation }) =>
            !saving && isDirty && currentLocation.pathname !== nextLocation.pathname
    );

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

    const fetchConnections = async () => {
        try {
            const data = await getConnections();
            setMyConnections(data || []);
        } catch (error) {
            console.error("Failed to fetch connections", error);
        }
    };

    const handleShare = async () => {
        if (!id || id === "new") {
            toast.error("Please save the note first.");
            return;
        }

        let targetGroupId = selectedGroupId;

        setSharing(true);
        try {
            if (selectedConnectionId) {
                // Get or Create Direct Group
                const group = await getOrCreateDirectGroup(selectedConnectionId);
                targetGroupId = group.id;
            }

            if (!targetGroupId) return;

            await shareItem(targetGroupId, { note_id: id, comment: shareComment });
            toast.success("Note shared!");
            setSelectedGroupId("");
            setSelectedConnectionId("");
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
            const res = await getBiblePassage(passageRef, passageVersion);
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
        const html = `<blockquote><p><strong>${passageRef} (${passageVersion})</strong></p>${bibleText}</blockquote><p></p>`;

        editorRef.current.chain().insertContent(html).run();

        setPassageRef("");
        setBibleText("");
        setPassageDialogOpen(false);
    };

    const handleSearchVerses = async (q: string) => {
        if (!q || q.length === 0) {
            setVerses([]);
            return;
        }
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

    const handleInsertVerse = async (verse: MemoryVerse) => {
        if (!editorRef.current) return;
        setInsertingVerse(true);
        try {
             // For memory verses, we use the verse's version if available, or fetch using default
             const res = await getBiblePassage(verse.reference, verse.version || userVersion);
             // API returns "verse" (real) or "text" (mock). Support both.
             const text = res.verse || res.text || res.content || "Passage found but no text returned.";

             const html = `<blockquote><p><strong>${verse.reference} (${verse.version || userVersion})</strong></p>${text}</blockquote><p></p>`;
             editorRef.current.chain().insertContent(html).run();
             setVerseDialogOpen(false);
        } catch (error) {
            console.error("Failed to insert verse", error);
            toast.error("Failed to load verse text.");
        } finally {
            setInsertingVerse(false);
        }
    };

    const handleAskAI = async () => {
        setAskingAI(true);
        setAiResponse("");
        try {
            await askAIStream(content, aiPrompt, aiVersion, {
                onStart: () => setAiResponse(""),
                onChunk: (chunk) => setAiResponse(prev => prev + chunk),
                onDone: () => setAskingAI(false),
                onError: (e) => {
                    console.error(e);
                    toast.error("Error asking AI");
                    setAiResponse("Error asking AI.");
                    setAskingAI(false);
                }
            });
        } catch (e) {
            console.error(e);
            setAiResponse("Error asking AI.");
            setAskingAI(false);
        }
    };

    const handleAddAIResponse = () => {
        if (!aiResponse || !editorRef.current) return;

        const html = `<blockquote><p><em>Question: ${aiPrompt}</em></p>${aiResponse}</blockquote><p></p>`;

        editorRef.current.chain().insertContent(html).run();

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
                    {isDirty && !saving && <span className="text-sm text-yellow-600 hidden sm:inline">Unsaved changes</span>}

                    <Button variant={mode === "edit" ? "default" : "outline"} onClick={() => setMode("edit")} className="hidden md:inline-flex">Edit</Button>
                    <Button variant={mode === "preview" ? "default" : "outline"} onClick={() => setMode("preview")} className="hidden md:inline-flex">Preview</Button>
                    <Button onClick={handleSave} disabled={saving || !isDirty}>{saving ? "Saving..." : "Save"}</Button>

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
                                            fetchConnections();
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

             {/* Unsaved Changes Blocker Dialog */}
            {blocker.state === "blocked" && (
                <AlertDialog open={true}>
                    <AlertDialogContent>
                        <AlertDialogHeader>
                            <AlertDialogTitle>Unsaved Changes</AlertDialogTitle>
                            <AlertDialogDescription>
                                You have unsaved changes. Do you want to save them before leaving?
                            </AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                            <AlertDialogCancel onClick={() => blocker.reset()}>Cancel</AlertDialogCancel>
                            <Button variant="destructive" onClick={() => blocker.proceed()}>
                                Discard Changes
                            </Button>
                            <AlertDialogAction onClick={async (e) => {
                                e.preventDefault(); // Prevent closing immediately
                                const target = blocker.location;
                                const success = await saveNote(false); // Save without navigating
                                if (success && target) {
                                    navigate(target);
                                }
                            }}>
                                Save & Leave
                            </AlertDialogAction>
                        </AlertDialogFooter>
                    </AlertDialogContent>
                </AlertDialog>
            )}

            <Dialog open={deleteConfirmOpen} onOpenChange={setDeleteConfirmOpen}>
                <DialogContent aria-describedby={undefined}>
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
                        <DialogDescription>Search for a bible passage to add to your note.</DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4">
                        <div className="flex flex-col gap-2">
                            <BibleReferenceInput
                                multiline
                                allowMultiple
                                placeholder="e.g. John 3:16"
                                value={passageRef}
                                onChange={setPassageRef}
                                className="min-h-[100px]"
                            />
                            <div className="flex gap-2 items-center">
                                <BibleVersionSelector
                                    value={passageVersion}
                                    onChange={setPassageVersion}
                                    className="w-[150px]"
                                />
                                <Button onClick={handleFetchPassage} disabled={loadingPassage} className="flex-1">
                                    {loadingPassage ? "..." : "Search"}
                                </Button>
                            </div>
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
                        <DialogDescription>Search for a memory verse to insert.</DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4">
                        <Input
                            placeholder="Search by reference, pack or tags..."
                            value={verseSearch}
                            onChange={(e) => setVerseSearch(e.target.value)}
                        />
                        <div className="max-h-60 overflow-auto space-y-2">
                            {loadingVerses ? (
                                <div className="text-center text-sm text-muted-foreground">Loading...</div>
                            ) : verses.length === 0 ? (
                                <div className="text-center text-sm text-muted-foreground">
                                    {verseSearch.length > 0 ? "No verses found" : "Start typing to search..."}
                                </div>
                            ) : (
                                verses.map(v => (
                                    <div
                                        key={v.id}
                                        className={`p-2 border rounded hover:bg-muted cursor-pointer ${insertingVerse ? 'opacity-50 pointer-events-none' : ''}`}
                                        onClick={() => handleInsertVerse(v)}
                                    >
                                        <div className="flex justify-between items-start">
                                            <div className="font-semibold text-sm">
                                                {v.title ? `${v.title} (${v.reference})` : v.reference}
                                            </div>
                                            <span className="text-[10px] bg-secondary px-1.5 py-0.5 rounded text-secondary-foreground">
                                                {v.pack_title || "Unknown Pack"}
                                            </span>
                                        </div>
                                        {v.tags && v.tags.length > 0 && (
                                            <div className="flex gap-1 mt-1 flex-wrap">
                                                {v.tags.map(tag => (
                                                    <span key={tag} className="text-[10px] text-muted-foreground">#{tag}</span>
                                                ))}
                                            </div>
                                        )}
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
                        <DialogDescription>Ask the AI questions about your note content.</DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4">
                        <div className="flex flex-col gap-2">
                            <Textarea
                                placeholder="Ask a question..."
                                value={aiPrompt}
                                onChange={(e) => setAiPrompt(e.target.value)}
                            />
                            <div className="flex gap-2 items-center">
                                <BibleVersionSelector
                                    value={aiVersion}
                                    onChange={setAiVersion}
                                    className="w-[150px]"
                                />
                                <Button onClick={handleAskAI} disabled={askingAI || !content} className="flex-1">
                                    {askingAI ? "Thinking..." : "Ask"}
                                </Button>
                            </div>
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
                if(open) {
                    fetchMyGroups();
                    fetchConnections();
                }
            }}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Share Note</DialogTitle>
                        <DialogDescription>Share this note with your groups or connections.</DialogDescription>
                    </DialogHeader>

                    <Tabs defaultValue="groups" onValueChange={() => {
                        setSelectedGroupId("");
                        setSelectedConnectionId("");
                    }}>
                        <TabsList className="grid w-full grid-cols-2">
                            <TabsTrigger value="groups">Groups</TabsTrigger>
                            <TabsTrigger value="connections">Direct Message</TabsTrigger>
                        </TabsList>

                        <div className="py-4 space-y-4">
                            <TabsContent value="groups">
                                <select
                                    className="w-full p-2 border rounded bg-background"
                                    value={selectedGroupId}
                                    onChange={(e) => setSelectedGroupId(e.target.value)}
                                >
                                    <option value="">Select a Group...</option>
                                    {myGroups.map(g => (
                                        <option key={g.id} value={g.id}>{g.name}</option>
                                    ))}
                                </select>
                            </TabsContent>

                            <TabsContent value="connections">
                                <select
                                    className="w-full p-2 border rounded bg-background"
                                    value={selectedConnectionId}
                                    onChange={(e) => setSelectedConnectionId(e.target.value)}
                                >
                                    <option value="">Select a Connection...</option>
                                    {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
                                    {myConnections.filter((c: any) => c.status === 'accepted').map((c: any) => {
                                        // Need to identify which user is the 'other' one
                                        // But we don't have current user email easily here unless we use auth hook or syncUser data
                                        // Just display both emails or try to guess?
                                        // Connections list usually returns requester_email and receiver_email.
                                        // If I am requester, show receiver.
                                        // I will assume simple display for now: "requester <-> receiver" or just map all?
                                        // Better: The ConnectionsPage logic filters this.
                                        // I'll try to find the "other" email.
                                        // But I don't have 'user' object in this scope easily (useAuth hook is not used in top level... wait, it is NOT used in NoteEditor currently)
                                        // NoteEditor uses `syncUser`.
                                        // I'll just show the email that isn't null? Or maybe just render the object as string if I can't filter?
                                        // Actually `getConnections` returns { ... requester_email, receiver_email ... }.
                                        // I'll list both emails or just the ID.
                                        // Wait, I need to know which one is the OTHER.
                                        // I'll show: "Connection (ID: ...)" fallback?
                                        // No, that's bad UX.
                                        // I'll fetch user in useEffect or use `syncUser` result.
                                        // `syncUser` is called in useEffect. I can store user.
                                        // I'll add `currentUser` state.
                                        return (
                                            <option key={c.id} value={c.requester_email === currentUserEmail ? c.receiver_id : c.requester_id}>
                                                {c.requester_email === currentUserEmail ? c.receiver_email : c.requester_email}
                                            </option>
                                        );
                                    })}
                                </select>
                            </TabsContent>

                            <Input
                                placeholder="Add a comment (optional)..."
                                value={shareComment}
                                onChange={(e) => setShareComment(e.target.value)}
                            />

                            <Button onClick={handleShare} disabled={sharing || (!selectedGroupId && !selectedConnectionId)} className="w-full">
                                {sharing ? "Sharing..." : "Share Note"}
                            </Button>
                        </div>
                    </Tabs>
                </DialogContent>
            </Dialog>

            <input
                className="text-4xl font-bold w-full mb-4 p-2 border-b outline-none bg-transparent text-foreground placeholder:text-muted-foreground"
                placeholder="Title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
            />

            <div className="mb-4">
                <TagInput value={tags} onChange={setTags} suggestions={suggestions} />
            </div>

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