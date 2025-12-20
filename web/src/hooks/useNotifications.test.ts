/* eslint-disable @typescript-eslint/no-explicit-any */
import { renderHook, waitFor, act } from "@testing-library/react";
import { useNotifications } from "./useNotifications";
import { vi, describe, it, expect, beforeEach } from "vitest";

// Mock dependencies
vi.mock("@/lib/firebase", () => ({
  messaging: {}, // Truthy value for supported env
  auth: { currentUser: { uid: "test-user" } }
}));

vi.mock("./useAuth", () => ({
  useAuth: () => ({ user: { uid: "test-user" }, loading: false })
}));

vi.mock("firebase/messaging", () => ({
  getToken: vi.fn(),
  onMessage: vi.fn(() => vi.fn()) // returns unsubscribe
}));

vi.mock("@/services/api", () => ({
  registerDevice: vi.fn()
}));

// Mock sonner
vi.mock("sonner", () => ({
  toast: vi.fn()
}));

describe("useNotifications", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Mock global Notification
    global.Notification = {
      requestPermission: vi.fn(),
      permission: "default"
    } as any;

    // Mock navigator.serviceWorker
    Object.defineProperty(navigator, 'serviceWorker', {
        value: {
            ready: Promise.resolve({})
        },
        configurable: true,
        writable: true
    });
  });

  it("registers token immediately if permission is already granted", async () => {
    const { getToken } = await import("firebase/messaging");
    const { registerDevice } = await import("@/services/api");

    // Simulate already granted
    (global.Notification as any).permission = "granted";
    (getToken as any).mockResolvedValue("mock-fcm-token");

    const { result } = renderHook(() => useNotifications());

    await waitFor(() => {
      expect(getToken).toHaveBeenCalled();
    });

    expect(registerDevice).toHaveBeenCalledWith("mock-fcm-token");
    expect(result.current.permission).toBe("granted");
  });

  it("does not request permission or register automatically if permission is default", async () => {
    const { getToken } = await import("firebase/messaging");

    (global.Notification as any).permission = "default";

    renderHook(() => useNotifications());

    expect(global.Notification.requestPermission).not.toHaveBeenCalled();
    expect(getToken).not.toHaveBeenCalled();
  });

  it("requests permission and registers when triggered manually", async () => {
    const { getToken } = await import("firebase/messaging");
    const { registerDevice } = await import("@/services/api");

    (global.Notification.requestPermission as any).mockResolvedValue("granted");
    (getToken as any).mockResolvedValue("mock-fcm-token");

    const { result } = renderHook(() => useNotifications());

    await act(async () => {
      await result.current.requestPermission();
    });

    expect(global.Notification.requestPermission).toHaveBeenCalled();
    expect(getToken).toHaveBeenCalled();
    expect(registerDevice).toHaveBeenCalledWith("mock-fcm-token");
  });
});
