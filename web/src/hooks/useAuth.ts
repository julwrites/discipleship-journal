import { useEffect, useState } from "react";
import { User, onAuthStateChanged } from "firebase/auth";
import { auth } from "@/lib/firebase";

export function useAuth() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // If we are in development mode and using mock keys, we might need to handle the error
    // or just timeout the loading state if onAuthStateChanged never fires.
    // However, onAuthStateChanged should still fire or we should catch the initialization error.

    // As a fallback for tests where keys are invalid
    if (import.meta.env.MODE === 'test' || import.meta.env.VITE_FIREBASE_API_KEY === 'mock-key') {
         console.warn("Using mock auth keys, bypassing auth check after timeout.");
         const timer = setTimeout(() => {
             setLoading(false);
         }, 1000);
         return () => clearTimeout(timer);
    }

    const unsubscribe = onAuthStateChanged(auth, (user) => {
      setUser(user);
      setLoading(false);
    }, (error) => {
        console.error("Auth state change error:", error);
        setLoading(false);
    });

    return unsubscribe;
  }, []);

  return { user, loading };
}
