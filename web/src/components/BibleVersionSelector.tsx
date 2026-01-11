import { useState, useEffect, useCallback } from "react";
import { Check, ChevronsUpDown } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
    Command,
    CommandEmpty,
    CommandGroup,
    CommandInput,
    CommandItem,
    CommandList,
} from "@/components/ui/command";
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover";
import { getBibleVersions, BibleVersion } from "@/services/api";
import { useDebounce } from "@/hooks/useDebounce";

interface BibleVersionSelectorProps {
    value: string;
    onChange: (value: string) => void;
    placeholder?: string;
    className?: string;
}

export function BibleVersionSelector({ value, onChange, placeholder = "Select version...", className }: BibleVersionSelectorProps) {
    const [open, setOpen] = useState(false);
    const [versions, setVersions] = useState<BibleVersion[]>([]);
    const [loading, setLoading] = useState(false);
    const [search, setSearch] = useState("");
    const debouncedSearch = useDebounce(search, 300);

    // Initial load and search logic
    const loadVersions = useCallback(async (q: string) => {
        setLoading(true);
        try {
            // If we have a lot of versions, we might want to paginate, but for a dropdown
            // usually fetching the top 20 matches is enough.
            const res = await getBibleVersions({ name: q, limit: 20 });
            setVersions(res.data || []);
        } catch (error) {
            console.error("Failed to load versions", error);
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadVersions(debouncedSearch);
    }, [debouncedSearch, loadVersions]);

    // Handle display value
    // If we have the version in the list, use its abbreviation/name.
    // If not (e.g. initial load where 'value' is set but list is empty/loading), show value.
    const selectedVersion = versions.find(v => v.id === value || v.abbreviation === value)
        || (value ? { id: value, abbreviation: value, name: value, language: "" } : null);

    return (
        <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
                <Button
                    variant="outline"
                    role="combobox"
                    aria-expanded={open}
                    className={cn("w-full justify-between", className)}
                >
                    {selectedVersion ? selectedVersion.abbreviation : placeholder}
                    <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                </Button>
            </PopoverTrigger>
            <PopoverContent className="w-[300px] p-0">
                <Command shouldFilter={false}>
                    <CommandInput
                        placeholder="Search versions..."
                        value={search}
                        onValueChange={setSearch}
                    />
                    <CommandList>
                        {loading && <div className="py-2 text-center text-sm text-muted-foreground">Loading...</div>}
                        {!loading && versions.length === 0 && <CommandEmpty>No versions found.</CommandEmpty>}
                        <CommandGroup>
                            {versions.map((version) => (
                                <CommandItem
                                    key={version.id}
                                    value={version.id}
                                    onSelect={() => {
                                        // The value passed to onSelect is usually lowercased by Command
                                        // So we use the version object
                                        onChange(version.id);
                                        setOpen(false);
                                    }}
                                >
                                    <Check
                                        className={cn(
                                            "mr-2 h-4 w-4",
                                            value === version.id ? "opacity-100" : "opacity-0"
                                        )}
                                    />
                                    <div className="flex flex-col">
                                        <span className="font-medium">{version.abbreviation}</span>
                                        <span className="text-xs text-muted-foreground">{version.name} - {version.language}</span>
                                    </div>
                                </CommandItem>
                            ))}
                        </CommandGroup>
                    </CommandList>
                </Command>
            </PopoverContent>
        </Popover>
    );
}
