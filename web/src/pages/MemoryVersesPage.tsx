import { useEffect, useState } from "react";
import { useDebounce } from "@/hooks/useDebounce";
import { searchMemoryVerses, createMemoryVerse, getMemoryVersePacks, MemoryVerse } from "@/services/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent } from "@/components/ui/card";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger, DialogFooter } from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { toast } from "sonner";
import { Search, Plus, BookOpen, Tag } from "lucide-react";

export default function MemoryVersesPage() {
    const [searchQuery, setSearchQuery] = useState("");
    const debouncedSearch = useDebounce(searchQuery, 500);
    const [verses, setVerses] = useState<MemoryVerse[]>([]);
    const [packs, setPacks] = useState<string[]>([]);
    const [loading, setLoading] = useState(false);
    const [isCreateOpen, setIsCreateOpen] = useState(false);

    // Create Form State
    const [newVerse, setNewVerse] = useState<MemoryVerse>({
        pack_name: "My Verses",
        reference: "",
        text: "",
        version: "ESV",
        tags: []
    });
    const [tagInput, setTagInput] = useState("");

    useEffect(() => {
        loadPacks();
        doSearch("");
    }, []);

    useEffect(() => {
        doSearch(debouncedSearch);
    }, [debouncedSearch]);

    const loadPacks = async () => {
        try {
            const res = await getMemoryVersePacks();
            setPacks(res.data || []);
        } catch (error) {
            console.error(error);
        }
    };

    const doSearch = async (query: string) => {
        setLoading(true);
        try {
            const res = await searchMemoryVerses(query);
            setVerses(res.data || []);
        } catch (error) {
            toast.error("Failed to load verses");
        } finally {
            setLoading(false);
        }
    };

    const handleCreate = async () => {
        if (!newVerse.reference || !newVerse.text) {
            toast.error("Reference and Text are required");
            return;
        }
        try {
            await createMemoryVerse(newVerse);
            toast.success("Verse created!");
            setIsCreateOpen(false);
            setNewVerse({
                pack_name: "My Verses",
                reference: "",
                text: "",
                version: "ESV",
                tags: []
            });
            doSearch(debouncedSearch);
        } catch (error) {
            toast.error("Failed to create verse");
        }
    };

    const addTag = () => {
        if (tagInput.trim()) {
            setNewVerse(prev => ({
                ...prev,
                tags: [...prev.tags, tagInput.trim()]
            }));
            setTagInput("");
        }
    };

    return (
        <div className="p-4 md:p-8 max-w-4xl mx-auto space-y-6">
            <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
                <div>
                    <h1 className="text-3xl font-bold">Scripture Memory</h1>
                    <p className="text-muted-foreground">Memorize and meditate on God's Word</p>
                </div>
                <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
                    <DialogTrigger asChild>
                        <Button>
                            <Plus className="mr-2 h-4 w-4" />
                            Add Verse
                        </Button>
                    </DialogTrigger>
                    <DialogContent>
                        <DialogHeader>
                            <DialogTitle>Add Memory Verse</DialogTitle>
                        </DialogHeader>
                        <div className="space-y-4 py-4">
                            <div className="grid gap-2">
                                <Label>Pack Name</Label>
                                <Input
                                    value={newVerse.pack_name}
                                    onChange={e => setNewVerse({...newVerse, pack_name: e.target.value})}
                                    placeholder="e.g. My Verses"
                                />
                            </div>
                            <div className="grid gap-2">
                                <Label>Reference</Label>
                                <Input
                                    value={newVerse.reference}
                                    onChange={e => setNewVerse({...newVerse, reference: e.target.value})}
                                    placeholder="e.g. John 3:16"
                                />
                            </div>
                            <div className="grid gap-2">
                                <Label>Scripture Text</Label>
                                <Textarea
                                    value={newVerse.text}
                                    onChange={e => setNewVerse({...newVerse, text: e.target.value})}
                                    placeholder="For God so loved..."
                                    rows={4}
                                />
                            </div>
                            <div className="grid gap-2">
                                <Label>Version</Label>
                                <Select
                                    value={newVerse.version}
                                    onValueChange={v => setNewVerse({...newVerse, version: v})}
                                >
                                    <SelectTrigger>
                                        <SelectValue />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="ESV">ESV</SelectItem>
                                        <SelectItem value="NIV">NIV</SelectItem>
                                        <SelectItem value="KJV">KJV</SelectItem>
                                        <SelectItem value="NASB">NASB</SelectItem>
                                    </SelectContent>
                                </Select>
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
                                    {newVerse.tags.map((tag, i) => (
                                        <span key={i} className="bg-secondary text-secondary-foreground px-2 py-1 rounded text-xs flex items-center gap-1">
                                            {tag}
                                            <button
                                                onClick={() => setNewVerse(prev => ({...prev, tags: prev.tags.filter((_, idx) => idx !== i)}))}
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
                            <Button variant="outline" onClick={() => setIsCreateOpen(false)}>Cancel</Button>
                            <Button onClick={handleCreate}>Save Verse</Button>
                        </DialogFooter>
                    </DialogContent>
                </Dialog>
            </div>

            <div className="relative">
                <Search className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                <Input
                    placeholder="Search verses, references, or packs..."
                    className="pl-9"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                />
            </div>

            {loading ? (
                <div className="text-center py-8">Loading...</div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {verses.length === 0 ? (
                        <div className="col-span-full text-center py-8 text-muted-foreground">
                            No verses found. Try adding one!
                        </div>
                    ) : (
                        verses.map((verse) => (
                            <Card key={verse.id} className="bg-card hover:bg-muted/50 transition-colors">
                                <CardContent className="p-4 space-y-3">
                                    <div className="flex justify-between items-start">
                                        <div className="font-semibold flex items-center gap-2">
                                            <BookOpen className="h-4 w-4 text-primary" />
                                            {verse.reference}
                                        </div>
                                        <span className="text-xs text-muted-foreground bg-secondary px-2 py-1 rounded">
                                            {verse.version}
                                        </span>
                                    </div>
                                    <p className="text-sm text-foreground/90 italic">
                                        "{verse.text}"
                                    </p>
                                    <div className="flex items-center justify-between text-xs text-muted-foreground pt-2 border-t">
                                        <div className="flex items-center gap-2">
                                            <Tag className="h-3 w-3" />
                                            <span>{verse.pack_name}</span>
                                        </div>
                                        {verse.tags && verse.tags.length > 0 && (
                                            <div className="flex gap-1">
                                                {verse.tags.slice(0, 3).map(t => (
                                                    <span key={t} className="bg-primary/10 text-primary px-1.5 rounded">
                                                        {t}
                                                    </span>
                                                ))}
                                            </div>
                                        )}
                                    </div>
                                </CardContent>
                            </Card>
                        ))
                    )}
                </div>
            )}
        </div>
    );
}
