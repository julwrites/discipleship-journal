import { createContext, useContext, useEffect, useState } from "react"

export type Theme = "default" | "serene" | "elegant"
export type Mode = "light" | "dark" | "system"

type ThemeProviderProps = {
  children: React.ReactNode
  defaultTheme?: Theme
  defaultMode?: Mode
  storageKey?: string
}

type ThemeProviderState = {
  theme: Theme
  mode: Mode
  setTheme: (theme: Theme) => void
  setMode: (mode: Mode) => void
}

const initialState: ThemeProviderState = {
  theme: "default",
  mode: "system",
  setTheme: () => null,
  setMode: () => null,
}

const ThemeProviderContext = createContext<ThemeProviderState>(initialState)

export function ThemeProvider({
  children,
  defaultTheme = "default",
  defaultMode = "system",
  storageKey = "ui-theme",
  ...props
}: ThemeProviderProps) {
  const [theme, setTheme] = useState<Theme>(
    () => (localStorage.getItem(`${storageKey}-theme`) as Theme) || defaultTheme
  )
  const [mode, setMode] = useState<Mode>(
    () => (localStorage.getItem(`${storageKey}-mode`) as Mode) || defaultMode
  )

  useEffect(() => {
    const root = window.document.documentElement

    // Remove old theme classes
    root.classList.remove("theme-default", "theme-serene", "theme-elegant")
    // Add new theme class
    root.classList.add(`theme-${theme}`)

    // Handle dark mode
    root.classList.remove("light", "dark")

    if (mode === "system") {
      const systemTheme = window.matchMedia("(prefers-color-scheme: dark)")
        .matches
        ? "dark"
        : "light"

      root.classList.add(systemTheme)
      return
    }

    root.classList.add(mode)
  }, [theme, mode])

  const value = {
    theme,
    mode,
    setTheme: (theme: Theme) => {
      localStorage.setItem(`${storageKey}-theme`, theme)
      setTheme(theme)
    },
    setMode: (mode: Mode) => {
      localStorage.setItem(`${storageKey}-mode`, mode)
      setMode(mode)
    },
  }

  return (
    <ThemeProviderContext.Provider {...props} value={value}>
      {children}
    </ThemeProviderContext.Provider>
  )
}

// eslint-disable-next-line react-refresh/only-export-components
export const useTheme = () => {
  const context = useContext(ThemeProviderContext)

  if (context === undefined)
    throw new Error("useTheme must be used within a ThemeProvider")

  return context
}
