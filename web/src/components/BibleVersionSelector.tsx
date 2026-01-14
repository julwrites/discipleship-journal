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
    // Support various API response formats (abbreviation, code, id, version)
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const getVersionCode = (v: any) => v.abbreviation || v.code || v.id || v.version || v.name;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const getUniqueId = (v: any) => v.id || v.abbreviation || v.code || v.version || v.name;

    // If we have the version in the list, use its code.
    // If not (e.g. initial load where 'value' is set but list is empty/loading), show value.
    const selectedVersion = versions.find(v => {
        const code = getVersionCode(v);
        return code && value && code.toLowerCase() === value.toLowerCase();
    }) || (value ? { id: value, abbreviation: value, name: value, language: "" } : null);

    const displayLabel = selectedVersion ? getVersionCode(selectedVersion) : placeholder;

    return (
        <Popover open={open} onOpenChange={(isOpen) => {
            setOpen(isOpen);
            if (!isOpen) setSearch(""); // Reset search on close
        }}>
            <PopoverTrigger asChild>
                <Button
                    variant="outline"
                    role="combobox"
                    aria-expanded={open}
                    className={cn("w-full justify-between", className)}
                >
                    {displayLabel}
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
                            {versions.map((version) => {
                                const code = getVersionCode(version);
                                const uniqueId = getUniqueId(version);
                                
                                if (!code || !uniqueId) return null;

                                const isSelected = value && code.toLowerCase() === value.toLowerCase();
                                
                                return (
                                    <CommandItem
                                        key={uniqueId}
                                        value={uniqueId.toLowerCase()} 
                                        onSelect={() => {
                                            console.log("BibleVersionSelector: onSelect", { code, uniqueId, version });
                                            if (code) {
                                                onChange(code);
                                            }
                                            setOpen(false);
                                        }}
                                    >
                                        {isSelected ? (
                                            <Check className={cn("mr-2 h-4 w-4")} />
                                        ) : (
                                            <div className="mr-2 h-4 w-4" />
                                        )}
                                        <div className="flex flex-col">
                                            <span className="font-medium">{code}</span>
                                            <span className="text-xs text-muted-foreground">{version.name} - {version.language}</span>
                                        </div>
                                    </CommandItem>
                                );
                            })}
                        </CommandGroup>
                    </CommandList>
                </Command>
            </PopoverContent>
        </Popover>
    );
}