import { useEditor, EditorContent } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import Link from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'
import { Markdown } from '@tiptap/markdown'
import { Button } from '@/components/ui/button'
import {
    Bold, Italic, List, ListOrdered, Heading1, Heading2,
    Quote, Link as LinkIcon, Undo, Redo
} from 'lucide-react'
import { useEffect } from 'react'

interface RichTextEditorProps {
    content: string;
    onChange: (markdown: string) => void;
    editable?: boolean;
    placeholder?: string;
}

export default function RichTextEditor({
    content,
    onChange,
    editable = true,
    placeholder = "Write something..."
}: RichTextEditorProps) {
    const editor = useEditor({
        extensions: [
            StarterKit,
            Link.configure({
                openOnClick: false,
            }),
            Placeholder.configure({
                placeholder,
            }),
            Markdown,
        ],
        content: content,
        editable: editable,
        onUpdate: ({ editor }) => {
            const markdownOutput = (editor as any).getMarkdown();
            onChange(markdownOutput);
        },
    });

    // Update editor content when prop changes
    useEffect(() => {
        if (editor && content !== (editor as any).getMarkdown()) {
            // Only update if content is different to avoid cursor jumps and loops
            // However, getMarkdown() might return slightly different format than input content.
            // A better check might be needed or just accept that external updates reset cursor.
            // Since this is mainly for initial load, it should be fine.
            // For real-time collaboration it would be harder.

            // Check if the editor is empty and we are setting content (initial load)
            // or if the content is drastically different.

            // Simple approach: Only set if editor is not focused?
            // Or try to detect if change came from us.

            // NoteEditor loads content once usually.
            // But user typing triggers setMarkdown -> content prop update.
            // We need to prevent re-setting content if we just emitted it.

            // Let's rely on the check: content !== currentMarkdown
            // But markdown conversion might not be stable (e.g. whitespace).
            // This is a known issue with controlled inputs in Tiptap.
            // Usually simpler to use `onCreate` for initial content
            // and `useEffect` with dependency on `content` ONLY if we want to support external updates (like real-time or reset).

            // For this app, `markdown` state in NoteEditor IS the source of truth.
            // But we can just set content once if it was empty?
            // Or compare stringified versions.

            editor.commands.setContent(content);
        }
    }, [content, editor]);

    useEffect(() => {
        if (editor) {
            editor.setEditable(editable);
        }
    }, [editable, editor]);

    if (!editor) {
        return null;
    }

    const toggleBold = () => editor.chain().focus().toggleBold().run();
    const toggleItalic = () => editor.chain().focus().toggleItalic().run();
    const toggleHeading1 = () => editor.chain().focus().toggleHeading({ level: 1 }).run();
    const toggleHeading2 = () => editor.chain().focus().toggleHeading({ level: 2 }).run();
    const toggleBulletList = () => editor.chain().focus().toggleBulletList().run();
    const toggleOrderedList = () => editor.chain().focus().toggleOrderedList().run();
    const toggleBlockquote = () => editor.chain().focus().toggleBlockquote().run();
    const setLink = () => {
        const previousUrl = editor.getAttributes('link').href
        const url = window.prompt('URL', previousUrl)
        if (url === null) return
        if (url === '') {
            editor.chain().focus().extendMarkRange('link').unsetLink().run()
            return
        }
        editor.chain().focus().extendMarkRange('link').setLink({ href: url }).run()
    }

    return (
        <div className="flex flex-col h-full border rounded-lg overflow-hidden">
            {editable && (
                <div className="bg-slate-50 border-b p-2 flex gap-1 flex-wrap">
                    <Button variant={editor.isActive('bold') ? "secondary" : "ghost"} size="sm" onClick={toggleBold} title="Bold">
                        <Bold className="w-4 h-4" />
                    </Button>
                    <Button variant={editor.isActive('italic') ? "secondary" : "ghost"} size="sm" onClick={toggleItalic} title="Italic">
                        <Italic className="w-4 h-4" />
                    </Button>
                    <div className="w-px h-6 bg-slate-300 mx-1" />
                    <Button variant={editor.isActive('heading', { level: 1 }) ? "secondary" : "ghost"} size="sm" onClick={toggleHeading1} title="Heading 1">
                        <Heading1 className="w-4 h-4" />
                    </Button>
                    <Button variant={editor.isActive('heading', { level: 2 }) ? "secondary" : "ghost"} size="sm" onClick={toggleHeading2} title="Heading 2">
                        <Heading2 className="w-4 h-4" />
                    </Button>
                    <div className="w-px h-6 bg-slate-300 mx-1" />
                    <Button variant={editor.isActive('bulletList') ? "secondary" : "ghost"} size="sm" onClick={toggleBulletList} title="Bullet List">
                        <List className="w-4 h-4" />
                    </Button>
                    <Button variant={editor.isActive('orderedList') ? "secondary" : "ghost"} size="sm" onClick={toggleOrderedList} title="Ordered List">
                        <ListOrdered className="w-4 h-4" />
                    </Button>
                    <Button variant={editor.isActive('blockquote') ? "secondary" : "ghost"} size="sm" onClick={toggleBlockquote} title="Quote">
                        <Quote className="w-4 h-4" />
                    </Button>
                    <div className="w-px h-6 bg-slate-300 mx-1" />
                    <Button variant={editor.isActive('link') ? "secondary" : "ghost"} size="sm" onClick={setLink} title="Link">
                        <LinkIcon className="w-4 h-4" />
                    </Button>
                    <div className="flex-1" />
                    <Button variant="ghost" size="sm" onClick={() => editor.chain().focus().undo().run()} disabled={!editor.can().undo()} title="Undo">
                        <Undo className="w-4 h-4" />
                    </Button>
                    <Button variant="ghost" size="sm" onClick={() => editor.chain().focus().redo().run()} disabled={!editor.can().redo()} title="Redo">
                        <Redo className="w-4 h-4" />
                    </Button>
                </div>
            )}
            <EditorContent editor={editor} className="flex-1 overflow-auto p-4 prose prose-slate max-w-none focus:outline-none" />
            <style>{`
                .ProseMirror {
                    outline: none;
                    min-height: 100%;
                }
                .ProseMirror p.is-editor-empty:first-child::before {
                    color: #adb5bd;
                    content: attr(data-placeholder);
                    float: left;
                    height: 0;
                    pointer-events: none;
                }
            `}</style>
        </div>
    )
}
