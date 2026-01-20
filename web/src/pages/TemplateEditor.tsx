import { useEffect, useState, useCallback } from "react";
import { useNavigate, useParams, useBlocker } from "react-router-dom";
import { getTemplate, createTemplate, updateTemplate, StudyTemplate, TemplateField } from "@/services/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { toast } from "sonner";
import { Trash2, Plus, ArrowLeft, Info } from "lucide-react";
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

    const [initialState, setInitialState] = useState("");
    const [saving, setSaving] = useState(false);

    // Derived state for dirty check
    const currentState = JSON.stringify({ title, description, isPublic, systemPrompt, fields });
    const isDirty = initialState !== "" && currentState !== initialState;

    useEffect(() => {
        if (isEdit) {
            loadTemplate(id!);
        } else {
            // Set initial state for new template
            setInitialState(JSON.stringify({ title: "", description: "", isPublic: false, systemPrompt: "", fields: [] }));
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

            let sysPrompt = "";
            if (data.prompts && typeof data.prompts === 'object') {
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
                sysPrompt = (data.prompts as any).system || "";
            }
            setSystemPrompt(sysPrompt);

            setInitialState(JSON.stringify({
                title: data.title,
                description: data.description,
                isPublic: data.is_public,
                systemPrompt: sysPrompt,
                fields: data.fields || []
            }));
        } catch {
            toast.error("Failed to load template");
            navigate("/templates");
        } finally {
            setLoading(false);
        }
    };

    const saveTemplate = useCallback(async (shouldNavigate = true) => {
        if (!title) {
            toast.error("Title is required");
            return false;
        }

        setSaving(true);
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

            // Update initial state to match current, so isDirty becomes false
            setInitialState(JSON.stringify({ title, description, isPublic, systemPrompt, fields }));

            if (shouldNavigate) {
                navigate("/templates");
            }
            return true;
        } catch {
            toast.error("Failed to save template");
            return false;
        } finally {
            setSaving(false);
        }
    }, [title, description, isPublic, fields, systemPrompt, id, isEdit, navigate]);

    const handleSave = () => saveTemplate(true);

    const blocker = useBlocker(
        ({ currentLocation, nextLocation }) =>
            !saving && isDirty && currentLocation.pathname !== nextLocation.pathname
    );

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

            <div className="bg-blue-50 dark:bg-blue-950/30 border border-blue-200 dark:border-blue-900 rounded-lg p-4 flex items-start gap-3">
                <Info className="h-5 w-5 text-blue-600 dark:text-blue-400 shrink-0 mt-0.5" />
                <div className="text-sm text-blue-800 dark:text-blue-300">
                    <p className="font-semibold mb-1">How Templates Work</p>
                    <p>
                        Templates allow you to create reusable prompts for the AI. You can define variables (Input Fields)
                        that you'll fill in when you use the template. These variables are then inserted into the System Instructions.
                    </p>
                </div>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle>AI Prompt</CardTitle>
                    <CardDescription>Define how the AI should behave.</CardDescription>
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
                        <div className="text-xs text-muted-foreground space-y-1">
                            <p>This prompt guides the AI. You can use placeholders for the input fields defined below.</p>
                            <p>Use the format <code className="bg-muted px-1 py-0.5 rounded">{'{{key}}'}</code> to insert a field value. For example, if you have a field with key "audience", use <code className="bg-muted px-1 py-0.5 rounded">{'{{audience}}'}</code> in your prompt.</p>
                        </div>
                    </div>
                </CardContent>
            </Card>

            <Card>
                <CardHeader className="flex flex-row items-center justify-between">
                    <div>
                        <CardTitle>Input Fields</CardTitle>
                        <CardDescription>Variables the user will fill in when using the template.</CardDescription>
                    </div>
                    <Button size="sm" variant="outline" onClick={addField}>
                        <Plus className="mr-2 h-4 w-4" /> Add Field
                    </Button>
                </CardHeader>
                <CardContent className="space-y-4">
                    {fields.length === 0 ? <p className="text-sm text-muted-foreground text-center py-4">No input fields defined. Click "Add Field" to create variables.</p> :
                        fields.map((field, idx) => (
                            <div key={idx} className="flex gap-2 items-start border p-3 rounded bg-muted/20">
                                <div className="grid gap-2 flex-1">
                                    <div className="flex gap-2">
                                        <div className="flex-1">
                                            <Label className="text-xs text-muted-foreground mb-1 block">Label</Label>
                                            <Input
                                                placeholder="e.g. Target Audience"
                                                value={field.label}
                                                onChange={e => updateField(idx, "label", e.target.value)}
                                            />
                                        </div>
                                        <div className="flex-1">
                                            <Label className="text-xs text-muted-foreground mb-1 block">Key (for use in prompt)</Label>
                                            <Input
                                                placeholder="e.g. audience"
                                                value={field.key}
                                                onChange={e => updateField(idx, "key", e.target.value)}
                                            />
                                        </div>
                                    </div>
                                    <div>
                                        <Label className="text-xs text-muted-foreground mb-1 block">Placeholder (Optional)</Label>
                                        <Input
                                            placeholder="Example value to guide the user..."
                                            value={field.placeholder || ""}
                                            onChange={e => updateField(idx, "placeholder", e.target.value)}
                                        />
                                    </div>
                                </div>
                                <Button variant="ghost" size="icon" className="text-destructive mt-6" onClick={() => removeField(idx)}>
                                    <Trash2 className="h-4 w-4" />
                                </Button>
                            </div>
                        ))
                    }
                </CardContent>
            </Card>

            <div className="flex justify-end gap-2">
                <Button variant="outline" onClick={() => navigate("/templates")}>Cancel</Button>
                <Button onClick={handleSave} disabled={saving}>{saving ? "Saving..." : "Save Template"}</Button>
            </div>

            {blocker.state === "blocked" && (
                <AlertDialog open={true}>
                    <AlertDialogContent>
                        <AlertDialogHeader>
                            <AlertDialogTitle>Unsaved Changes</AlertDialogTitle>
                            <AlertDialogDescription>
                                You have unsaved changes. Do you want to save them before leaving?
                            </AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                            <AlertDialogCancel onClick={() => blocker.reset()}>Cancel</AlertDialogCancel>
                            <Button variant="destructive" onClick={() => blocker.proceed()}>
                                Discard Changes
                            </Button>
                            <AlertDialogAction onClick={async (e) => {
                                e.preventDefault();
                                const target = blocker.location;
                                const success = await saveTemplate(false);
                                if (success && target) {
                                    navigate(target);
                                }
                            }}>
                                Save & Leave
                            </AlertDialogAction>
                        </AlertDialogFooter>
                    </AlertDialogContent>
                </AlertDialog>
            )}
        </div>
    );
}
