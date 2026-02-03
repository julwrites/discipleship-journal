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
    if (isOpen && reference) {
      loadPassage();
    }
  }, [isOpen, reference]);

  const loadPassage = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await getBiblePassage(reference, defaultVersion);
      setPassageText(data.text || data.verse);
      setPassageRef(data.reference || reference);
      setVersion(data.version || defaultVersion);
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
          <DialogDescription>
             {version}
          </DialogDescription>
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
            <div className="prose dark:prose-invert text-sm leading-relaxed whitespace-pre-wrap">
              {passageText}
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
