import { useEffect, useState, useCallback } from "react";
import { getVersePacks, createVersePack, VersePack, getPackDetails, deletePack, createVerseInPack, MemoryVerse, clonePack, syncUser } from "@/services/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger, DialogFooter } from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { toast } from "sonner";
import { Plus, BookOpen, Trash2, ArrowLeft, Copy } from "lucide-react";
import { useNavigate, useParams, Link } from "react-router-dom";
import { BibleVersionSelector } from "@/components/BibleVersionSelector";

export default function MemoryVersesPage() {
    return (
        <div className="p-4 md:p-8 max-w-4xl mx-auto space-y-6">
            <Link to="/">
                <Button variant="ghost" className="pl-0">
                    <ArrowLeft className="mr-2 h-4 w-4" />
                    Back to Dashboard
                </Button>
            </Link>

            <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
                <div>
                    <h1 className="text-3xl font-bold">Scripture Memory</h1>
                    <p className="text-muted-foreground">Memorize and meditate on God's Word</p>
                </div>
            </div>

            <Tabs defaultValue="my-packs" className="w-full">
                <TabsList>
                    <TabsTrigger value="my-packs">My Packs</TabsTrigger>
                    <TabsTrigger value="system-packs">System Packs</TabsTrigger>
                </TabsList>
                <TabsContent value="my-packs" className="mt-6">
                    <PackList type="user" />
                </TabsContent>
                <TabsContent value="system-packs" className="mt-6">
                    <PackList type="system" />
                </TabsContent>
            </Tabs>
        </div>
    );
}

function PackList({ type }: { type: "system" | "user" }) {
    const [packs, setPacks] = useState<VersePack[]>([]);
    const [loading, setLoading] = useState(false);
    const [isCreateOpen, setIsCreateOpen] = useState(false);
    const [newPackTitle, setNewPackTitle] = useState("");
    const navigate = useNavigate();

    const loadPacks = useCallback(async () => {
        setLoading(true);
        try {
            const res = await getVersePacks(type);
            setPacks(res.data || []);
        } catch {
            toast.error("Failed to load packs");
        } finally {
            setLoading(false);
        }
    }, [type]);

    useEffect(() => {
        loadPacks();
    }, [loadPacks]);

    const handleCreatePack = async () => {
        if (!newPackTitle.trim()) return;
        try {
            await createVersePack(newPackTitle);
            toast.success("Pack created");
            setIsCreateOpen(false);
            setNewPackTitle("");
            loadPacks();
        } catch {
            toast.error("Failed to create pack");
        }
    };

    return (
        <div className="space-y-6">
            {type === "user" && (
                <div className="flex justify-end">
                    <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
                        <DialogTrigger asChild>
                            <Button>
                                <Plus className="mr-2 h-4 w-4" />
                                Create Pack
                            </Button>
                        </DialogTrigger>
                        <DialogContent>
                            <DialogHeader>
                                <DialogTitle>Create New Pack</DialogTitle>
                            </DialogHeader>
                            <div className="py-4">
                                <Label>Pack Title</Label>
                                <Input
                                    value={newPackTitle}
                                    onChange={e => setNewPackTitle(e.target.value)}
                                    placeholder="e.g. My Favorites"
                                    className="mt-2"
                                />
                            </div>
                            <DialogFooter>
                                <Button onClick={handleCreatePack}>Create</Button>
                            </DialogFooter>
                        </DialogContent>
                    </Dialog>
                </div>
            )}

            {loading ? (
                <div>Loading...</div>
            ) : packs.length === 0 ? (
                <div className="text-center py-10 text-muted-foreground">
                    No packs found.
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {packs.map(pack => (
                        <Card
                            key={pack.id}
                            className="cursor-pointer hover:border-primary/50 transition-colors"
                            onClick={() => navigate(`/memory-verses/${pack.id}`)}
                        >
                            <CardHeader>
                                <CardTitle className="text-lg flex justify-between items-start">
                                    <span>{pack.title}</span>
                                    {pack.identifier && (
                                        <span className="text-xs bg-muted px-2 py-1 rounded text-muted-foreground">
                                            {pack.identifier}
                                        </span>
                                    )}
                                </CardTitle>
                                <CardDescription>
                                    {pack.verse_count} verses
                                </CardDescription>
                            </CardHeader>
                            {pack.description && (
                                <CardContent>
                                    <p className="text-sm text-muted-foreground line-clamp-2">
                                        {pack.description}
                                    </p>
                                </CardContent>
                            )}
                        </Card>
                    ))}
                </div>
            )}
        </div>
    );
}

