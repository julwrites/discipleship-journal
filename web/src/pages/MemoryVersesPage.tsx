import { useEffect, useState, useCallback } from "react";
import { getVersePacks, createVersePack, VersePack, getPackDetails, deletePack, createVerseInPack, MemoryVerse, clonePack, syncUser, updateMemoryVerse, deleteMemoryVerse, getBiblePassage, setVersePreference, removeVersePreference, setVersePreferencesBatch } from "@/services/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogTrigger, DialogFooter } from "@/components/ui/dialog";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { toast } from "sonner";
import { Plus, BookOpen, Trash2, ArrowLeft, Copy, Pencil, Settings } from "lucide-react";
import { useNavigate, useParams, Link } from "react-router-dom";
import { BibleVersionSelector } from "@/components/BibleVersionSelector";
import { cn } from "@/lib/utils";
import { BibleReferenceInput } from "@/components/BibleReferenceInput";

export default function MemoryVersesPage() {
    const navigate = useNavigate();
    return (
        <div className="p-4 md:p-8 max-w-4xl mx-auto space-y-6">
            <div>
                <div className="flex items-center gap-3 mb-2">
                    <Button variant="ghost" size="icon" onClick={() => navigate("/")} aria-label="Go back">
                        <ArrowLeft className="h-5 w-5" />
                    </Button>
                    <h1 className="text-3xl font-bold">Scripture Memory</h1>
                </div>
                <p className="text-muted-foreground ml-12">Memorize and meditate on God's Word</p>
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
                                <DialogDescription>Enter a title for your new verse pack.</DialogDescription>
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
        title: "",
        version: "ESV",
        tags: []
    });
    const [tagInput, setTagInput] = useState("");

    // Edit Verse Form
    const [isEditOpen, setIsEditOpen] = useState(false);
    const [editVerseData, setEditVerseData] = useState<Partial<MemoryVerse>>({});
    const [editTagInput, setEditTagInput] = useState("");

    // View Verse State
    const [viewVerseOpen, setViewVerseOpen] = useState(false);
    const [viewVerseData, setViewVerseData] = useState<{ reference: string, text: string, version: string } | null>(null);
    const [loadingVerseText, setLoadingVerseText] = useState(false);

    // Clone Form
    const [isCloneOpen, setIsCloneOpen] = useState(false);
    const [cloneTitle, setCloneTitle] = useState("");
    const [cloneUseUserDefault, setCloneUseUserDefault] = useState(true);

    const [userDefaultVersion, setUserDefaultVersion] = useState<string | null>(null);

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
                setUserDefaultVersion(u.settings.bible_version);
            }
        }).catch(console.error);
    }, []);

    const handleApplyDefaultToAll = async () => {
        if (!id || !userDefaultVersion) return;
        if (!confirm(`This will set your preference for all verses in this pack to ${userDefaultVersion}. Continue?`)) return;

        try {
            const verseIds = verses
                .map(v => v.id)
                .filter((id): id is string => !!id);

            if (verseIds.length > 0) {
                await setVersePreferencesBatch(verseIds, userDefaultVersion);
            }
            toast.success(`Applied ${userDefaultVersion} to all verses`);
            loadDetails();
        } catch {
            toast.error("Failed to apply preferences");
        }
    };

    const handleAddVerse = async () => {
        if (!id || !newVerse.reference) return;
        try {
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            await createVerseInPack(id, newVerse as any);
            toast.success("Verse added");
            setIsAddOpen(false);
            setNewVerse(prev => ({ ...prev, reference: "", title: "", tags: [] })); // Keep version
            loadDetails();
        } catch {
            toast.error("Failed to add verse");
        }
    };

    const handleUpdateVerse = async () => {
        if (!editVerseData.id || !editVerseData.reference) return;
        try {
            if (isMyPack) {
                await updateMemoryVerse(editVerseData.id, editVerseData);
                toast.success("Verse updated");
            } else {
                if (editVerseData.version) {
                    await setVersePreference(editVerseData.id, editVerseData.version);
                    toast.success("Verse preference saved");
                }
            }
            setIsEditOpen(false);
            loadDetails();
        } catch {
            toast.error("Failed to update verse");
        }
    };

    const handleResetPreference = async () => {
        if (!editVerseData.id) return;
        try {
            await removeVersePreference(editVerseData.id);
            toast.success("Reset to default version");
            setIsEditOpen(false);
            loadDetails();
        } catch {
            toast.error("Failed to reset preference");
        }
    };

    const handleDeleteVerse = async (verseId: string) => {
        if (!confirm("Are you sure you want to delete this verse?")) return;
        try {
            await deleteMemoryVerse(verseId);
            toast.success("Verse deleted");
            loadDetails();
        } catch {
            toast.error("Failed to delete verse");
        }
    };

    const handleViewVerse = async (verse: MemoryVerse) => {
        setLoadingVerseText(true);
        setViewVerseOpen(true);
        setViewVerseData({ reference: verse.reference, version: verse.version, text: "Loading..." });
        try {
            const res = await getBiblePassage(verse.reference, verse.version);
            const text = res.verse || res.text || res.content || "Text not found.";
            setViewVerseData({ reference: verse.reference, version: verse.version, text });
        } catch {
            setViewVerseData(prev => prev ? { ...prev, text: "Failed to load verse text." } : null);
        } finally {
            setLoadingVerseText(false);
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
            await clonePack(id, cloneTitle || undefined, cloneUseUserDefault);
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

    const addEditTag = () => {
        if (editTagInput.trim()) {
            setEditVerseData(prev => ({
                ...prev,
                tags: [...(prev.tags || []), editTagInput.trim()]
            }));
            setEditTagInput("");
        }
    };

    const openEdit = (e: React.MouseEvent, v: MemoryVerse) => {
        e.stopPropagation(); // Prevent card click
        setEditVerseData({ ...v });
        setIsEditOpen(true);
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
                        <>
                            {userDefaultVersion && (
                                <Button variant="outline" onClick={handleApplyDefaultToAll}>
                                    Set all to {userDefaultVersion}
                                </Button>
                            )}
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
                                    <DialogDescription>Save a copy of this pack to your library.</DialogDescription>
                                </DialogHeader>
                                <div className="py-4 space-y-4">
                                    <div>
                                        <Label>New Title (Optional)</Label>
                                        <Input
                                            value={cloneTitle}
                                            onChange={e => setCloneTitle(e.target.value)}
                                            placeholder={pack.title}
                                            className="mt-2"
                                        />
                                    </div>
                                    {userDefaultVersion && (
                                        <div className="flex items-center space-x-2">
                                            <Checkbox
                                                id="use-default"
                                                checked={cloneUseUserDefault}
                                                onCheckedChange={(c) => setCloneUseUserDefault(!!c)}
                                            />
                                            <Label htmlFor="use-default" className="font-normal cursor-pointer">
                                                Apply my default version ({userDefaultVersion}) to all verses
                                            </Label>
                                        </div>
                                    )}
                                </div>
                                <DialogFooter>
                                    <Button onClick={handleClonePack}>Save</Button>
                                </DialogFooter>
                            </DialogContent>
                        </Dialog>
                        </>
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
                                        <DialogDescription>Add a new memory verse to this pack.</DialogDescription>
                                    </DialogHeader>
                                    <div className="space-y-4 py-4">
                                        <div className="grid gap-2">
                                            <Label>Title (Optional)</Label>
                                            <Input
                                                value={newVerse.title || ""}
                                                onChange={e => setNewVerse({...newVerse, title: e.target.value})}
                                                placeholder="e.g. God's Love"
                                            />
                                        </div>
                                        <div className="grid gap-2">
                                            <Label>Reference</Label>
                                            <BibleReferenceInput
                                                value={newVerse.reference || ""}
                                                onChange={val => setNewVerse({...newVerse, reference: val})}
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

            {/* Edit Verse Dialog */}
            <Dialog open={isEditOpen} onOpenChange={setIsEditOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>{isMyPack ? "Edit Verse" : "Verse Preferences"}</DialogTitle>
                        <DialogDescription>
                            {isMyPack ? "Make changes to your memory verse." : "Customize the Bible version for this verse."}
                        </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4 py-4">
                        {isMyPack && (
                            <>
                                <div className="grid gap-2">
                                    <Label>Title (Optional)</Label>
                                    <Input
                                        value={editVerseData.title || ""}
                                        onChange={e => setEditVerseData({...editVerseData, title: e.target.value})}
                                        placeholder="e.g. God's Love"
                                    />
                                </div>
                                <div className="grid gap-2">
                                    <Label>Reference</Label>
                                    <BibleReferenceInput
                                        value={editVerseData.reference || ""}
                                        onChange={val => setEditVerseData({...editVerseData, reference: val})}
                                        placeholder="e.g. John 3:16"
                                    />
                                </div>
                            </>
                        )}
                        <div className="grid gap-2">
                            <Label>Version</Label>
                            <BibleVersionSelector
                                value={editVerseData.version || "ESV"}
                                onChange={v => setEditVerseData({...editVerseData, version: v})}
                            />
                            {!isMyPack && editVerseData.version_source === 'override' && (
                                <p className="text-xs text-muted-foreground">
                                    You have customized this version.
                                </p>
                            )}
                            {!isMyPack && (
                                <Link to="/settings" className="text-xs text-primary hover:underline flex items-center gap-1 mt-1">
                                    <Settings className="h-3 w-3" /> Manage default version
                                </Link>
                            )}
                        </div>
                        {isMyPack && (
                            <div className="grid gap-2">
                                <Label>Tags</Label>
                                <div className="flex gap-2">
                                    <Input
                                        value={editTagInput}
                                        onChange={e => setEditTagInput(e.target.value)}
                                        onKeyDown={e => e.key === 'Enter' && addEditTag()}
                                        placeholder="Add tag..."
                                    />
                                    <Button type="button" variant="secondary" onClick={addEditTag}>Add</Button>
                                </div>
                                <div className="flex flex-wrap gap-2 mt-2">
                                    {editVerseData.tags?.map((tag, i) => (
                                        <span key={i} className="bg-secondary text-secondary-foreground px-2 py-1 rounded text-xs flex items-center gap-1">
                                            {tag}
                                            <button
                                                onClick={() => setEditVerseData(prev => ({...prev, tags: prev.tags?.filter((_, idx) => idx !== i)}))}
                                                className="hover:text-destructive"
                                            >
                                                ×
                                            </button>
                                        </span>
                                    ))}
                                </div>
                            </div>
                        )}
                    </div>
                    <DialogFooter className="flex-col sm:flex-row gap-2">
                        {!isMyPack && editVerseData.version_source === 'override' && (
                            <Button variant="outline" onClick={handleResetPreference} type="button">
                                Reset to Default
                            </Button>
                        )}
                        <Button onClick={handleUpdateVerse}>Save Changes</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            {/* View Verse Dialog */}
            <Dialog open={viewVerseOpen} onOpenChange={setViewVerseOpen}>
                <DialogContent className="max-w-2xl max-h-[80vh] overflow-hidden flex flex-col">
                    <DialogHeader>
                        <DialogTitle>{viewVerseData?.reference}</DialogTitle>
                        <CardDescription>{viewVerseData?.version}</CardDescription>
                        <DialogDescription className="sr-only">Verse text display</DialogDescription>
                    </DialogHeader>
                    <div className="flex-1 overflow-auto p-4 bg-muted/20 rounded-md mt-2">
                        {loadingVerseText ? (
                            <div className="text-center py-8 text-muted-foreground">Loading text...</div>
                        ) : (
                            <div className="prose dark:prose-invert max-w-none break-words">
                                <div dangerouslySetInnerHTML={{ __html: viewVerseData?.text || "" }} />
                            </div>
                        )}
                    </div>
                    <DialogFooter>
                        <Button onClick={() => setViewVerseOpen(false)}>Close</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            <div className="grid gap-4">
                {verses.map(verse => (
                    <Card
                        key={verse.id}
                        className="cursor-pointer hover:bg-muted/50 transition-colors"
                        onClick={() => handleViewVerse(verse)}
                    >
                        <CardContent className="p-4 flex justify-between items-center">
                            <div className="flex items-center gap-3">
                                <BookOpen className="h-5 w-5 text-primary" />
                                <div>
                                    <div className="font-semibold text-lg">
                                        {verse.title ? (
                                            <>
                                                {verse.title}
                                                <span className="text-muted-foreground font-normal ml-2 text-sm">{verse.reference}</span>
                                            </>
                                        ) : verse.reference}
                                    </div>
                                    <div className="text-xs text-muted-foreground flex gap-2 mt-1">
                                        <span
                                            className={cn("px-1.5 rounded text-xs font-medium",
                                                verse.version_source === 'override' ? "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-100" :
                                                verse.version_source === 'user_default' ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-100" :
                                                "bg-muted text-muted-foreground"
                                            )}
                                            title={
                                                verse.version_source === 'override' ? "Your custom preference" :
                                                verse.version_source === 'user_default' ? "Your default setting" :
                                                "Original version"
                                            }
                                        >
                                            {verse.version}
                                        </span>
                                        {verse.tags?.map(t => (
                                            <span key={t}>#{t}</span>
                                        ))}
                                    </div>
                                </div>
                            </div>
                            <div className="flex gap-2">
                                <Button variant="ghost" size="icon" onClick={(e) => openEdit(e, verse)}>
                                    <Pencil className="h-4 w-4" />
                                </Button>
                                {isMyPack && (
                                    <Button
                                        variant="ghost"
                                        size="icon"
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            if(verse.id) handleDeleteVerse(verse.id);
                                        }}
                                        className="text-destructive"
                                    >
                                        <Trash2 className="h-4 w-4" />
                                    </Button>
                                )}
                            </div>
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
