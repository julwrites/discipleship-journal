import { useState, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { auth } from "@/lib/firebase";
import { useAuthState } from "react-firebase-hooks/auth";
import { useDebounce } from "@/hooks/useDebounce";
import { ArrowLeft } from "lucide-react";
import {
    getConnections,
    searchUsers as apiSearchUsers,
    sendConnectionRequest,
    respondToConnectionRequest
} from "@/services/api";

interface User {
    id: string;
    email: string;
    username?: string;
}

interface Connection {
    id: string;
    requester_id: string;
    receiver_id: string;
    status: string;
    requester_email?: string;
    receiver_email?: string;
}

export default function ConnectionsPage() {
    const navigate = useNavigate();
    const [user] = useAuthState(auth);
    const [searchQuery, setSearchQuery] = useState("");
    const debouncedSearch = useDebounce(searchQuery, 500);
    const [searchResults, setSearchResults] = useState<User[]>([]);
    const [connections, setConnections] = useState<Connection[]>([]);
    const [loading, setLoading] = useState(false);

    const fetchConnections = useCallback(async () => {
        try {
            const data = await getConnections();
            setConnections(data || []);
        } catch (error) {
            console.error("Failed to fetch connections", error);
        }
    }, []);

    useEffect(() => {
        if (user) {
            fetchConnections();
        }
    }, [user, fetchConnections]);

    useEffect(() => {
        let ignore = false;
        const search = async () => {
            if (debouncedSearch.length < 3) {
                setSearchResults([]);
                return;
            }

            setLoading(true);
            try {
                const data = await apiSearchUsers(debouncedSearch);
                if (!ignore) setSearchResults(data || []);
            } catch (error) {
                if (!ignore) console.error("Search failed", error);
            } finally {
                if (!ignore) setLoading(false);
            }
        };
        search();
        return () => { ignore = true; };
    }, [debouncedSearch]);

    const sendRequest = async (receiverEmail: string) => {
        try {
            await sendConnectionRequest(receiverEmail);
            toast.success("Request sent!");
            fetchConnections();
        } catch (error) {
            console.error("Request failed", error);
            let message = "Unknown error";
            if (error instanceof Error) {
                message = error.message;
            }
            toast.error("Failed: " + message);
        }
    };

    const respondToRequest = async (id: string, action: "accept" | "reject") => {
        try {
            await respondToConnectionRequest(id, action);
            fetchConnections();
        } catch (error) {
            console.error("Response failed", error);
        }
    };

    return (
        <div className="p-8 max-w-4xl mx-auto space-y-6">
            <div className="flex items-center gap-3">
                <Button variant="ghost" size="icon" onClick={() => navigate("/")}>
                    <ArrowLeft className="h-5 w-5" />
                </Button>
                <h1 className="text-3xl font-bold">Connections</h1>
            </div>

            <Tabs defaultValue="connections">
                <TabsList>
                    <TabsTrigger value="connections">My Connections</TabsTrigger>
                    <TabsTrigger value="find">Find People</TabsTrigger>
                </TabsList>

                <TabsContent value="connections" className="space-y-4">
                    <h2 className="text-xl font-semibold mt-4">Pending Requests</h2>
                    {connections.filter(c => c.status === 'pending' && c.receiver_email === user?.email).map(c => (
                         <Card key={c.id}>
                            <CardContent className="flex justify-between items-center p-4">
                                <div>
                                    <p className="font-medium">{c.requester_email}</p>
                                    <p className="text-sm text-muted-foreground">Wants to connect</p>
                                </div>
                                <div className="space-x-2">
                                    <Button size="sm" onClick={() => respondToRequest(c.id, "accept")}>Accept</Button>
                                    <Button size="sm" variant="outline" onClick={() => respondToRequest(c.id, "reject")}>Reject</Button>
                                </div>
                            </CardContent>
                         </Card>
                    ))}
                    {connections.filter(c => c.status === 'pending' && c.requester_email === user?.email).length > 0 && (
                        <div className="mt-4">
                            <h3 className="text-sm font-medium text-muted-foreground uppercase">Sent Requests</h3>
                             {connections.filter(c => c.status === 'pending' && c.requester_email === user?.email).map(c => (
                                <div key={c.id} className="p-2 border-b">
                                    To: {c.receiver_email} (Pending)
                                </div>
                            ))}
                        </div>
                    )}

                    <h2 className="text-xl font-semibold mt-8">My Network</h2>
                    {connections.filter(c => c.status === 'accepted').map(c => {
                        const otherEmail = c.requester_email === user?.email ? c.receiver_email : c.requester_email;
                        return (
                            <Card key={c.id}>
                                <CardContent className="p-4">
                                    <p className="font-medium">{otherEmail}</p>
                                </CardContent>
                            </Card>
                        )
                    })}
                     {connections.filter(c => c.status === 'accepted').length === 0 && (
                        <p className="text-muted-foreground">No connections yet.</p>
                    )}
                </TabsContent>

                <TabsContent value="find" className="space-y-4">
                    <div className="flex gap-2">
                        <Input
                            placeholder="Search by email or username..."
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                        />
                        {/* Search button removed as it's auto-debounced */}
                    </div>
                    {loading && <p className="text-sm text-muted-foreground">Searching...</p>}

                    <div className="space-y-2">
                        {searchResults.map(u => (
                            <Card key={u.id}>
                                <CardContent className="flex justify-between items-center p-4">
                                    <div>
                                        <p className="font-medium">{u.username || u.email}</p>
                                        {u.username && <p className="text-sm text-muted-foreground">{u.email}</p>}
                                    </div>
                                    <Button size="sm" onClick={() => sendRequest(u.email)}>Connect</Button>
                                </CardContent>
                            </Card>
                        ))}
                    </div>
                </TabsContent>
            </Tabs>
        </div>
    );
}
