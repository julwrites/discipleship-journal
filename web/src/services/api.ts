import { auth } from "@/lib/firebase";

const API_URL = import.meta.env.GCP_API_URL || "http://localhost:8080/api";

async function getHeaders() {
  const shouldBypassAuth = import.meta.env.FIREBASE_API_KEY === 'mock-key';
  let token = null;

  if (shouldBypassAuth) {
      const mockUserJson = localStorage.getItem('E2E_TEST_USER');
      if (mockUserJson) {
          token = "mock-token";
      }
  } else {
      token = await auth?.currentUser?.getIdToken();
  }

  return {
    "Content-Type": "application/json",
    "Authorization": `Bearer ${token}`,
  };
}

// --- Notes ---

export interface NoteFilter {
  search?: string;
  tag?: string;
  startDate?: Date;
  endDate?: Date;
  sortBy?: "updated_at" | "created_at" | "title";
  sortOrder?: "asc" | "desc";
}

export interface Tag {
    id: string;
    name: string;
    created_at?: string;
}

export async function fetchNotes(page = 1, limit = 20, filter: NoteFilter | string = {}) {
  const headers = await getHeaders();
  const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
  });

  if (typeof filter === "string") {
      if (filter) params.append("q", filter);
  } else {
      if (filter.search) params.append("q", filter.search);
      if (filter.tag) params.append("tag", filter.tag);
      if (filter.startDate) params.append("startDate", filter.startDate.toISOString());

      // Fix: Adjust endDate to be the end of the day (23:59:59.999)
      if (filter.endDate) {
          const eod = new Date(filter.endDate);
          eod.setHours(23, 59, 59, 999);
          params.append("endDate", eod.toISOString());
      }

      if (filter.sortBy) params.append("sortBy", filter.sortBy);
      if (filter.sortOrder) params.append("sortOrder", filter.sortOrder);
  }

  const res = await fetch(`${API_URL}/notes?${params.toString()}`, { headers });
  if (!res.ok) throw new Error("Failed to fetch notes");
  return res.json();
}

export async function createNote(title: string, content: string | Record<string, unknown>, tags: string[] = []) {
  const headers = await getHeaders();
  const res = await fetch(`${API_URL}/notes`, {
    method: "POST",
    headers,
    body: JSON.stringify({ title, content, tags }),
  });
  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    if (data.errors) {
      const messages = Object.values(data.errors).join(", ");
      throw new Error(`Validation failed: ${messages}`);
    }
    throw new Error("Failed to create note");
  }
  return res.json();
}

export async function getNote(id: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/notes/${id}`, { headers });
    if (!res.ok) throw new Error("Failed to load note");
    return res.json();
}

export async function updateNote(id: string, title: string, content: string | Record<string, unknown>, tags: string[] = []) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/notes/${id}`, {
        method: "PUT",
        headers,
        body: JSON.stringify({ title, content, tags })
    });
    if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        if (data.errors) {
            const messages = Object.values(data.errors).join(", ");
            throw new Error(`Validation failed: ${messages}`);
        }
        throw new Error("Failed to update note");
    }
}

export async function deleteNote(id: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/notes/${id}`, {
        method: "DELETE",
        headers,
    });
    if (!res.ok) throw new Error("Failed to delete note");
}

export async function getTags() {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/tags`, { headers });
    if (!res.ok) throw new Error("Failed to fetch tags");
    return res.json();
}

export async function createTag(name: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/tags`, {
        method: "POST",
        headers,
        body: JSON.stringify({ name }),
    });
    if (!res.ok) throw new Error("Failed to create tag");
    return res.json();
}

export async function deleteTag(id: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/tags/${id}`, {
        method: "DELETE",
        headers,
    });
    if (!res.ok) throw new Error("Failed to delete tag");
}

// --- User ---

export async function syncUser() {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/users/me`, { method: "POST", headers });
    if (!res.ok) throw new Error("Failed to sync user");
    return res.json();
}

export async function updateUser(data: { username?: string; bible_version?: string }) {
    const headers = await getHeaders();
    
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const payload: any = {};
    if (data.username !== undefined) payload.username = data.username;
    
    // Nest bible_version under settings to match backend UpdateUserRequest structure
    if (data.bible_version !== undefined) {
        payload.settings = { bible_version: data.bible_version };
    }

    const res = await fetch(`${API_URL}/users/me`, {
        method: "PUT",
        headers,
        body: JSON.stringify(payload),
    });
    if (!res.ok) throw new Error("Failed to update user");
}

export async function searchUsers(query: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/users/search?q=${encodeURIComponent(query)}`, { headers });
    if (!res.ok) throw new Error("Failed to search users");
    return res.json();
}

