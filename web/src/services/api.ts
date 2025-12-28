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

export async function createNote(title: string, content: Record<string, unknown>) {
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

export async function updateNote(id: string, title: string, content: Record<string, unknown>) {
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
    const res = await fetch(`${API_URL}/users/me`, {
        method: "PUT",
        headers,
        body: JSON.stringify(data),
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

export async function getBiblePassage(ref: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/bible/passage?ref=${encodeURIComponent(ref)}`, { headers });
    if (!res.ok) throw new Error("Failed to fetch passage");
    return res.json();
}

export async function chatWithAI(passage: string, themes: string[], prompt: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/chat`, {
        method: "POST",
        headers,
        body: JSON.stringify({ passage, themes, prompt }),
    });
    if (!res.ok) throw new Error("Failed to chat with AI");
    return res.json();
}

export async function askAI(context: string, prompt: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/ai/ask`, {
        method: "POST",
        headers,
        body: JSON.stringify({ context, prompt }),
    });
    if (!res.ok) throw new Error("Failed to ask AI");
    return res.json();
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
    if (!res.ok) throw new Error("Failed to fetch shared notes");
    return res.json();
}

export async function shareNote(groupId: string, noteId: string, comment: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/${groupId}/shares`, {
        method: "POST",
        headers,
        body: JSON.stringify({ note_id: noteId, comment }),
    });
    if (!res.ok) throw new Error("Failed to share note");
    return res.json();
}

export async function getSharedNote(groupId: string, shareId: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/${groupId}/shares/${shareId}`, { headers });
    if (!res.ok) throw new Error("Failed to fetch shared note details");
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
