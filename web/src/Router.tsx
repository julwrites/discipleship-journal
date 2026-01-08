import { createBrowserRouter, Navigate, Outlet } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import LoginPage from "@/pages/Login";
import Dashboard from "@/pages/Dashboard";
import NoteEditor from "@/pages/NoteEditor";
import ChatPage from "@/pages/ChatPage";
import ConnectionsPage from "@/pages/ConnectionsPage";
import GroupsPage from "@/pages/GroupsPage";
import ReadingPlansPage from "@/pages/ReadingPlansPage";
import ReadingPlanDetail from "@/pages/ReadingPlanDetail";
import Settings from "@/pages/Settings";
import PrivacyPolicy from "@/pages/PrivacyPolicy";
import TermsOfService from "@/pages/TermsOfService";
import PublicLayout from "@/components/PublicLayout";

function AuthGuard() {
  const { user, loading } = useAuth();
  if (loading) return <div className="flex h-screen items-center justify-center">Loading...</div>;
  if (!user) return <Navigate to="/login" replace />;
  return <Outlet />;
}

function GuestGuard() {
  const { user, loading } = useAuth();
  if (loading) return <div className="flex h-screen items-center justify-center">Loading...</div>;
  if (user) return <Navigate to="/" replace />;
  return <Outlet />;
}

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
    ],
  },
  {
    path: "*",
    element: <Navigate to="/" />,
  },
]);