// --- Bible & AI ---

export async function getBiblePassage(ref: string, version?: string) {
    const headers = await getHeaders();
    const params = new URLSearchParams({ ref });
    if (version) params.append("version", version);
    const res = await fetch(`${API_URL}/bible/passage?${params.toString()}`, { headers });
    if (!res.ok) throw new Error("Failed to fetch passage");
    return res.json();
}

// Stream Handlers
type StreamCallback = {
    onStart?: (noteId: string) => void;
    onChunk?: (data: string) => void;
    onDone?: () => void;
    onError?: (err: Error) => void;
};

async function handleStreamRequest(url: string, body: object, callbacks: StreamCallback) {
    try {
        const headers = await getHeaders();
        const res = await fetch(url, {
            method: "POST",
            headers,
            body: JSON.stringify(body),
        });

        if (!res.ok) {
            throw new Error(`Stream request failed: ${res.status} ${res.statusText}`);
        }

        if (!res.body) {
            throw new Error("ReadableStream not supported by browser");
        }

        const reader = res.body.getReader();
        const decoder = new TextDecoder();
        let buffer = "";

        while (true) {
            const { value, done } = await reader.read();
            if (done) break;

            buffer += decoder.decode(value, { stream: true });
            const lines = buffer.split("\n\n");
            // Keep the last partial line in the buffer
            buffer = lines.pop() || "";

            for (const line of lines) {
                if (line.startsWith("event: ")) {
                    const parts = line.split("\n");
                    const eventLine = parts.find(l => l.startsWith("event: "));
                    const dataLine = parts.find(l => l.startsWith("data: "));

                    if (eventLine) {
                        const event = eventLine.substring(7).trim();
                        const dataStr = dataLine ? dataLine.substring(6).trim() : "{}";

                        try {
                            const data = dataStr ? JSON.parse(dataStr) : {};

                            if (event === "start" && callbacks.onStart) {
                                callbacks.onStart(dataStr); // Note ID is sent as raw string/json
                            } else if (event === "chunk" && callbacks.onChunk) {
                                callbacks.onChunk(data.response || "");
                            } else if (event === "done" && callbacks.onDone) {
                                callbacks.onDone();
                            } else if (event === "error" && callbacks.onError) {
                                callbacks.onError(new Error(data.error || "Unknown stream error"));
                            }
                        } catch (e) {
                             // If parse fails, it might be raw string for 'start' (which sends note ID directly)
                             // Actually backend sends `fmt.Fprintf(w, "event: start\ndata: %s\n\n", note.ID)`
                             // So dataStr IS the note ID, not JSON.
                             if (event === "start" && callbacks.onStart) {
                                 callbacks.onStart(dataStr);
                             } else {
                                 console.warn("Failed to parse SSE data", dataStr, e);
                             }
                        }
                    }
                }
            }
        }
    } catch (e) {
        if (callbacks.onError) callbacks.onError(e instanceof Error ? e : new Error(String(e)));
    }
}

export async function chatWithAI(passage: string, themes: string[], prompt: string, version?: string) {
    // Deprecated wrapper around stream, or keep old behavior?
    // Old behavior expects Promise<{response: string}>.
    // If backend now returns SSE, this will fail if we just await res.json().
    // We must update this to accumulate the stream and return full response for backward compatibility if we want.
    // But since I updated the backend to return SSE, calling res.json() will crash or be weird.
    // So I MUST update usages or wrap it.
    // I'll wrap it to wait for full stream.
    return new Promise((resolve, reject) => {
        let fullResponse = "";
        handleStreamRequest(`${API_URL}/chat`, { passage, themes, prompt, version }, {
            onChunk: (c) => fullResponse += c, // The chunk IS the full text in current implementation?
            // Actually backend sends `data, _ := json.Marshal(map[string]string{"response": answer})`.
            // And handleStreamRequest parses that JSON and passes `data.response`.
            // Backend sends full answer in one chunk currently.
            // But if it sent real chunks, we'd concat.
            // Since backend currently sends full answer in one chunk (see ChatHandler),
            // fullResponse will be the full text.
            onDone: () => resolve({ response: fullResponse }),
            onError: reject
        });
    });
}

