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
    shareNote
} from "@/services/api";
import ReactMarkdown from "react-markdown";
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
import { useDebounce } from "@/hooks/useDebounce";
import {
    DropdownMenu,
    DropdownMenuTrigger,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator
} from "@/components/ui/dropdown-menu";
import { MoreVertical, Book, Sparkles, Share2, Trash2 } from "lucide-react";

export default function NoteEditor() {
    const { id } = useParams();
    const navigate = useNavigate();
    const [title, setTitle] = useState("");
    const [markdown, setMarkdown] = useState("");
    const [mode, setMode] = useState<"edit" | "preview">("edit");
    const [saving, setSaving] = useState(false);
    const [lastSaved, setLastSaved] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);

    const editorRef = useRef<Editor | null>(null);

    // Bible Passage State
    const [passageRef, setPassageRef] = useState("");
    const debouncedPassageRef = useDebounce(passageRef, 500);
    const [bibleText, setBibleText] = useState("");
    const [loadingPassage, setLoadingPassage] = useState(false);
    const [passageDialogOpen, setPassageDialogOpen] = useState(false);

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
                setMarkdown(note.content.markdown || "");
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

    const handleSave = async (manual = true) => {
        setSaving(true);
        setSaveError(false);
        try {
            const content = { markdown };
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
            if (title || markdown) {
                handleSave(false);
            }
        }, 2000); // 2 second debounce

        return () => clearTimeout(timer);
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [title, markdown, id]);

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
        } catch (e) {
            console.error(e);
            setBibleText("Error fetching passage.");
        } finally {
            setLoadingPassage(false);
        }
    };

    useEffect(() => {
        if (debouncedPassageRef && debouncedPassageRef.length > 2) {
            handleFetchPassage();
        } else if (!debouncedPassageRef) {
            setBibleText("");
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [debouncedPassageRef]);

    const handleAddPassage = () => {
        if (!editorRef.current) return;

        // Build HTML for the passage
        const formattedText = bibleText.split('\n').map(line => `<p>${line}</p>`).join('');
        const html = `<blockquote><p><strong>${passageRef}</strong></p>${formattedText}</blockquote><p></p>`;

        editorRef.current.chain().focus().insertContent(html).run();

        setPassageRef("");
        setBibleText("");
        setPassageDialogOpen(false);
    };

    const handleAskAI = async () => {
        setAskingAI(true);
        setAiResponse("");
        try {
            const res = await askAI(markdown, aiPrompt);
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

        const formattedResponse = aiResponse.split('\n').map(line => `<p>${line}</p>`).join('');
        const html = `<blockquote><p><strong>AI Response</strong></p>${formattedResponse}<p><em>Question: ${aiPrompt}</em></p></blockquote><p></p>`;

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
                         <Button variant="outline" onClick={() => setPassageDialogOpen(true)}>Add Scripture</Button>
                         <Button variant="outline" onClick={() => setAiDialogOpen(true)}>Ask AI</Button>

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
                        <div className="flex gap-2">
                            <Input
                                placeholder="e.g. John 3:16"
                                value={passageRef}
                                onChange={(e) => setPassageRef(e.target.value)}
                            />
                            <Button onClick={handleFetchPassage} disabled={loadingPassage}>
                                {loadingPassage ? "..." : "Search"}
                            </Button>
                        </div>
                        {bibleText && (
                            <div className="p-2 bg-muted border rounded max-h-40 overflow-auto text-sm italic">
                                <div className="prose prose-sm dark:prose-invert max-w-none">
                                    <ReactMarkdown>
                                        {bibleText}
                                    </ReactMarkdown>
                                </div>
                            </div>
                        )}
                        {bibleText && (
                            <Button onClick={handleAddPassage} className="w-full">Insert into Note</Button>
                        )}
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
                            <Input
                                placeholder="Ask a question..."
                                value={aiPrompt}
                                onChange={(e) => setAiPrompt(e.target.value)}
                            />
                            <Button onClick={handleAskAI} disabled={askingAI || !markdown}>
                                {askingAI ? "Thinking..." : "Ask"}
                            </Button>
                        </div>
                        {aiResponse && (
                            <>
                                <div className="p-4 bg-muted border rounded max-h-60 overflow-auto text-sm">
                                    <p className="font-semibold mb-2">Answer:</p>
                                    <ReactMarkdown>{aiResponse}</ReactMarkdown>
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
                            initialContent={markdown}
                            onChange={setMarkdown}
                            editable={true}
                            onEditorReady={(editor) => { editorRef.current = editor; }}
                        />
                    )
                ) : (
                    <div className="flex-1 border rounded-lg overflow-auto p-4 prose dark:prose-invert max-w-none bg-muted">
                        <ReactMarkdown>{markdown}</ReactMarkdown>
                    </div>
                )}
            </div>
        </div>
    );
}
