import { useState, useEffect } from "react";
import { StudyTemplate, generateFromTemplate, createNote, syncUser } from "@/services/api";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { useNavigate } from "react-router-dom";
import { BibleReferenceInput } from "@/components/BibleReferenceInput";
import { BibleVersionSelector } from "@/components/BibleVersionSelector";

interface UseTemplateDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    template: StudyTemplate;
}

export default function UseTemplateDialog({ open, onOpenChange, template }: UseTemplateDialogProps) {
    const navigate = useNavigate();
    const [inputs, setInputs] = useState<Record<string, string>>({});
    const [generating, setGenerating] = useState(false);

    const [userPassages, setUserPassages] = useState("");
    const [userVersion, setUserVersion] = useState("ESV");

    useEffect(() => {
        if (open) {
            syncUser().then(u => {
                if (u.settings?.bible_version) {
                    setUserVersion(u.settings.bible_version);
                }
            }).catch(console.error);
            // Reset inputs
            setInputs({});
            setUserPassages("");
        }
    }, [open]);

    const handleInputChange = (key: string, value: string) => {
        setInputs(prev => ({ ...prev, [key]: value }));
    };

    const handleGenerate = async () => {
        setGenerating(true);
        try {
            const passagesList = template.allow_user_passages
                ? userPassages.split(',').map(s => s.trim()).filter(Boolean)
                : undefined;

            // Use specified version or fallback to user selection (if not enforced by template)
            const versionToUse = template.required_version || userVersion;

            const res = await generateFromTemplate(template.id!, inputs, passagesList, versionToUse);

            // Create a new note with the content
            const noteTitle = `${template.title} - ${new Date().toLocaleDateString()}`;
            const noteRes = await createNote(noteTitle, res.content);

            toast.success("Study generated!");
            onOpenChange(false);
            navigate(`/notes/${noteRes.id}`);
        } catch (error) {
            console.error(error);
            toast.error("Failed to generate study");
        } finally {
            setGenerating(false);
        }
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-lg">
                <DialogHeader>
                    <DialogTitle>Use Template: {template.title}</DialogTitle>
                </DialogHeader>
                <div className="space-y-4 py-4 max-h-[60vh] overflow-y-auto px-1">

                    {/* Bible References Section */}
                    {template.allow_user_passages && (
                         <div className="grid gap-2">
                            <Label>Bible Passages (comma separated)</Label>
                            <BibleReferenceInput
                                value={userPassages}
                                onChange={setUserPassages}
                                allowMultiple
                                placeholder="e.g. John 3:16, Psalm 23"
                            />
                        </div>
                    )}

                    {/* Version Selection */}
                    {!template.required_version && (
                        <div className="grid gap-2">
                            <Label>Bible Version</Label>
                            <BibleVersionSelector
                                value={userVersion}
                                onChange={setUserVersion}
                            />
                        </div>
                    )}
                    {template.required_version && (
                        <div className="text-sm text-muted-foreground">
                            Using required version: <strong>{template.required_version}</strong>
                        </div>
                    )}

                    <div className="border-t my-2" />

                    {template.fields.length === 0 ? <p className="text-sm text-muted-foreground">No additional inputs required.</p> :
                        template.fields.map((field, idx) => (
                            <div key={idx} className="grid gap-2">
                                <Label>{field.label}</Label>
                                {field.type === 'textarea' ? (
                                    <Textarea
                                        placeholder={field.placeholder}
                                        value={inputs[field.key] || ""}
                                        onChange={e => handleInputChange(field.key, e.target.value)}
                                    />
                                ) : (
                                    <Input
                                        placeholder={field.placeholder}
                                        value={inputs[field.key] || ""}
                                        onChange={e => handleInputChange(field.key, e.target.value)}
                                    />
                                )}
                            </div>
                        ))
                    }
                </div>
                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)}>Cancel</Button>
                    <Button onClick={handleGenerate} disabled={generating}>
                        {generating ? "Generating..." : "Generate Study"}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}
