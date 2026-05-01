import { useEffect, useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { getBiblePassage } from "@/services/api";
import { Button } from "@/components/ui/button";
import { BibleVersionSelector } from "@/components/BibleVersionSelector";

interface BiblePassageDialogProps {
  reference: string;
  isOpen: boolean;
  onClose: () => void;
  defaultVersion?: string;
}

export function BiblePassageDialog({ reference, isOpen, onClose, defaultVersion = "ESV" }: BiblePassageDialogProps) {
  const [loading, setLoading] = useState(false);
  const [passageText, setPassageText] = useState("");
  const [passageRef, setPassageRef] = useState("");
  const [version, setVersion] = useState(defaultVersion);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (isOpen) {
        loadPassage();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isOpen, reference, version]);

  const loadPassage = async () => {
    if (!reference) return;
    setLoading(true);
    setError(null);
    try {
      const data = await getBiblePassage(reference, version);
      setPassageText(data.text || data.verse);
      setPassageRef(data.reference || reference);
    } catch (err) {
      console.error(err);
      setError("Failed to load bible passage.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="max-w-md max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{passageRef || reference}</DialogTitle>
          <DialogDescription className="sr-only">
             Read Bible passage {reference} in {version}
          </DialogDescription>
          <div className="mt-2 flex items-center justify-between gap-2 text-sm text-muted-foreground">
             <span>Reading from:</span>
             <BibleVersionSelector
                value={version}
                onChange={setVersion}
                className="w-[120px] h-8"
             />
          </div>
        </DialogHeader>

        <div className="mt-4">
          {loading ? (
            <div className="flex justify-center py-8 text-muted-foreground">
               Loading...
            </div>
          ) : error ? (
            <div className="text-destructive text-center py-4">
              {error}
              <div className="mt-2">
                <Button variant="outline" size="sm" onClick={loadPassage}>Retry</Button>
              </div>
            </div>
          ) : (
            <div className="prose dark:prose-invert text-sm leading-relaxed">
              <div dangerouslySetInnerHTML={{ __html: passageText }} />
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
