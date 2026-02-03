import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import ReadingPlanDetail from "./ReadingPlanDetail";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { describe, it, expect, vi, beforeEach, Mock } from "vitest";
import * as api from "@/services/api";

// Mock API
vi.mock("@/services/api", () => ({
    getReadingPlan: vi.fn(),
    getPlanProgress: vi.fn(),
    markPlanDayComplete: vi.fn(),
    unmarkPlanDayComplete: vi.fn(),
    syncUser: vi.fn(),
    getBiblePassage: vi.fn(),
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

describe("ReadingPlanDetail", () => {
    const planId = "test-plan-id";
    const mockPlan = {
        id: planId,
        title: "Test Plan",
        description: "Test Description",
        days: [
            { id: "d1", day_number: 1, passage: "Gen 1" },
            { id: "d2", day_number: 2, passage: "Gen 2" },
        ],
    };

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders plan details", async () => {
        (api.getReadingPlan as Mock).mockResolvedValue(mockPlan);
        (api.getPlanProgress as Mock).mockResolvedValue({ completed_days: [] });

        render(
            <MemoryRouter initialEntries={[`/reading-plans/${planId}`]}>
                <Routes>
                    <Route path="/reading-plans/:id" element={<ReadingPlanDetail />} />
                </Routes>
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(screen.getByText("Test Plan")).toBeInTheDocument();
        });
        expect(screen.getByText("Test Description")).toBeInTheDocument();
        expect(screen.getByText("Day 1")).toBeInTheDocument();
        expect(screen.getByText("Gen 1")).toBeInTheDocument();
    });

    it("shows completed status", async () => {
        (api.getReadingPlan as Mock).mockResolvedValue(mockPlan);
        (api.getPlanProgress as Mock).mockResolvedValue({ completed_days: [1] });

        render(
            <MemoryRouter initialEntries={[`/reading-plans/${planId}`]}>
                <Routes>
                    <Route path="/reading-plans/:id" element={<ReadingPlanDetail />} />
                </Routes>
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(screen.getByText("1 / 2 Days Completed")).toBeInTheDocument();
        });

        // Day 1 should be disabled/completed
        // We look for the CheckCircle icon or disabled button
        // The implementation renders a disabled button with ghost variant and CheckCircle

        // Find button for Day 1. It is the one that has CheckCircle inside?
        // Or simpler, look for "Mark Complete" which should be present for Day 2 but NOT Day 1

        expect(screen.queryAllByText("Mark Complete")).toHaveLength(1); // Only for day 2
    });

    it("handles marking day as complete", async () => {
        const user = userEvent.setup();
        (api.getReadingPlan as Mock).mockResolvedValue(mockPlan);
        (api.getPlanProgress as Mock).mockResolvedValue({ completed_days: [] });
        (api.markPlanDayComplete as Mock).mockResolvedValue({});

        render(
            <MemoryRouter initialEntries={[`/reading-plans/${planId}`]}>
                <Routes>
                    <Route path="/reading-plans/:id" element={<ReadingPlanDetail />} />
                </Routes>
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(screen.getByText("Test Plan")).toBeInTheDocument();
        });

        const markButtons = screen.getAllByText("Mark Complete");
        await user.click(markButtons[0]); // Mark Day 1

        await waitFor(() => {
            expect(api.markPlanDayComplete).toHaveBeenCalledWith(planId, 1);
        });

        // Should optimistically update UI
        expect(screen.getByText("1 / 2 Days Completed")).toBeInTheDocument();
    });
});
