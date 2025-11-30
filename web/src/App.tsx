import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import LoginPage from "@/pages/Login";
import Dashboard from "@/pages/Dashboard";
import NoteEditor from "@/pages/NoteEditor";
import ChatPage from "@/pages/ChatPage";
import Settings from "@/pages/Settings";

function App() {
  const { user, loading } = useAuth();

  if (loading) {
    return <div className="flex h-screen items-center justify-center">Loading...</div>;
  }

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={user ? <Navigate to="/" /> : <LoginPage />} />

        <Route path="/" element={user ? <Dashboard /> : <Navigate to="/login" />} />
        <Route path="/settings" element={user ? <Settings /> : <Navigate to="/login" />} />
        <Route path="/notes/new" element={user ? <NoteEditor /> : <Navigate to="/login" />} />
        <Route path="/notes/:id" element={user ? <NoteEditor /> : <Navigate to="/login" />} />
        <Route path="/chat" element={user ? <ChatPage /> : <Navigate to="/login" />} />

        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