export function VersePackDetail() {
    const { id } = useParams();
    const navigate = useNavigate();
    const [pack, setPack] = useState<VersePack | null>(null);
    const [verses, setVerses] = useState<MemoryVerse[]>([]);
    const [loading, setLoading] = useState(true);
    const [isAddOpen, setIsAddOpen] = useState(false);

    // Add Verse Form
    const [newVerse, setNewVerse] = useState<Partial<MemoryVerse>>({
        reference: "",
        version: "ESV",
        tags: []
    });
    const [tagInput, setTagInput] = useState("");

    // Clone Form
    const [isCloneOpen, setIsCloneOpen] = useState(false);
    const [cloneTitle, setCloneTitle] = useState("");

    const loadDetails = useCallback(async () => {
        if (!id) return;
        setLoading(true);
        try {
            const res = await getPackDetails(id);
            setPack(res.pack);
            setVerses(res.verses || []);
        } catch {
            toast.error("Failed to load pack details");
            navigate("/memory-verses");
        } finally {
            setLoading(false);
        }
    }, [id, navigate]);

    useEffect(() => {
        if (id) loadDetails();
    }, [id, loadDetails]);

    // Pre-load default version
    useEffect(() => {
        syncUser().then(u => {
            if (u.settings?.bible_version) {
                setNewVerse(prev => ({ ...prev, version: u.settings.bible_version }));
            }
        }).catch(console.error);
    }, []);

    const handleAddVerse = async () => {
        if (!id || !newVerse.reference) return;
        try {
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            await createVerseInPack(id, newVerse as any);
            toast.success("Verse added");
            setIsAddOpen(false);
            setNewVerse(prev => ({ ...prev, reference: "", tags: [] })); // Keep version
            loadDetails();
        } catch {
            toast.error("Failed to add verse");
        }
    };

    const handleDeletePack = async () => {
        if (!id || !confirm("Are you sure you want to delete this pack?")) return;
        try {
            await deletePack(id);
            toast.success("Pack deleted");
            navigate("/memory-verses");
        } catch {
            toast.error("Failed to delete pack");
        }
    };

    const handleClonePack = async () => {
        if (!id) return;
        try {
            await clonePack(id, cloneTitle || undefined);
            toast.success("Pack cloned to your library");
            setIsCloneOpen(false);
            navigate("/memory-verses"); // Or navigate to new pack?
        } catch {
            toast.error("Failed to clone pack");
        }
    };

    const addTag = () => {
        if (tagInput.trim()) {
            setNewVerse(prev => ({
                ...prev,
                tags: [...(prev.tags || []), tagInput.trim()]
            }));
            setTagInput("");
        }
    };

    const isMyPack = pack?.user_id !== undefined && pack?.user_id !== null;

    if (loading) return <div className="p-8 text-center">Loading...</div>;
    if (!pack) return null;

    return (
        <div className="p-4 md:p-8 max-w-4xl mx-auto space-y-6">
            <Button variant="ghost" onClick={() => navigate("/memory-verses")} className="pl-0">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Back to Packs
            </Button>

            <div className="flex flex-col md:flex-row justify-between gap-4">
                <div>
                    <div className="flex items-center gap-3">
                        <h1 className="text-3xl font-bold">{pack.title}</h1>
                        {pack.identifier && (
                            <span className="bg-primary/10 text-primary px-2 py-1 rounded text-sm font-medium">
                                {pack.identifier}
                            </span>
                        )}
                    </div>
                    {pack.description && (
                        <p className="text-muted-foreground mt-2">{pack.description}</p>
                    )}
                </div>
                <div className="flex gap-2">
                    {!isMyPack ? (
                        <Dialog open={isCloneOpen} onOpenChange={setIsCloneOpen}>
                            <DialogTrigger asChild>
                                <Button variant="secondary">
                                    <Copy className="mr-2 h-4 w-4" />
                                    Save to My Packs
                                </Button>
                            </DialogTrigger>
                            <DialogContent>
                                <DialogHeader>
                                    <DialogTitle>Save Pack</DialogTitle>
                                </DialogHeader>
                                <div className="py-4">
                                    <Label>New Title (Optional)</Label>
                                    <Input
                                        value={cloneTitle}
                                        onChange={e => setCloneTitle(e.target.value)}
                                        placeholder={pack.title}
                                        className="mt-2"
                                    />
                                </div>
                                <DialogFooter>
                                    <Button onClick={handleClonePack}>Save</Button>
                                </DialogFooter>
                            </DialogContent>
                        </Dialog>
                    ) : (
                        <>
                            <Dialog open={isAddOpen} onOpenChange={setIsAddOpen}>
                                <DialogTrigger asChild>
                                    <Button>
                                        <Plus className="mr-2 h-4 w-4" />
                                        Add Verse
                                    </Button>
                                </DialogTrigger>
                                <DialogContent>
                                    <DialogHeader>
                                        <DialogTitle>Add Verse</DialogTitle>
                                    </DialogHeader>
                                    <div className="space-y-4 py-4">
                                        <div className="grid gap-2">
                                            <Label>Reference</Label>
                                            <Input
                                                value={newVerse.reference}
                                                onChange={e => setNewVerse({...newVerse, reference: e.target.value})}
                                                placeholder="e.g. John 3:16"
                                            />
                                        </div>
                                        <div className="grid gap-2">
                                            <Label>Version</Label>
                                            <BibleVersionSelector
                                                value={newVerse.version || "ESV"}
                                                onChange={v => setNewVerse({...newVerse, version: v})}
                                            />
                                        </div>
                                        <div className="grid gap-2">
                                            <Label>Tags</Label>
                                            <div className="flex gap-2">
                                                <Input
                                                    value={tagInput}
                                                    onChange={e => setTagInput(e.target.value)}
                                                    onKeyDown={e => e.key === 'Enter' && addTag()}
                                                    placeholder="Add tag..."
                                                />
                                                <Button type="button" variant="secondary" onClick={addTag}>Add</Button>
                                            </div>
                                            <div className="flex flex-wrap gap-2 mt-2">
                                                {newVerse.tags?.map((tag, i) => (
                                                    <span key={i} className="bg-secondary text-secondary-foreground px-2 py-1 rounded text-xs flex items-center gap-1">
                                                        {tag}
                                                        <button
                                                            onClick={() => setNewVerse(prev => ({...prev, tags: prev.tags?.filter((_, idx) => idx !== i)}))}
                                                            className="hover:text-destructive"
                                                        >
                                                            ×
                                                        </button>
                                                    </span>
                                                ))}
                                            </div>
                                        </div>
                                    </div>
                                    <DialogFooter>
                                        <Button onClick={handleAddVerse}>Add Verse</Button>
                                    </DialogFooter>
                                </DialogContent>
                            </Dialog>

                            <Button variant="destructive" size="icon" onClick={handleDeletePack}>
                                <Trash2 className="h-4 w-4" />
                            </Button>
                        </>
                    )}
                </div>
            </div>

            <div className="grid gap-4">
                {verses.map(verse => (
                    <Card key={verse.id}>
                        <CardContent className="p-4 flex justify-between items-center">
                            <div className="flex items-center gap-3">
                                <BookOpen className="h-5 w-5 text-primary" />
                                <div>
                                    <div className="font-semibold text-lg">{verse.reference}</div>
                                    <div className="text-xs text-muted-foreground flex gap-2">
                                        <span className="bg-muted px-1.5 rounded">{verse.version}</span>
                                        {verse.tags?.map(t => (
                                            <span key={t}>#{t}</span>
                                        ))}
                                    </div>
                                </div>
                            </div>
                            {/* Actions for verse could go here (e.g. Delete if my pack) */}
                        </CardContent>
                    </Card>
                ))}
                {verses.length === 0 && (
                    <div className="text-center py-8 text-muted-foreground bg-muted/20 rounded-lg">
                        No verses in this pack yet.
                    </div>
                )}
            </div>
        </div>
    );
}
