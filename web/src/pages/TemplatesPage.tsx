import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { getMyTemplates, getPublicTemplates, deleteTemplate, cloneTemplate, StudyTemplate } from "@/services/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { toast } from "sonner";
import { Plus, Edit, Trash2, Copy, Play } from "lucide-react";
import UseTemplateDialog from "./UseTemplateDialog";

export default function TemplatesPage() {
    const navigate = useNavigate();
    const [myTemplates, setMyTemplates] = useState<StudyTemplate[]>([]);
    const [publicTemplates, setPublicTemplates] = useState<StudyTemplate[]>([]);
    const [loading, setLoading] = useState(true);
    const [activeTab, setActiveTab] = useState("my");
    const [selectedTemplate, setSelectedTemplate] = useState<StudyTemplate | null>(null);
    const [useDialogOpen, setUseDialogOpen] = useState(false);

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        setLoading(true);
        try {
            const [myRes, pubRes] = await Promise.all([
                getMyTemplates(),
                getPublicTemplates()
            ]);
            setMyTemplates(myRes.data || []);
            setPublicTemplates(pubRes.data || []);
        } catch {
            toast.error("Failed to load templates");
        } finally {
            setLoading(false);
        }
    };

    const handleDelete = async (id: string) => {
        if (!confirm("Are you sure?")) return;
        try {
            await deleteTemplate(id);
            toast.success("Template deleted");
            loadData();
        } catch {
            toast.error("Failed to delete");
        }
    };

    const handleClone = async (id: string) => {
        try {
            await cloneTemplate(id);
            toast.success("Template cloned to your library");
            loadData();
            setActiveTab("my");
        } catch {
            toast.error("Failed to clone");
        }
    };

    const handleUse = (template: StudyTemplate) => {
        setSelectedTemplate(template);
        setUseDialogOpen(true);
    };

    return (
        <div className="p-4 md:p-8 max-w-5xl mx-auto space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold">Study Templates</h1>
                    <p className="text-muted-foreground">Create reusable AI prompts and study structures.</p>
                </div>
                <Link to="/templates/new">
                    <Button>
                        <Plus className="mr-2 h-4 w-4" /> New Template
                    </Button>
                </Link>
            </div>

            <Tabs value={activeTab} onValueChange={setActiveTab}>
                <TabsList>
                    <TabsTrigger value="my">My Templates</TabsTrigger>
                    <TabsTrigger value="public">Public Library</TabsTrigger>
                </TabsList>

                <TabsContent value="my" className="mt-4">
                    {loading ? <div>Loading...</div> : (
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            {myTemplates.length === 0 ? <p className="text-muted-foreground">No templates yet.</p> :
                                myTemplates.map(t => (
                                    <TemplateCard
                                        key={t.id}
                                        template={t}
                                        isOwner={true}
                                        onDelete={() => handleDelete(t.id!)}
                                        onUse={() => handleUse(t)}
                                        onEdit={() => navigate(`/templates/${t.id}/edit`)}
                                    />
                                ))
                            }
                        </div>
                    )}
                </TabsContent>

                <TabsContent value="public" className="mt-4">
                    {loading ? <div>Loading...</div> : (
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            {publicTemplates.length === 0 ? <p className="text-muted-foreground">No public templates available.</p> :
                                publicTemplates.map(t => (
                                    <TemplateCard
                                        key={t.id}
                                        template={t}
                                        isOwner={false}
                                        onClone={() => handleClone(t.id!)}
                                        onUse={() => handleUse(t)}
                                    />
                                ))
                            }
                        </div>
                    )}
                </TabsContent>
            </Tabs>

            {selectedTemplate && (
                <UseTemplateDialog
                    open={useDialogOpen}
                    onOpenChange={setUseDialogOpen}
                    template={selectedTemplate}
                />
            )}
        </div>
    );
}

function TemplateCard({ template, isOwner, onDelete, onClone, onUse, onEdit }: {
    template: StudyTemplate,
    isOwner: boolean,
    onDelete?: () => void,
    onClone?: () => void,
    onUse: () => void,
    onEdit?: () => void
}) {
    return (
        <Card className="flex flex-col">
            <CardHeader>
                <CardTitle className="text-xl truncate">{template.title}</CardTitle>
                <CardDescription className="line-clamp-2">{template.description || "No description"}</CardDescription>
            </CardHeader>
            <CardContent className="flex-1">
                <div className="text-xs text-muted-foreground">
                    <span className="font-semibold">Fields:</span> {template.fields.map(f => f.label).join(", ")}
                </div>
            </CardContent>
            <CardFooter className="flex justify-between gap-2 border-t pt-4">
                <div className="flex gap-2">
                    <Button variant="default" size="sm" onClick={onUse}>
                        <Play className="mr-2 h-3 w-3" /> Use
                    </Button>
                    {!isOwner && onClone && (
                        <Button variant="secondary" size="sm" onClick={onClone}>
                            <Copy className="mr-2 h-3 w-3" /> Clone
                        </Button>
                    )}
                </div>
                {isOwner && (
                    <div className="flex gap-2">
                        <Button variant="ghost" size="icon" onClick={onEdit}>
                            <Edit className="h-4 w-4" />
                        </Button>
                        <Button variant="ghost" size="icon" className="text-destructive" onClick={onDelete}>
                            <Trash2 className="h-4 w-4" />
                        </Button>
                    </div>
                )}
            </CardFooter>
        </Card>
    );
}
