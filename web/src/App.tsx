import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import LoginPage from "@/pages/Login";
import Dashboard from "@/pages/Dashboard";
import NoteEditor from "@/pages/NoteEditor";
import ChatPage from "@/pages/ChatPage";
import ConnectionsPage from "@/pages/ConnectionsPage";
import GroupsPage from "@/pages/GroupsPage";
import Settings from "@/pages/Settings";
import ErrorBoundary from "@/components/ErrorBoundary";

function App() {
  const { user, loading } = useAuth();

  if (loading) {
    return <div className="flex h-screen items-center justify-center">Loading...</div>;
  }

  return (
    <ErrorBoundary>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={user ? <Navigate to="/" /> : <LoginPage />} />

          <Route path="/" element={user ? <Dashboard /> : <Navigate to="/login" />} />
          <Route path="/settings" element={user ? <Settings /> : <Navigate to="/login" />} />
          <Route path="/notes/:id" element={user ? <NoteEditor /> : <Navigate to="/login" />} />
          <Route path="/chat" element={user ? <ChatPage /> : <Navigate to="/login" />} />
          <Route path="/connections" element={user ? <ConnectionsPage /> : <Navigate to="/login" />} />
          <Route path="/groups" element={user ? <GroupsPage /> : <Navigate to="/login" />} />

          <Route path="*" element={<Navigate to="/" />} />
        </Routes>
      </BrowserRouter>
    </ErrorBoundary>
  );
}

export default App;
