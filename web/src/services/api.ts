import { auth } from "@/lib/firebase";

// Ensure API_URL always ends with /api, but prevent double /api if VITE_API_URL already has it
const BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";
export const API_URL = BASE_URL.endsWith("/api") ? BASE_URL : `${BASE_URL}/api`;

export async function getHeaders() {
  const token = await auth.currentUser?.getIdToken();
  return {
    "Content-Type": "application/json",
    "Authorization": `Bearer ${token}`,
  };
}

export async function fetchNotes() {
  const headers = await getHeaders();
  const res = await fetch(`${API_URL}/notes`, { headers });
  if (!res.ok) throw new Error("Failed to fetch notes");
  return res.json();
}

export async function getNote(id: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/notes/${id}`, { headers });
    if (!res.ok) throw new Error("Failed to load note");
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

export async function updateNote(id: string, title: string, content: Record<string, unknown>) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/notes/${id}`, {
        method: "PUT",
        headers,
        body: JSON.stringify({ title, content })
    });
    if (!res.ok) throw new Error("Failed to update note");
}

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

export async function searchUsers(query: string) {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/users/search?q=${encodeURIComponent(query)}`, { headers });
    if (!res.ok) throw new Error("Failed to search users");
    return res.json();
}
