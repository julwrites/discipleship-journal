import { useState, useEffect, useCallback } from "react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger, DialogFooter, DialogDescription } from "@/components/ui/dialog";
import { Textarea } from "@/components/ui/textarea";
import { useAuth } from "@/hooks/useAuth";
import { useDebounce } from "@/hooks/useDebounce";
import { ChevronDown, ChevronUp, UserPlus, Trash2, Copy, BookOpen, ArrowLeft } from "lucide-react";
import {
    getGroups,
    searchGroups,
    createGroup,
    joinGroup,
    leaveGroup,
    getGroupMembers,
    getGroupShares,
    getSharedItem,
    getConnections,
    addGroupMember,
    removeGroupMember,
    clonePack
} from "@/services/api";
import { useNavigate, useSearchParams } from "react-router-dom";
import { Label } from "@/components/ui/label";

interface Group {
    id: string;
    name: string;
    description: string;
    role?: string;
}

interface GroupMember {
    user_id: string;
    display_name: string;
    email: string;
    role: string;
}

interface SharedItem {
    id: string; // shareId
    note_id?: string;
    verse_pack_id?: string;
    title: string;
    subtitle?: string;
    type: "note" | "verse_pack";
    shared_by: string;
    shared_at: string;
    comment: string;
}

interface Connection {
    id: string;
    requester_id: string;
    receiver_id: string;
    status: string;
    requester_email?: string;
    receiver_email?: string;
}

