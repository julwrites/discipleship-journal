import { useEffect } from 'react';
import { messaging } from '../lib/firebase';
import { getToken, onMessage } from 'firebase/messaging';
import { api } from '../services/api';

import { Messaging } from 'firebase/messaging';

export function useNotifications() {
  useEffect(() => {
    // Skip if messaging is not initialized (e.g. mock mode or unsupported browser)
    // Also skip if running in E2E mock mode explicitly
    if (!messaging || import.meta.env.VITE_FIREBASE_API_KEY === 'mock-key') return;

    const msg = messaging as Messaging;

    const registerToken = async () => {
      try {
        const permission = await Notification.requestPermission();
        if (permission === 'granted') {
          const token = await getToken(msg, {
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
    const unsubscribe = onMessage(msg, (payload) => {
      console.log('Message received. ', payload);
      if (payload.notification) {
          // Show native notification for now, or could use Sonner
          new Notification(payload.notification.title || "New Notification", {
              body: payload.notification.body
          });
      }
    });

    return () => unsubscribe();
  }, []);
}
