import { auth } from "@/lib/firebase";
import { GoogleAuthProvider, signInWithPopup, AuthError } from "firebase/auth";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { useState } from "react";

export default function LoginPage() {
  const [error, setError] = useState<string | null>(null);

  const handleLogin = async () => {
    if (!auth) {
      console.error("Cannot sign in: Auth not initialized");
      setError("Auth not initialized");
      return;
    }
    const provider = new GoogleAuthProvider();
    try {
      await signInWithPopup(auth, provider);
    } catch (e: unknown) {
      console.error("Login failed", e);
      let message = "Login failed. Please try again.";
      const error = e as AuthError;

      if (error?.code === "auth/unauthorized-domain") {
        message = "This domain is not authorized for authentication. Please add it to the Firebase Console > Authentication > Settings > Authorized Domains.";
      } else if (error?.code === "auth/popup-closed-by-user") {
        message = "Login popup was closed before completing sign in.";
      } else if (error?.message) {
        message = error.message;
      }

      setError(message);
      toast.error(message);
    }
  };

  if (!auth) {
    return (
      <div className="flex h-screen items-center justify-center bg-gray-50">
        <div className="text-center space-y-4">
          <h1 className="text-3xl font-bold text-red-600">Configuration Error</h1>
          <p className="text-gray-600">
            Firebase authentication is not initialized. Please check your deployment configuration and environment variables.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-screen items-center justify-center bg-gray-50">
      <div className="text-center space-y-4">
        <h1 className="text-3xl font-bold">Discipleship Journal</h1>
        <p className="text-gray-600">Sign in to continue</p>

        {error && (
            <div className="max-w-md mx-auto p-4 bg-red-50 border border-red-200 rounded-md text-red-700 text-sm">
                {error}
            </div>
        )}

        <Button onClick={handleLogin}>Sign in with Google</Button>
      </div>
    </div>
  );
}
