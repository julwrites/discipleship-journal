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

    // Fallback for E2E testing with mock keys.
    // This bypasses Firebase auth initialization which fails with invalid keys.
    // We skip this for unit tests (MODE === 'test') as they mock onAuthStateChanged directly.
    const shouldBypassAuth = import.meta.env.VITE_FIREBASE_API_KEY === 'mock-key' && import.meta.env.MODE !== 'test';

    if (shouldBypassAuth) {
         console.warn("Using mock auth keys, bypassing auth check after timeout.");

         // Check for simulated user in localStorage for E2E tests
         const mockUserJson = localStorage.getItem('E2E_TEST_USER');
         if (mockUserJson) {
             try {
                 const mockUser = JSON.parse(mockUserJson) as User;
                 // Mock the getIdToken method which is required by some components/logic
                 mockUser.getIdToken = async () => "mock-token";

                 // Avoid synchronous setState in effect
                 setTimeout(() => {
                    setUser(mockUser);
                    setLoading(false);
                 }, 0);
                 return;
             } catch (e) {
                 console.error("Failed to parse E2E_TEST_USER", e);
             }
         }

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
