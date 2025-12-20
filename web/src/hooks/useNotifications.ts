import { useEffect, useState, useCallback } from "react";
import { getToken, onMessage } from "firebase/messaging";
import { messaging } from "@/lib/firebase";
import { registerDevice } from "@/services/api";
import { toast } from "sonner";
import { useAuth } from "./useAuth";

export function useNotifications() {
  const { user } = useAuth();
  const [token, setToken] = useState<string | null>(null);
  const [permission, setPermission] = useState<NotificationPermission>(
    typeof Notification !== "undefined" ? Notification.permission : "default"
  );
  const isSupported = !!messaging;

  const registerToken = useCallback(async () => {
    if (!messaging || !user) return;
    try {
        let registration;
        if ('serviceWorker' in navigator) {
            registration = await navigator.serviceWorker.ready;
        }

        const currentToken = await getToken(messaging, {
            serviceWorkerRegistration: registration
        });

        if (currentToken) {
            setToken(currentToken);
            await registerDevice(currentToken);
        }
    } catch (error) {
        console.error("Error retrieving token:", error);
    }
  }, [user]);

  const requestPermission = useCallback(async () => {
    if (!messaging) return;
    try {
      const status = await Notification.requestPermission();
      setPermission(status);
      if (status === 'granted') {
          await registerToken();
      }
    } catch (error) {
      console.error("Error requesting permission:", error);
    }
  }, [registerToken]);

  useEffect(() => {
    if (!user || !messaging) return;

    // Only attempt registration if permission is already granted
    if (permission === "granted") {
        setTimeout(() => {
            registerToken();
        }, 0);
    }

    const unsubscribe = onMessage(messaging, (payload) => {
      toast(payload.notification?.title || "New Notification", {
        description: payload.notification?.body,
      });
    });

    return () => {
      unsubscribe();
    };
  }, [user, permission, registerToken]);

  return { token, permission, requestPermission, isSupported };
}
