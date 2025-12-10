import { API_URL, getHeaders } from "./api";

export interface Group {
    id: string;
    name: string;
    description: string;
    role?: string;
    created_by?: string;
}

export interface GroupMember {
    user_id: string;
    display_name: string;
    email: string;
    role: string;
    joined_at: string;
}

export interface SharedNote {
    id: string; // shareId
    group_id: string;
    note_id: string;
    title: string;
    shared_by: string;
    shared_at: string;
    comment: string;
    content?: Record<string, unknown>; // populated in details
}

export async function getMyGroups(): Promise<Group[]> {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups`, { headers });
    if (!res.ok) throw new Error("Failed to fetch groups");
    return res.json();
}

export async function searchGroups(query: string): Promise<Group[]> {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/search?q=${encodeURIComponent(query)}`, { headers });
    if (!res.ok) throw new Error("Failed to search groups");
    return res.json();
}

export async function createGroup(data: { name: string; description: string }): Promise<{ id: string }> {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups`, {
        method: "POST",
        headers,
        body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error("Failed to create group");
    return res.json();
}

export async function joinGroup(id: string): Promise<void> {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/${id}/join`, {
        method: "POST",
        headers,
    });
    if (!res.ok) {
        if (res.status === 409) return; // Already joined, treat as success
        throw new Error("Failed to join group");
    }
}

export async function leaveGroup(id: string): Promise<void> {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/${id}/leave`, {
        method: "DELETE",
        headers,
    });
    if (!res.ok) throw new Error("Failed to leave group");
}

export async function getGroupMembers(groupId: string): Promise<GroupMember[]> {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/${groupId}/members`, { headers });
    if (!res.ok) throw new Error("Failed to fetch members");
    return res.json();
}

export async function addGroupMember(groupId: string, userId: string): Promise<void> {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/${groupId}/members`, {
        method: "POST",
        headers,
        body: JSON.stringify({ user_id: userId }),
    });
    if (!res.ok) throw new Error("Failed to add member");
}

export async function removeGroupMember(groupId: string, userId: string): Promise<void> {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/${groupId}/members/${userId}`, {
        method: "DELETE",
        headers,
    });
    if (!res.ok) throw new Error("Failed to remove member");
}

export async function getGroupShares(groupId: string): Promise<SharedNote[]> {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/${groupId}/shares`, { headers });
    if (!res.ok) throw new Error("Failed to fetch shared notes");
    return res.json();
}

export async function shareNoteToGroup(groupId: string, noteId: string, comment: string): Promise<void> {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/${groupId}/shares`, {
        method: "POST",
        headers,
        body: JSON.stringify({ note_id: noteId, comment }),
    });
    if (!res.ok) throw new Error("Failed to share note");
}

export async function getSharedNote(groupId: string, shareId: string): Promise<SharedNote> {
    const headers = await getHeaders();
    const res = await fetch(`${API_URL}/groups/${groupId}/shares/${shareId}`, { headers });
    if (!res.ok) throw new Error("Failed to fetch shared note details");
    return res.json();
}
