import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Trash2, Share2, Sparkles } from "lucide-react";

export interface Note {
  id: string;
  title: string;
  updated_at: string;
}

interface NoteCardProps {
  note: Note;
  onDelete: (note: Note) => void;
  onShare: (note: Note) => void;
  onAskAI: (note: Note) => void;
}

export function NoteCard({ note, onDelete, onShare, onAskAI }: NoteCardProps) {
  return (
    <Link to={`/notes/${note.id}`} className="block group h-full">
      <div className="p-4 bg-card hover:bg-accent/50 text-card-foreground transition rounded-lg shadow border h-44 flex flex-col">
         {/* Main Content */}
         <div className="flex-1">
            <h3 className="font-semibold mb-2 line-clamp-2 leading-tight">{note.title || "Untitled Note"}</h3>
            <p className="text-muted-foreground text-xs mt-1">{new Date(note.updated_at).toLocaleDateString()}</p>
         </div>

         {/* Actions Footer */}
         <div className="flex justify-end gap-1 mt-auto pt-2">
            <Button
                variant="ghost"
                size="sm"
                className="h-8 w-8 p-0 text-muted-foreground hover:text-primary hover:bg-background/80"
                onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    onAskAI(note);
                }}
                title="Ask AI"
            >
                <Sparkles className="w-4 h-4" />
            </Button>
            <Button
                variant="ghost"
                size="sm"
                className="h-8 w-8 p-0 text-muted-foreground hover:text-primary hover:bg-background/80"
                onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    onShare(note);
                }}
                title="Share"
            >
                <Share2 className="w-4 h-4" />
            </Button>
             <Button
                variant="ghost"
                size="sm"
                className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive hover:bg-background/80"
                onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    onDelete(note);
                }}
                title="Delete"
            >
                <Trash2 className="w-4 h-4" />
            </Button>
         </div>
      </div>
    </Link>
  );
}