export async function chatWithAIStream(passage: string, themes: string[], prompt: string, version: string | undefined, callbacks: StreamCallback) {
    return handleStreamRequest(`${API_URL}/chat`, { passage, themes, prompt, version }, callbacks);
}

export async function askAI(context: string, prompt: string, version?: string) {
     return new Promise((resolve, reject) => {
        let fullResponse = "";
        handleStreamRequest(`${API_URL}/ai/ask`, { context, prompt, version }, {
            onChunk: (c) => fullResponse = c, // Backend sends full text currently
            onDone: () => resolve({ response: fullResponse }),
            onError: reject
        });
    });
}

export async function askAIStream(context: string, prompt: string, version: string | undefined, callbacks: StreamCallback) {
    return handleStreamRequest(`${API_URL}/ai/ask`, { context, prompt, version }, callbacks);
}

export interface BibleVersion {
    id: string;
    name: string;
    language: string;
    abbreviation: string;
    code?: string;
    version?: string;
    value?: string;
}

export async function getBibleVersions(params: { name?: string; language?: string; page?: number; limit?: number; sort?: "code" | "name" | "language" } = {}) {
    const headers = await getHeaders();
    const query = new URLSearchParams();
    if (params.name) query.append("name", params.name);
    if (params.language) query.append("language", params.language);
    if (params.page) query.append("page", params.page.toString());
    if (params.limit) query.append("limit", params.limit.toString());
    if (params.sort) query.append("sort", params.sort);

    const res = await fetch(`${API_URL}/bible/versions?${query.toString()}`, { headers });
    if (!res.ok) throw new Error("Failed to fetch bible versions");
    const data = await res.json();
    // Support both { data: [...] } and raw array [...]
    if (Array.isArray(data)) return { data };
    return data;
}

// --- Groups ---

export async function getGroups() {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups`, { headers });
    if (!res.ok) throw new Error("Failed to fetch groups");
    return res.json();
}

export async function searchGroups(query: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/search?q=${encodeURIComponent(query)}`, { headers });
    if (!res.ok) throw new Error("Failed to search groups");
    return res.json();
}

export async function createGroup(data: { name: string; description: string }) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups`, {
        method: "POST",
        headers,
        body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error("Failed to create group");
    return res.json();
}

export async function joinGroup(id: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/${id}/join`, {
        method: "POST",
        headers,
    });
    if (!res.ok) throw new Error("Failed to join group");
    return res.json();
}

export async function leaveGroup(id: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/${id}/leave`, {
        method: "DELETE",
        headers,
    });
    if (!res.ok) throw new Error("Failed to leave group");
    return res.json();
}

export async function getGroupMembers(groupId: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/${groupId}/members`, { headers });
    if (!res.ok) throw new Error("Failed to fetch group members");
    return res.json();
}

export async function addGroupMember(groupId: string, userId: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/${groupId}/members`, {
        method: "POST",
        headers,
        body: JSON.stringify({ user_id: userId }),
    });
    if (!res.ok) throw new Error("Failed to add member");
    return res.json();
}

export async function removeGroupMember(groupId: string, userId: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/${groupId}/members/${userId}`, {
        method: "DELETE",
        headers,
    });
    if (!res.ok) throw new Error("Failed to remove member");
    return res.json();
}

// --- Group Shares ---

export async function getGroupShares(groupId: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/${groupId}/shares`, { headers });
    if (!res.ok) throw new Error("Failed to fetch shared items");
    return res.json();
}

export async function shareItem(groupId: string, data: { note_id?: string; verse_pack_id?: string; comment: string }) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/${groupId}/shares`, {
        method: "POST",
        headers,
        body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error("Failed to share item");
    return res.json();
}

export async function getSharedItem(groupId: string, shareId: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/${groupId}/shares/${shareId}`, { headers });
    if (!res.ok) throw new Error("Failed to fetch shared item details");
    return res.json();
}

// --- Connections ---

export interface Connection {
    id: string;
    requester_id: string;
    receiver_id: string;
    status: string;
    requester_email?: string;
    receiver_email?: string;
    requester_username?: string;
    receiver_username?: string;
}

export async function getConnections() {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/connections`, { headers });
    if (!res.ok) throw new Error("Failed to fetch connections");
    return res.json();
}

