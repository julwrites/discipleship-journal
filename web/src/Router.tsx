import { createBrowserRouter, Navigate } from "react-router-dom";
import LoginPage from "@/pages/Login";
import Dashboard from "@/pages/Dashboard";
import NoteEditor from "@/pages/NoteEditor";
import ChatPage from "@/pages/ChatPage";
import ConnectionsPage from "@/pages/ConnectionsPage";
import GroupsPage from "@/pages/GroupsPage";
import ReadingPlansPage from "@/pages/ReadingPlansPage";
import ReadingPlanDetail from "@/pages/ReadingPlanDetail";
import MemoryVersesPage, { VersePackDetail } from "@/pages/MemoryVersesPage";
import TemplatesPage from "@/pages/TemplatesPage";
import TemplateEditor from "@/pages/TemplateEditor";
import Settings from "@/pages/Settings";
import PrivacyPolicy from "@/pages/PrivacyPolicy";
import TermsOfService from "@/pages/TermsOfService";
import PublicLayout from "@/components/PublicLayout";
import { AuthGuard, GuestGuard } from "@/components/AuthGuards";

export const router = createBrowserRouter([
  {
    element: <PublicLayout />,
    children: [
      {
        path: "/privacy-policy",
        element: <PrivacyPolicy />,
      },
      {
        path: "/terms-of-service",
        element: <TermsOfService />,
      },
    ],
  },
  {
    path: "/login",
    element: <GuestGuard />,
    children: [
      {
        path: "",
        element: <LoginPage />,
      },
    ],
  },
  {
    element: <AuthGuard />,
    children: [
      {
        path: "/",
        element: <Dashboard />,
      },
      {
        path: "/settings",
        element: <Settings />,
      },
      {
        path: "/notes/:id",
        element: <NoteEditor />,
      },
      {
        path: "/chat",
        element: <ChatPage />,
      },
      {
        path: "/connections",
        element: <ConnectionsPage />,
      },
      {
        path: "/groups",
        element: <GroupsPage />,
      },
      {
        path: "/reading-plans",
        element: <ReadingPlansPage />,
      },
      {
        path: "/reading-plans/:id",
        element: <ReadingPlanDetail />,
      },
      {
        path: "/memory-verses",
        element: <MemoryVersesPage />,
      },
      {
        path: "/memory-verses/:id",
        element: <VersePackDetail />,
      },
      {
        path: "/templates",
        element: <TemplatesPage />,
      },
      {
        path: "/templates/:id/edit",
        element: <TemplateEditor />,
      },
      {
        path: "/templates/new",
        element: <TemplateEditor />,
      },
    ],
  },
  {
    path: "*",
    element: <Navigate to="/" />,
  },
]);
