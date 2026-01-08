import { RouterProvider } from "react-router-dom";
import ErrorBoundary from "@/components/ErrorBoundary";
import { Toaster } from "@/components/ui/sonner";
import { useNotifications } from "@/hooks/useNotifications";
import { router } from "@/Router";

function App() {
  useNotifications();

  return (
    <ErrorBoundary>
      <Toaster />
      <RouterProvider router={router} />
    </ErrorBoundary>
  );
}

export default App;