export async function sendConnectionRequest(receiverEmailOrId: string, isId: boolean = false) {
    const headers = await getHeaders();
    const payload = isId ? { receiver_id: receiverEmailOrId } : { receiver_email: receiverEmailOrId };
    const res = await fetch(`${API_URL}/connections/request`, {
        method: "POST",
        headers,
        body: JSON.stringify(payload),
    });
    if (!res.ok) {
        const err = await res.text();
        throw new Error(err || "Failed to send connection request");
    }
    return res.json();
}

export async function unsubscribeFromPlan(id: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/reading-plans/${id}/subscribe`, {
        method: "DELETE",
        headers,
    });
    if (!res.ok) {
        throw new Error("Failed to unsubscribe from plan");
    }
}

export async function getOrCreateDirectGroup(partnerId: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/direct`, {
        method: "POST",
        headers,
        body: JSON.stringify({ partner_id: partnerId }),
    });
    if (!res.ok) throw new Error("Failed to get/create direct group");
    return res.json();
}

export async function respondToConnectionRequest(id: string, action: "accept" | "reject") {
    const headers = await getHeaders();
    const method = action === "reject" ? "DELETE" : "PUT";
    const res = await fetch(`${API_URL}/connections/${id}`, {
        method: method,
        headers,
    });
    if (!res.ok) throw new Error("Failed to respond to connection request");
    return res.json();
}

// --- Notifications ---

export async function registerDevice(token: string, deviceType: string = "web") {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/notifications/register`, {
        method: "POST",
        headers,
        body: JSON.stringify({ token, device_type: deviceType }),
    });
    if (!res.ok) throw new Error("Failed to register device");
    return res.json();
}

// --- Reading Plans ---

export async function getReadingPlans() {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/reading-plans`, { headers });
    if (!res.ok) throw new Error("Failed to fetch reading plans");
    return res.json();
}

export async function getReadingPlan(id: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/reading-plans/${id}`, { headers });
    if (!res.ok) throw new Error("Failed to fetch reading plan details");
    return res.json();
}

export async function subscribeToPlan(id: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/reading-plans/${id}/subscribe`, {
        method: "POST",
        headers,
    });
    if (!res.ok) {
        if (res.status === 409) throw new Error("Already subscribed");
        throw new Error("Failed to subscribe to plan");
    }
    return res.json();
}

export async function getMyReadingPlans() {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/my-reading-plans`, { headers });
    if (!res.ok) throw new Error("Failed to fetch my reading plans");
    return res.json();
}

export async function markPlanDayComplete(planId: string, dayNumber: number) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/my-reading-plans/${planId}/progress`, {
        method: "POST",
        headers,
        body: JSON.stringify({ day_number: dayNumber }),
    });
    if (!res.ok) throw new Error("Failed to mark day as complete");
    return res.json();
}

export async function unmarkPlanDayComplete(planId: string, dayNumber: number) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/my-reading-plans/${planId}/progress/${dayNumber}`, {
        method: "DELETE",
        headers,
    });
    if (!res.ok) throw new Error("Failed to unmark day as complete");
    return res.json();
}

export async function getPlanProgress(planId: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/my-reading-plans/${planId}/progress`, { headers });
    if (!res.ok) throw new Error("Failed to fetch plan progress");
    return res.json();
}

// --- Memory Verses ---

export interface VersePack {
    id: string;
    title: string;
    identifier?: string;
    description?: string;
    is_public: boolean;
    verse_count: number;
    user_id?: string;
}

export interface MemoryVerse {
    id?: string;
    verse_pack_id: string;
    reference: string;
    title?: string;
    version: string;
    version_source?: 'override' | 'user_default' | 'original';
    tags: string[];
    pack_title?: string;
}

export async function getVersePacks(type: "system" | "user" = "user") {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/verse-packs?type=${type}`, { headers });
    if (!res.ok) throw new Error("Failed to get packs");
    return res.json();
}

export async function createVersePack(title: string, description?: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/verse-packs`, {
        method: "POST",
        headers,
        body: JSON.stringify({ title, description }),
    });
    if (!res.ok) throw new Error("Failed to create pack");
    return res.json();
}

export async function getPackDetails(id: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/verse-packs/${id}`, { headers });
    if (!res.ok) throw new Error("Failed to get pack details");
    return res.json();
}

export async function createVerseInPack(packId: string, verse: Omit<MemoryVerse, "verse_pack_id">) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/verse-packs/${packId}/verses`, {
        method: "POST",
        headers,
        body: JSON.stringify(verse),
    });
    if (!res.ok) throw new Error("Failed to create verse");
    return res.json();
}

