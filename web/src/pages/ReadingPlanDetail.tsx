import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { getReadingPlan, markPlanDayComplete, getPlanProgress } from "@/services/api";
import { toast } from "sonner";
import { ArrowLeft, CheckCircle } from "lucide-react";

interface ReadingPlanDay {
    id: string;
    day_number: number;
    passage: string;
}

interface ReadingPlan {
    id: string;
    title: string;
    description: string;
    days: ReadingPlanDay[];
}

export default function ReadingPlanDetail() {
    const { id } = useParams<{ id: string }>();
    const [plan, setPlan] = useState<ReadingPlan | null>(null);
    const [completedDays, setCompletedDays] = useState<Set<number>>(new Set());
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (id) {
            loadData(id);
        }
    }, [id]);

    const loadData = async (planId: string) => {
        setLoading(true);
        try {
            const [planRes, progressRes] = await Promise.all([
                getReadingPlan(planId),
                getPlanProgress(planId)
            ]);
            setPlan(planRes);
            setCompletedDays(new Set(progressRes.completed_days));
        } catch (error) {
            console.error(error);
            toast.error("Failed to load reading plan details");
        } finally {
            setLoading(false);
        }
    };

    const handleMarkComplete = async (dayNumber: number) => {
        if (!id) return;
        try {
            // Optimistic update
            const newCompleted = new Set(completedDays);
            newCompleted.add(dayNumber);
            setCompletedDays(newCompleted);

            await markPlanDayComplete(id, dayNumber);
            toast.success(`Day ${dayNumber} completed!`);
        } catch (error) {
            console.error(error);
            toast.error("Failed to update progress");
            // Revert
            const reverted = new Set(completedDays);
            reverted.delete(dayNumber);
            setCompletedDays(reverted);
        }
    };

    if (loading) return <div className="p-8 text-center">Loading...</div>;
    if (!plan) return <div className="p-8 text-center">Plan not found</div>;

    return (
        <div className="p-4 md:p-8 max-w-4xl mx-auto">
            <Link to="/reading-plans">
                <Button variant="ghost" className="mb-4 pl-0 hover:bg-transparent hover:text-primary">
                    <ArrowLeft className="mr-2 h-4 w-4" />
                    Back to Plans
                </Button>
            </Link>

            <div className="mb-8">
                <h1 className="text-3xl font-bold mb-2">{plan.title}</h1>
                <p className="text-muted-foreground">{plan.description}</p>
                <div className="mt-4 flex items-center gap-2 text-sm font-medium">
                    <div className="bg-primary/10 text-primary px-3 py-1 rounded-full">
                        {completedDays.size} / {plan.days?.length || 0} Days Completed
                    </div>
                </div>
            </div>

            <div className="space-y-4">
                {plan.days?.map((day) => {
                    const isCompleted = completedDays.has(day.day_number);
                    return (
                        <Card key={day.id} className={`transition-colors ${isCompleted ? 'bg-muted/50' : 'bg-card'}`}>
                            <CardContent className="flex items-center justify-between p-4">
                                <div className="flex items-center gap-4">
                                    <div className={`h-8 w-8 rounded-full flex items-center justify-center font-bold text-sm ${isCompleted ? 'bg-green-100 text-green-700' : 'bg-secondary text-secondary-foreground'}`}>
                                        {day.day_number}
                                    </div>
                                    <div>
                                        <div className="font-medium">Day {day.day_number}</div>
                                        <div className="text-sm text-muted-foreground">{day.passage}</div>
                                    </div>
                                </div>
                                {isCompleted ? (
                                    <Button variant="ghost" size="icon" disabled className="text-green-600">
                                        <CheckCircle className="h-6 w-6" />
                                    </Button>
                                ) : (
                                    <Button variant="outline" size="sm" onClick={() => handleMarkComplete(day.day_number)}>
                                        Mark Complete
                                    </Button>
                                )}
                            </CardContent>
                        </Card>
                    );
                })}
            </div>
        </div>
    );
}
