import { useEffect, useState } from "react";
import { User, onAuthStateChanged } from "firebase/auth";
import { auth } from "@/lib/firebase";

export function useAuth() {
  const shouldBypassAuth = import.meta.env.VITE_FIREBASE_API_KEY === 'mock-key' && import.meta.env.MODE !== 'test';

  const [user, setUser] = useState<User | null>(() => {
    if (shouldBypassAuth) {
        const mockUserJson = localStorage.getItem('E2E_TEST_USER');
        if (mockUserJson) {
            try {
                const mockUser = JSON.parse(mockUserJson) as User;
                // Mock the getIdToken method which is required by some components/logic
                mockUser.getIdToken = async () => "mock-token";
                return mockUser;
            } catch (e) {
                console.error("Failed to parse E2E_TEST_USER", e);
            }
        }
    }
    return null;
  });

  // If we already have a user from mock init, we are not loading.
  // Otherwise, we start as loading.
  const [loading, setLoading] = useState<boolean>(() => {
      if (shouldBypassAuth && user) return false;
      return true;
  });

  useEffect(() => {
    if (shouldBypassAuth) {
         if (user) return; // Already loaded from localStorage

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
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  return { user, loading };
}