export async function clonePack(id: string, title?: string, useUserDefault?: boolean) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/verse-packs/${id}/clone`, {
        method: "POST",
        headers,
        body: JSON.stringify({ title, use_user_default: useUserDefault }),
    });
    if (!res.ok) throw new Error("Failed to clone pack");
    return res.json();
}

export async function deletePack(id: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/verse-packs/${id}`, {
        method: "DELETE",
        headers,
    });
    if (!res.ok) throw new Error("Failed to delete pack");
    return res.json();
}

export async function updateMemoryVerse(id: string, verse: Partial<MemoryVerse>) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/memory-verses/${id}`, {
        method: "PUT",
        headers,
        body: JSON.stringify(verse),
    });
    if (!res.ok) throw new Error("Failed to update verse");
    return res.json();
}

export async function deleteMemoryVerse(id: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/memory-verses/${id}`, {
        method: "DELETE",
        headers,
    });
    if (!res.ok) throw new Error("Failed to delete verse");
    return res.json();
}

export async function setVersePreference(verseId: string, version: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/memory-verses/${verseId}/preference`, {
        method: "PUT",
        headers,
        body: JSON.stringify({ version }),
    });
    if (!res.ok) throw new Error("Failed to set verse preference");
    return res.json();
}

export async function removeVersePreference(verseId: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/memory-verses/${verseId}/preference`, {
        method: "DELETE",
        headers,
    });
    if (!res.ok) throw new Error("Failed to remove verse preference");
    return res.json();
}

// Deprecated: Compatibility aliases
export const shareNote = (groupId: string, noteId: string, comment: string) => shareItem(groupId, { note_id: noteId, comment });
export const searchMemoryVerses = async (query: string = "") => {
    // This requires implementing a search endpoint in MemoryVerseService that searches verses across packs
    // For now, if we don't have it, we might error or mock empty.
    // BUT the requirement is to use the new system.
    // I will point this to a search endpoint I just added to the service interface
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/memory-verses?q=${encodeURIComponent(query)}`, { headers });
    if (!res.ok) throw new Error("Failed to search verses");
    return res.json();
};


// --- Templates ---

export interface TemplateField {
    key: string;
    label: string;
    type: "text" | "textarea" | "select";
    placeholder?: string;
}

export interface StudyTemplate {
    id?: string;
    creator_id?: string;
    title: string;
    description: string;
    structure: Record<string, unknown>;
    prompts: Record<string, unknown>;
    fields: TemplateField[];
    is_public: boolean;
    bible_references?: string[];
    allow_user_passages?: boolean;
    template_body?: string;
    required_version?: string;
}

export async function getMyTemplates() {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/templates`, { headers });
    if (!res.ok) throw new Error("Failed to fetch templates");
    return res.json();
}

export async function getPublicTemplates() {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/templates/public`, { headers });
    if (!res.ok) throw new Error("Failed to fetch public templates");
    return res.json();
}

export async function getTemplate(id: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/templates/${id}`, { headers });
    if (!res.ok) throw new Error("Failed to fetch template");
    return res.json();
}

export async function createTemplate(template: StudyTemplate) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/templates`, {
        method: "POST",
        headers,
        body: JSON.stringify(template),
    });
    if (!res.ok) throw new Error("Failed to create template");
    return res.json();
}

export async function updateTemplate(id: string, template: StudyTemplate) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/templates/${id}`, {
        method: "PUT",
        headers,
        body: JSON.stringify(template),
    });
    if (!res.ok) throw new Error("Failed to update template");
    return res.json();
}

export async function deleteTemplate(id: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/templates/${id}`, {
        method: "DELETE",
        headers,
    });
    if (!res.ok) throw new Error("Failed to delete template");
    return res.json();
}

export async function cloneTemplate(id: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/templates/${id}/clone`, {
        method: "POST",
        headers,
    });
    if (!res.ok) throw new Error("Failed to clone template");
    return res.json();
}

export async function generateFromTemplate(id: string, inputs: Record<string, string>, user_passages?: string[], user_version?: string) {
    const headers = await getHeaders();
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const payload: any = { inputs };
    if (user_passages) payload.user_passages = user_passages;
    if (user_version) payload.user_version = user_version;

    const res = await fetch(`${API_URL}/templates/${id}/generate`, {
        method: "POST",
        headers,
        body: JSON.stringify(payload),
    });
    if (!res.ok) throw new Error("Failed to generate content");
    return res.json();
}
