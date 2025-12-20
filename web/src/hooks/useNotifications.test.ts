import { renderHook, waitFor } from "@testing-library/react";
import { useNotifications } from "./useNotifications";
import { vi, describe, it, expect, beforeEach } from "vitest";

// Mock dependencies
vi.mock("@/lib/firebase", () => ({
  messaging: {}, // Truthy value
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

  it("requests permission and registers token if granted", async () => {
    const { getToken } = await import("firebase/messaging");
    const { registerDevice } = await import("@/services/api");

    (global.Notification.requestPermission as any).mockResolvedValue("granted");
    (getToken as any).mockResolvedValue("mock-fcm-token");

    const { result } = renderHook(() => useNotifications());

    expect(global.Notification.requestPermission).toHaveBeenCalled();

    await waitFor(() => {
      expect(result.current.permission).toBe("granted");
    });

    expect(getToken).toHaveBeenCalled();
    expect(registerDevice).toHaveBeenCalledWith("mock-fcm-token");
  });

  it("does not register token if permission denied", async () => {
    const { getToken } = await import("firebase/messaging");
    const { registerDevice } = await import("@/services/api");

    (global.Notification.requestPermission as any).mockResolvedValue("denied");

    const { result } = renderHook(() => useNotifications());

    await waitFor(() => {
      expect(result.current.permission).toBe("denied");
    });

    expect(getToken).not.toHaveBeenCalled();
    expect(registerDevice).not.toHaveBeenCalled();
  });
});