export default function GroupsPage() {
    const { user } = useAuth();
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    const initialGroupId = searchParams.get("id");

    const [myGroups, setMyGroups] = useState<Group[]>([]);

    const [searchResults, setSearchResults] = useState<Group[]>([]);
    const [searchQuery, setSearchQuery] = useState("");
    const debouncedSearchQuery = useDebounce(searchQuery, 500);

    const [isCreateOpen, setIsCreateOpen] = useState(false);
    const [newGroup, setNewGroup] = useState({ name: "", description: "" });
    const [expandedGroupId, setExpandedGroupId] = useState<string | null>(null);
    const [groupMembers, setGroupMembers] = useState<GroupMember[]>([]);
    const [groupShares, setGroupShares] = useState<SharedItem[]>([]);

    const [myConnections, setMyConnections] = useState<Connection[]>([]);
    const [isAddMemberOpen, setIsAddMemberOpen] = useState(false);
    const [activeGroupId, setActiveGroupId] = useState<string | null>(null);

    // Detailed view of shared item
    // Use proper typing or omit any if complex
    const [viewingSharedItem, setViewingSharedItem] = useState<SharedItem & { content?: { html?: string } | string } | null>(null);
    const [isClonePackOpen, setIsClonePackOpen] = useState(false);
    const [clonePackTitle, setClonePackTitle] = useState("");

    const fetchMyGroups = useCallback(async () => {
        if (!user) return;
        try {
            const data = await getGroups();
            setMyGroups(data || []);
        } catch (error) {
            console.error("Failed to fetch groups", error);
        }
    }, [user]);

    useEffect(() => {
        const load = async () => {
            await fetchMyGroups();
        };
        load();
    }, [fetchMyGroups]);

    // Auto-expand group from URL
    useEffect(() => {
        if (initialGroupId && myGroups.length > 0 && !expandedGroupId) {
            // Verify group exists in myGroups
            const exists = myGroups.find(g => g.id === initialGroupId);
            if (exists) {
                toggleGroupDetails(initialGroupId);
                // Clean URL
                navigate("/groups", { replace: true });
            }
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [initialGroupId, myGroups]); // Run once when groups load

    // Find Groups Search
    useEffect(() => {
        let ignore = false;
        const run = async () => {
            if (debouncedSearchQuery.length < 3) {
                setSearchResults([]);
                return;
            }
            try {
                const data = await searchGroups(debouncedSearchQuery);
                if (!ignore) setSearchResults(data || []);
            } catch (error) {
                if (!ignore) console.error("Search failed", error);
            }
        };
        run();
        return () => { ignore = true; };
    }, [debouncedSearchQuery]);

    const loadConnections = async () => {
        try {
            const data = await getConnections();
            setMyConnections(data || []);
        } catch (error) {
            console.error("Failed to load connections", error);
        }
    };


    const handleCreate = async () => {
        try {
            await createGroup(newGroup);
            setIsCreateOpen(false);
            setNewGroup({ name: "", description: "" });
            fetchMyGroups();
        } catch (error) {
            console.error("Create failed", error);
        }
    };

    const handleJoin = async (id: string) => {
        try {
            await joinGroup(id);
            toast.success("Joined group!");
            fetchMyGroups();
            if (debouncedSearchQuery.length >= 3) {
                 const data = await searchGroups(debouncedSearchQuery);
                 setSearchResults(data || []);
             }
        } catch (error) {
            console.error("Join failed", error);
        }
    };

    const handleLeave = async (id: string) => {
        if (!confirm("Are you sure you want to leave this group?")) return;
        try {
            await leaveGroup(id);
            fetchMyGroups();
        } catch (error) {
            console.error("Leave failed", error);
        }
    };

    const toggleGroupDetails = async (groupId: string) => {
        if (expandedGroupId === groupId) {
            setExpandedGroupId(null);
            setGroupMembers([]);
            return;
        }
        setExpandedGroupId(groupId);
        try {
            const members = await getGroupMembers(groupId);
            setGroupMembers(members || []);

            // Fetch shared notes
            const shares = await getGroupShares(groupId);
            setGroupShares(shares || []);
        } catch (error) {
            console.error("Failed to fetch details", error);
        }
    };

    const viewSharedItem = async (groupId: string, shareId: string) => {
        try {
            const data = await getSharedItem(groupId, shareId);
            setViewingSharedItem(data);
        } catch (error) {
            console.error("Failed to fetch shared item", error);
            toast.error("Failed to load details");
        }
    };

    const handleClonePack = async () => {
        if (!viewingSharedItem?.verse_pack_id) return;
        try {
            await clonePack(viewingSharedItem.verse_pack_id, clonePackTitle || undefined);
            toast.success("Pack cloned to your library!");
            setIsClonePackOpen(false);
            setViewingSharedItem(null);
            navigate("/memory-verses");
        } catch {
            toast.error("Failed to clone pack");
        }
    };

    const addMember = async (userId: string) => {
        if (!activeGroupId) return;
        try {
            await addGroupMember(activeGroupId, userId);
            setIsAddMemberOpen(false);
            toggleGroupDetails(activeGroupId); // Refresh members
        } catch (error) {
            console.error("Failed to add member", error);
            toast.error("Failed to add member. Are you connected?");
        }
    };

    const removeMember = async (groupId: string, userId: string) => {
        if (!confirm("Remove this member?")) return;
        try {
            await removeGroupMember(groupId, userId);
            toggleGroupDetails(groupId); // Refresh members
        } catch (error) {
            console.error("Failed to remove member", error);
        }
    };

    const getSharedContent = () => {
        if (!viewingSharedItem) return "No content";
        const content = viewingSharedItem.content;
        if (typeof content === 'string') return content;
        return content?.html || "No content";
    };

    return (
        <div className="p-4 md:p-8 max-w-4xl mx-auto space-y-6">
             {/* Shared Item Detail Dialog */}
             {viewingSharedItem && (
                 <Dialog open={!!viewingSharedItem} onOpenChange={(o) => !o && setViewingSharedItem(null)}>
                     <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
                         <DialogHeader>
                             <DialogTitle className="flex flex-col gap-1">
                                <span>{viewingSharedItem.title}</span>
                                {viewingSharedItem.subtitle && (
                                    <span className="text-sm font-normal text-muted-foreground bg-muted w-fit px-2 py-0.5 rounded">
                                        {viewingSharedItem.subtitle}
                                    </span>
                                )}
                             </DialogTitle>
                             <DialogDescription>Details of the shared item.</DialogDescription>
                         </DialogHeader>
                         <div className="space-y-4">
                             <div className="bg-muted p-3 rounded text-sm text-muted-foreground">
                                 <p><strong>Shared by:</strong> {viewingSharedItem.shared_by}</p>
                                 <p><strong>Date:</strong> {new Date(viewingSharedItem.shared_at).toLocaleString()}</p>
                                 {viewingSharedItem.comment && <p className="mt-1 italic">"{viewingSharedItem.comment}"</p>}
                             </div>

                             {viewingSharedItem.type === "note" ? (
                                 <div className="prose dark:prose-invert max-w-none">
                                     <div dangerouslySetInnerHTML={{ __html: getSharedContent() }} />
                                 </div>
                             ) : (
                                <div className="space-y-2">
                                    <p className="text-sm text-muted-foreground">This is a shared verse pack.</p>
                                    <div className="flex justify-end">
                                        <Dialog open={isClonePackOpen} onOpenChange={setIsClonePackOpen}>
                                            <DialogTrigger asChild>
                                                <Button>
                                                    <Copy className="mr-2 h-4 w-4" />
                                                    Save to My Packs
                                                </Button>
                                            </DialogTrigger>
                                            <DialogContent>
                                                <DialogHeader>
                                                    <DialogTitle>Save Pack</DialogTitle>
                                                    <DialogDescription>Save this verse pack to your library.</DialogDescription>
                                                </DialogHeader>
                                                <div className="py-4">
                                                    <Label>New Title (Optional)</Label>
                                                    <Input
                                                        value={clonePackTitle}
                                                        onChange={e => setClonePackTitle(e.target.value)}
                                                        placeholder={viewingSharedItem.title}
                                                        className="mt-2"
                                                    />
                                                </div>
                                                <DialogFooter>
                                                    <Button onClick={handleClonePack}>Save</Button>
                                                </DialogFooter>
                                            </DialogContent>
                                        </Dialog>
                                    </div>
                                </div>
                             )}
                         </div>
                     </DialogContent>
                 </Dialog>
             )}

            <div className="flex justify-between items-center">
                <div className="flex items-center gap-2">
                    <Button variant="ghost" size="icon" onClick={() => navigate("/")} aria-label="Go back">
                        <ArrowLeft className="h-6 w-6" />
                    </Button>
                    <h1 className="text-3xl font-bold">Groups</h1>
                </div>
                <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
                    <DialogTrigger asChild>
                        <Button>Create Group</Button>
                    </DialogTrigger>
                    <DialogContent>
                        <DialogHeader>
                            <DialogTitle>Create a New Group</DialogTitle>
                            <DialogDescription>Create a new group to share notes and verses.</DialogDescription>
                        </DialogHeader>
                        <div className="space-y-4">
                            <Input
                                placeholder="Group Name"
                                value={newGroup.name}
                                onChange={(e) => setNewGroup({ ...newGroup, name: e.target.value })}
                            />
                            <Textarea
                                placeholder="Description"
                                value={newGroup.description}
                                onChange={(e) => setNewGroup({ ...newGroup, description: e.target.value })}
                            />
                            <Button onClick={handleCreate} className="w-full">Create</Button>
                        </div>
                    </DialogContent>
                </Dialog>
            </div>

            <Tabs defaultValue="my-groups">
                <TabsList>
                    <TabsTrigger value="my-groups">My Groups</TabsTrigger>
                    <TabsTrigger value="find">Find Groups</TabsTrigger>
                </TabsList>

                <TabsContent value="my-groups" className="space-y-4">
                    {myGroups.length === 0 && <p className="text-muted-foreground">You haven't joined any groups yet.</p>}
                    <div className="space-y-4">
                        {myGroups.map(g => (
                            <Card key={g.id}>
                                <CardHeader className="cursor-pointer" onClick={() => toggleGroupDetails(g.id)}>
                                    <div className="flex justify-between items-center">
                                        <div>
                                            <CardTitle>{g.name}</CardTitle>
                                            <CardDescription>{g.description}</CardDescription>
                                        </div>
                                        <div className="flex items-center gap-2">
                                            <span className="text-sm font-medium bg-muted px-2 py-1 rounded capitalize">{g.role}</span>
                                            {expandedGroupId === g.id ? <ChevronUp size={20} /> : <ChevronDown size={20} />}
                                        </div>
                                    </div>
                                </CardHeader>
                                {expandedGroupId === g.id && (
                                    <CardContent>
                                        <Tabs defaultValue="shares">
                                            <TabsList className="mb-4">
                                                <TabsTrigger value="shares">Shared Items</TabsTrigger>
                                                <TabsTrigger value="members">Members</TabsTrigger>
                                            </TabsList>

                                            <TabsContent value="shares" className="space-y-4">
                                                 {groupShares.length === 0 && <p className="text-sm text-muted-foreground">No items shared yet.</p>}
                                                 {groupShares.map(s => (
                                                     <Card key={s.id} className="bg-muted/50 cursor-pointer hover:bg-muted transition" onClick={() => viewSharedItem(g.id, s.id)}>
                                                         <CardContent className="p-4">
                                                             <div className="flex justify-between items-start">
                                                                 <div>
                                                                     <div className="flex items-center gap-2">
                                                                        {s.type === "verse_pack" ? <BookOpen className="h-4 w-4 text-blue-500" /> : <span className="h-4 w-4">📝</span>}
                                                                        <h4 className="font-bold text-md">{s.title}</h4>
                                                                     </div>
                                                                     <p className="text-xs text-muted-foreground mt-1">Shared by {s.shared_by} on {new Date(s.shared_at).toLocaleDateString()}</p>
                                                                     {s.comment && <p className="text-sm mt-2 italic">"{s.comment}"</p>}
                                                                 </div>
                                                             </div>
                                                         </CardContent>
                                                     </Card>
                                                 ))}
                                            </TabsContent>

                                            <TabsContent value="members" className="space-y-4">
                                                <div className="flex justify-between items-center border-b pb-2">
                                                    <h3 className="font-semibold">Members</h3>
                                                    {g.role === 'admin' && (
                                                        <Dialog open={isAddMemberOpen} onOpenChange={(o) => {
                                                            setIsAddMemberOpen(o);
                                                            if (o) loadConnections();
                                                        }}>
                                                            <DialogTrigger asChild>
                                                                <Button size="sm" variant="outline" onClick={() => setActiveGroupId(g.id)}>
                                                                    <UserPlus size={16} className="mr-2" /> Add Member
                                                                </Button>
                                                            </DialogTrigger>
                                                            <DialogContent>
                                                                <DialogHeader>
                                                                    <DialogTitle>Add Member to {g.name}</DialogTitle>
                                                                    <DialogDescription>Add from your connections.</DialogDescription>
                                                                </DialogHeader>
                                                                <div className="space-y-4">
                                                                    <div className="space-y-2 max-h-60 overflow-y-auto">
                                                                        {myConnections.filter(c => c.status === 'accepted').length === 0 && (
                                                                            <p className="text-muted-foreground text-sm">No connections found. Go to Connections page to add friends.</p>
                                                                        )}
                                                                        {myConnections.filter(c => c.status === 'accepted').map(c => {
                                                                            const isRequester = c.requester_email === user?.email;
                                                                            const otherEmail = isRequester ? c.receiver_email : c.requester_email;
                                                                            const otherId = isRequester ? c.receiver_id : c.requester_id;

                                                                            // Check if already in group
                                                                            const isMember = groupMembers.some(m => m.user_id === otherId);

                                                                            return (
                                                                                <div key={c.id} className="flex justify-between items-center p-2 border rounded">
                                                                                    <div>
                                                                                        <p className="font-medium">{otherEmail}</p>
                                                                                    </div>
                                                                                    {isMember ? (
                                                                                        <Button size="sm" variant="secondary" disabled>Added</Button>
                                                                                    ) : (
                                                                                        <Button size="sm" onClick={() => addMember(otherId)}>Add</Button>
                                                                                    )}
                                                                                </div>
                                                                            );
                                                                        })}
                                                                    </div>
                                                                </div>
                                                            </DialogContent>
                                                        </Dialog>
                                                    )}
                                                </div>
                                                <div className="space-y-2">
                                                    {groupMembers.map(m => (
                                                        <div key={m.user_id} className="flex justify-between items-center">
                                                            <div className="flex items-center gap-2">
                                                                <div>
                                                                    <p className="font-medium">{m.display_name}</p>
                                                                    <p className="text-xs text-muted-foreground">{m.email}</p>
                                                                </div>
                                                                <span className="text-xs bg-muted px-1 rounded">{m.role}</span>
                                                            </div>
                                                            {g.role === 'admin' && m.role !== 'admin' && (
                                                                <Button size="icon" variant="ghost" className="text-destructive h-8 w-8" onClick={() => removeMember(g.id, m.user_id)} aria-label="Remove member">
                                                                    <Trash2 size={16} />
                                                                </Button>
                                                            )}
                                                        </div>
                                                    ))}
                                                </div>
                                                <div className="pt-4 border-t">
                                                    <Button variant="destructive" size="sm" onClick={() => handleLeave(g.id)}>Leave Group</Button>
                                                </div>
                                            </TabsContent>
                                        </Tabs>
                                    </CardContent>
                                )}
                            </Card>
                        ))}
                    </div>
                </TabsContent>

                <TabsContent value="find" className="space-y-4">
                    <div className="flex gap-2">
                        <Input
                            placeholder="Search groups..."
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                        />
                    </div>
                    <div className="grid gap-4 md:grid-cols-2">
                        {searchResults.map(g => (
                            <Card key={g.id}>
                                <CardHeader>
                                    <CardTitle>{g.name}</CardTitle>
                                    <CardDescription>{g.description}</CardDescription>
                                </CardHeader>
                                <CardContent>
                                    {g.role ? (
                                        <Button disabled variant="secondary" className="w-full">Already Member</Button>
                                    ) : (
                                        <Button onClick={() => handleJoin(g.id)} className="w-full">Join Group</Button>
                                    )}
                                </CardContent>
                            </Card>
                        ))}
                    </div>
                </TabsContent>
            </Tabs>
        </div>
    );
}
