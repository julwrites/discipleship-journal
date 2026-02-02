import { useEffect, useState, useRef } from "react";
import { useParams, Link, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { getReadingPlan, markPlanDayComplete, unmarkPlanDayComplete, getPlanProgress } from "@/services/api";
import { toast } from "sonner";
import { ArrowLeft, CheckCircle, Calendar, StickyNote } from "lucide-react";

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
    plan_type?: string;
}

export default function ReadingPlanDetail() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [plan, setPlan] = useState<ReadingPlan | null>(null);
    const [completedDays, setCompletedDays] = useState<Set<number>>(new Set());
    const [loading, setLoading] = useState(true);
    const todayRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (id) {
            loadData(id);
        }
    }, [id]);

    useEffect(() => {
        // Scroll to today if it's a calendar plan
        if (plan?.plan_type === 'calendar' && todayRef.current) {
            todayRef.current.scrollIntoView({ behavior: 'smooth', block: 'center' });
        }
    }, [plan]);

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

        const isCompleted = completedDays.has(dayNumber);
        const newCompleted = new Set(completedDays);

        if (isCompleted) {
            newCompleted.delete(dayNumber);
            setCompletedDays(newCompleted); // Optimistic

            try {
                await unmarkPlanDayComplete(id, dayNumber);
                toast.success(`Day ${dayNumber} unmarked`);
            } catch (error) {
                console.error(error);
                toast.error("Failed to unmark day");
                // Revert
                newCompleted.add(dayNumber);
                setCompletedDays(new Set(newCompleted));
            }
        } else {
            newCompleted.add(dayNumber);
            setCompletedDays(newCompleted); // Optimistic

            try {
                await markPlanDayComplete(id, dayNumber);
                toast.success(`Day ${dayNumber} completed!`);
            } catch (error) {
                console.error(error);
                toast.error("Failed to update progress");
                // Revert
                newCompleted.delete(dayNumber);
                setCompletedDays(new Set(newCompleted));
            }
        }
    };

    const handleCreateNote = (day: ReadingPlanDay) => {
        if (!plan) return;
        navigate("/notes/new", {
            state: {
                title: `Bible Reading: ${day.passage}`,
                tags: ["Bible Reading"],
                passageRef: day.passage,
                context: {
                    planId: plan.id,
                    planTitle: plan.title,
                    dayNumber: day.day_number
                }
            }
        });
    };

    const getDayOfYear = () => {
        const now = new Date();
        const start = new Date(now.getFullYear(), 0, 0);
        const diff = now.getTime() - start.getTime();
        const oneDay = 1000 * 60 * 60 * 24;
        return Math.floor(diff / oneDay);
    };

    const currentDayOfYear = getDayOfYear();

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
                <div className="flex items-start justify-between">
                    <div>
                        <h1 className="text-3xl font-bold mb-2">{plan.title}</h1>
                        <p className="text-muted-foreground">{plan.description}</p>
                    </div>
                    {plan.plan_type === 'calendar' && (
                        <Button
                            variant="outline"
                            onClick={() => todayRef.current?.scrollIntoView({ behavior: 'smooth', block: 'center' })}
                        >
                            <Calendar className="mr-2 h-4 w-4" />
                            Jump to Today
                        </Button>
                    )}
                </div>

                <div className="mt-4 flex items-center gap-2 text-sm font-medium">
                    <div className="bg-primary/10 text-primary px-3 py-1 rounded-full">
                        {completedDays.size} / {plan.days?.length || 0} Days Completed
                    </div>
                </div>
            </div>

            <div className="space-y-4">
                {plan.days?.map((day) => {
                    const isCompleted = completedDays.has(day.day_number);
                    const isToday = plan.plan_type === 'calendar' && day.day_number === currentDayOfYear;

                    return (
                        <div key={day.id} ref={isToday ? todayRef : null}>
                            <Card className={`transition-all ${isCompleted ? 'bg-muted/50' : 'bg-card'} ${isToday ? 'border-primary ring-1 ring-primary' : ''}`}>
                                <CardContent className="flex items-center justify-between p-4">
                                    <div className="flex items-center gap-4">
                                        <div className={`h-8 w-8 rounded-full flex items-center justify-center font-bold text-sm ${isCompleted ? 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300' : isToday ? 'bg-primary text-primary-foreground' : 'bg-secondary text-secondary-foreground'}`}>
                                            {day.day_number}
                                        </div>
                                        <div>
                                            <div className="flex items-center gap-2">
                                                <div className="font-medium">Day {day.day_number}</div>
                                                {isToday && <span className="text-xs bg-primary/20 text-primary px-1.5 rounded">Today</span>}
                                            </div>
                                            <div className="text-sm text-muted-foreground">{day.passage}</div>
                                        </div>
                                    </div>
                                    <div className="flex items-center gap-2">
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            onClick={() => handleCreateNote(day)}
                                            title="Create Note"
                                        >
                                            <StickyNote className="h-4 w-4" />
                                        </Button>
                                        {isCompleted ? (
                                            <Button
                                                variant="ghost"
                                                size="icon"
                                                onClick={() => handleMarkComplete(day.day_number)}
                                                className="text-green-600 dark:text-green-400 hover:text-destructive hover:bg-destructive/10"
                                                title="Unmark"
                                            >
                                                <CheckCircle className="h-6 w-6" />
                                            </Button>
                                        ) : (
                                            <Button variant="outline" size="sm" onClick={() => handleMarkComplete(day.day_number)}>
                                                Mark Complete
                                            </Button>
                                        )}
                                    </div>
                                </CardContent>
                            </Card>
                        </div>
                    );
                })}
            </div>
        </div>
    );
}
