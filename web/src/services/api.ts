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
      token = await auth.currentUser?.getIdToken();
  }

  return {
    "Content-Type": "application/json",
    "Authorization": `Bearer ${token}`,
  };
}

// --- Notes ---

export interface NoteFilter {
  search?: string;
  startDate?: Date;
  endDate?: Date;
  sortBy?: "updated_at" | "created_at" | "title";
  sortOrder?: "asc" | "desc";
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

export async function createNote(title: string, content: string | Record<string, unknown>) {
  const headers = await getHeaders();
  const res = await fetch(`${API_URL}/notes`, {
    method: "POST",
    headers,
    body: JSON.stringify({ title, content }),
  });
  if (!res.ok) throw new Error("Failed to create note");
  return res.json();
}

export async function getNote(id: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/notes/${id}`, { headers });
    if (!res.ok) throw new Error("Failed to load note");
    return res.json();
}

export async function updateNote(id: string, title: string, content: string | Record<string, unknown>) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/notes/${id}`, {
        method: "PUT",
        headers,
        body: JSON.stringify({ title, content })
    });
    if (!res.ok) throw new Error("Failed to update note");
}

export async function deleteNote(id: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/notes/${id}`, {
        method: "DELETE",
        headers,
    });
    if (!res.ok) throw new Error("Failed to delete note");
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

export async function chatWithAI(passage: string, themes: string[], prompt: string, version?: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/chat`, {
        method: "POST",
        headers,
        body: JSON.stringify({ passage, themes, prompt, version }),
    });
    if (!res.ok) throw new Error("Failed to chat with AI");
    return res.json();
}

export async function askAI(context: string, prompt: string, version?: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/ai/ask`, {
        method: "POST",
        headers,
        body: JSON.stringify({ context, prompt, version }),
    });
    if (!res.ok) throw new Error("Failed to ask AI");
    return res.json();
}

export interface BibleVersion {
    id: string;
    name: string;
    language: string;
    abbreviation: string;
    code?: string;
    version?: string;
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

export async function getConnections() {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/connections`, { headers });
    if (!res.ok) throw new Error("Failed to fetch connections");
    return res.json();
}

export async function sendConnectionRequest(receiverEmail: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/connections/request`, {
        method: "POST",
        headers,
        body: JSON.stringify({ receiver_email: receiverEmail }),
    });
    if (!res.ok) {
        const err = await res.text();
        throw new Error(err || "Failed to send connection request");
    }
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
    version: string;
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

export async function clonePack(id: string, title?: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/verse-packs/${id}/clone`, {
        method: "POST",
        headers,
        body: JSON.stringify({ title }),
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

export async function generateFromTemplate(id: string, inputs: Record<string, string>) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/templates/${id}/generate`, {
        method: "POST",
        headers,
        body: JSON.stringify({ inputs }),
    });
    if (!res.ok) throw new Error("Failed to generate content");
    return res.json();
}
