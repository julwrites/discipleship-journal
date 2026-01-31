import { initializeApp } from "firebase/app";
import { getAuth, setPersistence, browserLocalPersistence } from "firebase/auth";
import { getMessaging } from "firebase/messaging";

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
let messagingInstance;
let authReadyPromise: Promise<void> | null = null;

try {
  console.log("Initializing Firebase with config:", {
    projectId: firebaseConfig.projectId,
    authDomain: firebaseConfig.authDomain,
    hasApiKey: !!firebaseConfig.apiKey
  });
  appInstance = initializeApp(firebaseConfig);
  authInstance = getAuth(appInstance);

  // Set persistence for auth state and store promise for components to wait on
  try {
    authReadyPromise = setPersistence(authInstance, browserLocalPersistence)
      .then(() => {
        console.log("Firebase auth persistence set to browserLocalPersistence");
      })
      .catch((persistenceError) => {
        console.warn("Failed to set auth persistence:", persistenceError);
      });
  } catch (persistenceError) {
    console.warn("Failed to set auth persistence (sync error):", persistenceError);
    // Create a resolved promise so components don't wait forever
    authReadyPromise = Promise.resolve();
  }

  console.log("Firebase initialized successfully");

  try {
    if (typeof window !== 'undefined' && 'serviceWorker' in navigator) {
       messagingInstance = getMessaging(appInstance);
    }
  } catch (e) {
    console.warn("Firebase Messaging initialization failed:", e);
  }

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

      // Messaging is null in mock mode to prevent SDK crashes
      messagingInstance = null;
  } else {
      // For non-mock failures, set authInstance to null
      authInstance = null;
  }
  // Ensure authReadyPromise is set to avoid waiting forever
  authReadyPromise = Promise.resolve();
  console.error("Firebase initialization catch block activated. Auth instance might be mocked or null.");
}

export const app = appInstance;
export const auth = authInstance;
export const messaging = messagingInstance;
export { authReadyPromise };
