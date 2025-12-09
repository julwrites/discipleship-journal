import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { createNote, getBiblePassage, askAI } from "@/services/api";
import ReactMarkdown from "react-markdown";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"

// Since we are building iteratively, I will use createNote for everything first, then update for proper updates.
// But first, let's fix the API service to support Get/Update.

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api";
import { auth } from "@/lib/firebase";

async function getNote(id: string) {
    const token = await auth.currentUser?.getIdToken();
    const res = await fetch(`${API_URL}/notes/${id}`, {
        headers: { "Authorization": `Bearer ${token}` }
    });
    if (!res.ok) throw new Error("Failed to load note");
    return res.json();
}

async function updateNote(id: string, title: string, content: Record<string, unknown>) {
    const token = await auth.currentUser?.getIdToken();
    const res = await fetch(`${API_URL}/notes/${id}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({ title, content })
    });
    if (!res.ok) throw new Error("Failed to update note");
}

export default function NoteEditor() {
    const { id } = useParams();
    const navigate = useNavigate();
    const [title, setTitle] = useState("");
    const [markdown, setMarkdown] = useState("");
    const [mode, setMode] = useState<"edit" | "preview">("edit");
    const [saving, setSaving] = useState(false);

    // Bible Passage State
    const [passageRef, setPassageRef] = useState("");
    const [bibleText, setBibleText] = useState("");
    const [loadingPassage, setLoadingPassage] = useState(false);

    // AI State
    const [aiPrompt, setAiPrompt] = useState("");
    const [aiResponse, setAiResponse] = useState("");
    const [askingAI, setAskingAI] = useState(false);

    // Sharing State
    const [myGroups, setMyGroups] = useState<{ id: string, name: string }[]>([]);
    const [selectedGroupId, setSelectedGroupId] = useState("");
    const [shareComment, setShareComment] = useState("");
    const [sharing, setSharing] = useState(false);

    useEffect(() => {
        if (id && id !== "new") {
            getNote(id).then(note => {
                setTitle(note.title);
                setMarkdown(note.content.markdown || "");
            }).catch(console.error);
        }
    }, [id]);

    const handleSave = async () => {
        setSaving(true);
        try {
            const content = { markdown };
            if (id === "new") {
                const res = await createNote(title, content);
                navigate(`/notes/${res.id}`, { replace: true });
            } else if (id) {
                await updateNote(id, title, content);
            }
        } catch (e) {
            console.error(e);
            alert("Failed to save");
        } finally {
            setSaving(false);
        }
    };

    const fetchMyGroups = async () => {
        try {
            const token = await auth.currentUser?.getIdToken();
            const res = await fetch(`${API_URL}/groups`, {
                headers: { Authorization: `Bearer ${token}` }
            });
            if (res.ok) {
                const data = await res.json();
                setMyGroups(data || []);
            }
        } catch (error) {
            console.error("Failed to fetch groups", error);
        }
    };

    const handleShare = async () => {
        if (!id || id === "new") {
            alert("Please save the note first.");
            return;
        }
        if (!selectedGroupId) return;
        setSharing(true);
        try {
            const token = await auth.currentUser?.getIdToken();
            const res = await fetch(`${API_URL}/groups/${selectedGroupId}/shares`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`
                },
                body: JSON.stringify({ note_id: id, comment: shareComment })
            });
            if (res.ok) {
                alert("Note shared!");
                setSelectedGroupId("");
                setShareComment("");
            } else {
                alert("Failed to share.");
            }
        } catch (error) {
            console.error("Share failed", error);
        } finally {
            setSharing(false);
        }
    };

    const handleFetchPassage = async () => {
        setLoadingPassage(true);
        try {
            const res = await getBiblePassage(passageRef);
            // Assuming res.text or res.content. Adjust based on API response
            const text = res.text || res.content || "Passage found but no text returned.";
            setBibleText(text);
        } catch (e) {
            console.error(e);
            setBibleText("Error fetching passage.");
        } finally {
            setLoadingPassage(false);
        }
    };

    const handleAddPassage = () => {
        const newContent = `${markdown}\n\n> **${passageRef}**\n> ${bibleText}\n`;
        setMarkdown(newContent);
        setPassageRef("");
        setBibleText("");
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

    return (
        <div className="flex flex-col h-screen max-w-4xl mx-auto p-4">
            <div className="flex justify-between items-center mb-4">
                <Button variant="ghost" onClick={() => navigate("/")}>&larr; Back</Button>
                <div className="space-x-2 flex items-center">
                    <Dialog>
                        <DialogTrigger asChild>
                            <Button variant="outline">Add Scripture</Button>
                        </DialogTrigger>
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
                                    <div className="p-2 bg-slate-50 border rounded max-h-40 overflow-auto text-sm italic">
                                        {bibleText}
                                    </div>
                                )}
                                {bibleText && (
                                    <Button onClick={handleAddPassage} className="w-full">Insert into Note</Button>
                                )}
                            </div>
                        </DialogContent>
                    </Dialog>

                    <Dialog>
                        <DialogTrigger asChild>
                            <Button variant="outline">Ask AI</Button>
                        </DialogTrigger>
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
                                    <div className="p-4 bg-slate-50 border rounded max-h-60 overflow-auto text-sm">
                                        <p className="font-semibold mb-2">Answer:</p>
                                        <ReactMarkdown>{aiResponse}</ReactMarkdown>
                                    </div>
                                )}
                            </div>
                        </DialogContent>
                    </Dialog>

                    <Button variant={mode === "edit" ? "default" : "outline"} onClick={() => setMode("edit")}>Edit</Button>
                    <Button variant={mode === "preview" ? "default" : "outline"} onClick={() => setMode("preview")}>Preview</Button>
                    <Button onClick={handleSave} disabled={saving}>{saving ? "Saving..." : "Save"}</Button>

                    <Dialog onOpenChange={(open) => { if (open) fetchMyGroups(); }}>
                        <DialogTrigger asChild>
                            <Button variant="outline" disabled={!id || id === "new"}>Share</Button>
                        </DialogTrigger>
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
                </div>
            </div>

            <input
                className="text-4xl font-bold w-full mb-4 p-2 border-b outline-none"
                placeholder="Title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
            />

            <div className="flex-1 border rounded-lg overflow-hidden">
                {mode === "edit" ? (
                    <textarea
                        className="w-full h-full p-4 resize-none outline-none"
                        placeholder="Write your thoughts..."
                        value={markdown}
                        onChange={(e) => setMarkdown(e.target.value)}
                    />
                ) : (
                    <div className="p-4 prose prose-slate max-w-none overflow-auto h-full">
                        <ReactMarkdown>{markdown}</ReactMarkdown>
                    </div>
                )}
            </div>
        </div>
    );
}
