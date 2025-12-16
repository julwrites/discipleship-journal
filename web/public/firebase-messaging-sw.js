// Give the service worker access to Firebase Messaging.
// Note that you can only use Firebase Messaging here. Other Firebase libraries are not available in the service worker.
importScripts('https://www.gstatic.com/firebasejs/8.10.1/firebase-app.js');
importScripts('https://www.gstatic.com/firebasejs/8.10.1/firebase-messaging.js');

// Initialize the Firebase app in the service worker.
// Note: Service workers in production often need these keys injected during build.
// Since we are using Vite, we can't easily inject env vars into files in `public/` without a plugin or script.
// However, Firebase config is generally public.
// TODO: Replace these placeholders with actual values or use a build script to generate this file.

const firebaseConfig = {
  // These placeholders must be replaced by the build system or manually for production
  apiKey: "REPLACE_WITH_YOUR_KEY",
  projectId: "discipleship-journal-pwa",
  messagingSenderId: "REPLACE_WITH_SENDER_ID",
  appId: "REPLACE_WITH_APP_ID"
};

// Try to fetch config if available (e.g. served by app)
// fetch('/firebase-config.json').then(response => response.json()).then(config => firebase.initializeApp(config)).catch(() => firebase.initializeApp(firebaseConfig));

firebase.initializeApp(firebaseConfig);

// Retrieve an instance of Firebase Messaging so that it can handle background
// messages.
const messaging = firebase.messaging();

messaging.onBackgroundMessage((payload) => {
  console.log('[firebase-messaging-sw.js] Received background message ', payload);
  // Customize notification here
  const notificationTitle = payload.notification.title;
  const notificationOptions = {
    body: payload.notification.body,
    icon: '/pwa-192x192.png'
  };

  self.registration.showNotification(notificationTitle, notificationOptions);
});
