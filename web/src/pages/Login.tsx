import { auth } from "@/lib/firebase";
import { GoogleAuthProvider, signInWithPopup } from "firebase/auth";
import { Button } from "@/components/ui/button";

export default function LoginPage() {
  const handleLogin = async () => {
    if (!auth) {
      console.error("Cannot sign in: Auth not initialized");
      return;
    }
    const provider = new GoogleAuthProvider();
    try {
      await signInWithPopup(auth, provider);
    } catch (error) {
      console.error("Login failed", error);
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
        <Button onClick={handleLogin}>Sign in with Google</Button>
      </div>
    </div>
  );
}
