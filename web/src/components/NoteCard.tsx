import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Trash2, Share2, Sparkles } from "lucide-react";

export interface Note {
  id: string;
  title: string;
  updated_at: string;
  tags?: { id: string; name: string }[];
}

interface NoteCardProps {
  note: Note;
  onDelete: (note: Note) => void;
  onShare: (note: Note) => void;
  onAskAI: (note: Note) => void;
}

export function NoteCard({ note, onDelete, onShare, onAskAI }: NoteCardProps) {
  const navigate = useNavigate();

  return (
    <div
        onClick={() => navigate(`/notes/${note.id}`)}
        className="block group h-full cursor-pointer"
        role="link"
        tabIndex={0}
        onKeyDown={(e) => {
            if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                navigate(`/notes/${note.id}`);
            }
        }}
    >
      <div className="p-4 bg-card hover:bg-accent/50 text-card-foreground transition rounded-lg shadow border h-44 flex flex-col">
         {/* Main Content */}
         <div className="flex-1">
            <h3 className="font-semibold mb-2 line-clamp-2 leading-tight">{note.title || "Untitled Note"}</h3>
            <p className="text-muted-foreground text-xs mt-1">{new Date(note.updated_at).toLocaleDateString()}</p>
            {note.tags && note.tags.length > 0 && (
                <div className="flex flex-wrap gap-1 mt-2 overflow-hidden max-h-6">
                    {note.tags.map(tag => (
                        <span key={tag.id} className="inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-medium bg-secondary text-secondary-foreground">
                            {tag.name}
                        </span>
                    ))}
                </div>
            )}
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
                aria-label={`Ask AI about ${note.title ? note.title : "Untitled Note"}`}
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
                aria-label={`Share ${note.title ? note.title : "Untitled Note"}`}
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
                aria-label={`Delete ${note.title ? note.title : "Untitled Note"}`}
            >
                <Trash2 className="w-4 h-4" />
            </Button>
         </div>
      </div>
    </div>
  );
}
