export default function PrivacyPolicy() {
  return (
    <div className="prose dark:prose-invert max-w-none">
      <h1>Privacy Policy</h1>
      <p className="lead">Last updated: {new Date().toLocaleDateString()}</p>

      <p>
        At Discipleship Journal, we take your privacy seriously. This Privacy Policy explains how we collect, use, and protect your information when you use our application.
      </p>

      <h2>1. Information We Collect</h2>
      <p>We collect the following types of information:</p>
      <ul>
        <li><strong>Account Information:</strong> When you create an account, we collect your email address and authentication credentials (via Google or Email/Password provider).</li>
        <li><strong>User Content:</strong> Any data you enter into the application, including journal entries, notes, prayers, and group messages ("User Content").</li>
        <li><strong>Device Information:</strong> We may collect basic device information to optimize the user experience and for push notifications (if enabled).</li>
      </ul>

      <h2>2. How We Use Your Information</h2>
      <p>We use your information solely to provide and improve the Discipleship Journal service:</p>
      <ul>
        <li>To authenticate you and secure your account.</li>
        <li>To store and sync your journal entries across your devices.</li>
        <li>To facilitate sharing with groups you explicitly join and share content with.</li>
        <li>To provide AI-assisted features (using Bible context) upon your request.</li>
      </ul>

      <h2>3. Data Privacy and Security</h2>
      <p>
        <strong>Your data is private by default.</strong> We do not sell your personal data to third parties.
      </p>
      <ul>
        <li>Your notes and journal entries are stored securely on our servers.</li>
        <li>Content is only shared with other users if you explicitly choose to share it (e.g., via Groups).</li>
        <li>We employ industry-standard security measures to protect your data from unauthorized access.</li>
      </ul>

      <h2>4. Data Retention and Deletion</h2>
      <p>
        We retain your data for as long as your account is active. You may delete your notes or your entire account at any time through the application settings. Upon account deletion, your data will be permanently removed from our active databases.
      </p>

      <h2>5. Changes to This Policy</h2>
      <p>
        We may update this Privacy Policy from time to time. We will notify you of any changes by posting the new Privacy Policy on this page.
      </p>

      <h2>6. Contact Us</h2>
      <p>
        If you have any questions about this Privacy Policy, please contact us.
      </p>
    </div>
  );
}
