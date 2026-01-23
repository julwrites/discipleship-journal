import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";

export function AuthGuard() {
  const { user, loading } = useAuth();
  if (loading) return <div className="flex h-screen items-center justify-center">Loading...</div>;
  if (!user) return <Navigate to="/login" replace />;
  return <Outlet />;
}

export function GuestGuard() {
  const { user, loading } = useAuth();
  if (loading) return <div className="flex h-screen items-center justify-center">Loading...</div>;
  if (user) return <Navigate to="/" replace />;
  return <Outlet />;
}
