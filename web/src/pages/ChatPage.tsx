import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { chatWithAIStream, syncUser } from "@/services/api";
import { useNavigate } from "react-router-dom";
import { BibleVersionSelector } from "@/components/BibleVersionSelector";
import { BibleReferenceInput } from "@/components/BibleReferenceInput";

export default function ChatPage() {
    const navigate = useNavigate();
    const [passage, setPassage] = useState("");
    const [themes, setThemes] = useState("");
    const [prompt, setPrompt] = useState("");
    const [version, setVersion] = useState("ESV");
    const [loading, setLoading] = useState(false);
    const [response, setResponse] = useState<string | null>(null);

    useEffect(() => {
        syncUser().then(u => {
            if (u.settings?.bible_version) {
                setVersion(u.settings.bible_version);
            }
        }).catch(console.error);
    }, []);

    const handleChat = async () => {
        setLoading(true);
        setResponse("");
        try {
            const themeList = themes.split(",").map(t => t.trim()).filter(Boolean);
            await chatWithAIStream(passage, themeList, prompt, version, {
                onStart: () => setResponse(""),
                onChunk: (chunk) => setResponse(prev => (prev || "") + chunk),
                onDone: () => setLoading(false),
                onError: (e) => {
                    console.error(e);
                    alert("Chat failed: " + e.message);
                    setLoading(false);
                }
            });
        } catch (e) {
            console.error(e);
            alert("Chat failed");
            setLoading(false);
        }
    };

    return (
        <div className="p-8 max-w-2xl mx-auto space-y-6">
            <h1 className="text-2xl font-bold">Chat with Bible AI</h1>

            <div className="space-y-2">
                <label className="text-sm font-medium">Bible Version</label>
                <div className="max-w-[200px]">
                    <BibleVersionSelector
                        value={version}
                        onChange={setVersion}
                    />
                </div>
            </div>

            <div className="space-y-2">
                <label className="text-sm font-medium">Bible Passage(s)</label>
                <BibleReferenceInput
                    multiline
                    allowMultiple
                    placeholder="e.g. Romans 8, Psalm 23"
                    value={passage}
                    onChange={setPassage}
                />
            </div>

            <div className="space-y-2">
                <label className="text-sm font-medium">Themes / Keywords (Optional)</label>
                <Input
                    placeholder="e.g. grace, suffering, hope"
                    value={themes}
                    onChange={(e) => setThemes(e.target.value)}
                />
            </div>

            <div className="space-y-2">
                <label className="text-sm font-medium">Your Question / Prompt</label>
                <Textarea
                    placeholder="What does this say about..."
                    value={prompt}
                    onChange={(e) => setPrompt(e.target.value)}
                />
            </div>

            <Button onClick={handleChat} disabled={loading || !passage || !prompt} className="w-full">
                {loading ? "Thinking..." : "Ask AI"}
            </Button>

            {response && (
                <div className="p-4 bg-muted rounded-lg border mt-6">
                    <h3 className="font-semibold mb-2">AI Response:</h3>
                    <div className="prose dark:prose-invert max-w-none" dangerouslySetInnerHTML={{ __html: response }} />
                    <div className="mt-4 text-sm text-muted-foreground">
                        A journal note has been created with this conversation.
                        <Button variant="link" className="p-0 h-auto ml-1" onClick={() => navigate("/")}>Go to Dashboard</Button>
                    </div>
                </div>
            )}
        </div>
    );
}
