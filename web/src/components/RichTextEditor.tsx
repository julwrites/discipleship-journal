import { useEditor, EditorContent, Editor } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import Link from '@tiptap/extension-link'
import Image from '@tiptap/extension-image'
import Placeholder from '@tiptap/extension-placeholder'
import { Markdown } from '@tiptap/markdown'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import {
    Bold, Italic, List, ListOrdered, Heading1, Heading2,
    Quote, Link as LinkIcon, Undo, Redo, Image as ImageIcon
} from 'lucide-react'
import { useEffect, useState } from 'react'

interface RichTextEditorProps {
    initialContent: string;
    onChange: (markdown: string) => void;
    editable?: boolean;
    placeholder?: string;
    onEditorReady?: (editor: Editor) => void;
}

export default function RichTextEditor({
    initialContent,
    onChange,
    editable = true,
    placeholder = "Write something...",
    onEditorReady
}: RichTextEditorProps) {
    const editor = useEditor({
        extensions: [
            StarterKit,
            Link.configure({
                openOnClick: false,
            }),
            Image,
            Placeholder.configure({
                placeholder,
            }),
            Markdown,
        ],
        content: initialContent,
        editable: editable,
        onUpdate: ({ editor }) => {
            const markdownOutput = (editor as Editor & { getMarkdown: () => string }).getMarkdown();
            onChange(markdownOutput);
        },
        onCreate: ({ editor }) => {
            if (onEditorReady) {
                onEditorReady(editor);
            }
        }
    });

    useEffect(() => {
        if (editor) {
            editor.setEditable(editable);
        }
    }, [editable, editor]);

    // Note: We intentionally do NOT use useEffect to sync `initialContent` to editor.
    // The content is only seeded initially. Updates should happen via the editor instance.

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

    return (
        <div className="flex flex-col h-full border rounded-lg overflow-hidden">
            {editable && (
                <div className="bg-muted border-b p-2 flex gap-1 flex-wrap">
                    <Button variant={editor.isActive('bold') ? "secondary" : "ghost"} size="sm" onClick={toggleBold} title="Bold">
                        <Bold className="w-4 h-4" />
                    </Button>
                    <Button variant={editor.isActive('italic') ? "secondary" : "ghost"} size="sm" onClick={toggleItalic} title="Italic">
                        <Italic className="w-4 h-4" />
                    </Button>
                    <div className="w-px h-6 bg-border mx-1" />
                    <Button variant={editor.isActive('heading', { level: 1 }) ? "secondary" : "ghost"} size="sm" onClick={toggleHeading1} title="Heading 1">
                        <Heading1 className="w-4 h-4" />
                    </Button>
                    <Button variant={editor.isActive('heading', { level: 2 }) ? "secondary" : "ghost"} size="sm" onClick={toggleHeading2} title="Heading 2">
                        <Heading2 className="w-4 h-4" />
                    </Button>
                    <div className="w-px h-6 bg-border mx-1" />
                    <Button variant={editor.isActive('bulletList') ? "secondary" : "ghost"} size="sm" onClick={toggleBulletList} title="Bullet List">
                        <List className="w-4 h-4" />
                    </Button>
                    <Button variant={editor.isActive('orderedList') ? "secondary" : "ghost"} size="sm" onClick={toggleOrderedList} title="Ordered List">
                        <ListOrdered className="w-4 h-4" />
                    </Button>
                    <Button variant={editor.isActive('blockquote') ? "secondary" : "ghost"} size="sm" onClick={toggleBlockquote} title="Quote">
                        <Quote className="w-4 h-4" />
                    </Button>
                    <div className="w-px h-6 bg-border mx-1" />

                    <LinkPopover editor={editor} />
                    <ImagePopover editor={editor} />

                    <div className="flex-1" />
                    <Button variant="ghost" size="sm" onClick={() => editor.chain().focus().undo().run()} disabled={!editor.can().undo()} title="Undo">
                        <Undo className="w-4 h-4" />
                    </Button>
                    <Button variant="ghost" size="sm" onClick={() => editor.chain().focus().redo().run()} disabled={!editor.can().redo()} title="Redo">
                        <Redo className="w-4 h-4" />
                    </Button>
                </div>
            )}
            <EditorContent editor={editor} className="flex-1 overflow-auto p-4 prose dark:prose-invert max-w-none focus:outline-none" />
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
                .ProseMirror img {
                    display: block;
                    max-width: 100%;
                    height: auto;
                    border-radius: 0.5rem;
                }
                .ProseMirror img.ProseMirror-selectednode {
                    outline: 2px solid var(--ring);
                    outline-offset: 2px;
                }
            `}</style>
        </div>
    )
}

function LinkPopover({ editor }: { editor: Editor }) {
    const [url, setUrl] = useState('')
    const [open, setOpen] = useState(false)

    // Handle opening/closing state
    const onOpenChange = (newOpen: boolean) => {
        if (newOpen) {
            setUrl(editor.getAttributes('link').href || '')
        }
        setOpen(newOpen)
    }

    const setLink = () => {
        if (url === '') {
            editor.chain().focus().extendMarkRange('link').unsetLink().run()
        } else {
            editor.chain().focus().extendMarkRange('link').setLink({ href: url }).run()
        }
        setOpen(false)
    }

    return (
        <Popover open={open} onOpenChange={onOpenChange}>
            <PopoverTrigger asChild>
                <Button variant={editor.isActive('link') ? "secondary" : "ghost"} size="sm" title="Link">
                    <LinkIcon className="w-4 h-4" />
                </Button>
            </PopoverTrigger>
            <PopoverContent className="w-80">
                <div className="grid gap-4">
                    <div className="space-y-2">
                        <h4 className="font-medium leading-none">Set Link</h4>
                        <p className="text-sm text-muted-foreground">
                            Enter the URL for the link.
                        </p>
                    </div>
                    <div className="grid gap-2">
                        <div className="grid grid-cols-3 items-center gap-4">
                            <Label htmlFor="url">URL</Label>
                            <Input
                                id="url"
                                value={url}
                                onChange={(e) => setUrl(e.target.value)}
                                className="col-span-2 h-8"
                                onKeyDown={(e) => {
                                    if (e.key === 'Enter') setLink()
                                }}
                            />
                        </div>
                        <Button onClick={setLink}>Save</Button>
                    </div>
                </div>
            </PopoverContent>
        </Popover>
    )
}

function ImagePopover({ editor }: { editor: Editor }) {
    const [url, setUrl] = useState('')
    const [open, setOpen] = useState(false)

    const addImage = () => {
        if (url) {
            editor.chain().focus().setImage({ src: url }).run()
        }
        setOpen(false)
        setUrl('')
    }

    return (
        <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
                <Button variant={editor.isActive('image') ? "secondary" : "ghost"} size="sm" title="Image">
                    <ImageIcon className="w-4 h-4" />
                </Button>
            </PopoverTrigger>
            <PopoverContent className="w-80">
                <div className="grid gap-4">
                    <div className="space-y-2">
                        <h4 className="font-medium leading-none">Add Image</h4>
                        <p className="text-sm text-muted-foreground">
                            Enter the URL of the image.
                        </p>
                    </div>
                    <div className="grid gap-2">
                        <div className="grid grid-cols-3 items-center gap-4">
                            <Label htmlFor="img-url">URL</Label>
                            <Input
                                id="img-url"
                                value={url}
                                onChange={(e) => setUrl(e.target.value)}
                                className="col-span-2 h-8"
                                onKeyDown={(e) => {
                                    if (e.key === 'Enter') addImage()
                                }}
                            />
                        </div>
                        <Button onClick={addImage}>Add Image</Button>
                    </div>
                </div>
            </PopoverContent>
        </Popover>
    )
}
