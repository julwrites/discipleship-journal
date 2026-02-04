import { Note } from "@/components/NoteCard";

const CACHE_KEY_PREFIX = "dashboard_notes_cache_";
const CACHE_VERSION = "v1";
const STALE_THRESHOLD_MS = 1000 * 60 * 60 * 24; // 24 hours

interface NotesCacheData {
  version: string;
  timestamp: number;
  notes: Note[];
  hasMore: boolean;
}

function getCacheKey(userId: string) {
    return `${CACHE_KEY_PREFIX}${userId}`;
}

export function getCachedNotes(userId: string): { notes: Note[], hasMore: boolean } | null {
  if (!userId) return null;
  try {
    const key = getCacheKey(userId);
    const raw = localStorage.getItem(key);
    if (!raw) return null;

    const data: NotesCacheData = JSON.parse(raw);
    if (data.version !== CACHE_VERSION) {
       localStorage.removeItem(key);
       return null;
    }

    if (Date.now() - data.timestamp > STALE_THRESHOLD_MS) {
        localStorage.removeItem(key);
        return null;
    }

    return { notes: data.notes, hasMore: data.hasMore };
  } catch (e) {
    console.warn("Failed to read notes cache", e);
    return null;
  }
}

export function setCachedNotes(userId: string, notes: Note[], hasMore: boolean) {
  if (!userId) return;
  try {
    const key = getCacheKey(userId);
    const data: NotesCacheData = {
      version: CACHE_VERSION,
      timestamp: Date.now(),
      notes,
      hasMore
    };
    localStorage.setItem(key, JSON.stringify(data));
  } catch (e) {
    console.warn("Failed to write notes cache", e);
  }
}

export function clearNotesCache(userId: string) {
    if (!userId) return;
    localStorage.removeItem(getCacheKey(userId));
}
