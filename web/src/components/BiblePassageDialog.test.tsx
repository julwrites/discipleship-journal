import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach, Mock } from "vitest";
import { BiblePassageDialog } from "./BiblePassageDialog";
import * as api from "@/services/api";

// Mock API
vi.mock("@/services/api", () => ({
    getBiblePassage: vi.fn(),
    getBibleVersions: vi.fn(),
}));

// Mock ResizeObserver
global.ResizeObserver = class ResizeObserver {
    observe() {}
    unobserve() {}
    disconnect() {}
};

// Mock ScrollIntoView (cmdk uses it)
window.HTMLElement.prototype.scrollIntoView = vi.fn();

describe("BiblePassageDialog", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders passage content", async () => {
        (api.getBiblePassage as Mock).mockResolvedValue({
            text: "<p>In the beginning...</p>",
            reference: "Genesis 1:1",
            version: "ESV"
        });
        (api.getBibleVersions as Mock).mockResolvedValue({ data: [] });

        render(
            <BiblePassageDialog
                reference="Genesis 1:1"
                isOpen={true}
                onClose={() => {}}
                defaultVersion="ESV"
            />
        );

        await waitFor(() => {
            expect(screen.getByText("Genesis 1:1")).toBeInTheDocument();
        });
        expect(screen.getByText("In the beginning...")).toBeInTheDocument();
        expect(api.getBiblePassage).toHaveBeenCalledWith("Genesis 1:1", "ESV");
    });

    it("changes version and refetches passage", async () => {
        (api.getBiblePassage as Mock).mockResolvedValue({
            text: "<p>In the beginning...</p>",
            reference: "Genesis 1:1",
            version: "ESV"
        });
        (api.getBibleVersions as Mock).mockResolvedValue({
            data: [
                { id: "ESV", name: "English Standard Version", abbreviation: "ESV" },
                { id: "KJV", name: "King James Version", abbreviation: "KJV" }
            ]
        });

        render(
            <BiblePassageDialog
                reference="Genesis 1:1"
                isOpen={true}
                onClose={() => {}}
                defaultVersion="ESV"
            />
        );

        await waitFor(() => {
            expect(screen.getByText("Genesis 1:1")).toBeInTheDocument();
        });

        // Open version selector
        const selector = screen.getByRole("combobox");
        fireEvent.click(selector);

        // Select KJV
        const kjvOption = await screen.findByText("KJV");

        // Mock next call
        (api.getBiblePassage as Mock).mockResolvedValue({
            text: "<p>In the beginning God...</p>",
            reference: "Genesis 1:1",
            version: "KJV"
        });

        fireEvent.click(kjvOption);

        await waitFor(() => {
            expect(api.getBiblePassage).toHaveBeenCalledWith("Genesis 1:1", "KJV");
        });

        await waitFor(() => {
             expect(screen.getByText("In the beginning God...")).toBeInTheDocument();
        });
    });
});
