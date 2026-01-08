import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import LoginPage from "@/pages/Login";
import Dashboard from "@/pages/Dashboard";
import NoteEditor from "@/pages/NoteEditor";
import ChatPage from "@/pages/ChatPage";
import ConnectionsPage from "@/pages/ConnectionsPage";
import GroupsPage from "@/pages/GroupsPage";
import ReadingPlansPage from "@/pages/ReadingPlansPage";
import ReadingPlanDetail from "@/pages/ReadingPlanDetail";
import MemoryVersesPage from "@/pages/MemoryVersesPage";
import TemplatesPage from "@/pages/TemplatesPage";
import TemplateEditor from "@/pages/TemplateEditor";
import Settings from "@/pages/Settings";
import PrivacyPolicy from "@/pages/PrivacyPolicy";
import TermsOfService from "@/pages/TermsOfService";
import PublicLayout from "@/components/PublicLayout";
import ErrorBoundary from "@/components/ErrorBoundary";
import { Toaster } from "@/components/ui/sonner";
import { useNotifications } from "@/hooks/useNotifications";
import { router } from "@/Router";

function App() {
  useNotifications();

  return (
    <ErrorBoundary>
      <Toaster />
      <BrowserRouter>
        <Routes>
          <Route element={<PublicLayout />}>
            <Route path="/privacy-policy" element={<PrivacyPolicy />} />
            <Route path="/terms-of-service" element={<TermsOfService />} />
          </Route>

          <Route path="/login" element={user ? <Navigate to="/" /> : <LoginPage />} />

          <Route path="/" element={user ? <Dashboard /> : <Navigate to="/login" />} />
          <Route path="/settings" element={user ? <Settings /> : <Navigate to="/login" />} />
          <Route path="/notes/:id" element={user ? <NoteEditor /> : <Navigate to="/login" />} />
          <Route path="/chat" element={user ? <ChatPage /> : <Navigate to="/login" />} />
          <Route path="/connections" element={user ? <ConnectionsPage /> : <Navigate to="/login" />} />
          <Route path="/groups" element={user ? <GroupsPage /> : <Navigate to="/login" />} />
          <Route path="/reading-plans" element={user ? <ReadingPlansPage /> : <Navigate to="/login" />} />
          <Route path="/reading-plans/:id" element={user ? <ReadingPlanDetail /> : <Navigate to="/login" />} />
          <Route path="/memory-verses" element={user ? <MemoryVersesPage /> : <Navigate to="/login" />} />

          <Route path="/templates" element={user ? <TemplatesPage /> : <Navigate to="/login" />} />
          <Route path="/templates/:id/edit" element={user ? <TemplateEditor /> : <Navigate to="/login" />} />
          <Route path="/templates/new" element={user ? <TemplateEditor /> : <Navigate to="/login" />} />

          <Route path="*" element={<Navigate to="/" />} />
        </Routes>
      </BrowserRouter>
    </ErrorBoundary>
  );
}

export default App;
