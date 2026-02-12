import { Input } from "@/components/ui/input";
import { Textarea, TextareaProps } from "@/components/ui/textarea";
import { parseBibleReference } from "@/utils/bibleReference";

interface BibleReferenceInputProps extends Omit<React.ComponentProps<"input">, 'onChange' | 'onBlur'> {
    value: string;
    onChange: (value: string) => void;
    multiline?: boolean;
    allowMultiple?: boolean;
    rows?: number;
    onBlur?: (e: React.FocusEvent<HTMLInputElement | HTMLTextAreaElement>) => void;
}

export function BibleReferenceInput({ value, onChange, multiline, allowMultiple, onBlur, ...props }: BibleReferenceInputProps) {
    const handleBlur = (e: React.FocusEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const val = e.target.value;
        if (!val.trim()) {
            onBlur?.(e);
            return;
        }

        let normalized = val;

        if (allowMultiple) {
            // Split by comma
            const parts = val.split(',').map(p => p.trim()).filter(Boolean);
            const corrected = parts.map(p => parseBibleReference(p) || p);
            normalized = corrected.join(', ');
        } else {
            const parsed = parseBibleReference(val);
            // Only update if we successfully parsed a reference
            if (parsed) normalized = parsed;
        }

        if (normalized !== val) {
            onChange(normalized);
        }

        onBlur?.(e);
    };

    if (multiline) {
        return (
            <Textarea
                value={value}
                onChange={(e) => onChange(e.target.value)}
                onBlur={handleBlur}
                {...(props as unknown as TextareaProps)}
            />
        );
    }

    return (
        <Input
            value={value}
            onChange={(e) => onChange(e.target.value)}
            onBlur={handleBlur as React.FocusEventHandler<HTMLInputElement>}
            {...props}
        />
    );
}
