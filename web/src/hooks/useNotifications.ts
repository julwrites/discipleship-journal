import { useEffect } from 'react';
import { messaging } from '../lib/firebase';
import { getToken, onMessage } from 'firebase/messaging';
import { api } from '../services/api';

export function useNotifications() {
  useEffect(() => {
    if (!messaging) return;

    const registerToken = async () => {
      try {
        const permission = await Notification.requestPermission();
        if (permission === 'granted') {
          const token = await getToken(messaging, {
            vapidKey: import.meta.env.VITE_FIREBASE_VAPID_KEY
          });

          if (token) {
            // Send token to backend
            await api.post('/api/notifications/register', {
              token,
              device_type: 'web'
            });
            console.log("FCM Token registered");
          }
        }
      } catch (error) {
        console.error("Failed to register notification token", error);
      }
    };

    registerToken();

    // Handle foreground messages
    const unsubscribe = onMessage(messaging, (payload) => {
      console.log('Message received. ', payload);
      // You can use sonner or existing toast notification here
      // For now just log it or maybe show a native notification if document is hidden?
      // Typically foreground handling is for in-app toasts.
      if (payload.notification) {
          // If we had a toast library hook, we'd use it here.
          // Since this is a hook, we might need to pass a callback or context.
          // For MVP, simply logging.
          // Ideally we dispatch an event or use a global store.
          new Notification(payload.notification.title || "New Notification", {
              body: payload.notification.body
          });
      }
    });

    return () => unsubscribe();
  }, []);
}
