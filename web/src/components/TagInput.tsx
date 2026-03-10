import * as React from "react"
import { Check, ChevronsUpDown, X } from "lucide-react"

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"

interface TagInputProps {
  value: string[]
  onChange: (value: string[]) => void
  suggestions: string[] // List of available tag names
}

export function TagInput({ value = [], onChange, suggestions = [] }: TagInputProps) {
  const [open, setOpen] = React.useState(false)
  const [inputValue, setInputValue] = React.useState("")

  const handleSelect = (tag: string) => {
    if (!value.includes(tag)) {
        onChange([...value, tag])
    }
    setOpen(false)
    setInputValue("")
  }

  const handleCreate = () => {
      if (inputValue && !value.includes(inputValue)) {
          onChange([...value, inputValue])
          setOpen(false)
          setInputValue("")
      }
  }

  const handleRemove = (tagToRemove: string) => {
    onChange(value.filter((t) => t !== tagToRemove))
  }

  // Filter out tags that are already selected to avoid duplicates in the list
  // although showing them with checkmark is also fine.
  // Let's show all suggestions, but mark selected.

  return (
    <div className="flex flex-col gap-2">
        <div className="flex flex-wrap gap-2">
            {value.map((tag) => (
            <span
                key={tag}
                className="inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80"
            >
                {tag}
                <button
                    type="button"
                    className="ml-1 ring-offset-background rounded-full outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                    onKeyDown={(e) => {
                        if (e.key === "Enter") {
                            handleRemove(tag)
                        }
                    }}
                    onMouseDown={(e) => {
                        e.preventDefault()
                        e.stopPropagation()
                    }}
                    onClick={() => handleRemove(tag)}
                >
                    <X className="h-3 w-3 text-muted-foreground hover:text-foreground" />
                    <span className="sr-only">Remove {tag}</span>
                </button>
            </span>
            ))}
        </div>

      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            variant="outline"
            role="combobox"
            aria-expanded={open}
            className="w-full justify-between"
          >
            Add tag...
            <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-[200px] p-0" align="start">
          <Command>
            <CommandInput
                placeholder="Search tag..."
                value={inputValue}
                onValueChange={setInputValue}
                onKeyDown={(e) => {
                    if (e.key === 'Enter' && inputValue) {
                        e.preventDefault();
                        handleCreate();
                    }
                }}
            />
            <CommandList>
                <CommandEmpty>
                    {inputValue && (
                        <div className="p-2 text-sm cursor-pointer hover:bg-accent hover:text-accent-foreground" onClick={handleCreate}>
                            Create "{inputValue}"
                        </div>
                    )}
                    {!inputValue && "No tags found."}
                </CommandEmpty>
                <CommandGroup>
                {suggestions.map((suggestion) => (
                    <CommandItem
                        key={suggestion}
                        value={suggestion}
                        onSelect={() => handleSelect(suggestion)}
                    >
                    <Check
                        className={cn(
                        "mr-2 h-4 w-4",
                        value.includes(suggestion) ? "opacity-100" : "opacity-0"
                        )}
                    />
                    {suggestion}
                    </CommandItem>
                ))}
                </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
    </div>
  )
}
