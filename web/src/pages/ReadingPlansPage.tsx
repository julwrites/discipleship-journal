import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { getReadingPlans, getMyReadingPlans, subscribeToPlan } from "@/services/api";
import { toast } from "sonner";
import { ArrowLeft, BookOpen, Check, Calendar } from "lucide-react";

interface ReadingPlan {
    id: string;
    title: string;
    description: string;
    days: number;
}

interface UserReadingPlan {
    id: string;
    reading_plan_id: string;
    start_date: string;
    status: string;
    plan?: ReadingPlan;
}

export default function ReadingPlansPage() {
    const [allPlans, setAllPlans] = useState<ReadingPlan[]>([]);
    const [myPlans, setMyPlans] = useState<UserReadingPlan[]>([]);
    const [loading, setLoading] = useState(true);
    const [activeTab, setActiveTab] = useState("my-plans");

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        setLoading(true);
        try {
            const [plansRes, myPlansRes] = await Promise.all([
                getReadingPlans(),
                getMyReadingPlans()
            ]);
            setAllPlans(plansRes.data || []);
            setMyPlans(myPlansRes.data || []);

            // If user has no plans, default to browse
            if ((!myPlansRes.data || myPlansRes.data.length === 0) && plansRes.data && plansRes.data.length > 0) {
                 // Only switch if we are strictly initializing?
                 // Maybe better to leave it to user, but let's see.
                 // Actually, let's keep it simple. Default to my-plans.
            }
        } catch (error) {
            console.error(error);
            toast.error("Failed to load reading plans");
        } finally {
            setLoading(false);
        }
    };

    const handleSubscribe = async (planId: string) => {
        try {
            await subscribeToPlan(planId);
            toast.success("Subscribed to plan!");
            loadData(); // Reload to update lists
            setActiveTab("my-plans");
        } catch (error) {
            console.error(error);
            toast.error(error instanceof Error ? error.message : "Failed to subscribe");
        }
    };

    return (
        <div className="p-4 md:p-8 max-w-7xl mx-auto">
            <div className="mb-8">
                <Link to="/">
                    <Button variant="ghost" className="mb-4 pl-0 hover:bg-transparent hover:text-primary">
                        <ArrowLeft className="mr-2 h-4 w-4" />
                        Back to Dashboard
                    </Button>
                </Link>
                <h1 className="text-3xl font-bold flex items-center gap-2">
                    <BookOpen className="h-8 w-8" />
                    Bible Reading Plans
                </h1>
                <p className="text-muted-foreground mt-2">
                    Follow structured plans to read through the Bible.
                </p>
            </div>

            <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
                <TabsList className="mb-4">
                    <TabsTrigger value="my-plans">My Plans</TabsTrigger>
                    <TabsTrigger value="browse">Browse Plans</TabsTrigger>
                </TabsList>

                <TabsContent value="my-plans">
                    {loading ? (
                        <div className="text-center py-8">Loading...</div>
                    ) : myPlans.length === 0 ? (
                        <div className="text-center py-12 bg-muted/30 rounded-lg">
                            <h3 className="text-lg font-medium mb-2">No active plans</h3>
                            <p className="text-muted-foreground mb-4">You haven't subscribed to any reading plans yet.</p>
                            <Button variant="outline" onClick={() => setActiveTab("browse")}>
                                Browse Plans
                            </Button>
                        </div>
                    ) : (
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                            {myPlans.map((userPlan) => (
                                <Card key={userPlan.id} className="hover:shadow-md transition-shadow">
                                    <CardHeader>
                                        <CardTitle>{userPlan.plan?.title || "Unknown Plan"}</CardTitle>
                                        <CardDescription>
                                            Started on {new Date(userPlan.start_date).toLocaleDateString()}
                                        </CardDescription>
                                    </CardHeader>
                                    <CardContent>
                                        <div className="flex justify-between items-center mb-4">
                                            <div className="text-sm text-muted-foreground">
                                                {userPlan.plan?.days} Days
                                            </div>
                                            <div className="flex items-center text-sm font-medium text-primary">
                                                <span className={`px-2 py-1 rounded-full text-xs ${userPlan.status === 'completed' ? 'bg-green-100 text-green-800' : 'bg-blue-100 text-blue-800'}`}>
                                                    {userPlan.status.toUpperCase()}
                                                </span>
                                            </div>
                                        </div>
                                        <Link to={`/reading-plans/${userPlan.plan?.id}`}>
                                            <Button className="w-full">Continue Reading</Button>
                                        </Link>
                                    </CardContent>
                                </Card>
                            ))}
                        </div>
                    )}
                </TabsContent>

                <TabsContent value="browse">
                    {loading ? (
                        <div className="text-center py-8">Loading...</div>
                    ) : (
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                            {allPlans.map((plan) => {
                                const isSubscribed = myPlans.some(mp => mp.reading_plan_id === plan.id && mp.status === 'active');
                                return (
                                    <Card key={plan.id} className="flex flex-col">
                                        <CardHeader>
                                            <CardTitle>{plan.title}</CardTitle>
                                            <CardDescription className="line-clamp-2">{plan.description}</CardDescription>
                                        </CardHeader>
                                        <CardContent className="mt-auto">
                                            <div className="flex items-center gap-2 text-sm text-muted-foreground mb-4">
                                                <Calendar className="h-4 w-4" />
                                                {plan.days} Days
                                            </div>
                                            {isSubscribed ? (
                                                <Button variant="secondary" className="w-full" disabled>
                                                    <Check className="mr-2 h-4 w-4" /> Subscribed
                                                </Button>
                                            ) : (
                                                <Button className="w-full" onClick={() => handleSubscribe(plan.id)}>
                                                    Subscribe
                                                </Button>
                                            )}
                                        </CardContent>
                                    </Card>
                                );
                            })}
                        </div>
                    )}
                </TabsContent>
            </Tabs>
        </div>
    );
}
