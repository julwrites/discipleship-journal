import { useEffect, useState } from "react";
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

  useEffect(() => {
    if (!user || !messaging) return;

    const requestPermission = async () => {
      try {
        const status = await Notification.requestPermission();
        setPermission(status);

        if (status === "granted") {
          let registration;
          if ('serviceWorker' in navigator) {
             registration = await navigator.serviceWorker.ready;
          }

          const currentToken = await getToken(messaging, {
              serviceWorkerRegistration: registration
          });

          if (currentToken) {
            setToken(currentToken);
            // Only register if token changed? Ideally backend handles idempotency.
            await registerDevice(currentToken);
            console.log("Device registered for notifications");
          }
        }
      } catch (error) {
        console.error("An error occurred while retrieving token: ", error);
      }
    };

    requestPermission();

    const unsubscribe = onMessage(messaging, (payload) => {
      console.log("Foreground message received:", payload);
      toast(payload.notification?.title || "New Notification", {
        description: payload.notification?.body,
      });
    });

    return () => {
      unsubscribe();
    };
  }, [user]);

  return { token, permission };
}
