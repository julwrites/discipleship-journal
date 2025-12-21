import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import ReadingPlansPage from "./ReadingPlansPage";
import { BrowserRouter } from "react-router-dom";
import { describe, it, expect, vi, beforeEach } from "vitest";
import * as api from "@/services/api";

// Mock API
vi.mock("@/services/api", () => ({
    getReadingPlans: vi.fn(),
    getMyReadingPlans: vi.fn(),
    subscribeToPlan: vi.fn(),
}));

// Mock Auth
vi.mock("@/hooks/useAuth", () => ({
    useAuth: () => ({ user: { uid: "test-user" } }),
}));

// Mock sonner
vi.mock("sonner", () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
    },
}));

// Mock matchMedia for Tabs (Radix UI)
beforeEach(() => {
    Object.defineProperty(window, 'matchMedia', {
        writable: true,
        value: vi.fn().mockImplementation(query => ({
            matches: false,
            media: query,
            onchange: null,
            addListener: vi.fn(), // deprecated
            removeListener: vi.fn(), // deprecated
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
            dispatchEvent: vi.fn(),
        })),
    });
});

describe("ReadingPlansPage", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders loading state", () => {
        (api.getReadingPlans as any).mockReturnValue(new Promise(() => {}));
        (api.getMyReadingPlans as any).mockReturnValue(new Promise(() => {}));

        render(
            <BrowserRouter>
                <ReadingPlansPage />
            </BrowserRouter>
        );

        expect(screen.getByText("Loading...")).toBeInTheDocument();
    });

    it("renders empty state for my plans", async () => {
        (api.getReadingPlans as any).mockResolvedValue({ data: [] });
        (api.getMyReadingPlans as any).mockResolvedValue({ data: [] });

        render(
            <BrowserRouter>
                <ReadingPlansPage />
            </BrowserRouter>
        );

        await waitFor(() => {
            expect(screen.queryByText("Loading...")).not.toBeInTheDocument();
        });

        expect(screen.getByText("No active plans")).toBeInTheDocument();
    });

    it("renders available plans", async () => {
        const mockPlans = [
            { id: "1", title: "Test Plan", description: "Test Desc", days: 30 },
        ];
        (api.getReadingPlans as any).mockResolvedValue({ data: mockPlans });
        (api.getMyReadingPlans as any).mockResolvedValue({ data: [] });

        render(
            <BrowserRouter>
                <ReadingPlansPage />
            </BrowserRouter>
        );

        await waitFor(() => {
            expect(screen.queryByText("Loading...")).not.toBeInTheDocument();
        });

        // Click Browse Plans tab
        const user = userEvent.setup();
        const browseTab = screen.getByRole("tab", { name: "Browse Plans" });
        await user.click(browseTab);

        await waitFor(() => {
            expect(screen.getByText("Test Plan")).toBeInTheDocument();
        });
        expect(screen.getByText("Test Desc")).toBeInTheDocument();
        expect(screen.getByText("Subscribe")).toBeInTheDocument();
    });

    it("handles subscription", async () => {
        const user = userEvent.setup();
        const mockPlans = [
            { id: "1", title: "Test Plan", description: "Test Desc", days: 30 },
        ];
        (api.getReadingPlans as any).mockResolvedValue({ data: mockPlans });
        (api.getMyReadingPlans as any).mockResolvedValue({ data: [] });
        (api.subscribeToPlan as any).mockResolvedValue({});

        render(
            <BrowserRouter>
                <ReadingPlansPage />
            </BrowserRouter>
        );

        await waitFor(() => {
            expect(screen.queryByText("Loading...")).not.toBeInTheDocument();
        });

        // Click Browse Plans
        const browseTab = screen.getByRole("tab", { name: "Browse Plans" });
        await user.click(browseTab);

        // Wait for plans to render
        await waitFor(() => expect(screen.getByText("Subscribe")).toBeInTheDocument());

        await user.click(screen.getByText("Subscribe"));

        await waitFor(() => {
            expect(api.subscribeToPlan).toHaveBeenCalledWith("1");
        });
    });
});
