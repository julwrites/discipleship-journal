import { initializeApp } from "firebase/app";
import { getAuth } from "firebase/auth";

const firebaseConfig = {
  apiKey: import.meta.env.FIREBASE_API_KEY,
  authDomain: import.meta.env.FIREBASE_AUTH_DOMAIN,
  projectId: import.meta.env.FIREBASE_PROJECT_ID,
  storageBucket: import.meta.env.FIREBASE_STORAGE_BUCKET,
  messagingSenderId: import.meta.env.FIREBASE_MESSAGING_SENDER_ID,
  appId: import.meta.env.FIREBASE_APP_ID,
  measurementId: import.meta.env.FIREBASE_MEASUREMENT_ID
};

// Conditional initialization to prevent crashes during mock testing if keys are invalid
let appInstance;
let authInstance;

try {
  appInstance = initializeApp(firebaseConfig);
  authInstance = getAuth(appInstance);
} catch (e) {
  console.warn("Firebase initialization failed (expected during mock testing):", e);
  // Provide a dummy auth object if needed, or rely on handling the error where it's used.
  // However, most components import 'auth' directly.
  // We can return a mock-like object or null, but type safety is key.
  // For now, if we are in a mock environment (detected by check above), we might want to return a dummy.

  if (import.meta.env.FIREBASE_API_KEY === 'mock-key') {
      // Mocking minimal auth object to prevent crash on import
      authInstance = {
          currentUser: null,
          onAuthStateChanged: () => () => {},
          // Add other methods as needed or keep it minimal
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      } as any;
  }
}

export const app = appInstance;
export const auth = authInstance;
