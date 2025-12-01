import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { auth } from "@/lib/firebase";
import { useAuthState } from "react-firebase-hooks/auth";
import { useCallback } from "react";

interface User {
    id: string;
    email: string;
    display_name: string;
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
    const [user] = useAuthState(auth);
    const [searchQuery, setSearchQuery] = useState("");
    const [searchResults, setSearchResults] = useState<User[]>([]);
    const [connections, setConnections] = useState<Connection[]>([]);
    const [loading, setLoading] = useState(false);

    const fetchConnections = useCallback(async () => {
        try {
            const token = await user?.getIdToken();
            const res = await fetch(`${import.meta.env.VITE_API_URL}/api/connections`, {
                headers: { Authorization: `Bearer ${token}` }
            });
            if (res.ok) {
                const data = await res.json();
                setConnections(data || []);
            }
        } catch (error) {
            console.error("Failed to fetch connections", error);
        }
    }, [user]);

    useEffect(() => {
        if (user) {
            fetchConnections();
        }
    }, [user, fetchConnections]);

    const handleSearch = async () => {
        if (searchQuery.length < 3) return;
        setLoading(true);
        try {
            const token = await user?.getIdToken();
            const res = await fetch(`${import.meta.env.VITE_API_URL}/api/users/search?q=${searchQuery}`, {
                headers: { Authorization: `Bearer ${token}` }
            });
            if (res.ok) {
                const data = await res.json();
                setSearchResults(data || []);
            }
        } catch (error) {
            console.error("Search failed", error);
        } finally {
            setLoading(false);
        }
    };

    const sendRequest = async (receiverEmail: string) => {
        try {
            const token = await user?.getIdToken();
            const res = await fetch(`${import.meta.env.VITE_API_URL}/api/connections/request`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`
                },
                body: JSON.stringify({ receiver_email: receiverEmail })
            });
            if (res.ok) {
                alert("Request sent!");
                fetchConnections();
            } else {
                const err = await res.text();
                alert("Failed: " + err);
            }
        } catch (error) {
            console.error("Request failed", error);
        }
    };

    const respondToRequest = async (id: string, action: "accept" | "reject") => {
        try {
            const token = await user?.getIdToken();
            const method = action === "reject" ? "DELETE" : "PUT";
            const url = `${import.meta.env.VITE_API_URL}/api/connections/${id}?action=${action}`;

            const res = await fetch(url, {
                method: method,
                headers: { Authorization: `Bearer ${token}` }
            });
            if (res.ok) {
                fetchConnections();
            }
        } catch (error) {
            console.error("Response failed", error);
        }
    };

    return (
        <div className="p-8 max-w-4xl mx-auto space-y-6">
            <h1 className="text-3xl font-bold">Connections</h1>

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
                                    <p className="text-sm text-gray-500">Wants to connect</p>
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
                            <h3 className="text-sm font-medium text-gray-500 uppercase">Sent Requests</h3>
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
                        <p className="text-gray-500">No connections yet.</p>
                    )}
                </TabsContent>

                <TabsContent value="find" className="space-y-4">
                    <div className="flex gap-2">
                        <Input
                            placeholder="Search by email or name..."
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                        />
                        <Button onClick={handleSearch} disabled={loading}>Search</Button>
                    </div>

                    <div className="space-y-2">
                        {searchResults.map(u => (
                            <Card key={u.id}>
                                <CardContent className="flex justify-between items-center p-4">
                                    <div>
                                        <p className="font-medium">{u.display_name || "User"}</p>
                                        <p className="text-sm text-gray-500">{u.email}</p>
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
