import { useState } from "react";
import { StudyTemplate, generateFromTemplate, createNote } from "@/services/api";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { useNavigate } from "react-router-dom";

interface UseTemplateDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    template: StudyTemplate;
}

export default function UseTemplateDialog({ open, onOpenChange, template }: UseTemplateDialogProps) {
    const navigate = useNavigate();
    const [inputs, setInputs] = useState<Record<string, string>>({});
    const [generating, setGenerating] = useState(false);

    const handleInputChange = (key: string, value: string) => {
        setInputs(prev => ({ ...prev, [key]: value }));
    };

    const handleGenerate = async () => {
        setGenerating(true);
        try {
            const res = await generateFromTemplate(template.id!, inputs);

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
                    {template.fields.length === 0 ? <p>No inputs required. Click Generate to proceed.</p> :
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
