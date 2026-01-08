import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getTemplate, createTemplate, updateTemplate, StudyTemplate, TemplateField } from "@/services/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { toast } from "sonner";
import { Trash2, Plus, ArrowLeft } from "lucide-react";

export default function TemplateEditor() {
    const { id } = useParams();
    const navigate = useNavigate();
    const isEdit = id && id !== "new";
    const [loading, setLoading] = useState(isEdit);

    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [isPublic, setIsPublic] = useState(false);
    const [systemPrompt, setSystemPrompt] = useState("");
    const [fields, setFields] = useState<TemplateField[]>([]);

    useEffect(() => {
        if (isEdit) {
            loadTemplate(id!);
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [id]);

    const loadTemplate = async (tmplId: string) => {
        try {
            const data = await getTemplate(tmplId);
            setTitle(data.title);
            setDescription(data.description);
            setIsPublic(data.is_public);
            setFields(data.fields || []);

            if (data.prompts && typeof data.prompts === 'object') {
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
                setSystemPrompt((data.prompts as any).system || "");
            }
        } catch {
            toast.error("Failed to load template");
            navigate("/templates");
        } finally {
            setLoading(false);
        }
    };

    const handleSave = async () => {
        if (!title) return toast.error("Title is required");

        const templateData: StudyTemplate = {
            title,
            description,
            is_public: isPublic,
            fields,
            prompts: { system: systemPrompt },
            structure: {} // Default empty for now
        };

        try {
            if (isEdit) {
                await updateTemplate(id!, templateData);
                toast.success("Template updated");
            } else {
                await createTemplate(templateData);
                toast.success("Template created");
            }
            navigate("/templates");
        } catch {
            toast.error("Failed to save template");
        }
    };

    const addField = () => {
        setFields([...fields, { key: "", label: "", type: "text", placeholder: "" }]);
    };

    const updateField = (index: number, key: keyof TemplateField, value: string) => {
        const newFields = [...fields];
        // @ts-expect-error Dynamic access
        newFields[index][key] = value;
        setFields(newFields);
    };

    const removeField = (index: number) => {
        setFields(fields.filter((_, i) => i !== index));
    };

    if (loading) return <div className="p-8">Loading...</div>;

    return (
        <div className="p-4 md:p-8 max-w-3xl mx-auto space-y-6">
            <div className="flex items-center gap-4">
                <Button variant="ghost" onClick={() => navigate("/templates")}>
                    <ArrowLeft className="h-4 w-4" />
                </Button>
                <h1 className="text-2xl font-bold">{isEdit ? "Edit Template" : "New Template"}</h1>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle>Basic Info</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                    <div className="grid gap-2">
                        <Label>Title</Label>
                        <Input value={title} onChange={e => setTitle(e.target.value)} placeholder="e.g. Inductive Study" />
                    </div>
                    <div className="grid gap-2">
                        <Label>Description</Label>
                        <Textarea value={description} onChange={e => setDescription(e.target.value)} placeholder="Brief description..." />
                    </div>
                    <div className="flex items-center gap-2">
                        <Checkbox
                            id="public"
                            checked={isPublic}
                            onCheckedChange={(c: boolean | string) => setIsPublic(!!c)}
                        />
                        <Label htmlFor="public">Make Public</Label>
                    </div>
                </CardContent>
            </Card>

            <Card>
                <CardHeader>
                    <CardTitle>AI Prompt</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="grid gap-2">
                        <Label>System Instructions</Label>
                        <Textarea
                            value={systemPrompt}
                            onChange={e => setSystemPrompt(e.target.value)}
                            placeholder="You are a Bible teacher. When the user provides a passage, analyze it by..."
                            rows={6}
                        />
                        <p className="text-xs text-muted-foreground">This prompt guides the AI on how to process the user inputs.</p>
                    </div>
                </CardContent>
            </Card>

            <Card>
                <CardHeader className="flex flex-row items-center justify-between">
                    <CardTitle>Input Fields</CardTitle>
                    <Button size="sm" variant="outline" onClick={addField}>
                        <Plus className="mr-2 h-4 w-4" /> Add Field
                    </Button>
                </CardHeader>
                <CardContent className="space-y-4">
                    {fields.length === 0 ? <p className="text-sm text-muted-foreground text-center">No input fields defined.</p> :
                        fields.map((field, idx) => (
                            <div key={idx} className="flex gap-2 items-start border p-3 rounded bg-muted/20">
                                <div className="grid gap-2 flex-1">
                                    <div className="flex gap-2">
                                        <Input
                                            placeholder="Label (e.g. Target Audience)"
                                            value={field.label}
                                            onChange={e => updateField(idx, "label", e.target.value)}
                                        />
                                        <Input
                                            placeholder="Key (e.g. audience)"
                                            value={field.key}
                                            onChange={e => updateField(idx, "key", e.target.value)}
                                        />
                                    </div>
                                    <Input
                                        placeholder="Placeholder text..."
                                        value={field.placeholder || ""}
                                        onChange={e => updateField(idx, "placeholder", e.target.value)}
                                    />
                                </div>
                                <Button variant="ghost" size="icon" className="text-destructive mt-1" onClick={() => removeField(idx)}>
                                    <Trash2 className="h-4 w-4" />
                                </Button>
                            </div>
                        ))
                    }
                </CardContent>
            </Card>

            <div className="flex justify-end gap-2">
                <Button variant="outline" onClick={() => navigate("/templates")}>Cancel</Button>
                <Button onClick={handleSave}>Save Template</Button>
            </div>
        </div>
    );
}
