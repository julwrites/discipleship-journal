import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getTags, createTag, deleteTag, Tag } from "@/services/api";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";
import { Plus, Trash2, ArrowLeft, Tag as TagIcon } from "lucide-react";
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";

export default function TagsPage() {
    const navigate = useNavigate();
    const [tags, setTags] = useState<Tag[]>([]);
    const [loading, setLoading] = useState(true);
    const [createDialogOpen, setCreateDialogOpen] = useState(false);
    const [newTagName, setNewTagName] = useState("");
    const [creating, setCreating] = useState(false);

    // Delete state
    const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false);
    const [tagToDelete, setTagToDelete] = useState<Tag | null>(null);
    const [deleting, setDeleting] = useState(false);

    useEffect(() => {
        loadTags();
    }, []);

    const loadTags = async () => {
        setLoading(true);
        try {
            const data = await getTags();
            setTags(data || []);
        } catch (e) {
            console.error(e);
            toast.error("Failed to load tags");
        } finally {
            setLoading(false);
        }
    };

    const handleCreate = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!newTagName.trim()) return;

        setCreating(true);
        try {
            await createTag(newTagName.trim());
            toast.success("Tag created");
            setNewTagName("");
            setCreateDialogOpen(false);
            loadTags();
        } catch (e) {
            console.error(e);
            toast.error("Failed to create tag");
        } finally {
            setCreating(false);
        }
    };

    const confirmDelete = async () => {
        if (!tagToDelete) return;
        setDeleting(true);
        try {
            await deleteTag(tagToDelete.id);
            toast.success("Tag deleted");
            setTags(tags.filter(t => t.id !== tagToDelete.id));
            setDeleteConfirmOpen(false);
            setTagToDelete(null);
        } catch (e) {
            console.error(e);
            toast.error("Failed to delete tag");
        } finally {
            setDeleting(false);
        }
    };

    return (
        <div className="p-4 md:p-8 max-w-5xl mx-auto space-y-6">
            <div>
                <div className="flex justify-between items-center mb-2">
                    <div className="flex items-center gap-3">
                        <Button variant="ghost" size="icon" onClick={() => navigate("/")} aria-label="Go back">
                            <ArrowLeft className="h-5 w-5" />
                        </Button>
                        <h1 className="text-3xl font-bold">Tags</h1>
                    </div>

                    <Dialog open={createDialogOpen} onOpenChange={setCreateDialogOpen}>
                        <DialogTrigger asChild>
                            <Button>
                                <Plus className="mr-2 h-4 w-4" /> New Tag
                            </Button>
                        </DialogTrigger>
                        <DialogContent>
                            <form onSubmit={handleCreate}>
                                <DialogHeader>
                                    <DialogTitle>Create New Tag</DialogTitle>
                                    <DialogDescription>
                                        Add a new tag to organize your notes.
                                    </DialogDescription>
                                </DialogHeader>
                                <div className="grid gap-4 py-4">
                                    <div className="grid grid-cols-4 items-center gap-4">
                                        <Label htmlFor="name" className="text-right">
                                            Name
                                        </Label>
                                        <Input
                                            id="name"
                                            value={newTagName}
                                            onChange={(e) => setNewTagName(e.target.value)}
                                            className="col-span-3"
                                            placeholder="e.g., Prayer, Study, Goals"
                                        />
                                    </div>
                                </div>
                                <DialogFooter>
                                    <Button type="submit" disabled={creating || !newTagName.trim()}>
                                        {creating ? "Creating..." : "Create Tag"}
                                    </Button>
                                </DialogFooter>
                            </form>
                        </DialogContent>
                    </Dialog>
                </div>
                <p className="text-muted-foreground ml-12">Manage your tags to better organize your journal entries.</p>
            </div>

            {loading ? (
                <div>Loading tags...</div>
            ) : (
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                    {tags.length === 0 ? (
                        <p className="text-muted-foreground col-span-full">No tags created yet.</p>
                    ) : (
                        tags.map(tag => (
                            <Card key={tag.id} className="flex flex-row items-center justify-between p-4 h-16">
                                <div className="flex items-center gap-2 overflow-hidden">
                                    <TagIcon className="h-4 w-4 text-muted-foreground shrink-0" />
                                    <span className="font-medium truncate">{tag.name}</span>
                                </div>
                                <Button
                                    variant="ghost"
                                    size="icon"
                                    className="text-muted-foreground hover:text-destructive shrink-0"
                                    onClick={() => {
                                        setTagToDelete(tag);
                                        setDeleteConfirmOpen(true);
                                    }}
                                    aria-label={`Delete tag ${tag.name}`}
                                >
                                    <Trash2 className="h-4 w-4" />
                                    <span className="sr-only">Delete {tag.name}</span>
                                </Button>
                            </Card>
                        ))
                    )}
                </div>
            )}

            <AlertDialog open={deleteConfirmOpen} onOpenChange={setDeleteConfirmOpen}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Delete Tag</AlertDialogTitle>
                        <AlertDialogDescription>
                            Are you sure you want to delete the tag "{tagToDelete?.name}"?
                            This will remove the tag from all notes, but the notes themselves will not be deleted.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel onClick={() => setTagToDelete(null)}>Cancel</AlertDialogCancel>
                        <AlertDialogAction onClick={confirmDelete} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
                            {deleting ? "Deleting..." : "Delete"}
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </div>
    );
}
