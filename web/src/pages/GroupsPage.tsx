import { useState, useEffect, useCallback } from "react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Textarea } from "@/components/ui/textarea";
import { useAuth } from "@/hooks/useAuth";
import { useDebounce } from "@/hooks/useDebounce";
import { ChevronDown, ChevronUp, UserPlus, Trash2 } from "lucide-react";
import ReactMarkdown from "react-markdown";
import {
    getGroups,
    searchGroups,
    createGroup,
    joinGroup,
    leaveGroup,
    getGroupMembers,
    getGroupShares,
    getSharedNote,
    searchUsers as apiSearchUsers,
    addGroupMember,
    removeGroupMember
} from "@/services/api";

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

interface SharedNote {
    id: string; // shareId
    note_id: string;
    title: string;
    shared_by: string;
    shared_at: string;
    comment: string;
}

interface UserSearchResult {
    id: string;
    email: string;
    display_name: string;
    username?: string;
}

export default function GroupsPage() {
    const { user } = useAuth();
    const [myGroups, setMyGroups] = useState<Group[]>([]);

    const [searchResults, setSearchResults] = useState<Group[]>([]);
    const [searchQuery, setSearchQuery] = useState("");
    const debouncedSearchQuery = useDebounce(searchQuery, 500);

    const [isCreateOpen, setIsCreateOpen] = useState(false);
    const [newGroup, setNewGroup] = useState({ name: "", description: "" });
    const [expandedGroupId, setExpandedGroupId] = useState<string | null>(null);
    const [groupMembers, setGroupMembers] = useState<GroupMember[]>([]);
    const [groupShares, setGroupShares] = useState<SharedNote[]>([]);

    const [userSearchQuery, setUserSearchQuery] = useState("");
    const debouncedUserSearchQuery = useDebounce(userSearchQuery, 500);
    const [userSearchResults, setUserSearchResults] = useState<UserSearchResult[]>([]);

    const [isAddMemberOpen, setIsAddMemberOpen] = useState(false);
    const [activeGroupId, setActiveGroupId] = useState<string | null>(null);
    const [viewingSharedNote, setViewingSharedNote] = useState<SharedNote & { content: { markdown?: string } } | null>(null);

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

    // Add Member User Search
    useEffect(() => {
        let ignore = false;
        const run = async () => {
            if (debouncedUserSearchQuery.length < 3) {
                setUserSearchResults([]);
                return;
            }
            try {
                const data = await apiSearchUsers(debouncedUserSearchQuery);
                if (!ignore) setUserSearchResults(data || []);
            } catch (error) {
                if (!ignore) console.error("User search failed", error);
            }
        };
        run();
        return () => { ignore = true; };
    }, [debouncedUserSearchQuery]);


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
            // Refresh search results to show updated role
            fetchMyGroups();
            // Re-trigger search if needed, but since it's debounced, manual call to searchGroups might be needed if we want immediate update
            // For now, assume fetchMyGroups handles the "my groups" part, and search results might need refresh.
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

    const viewSharedNote = async (groupId: string, shareId: string) => {
        try {
            const data = await getSharedNote(groupId, shareId);
            setViewingSharedNote(data);
        } catch (error) {
            console.error("Failed to fetch shared note", error);
        }
    };

    const addMember = async (userId: string) => {
        if (!activeGroupId) return;
        try {
            await addGroupMember(activeGroupId, userId);
            setIsAddMemberOpen(false);
            setUserSearchQuery("");
            setUserSearchResults([]);
            toggleGroupDetails(activeGroupId); // Refresh members
        } catch (error) {
            console.error("Failed to add member", error);
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

    return (
        <div className="p-8 max-w-4xl mx-auto space-y-6">
             {viewingSharedNote && (
                 <Dialog open={!!viewingSharedNote} onOpenChange={(o) => !o && setViewingSharedNote(null)}>
                     <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
                         <DialogHeader>
                             <DialogTitle>{viewingSharedNote.title}</DialogTitle>
                         </DialogHeader>
                         <div className="space-y-4">
                             <div className="bg-muted p-3 rounded text-sm text-muted-foreground">
                                 <p><strong>Shared by:</strong> {viewingSharedNote.shared_by}</p>
                                 <p><strong>Date:</strong> {new Date(viewingSharedNote.shared_at).toLocaleString()}</p>
                                 {viewingSharedNote.comment && <p className="mt-1 italic">"{viewingSharedNote.comment}"</p>}
                             </div>
                             <div className="prose dark:prose-invert max-w-none">
                                 <ReactMarkdown>{viewingSharedNote.content?.markdown || "No content"}</ReactMarkdown>
                             </div>
                         </div>
                     </DialogContent>
                 </Dialog>
             )}

            <div className="flex justify-between items-center">
                <h1 className="text-3xl font-bold">Groups</h1>
                <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
                    <DialogTrigger asChild>
                        <Button>Create Group</Button>
                    </DialogTrigger>
                    <DialogContent>
                        <DialogHeader>
                            <DialogTitle>Create a New Group</DialogTitle>
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
                                                <TabsTrigger value="shares">Shared Notes</TabsTrigger>
                                                <TabsTrigger value="members">Members</TabsTrigger>
                                            </TabsList>

                                            <TabsContent value="shares" className="space-y-4">
                                                 {groupShares.length === 0 && <p className="text-sm text-muted-foreground">No notes shared yet.</p>}
                                                 {groupShares.map(s => (
                                                     <Card key={s.id} className="bg-muted/50 cursor-pointer hover:bg-muted transition" onClick={() => viewSharedNote(g.id, s.id)}>
                                                         <CardContent className="p-4">
                                                             <div className="flex justify-between items-start">
                                                                 <div>
                                                                     <h4 className="font-bold text-md">{s.title}</h4>
                                                                     <p className="text-xs text-muted-foreground">Shared by {s.shared_by} on {new Date(s.shared_at).toLocaleDateString()}</p>
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
                                                        <Dialog open={isAddMemberOpen} onOpenChange={setIsAddMemberOpen}>
                                                            <DialogTrigger asChild>
                                                                <Button size="sm" variant="outline" onClick={() => setActiveGroupId(g.id)}>
                                                                    <UserPlus size={16} className="mr-2" /> Add Member
                                                                </Button>
                                                            </DialogTrigger>
                                                            <DialogContent>
                                                                <DialogHeader>
                                                                    <DialogTitle>Add Member to {g.name}</DialogTitle>
                                                                </DialogHeader>
                                                                <div className="space-y-4">
                                                                    <div className="flex gap-2">
                                                                        <Input
                                                                            placeholder="Search by email, name, or username"
                                                                            value={userSearchQuery}
                                                                            onChange={(e) => setUserSearchQuery(e.target.value)}
                                                                        />
                                                                        {/* Button removed */}
                                                                    </div>
                                                                    <div className="space-y-2 max-h-60 overflow-y-auto">
                                                                        {userSearchResults.map(u => (
                                                                            <div key={u.id} className="flex justify-between items-center p-2 border rounded">
                                                                                <div>
                                                                                    <p className="font-medium">{u.display_name}</p>
                                                                                    <p className="text-xs text-muted-foreground">@{u.username || 'unknown'} • {u.email}</p>
                                                                                </div>
                                                                                <Button size="sm" onClick={() => addMember(u.id)}>Add</Button>
                                                                            </div>
                                                                        ))}
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
                                                                <Button size="icon" variant="ghost" className="text-destructive h-8 w-8" onClick={() => removeMember(g.id, m.user_id)}>
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
                         {/* Button removed */}
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
