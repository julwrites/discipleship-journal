import { auth } from "@/lib/firebase";
import {
  GoogleAuthProvider,
  signInWithRedirect,
  createUserWithEmailAndPassword,
  signInWithEmailAndPassword,
  getRedirectResult,
  sendSignInLinkToEmail,
  isSignInWithEmailLink,
  signInWithEmailLink,
  AuthError
} from "firebase/auth";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Link } from "react-router-dom";
import { toast } from "sonner";
import { useState, useEffect } from "react";
import { Book, Shield, Users, CheckCircle2 } from "lucide-react";

export default function LoginPage() {
  const [isLoading, setIsLoading] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [magicLinkEmail, setMagicLinkEmail] = useState("");
  const [activeTab, setActiveTab] = useState("login");

  useEffect(() => {
    console.log("Login page useEffect running, auth:", !!auth, "URL:", window.location.href, "Search:", window.location.search, "Hash:", window.location.hash);
    if (!auth) return;

    // Debug: check for any OAuth parameters in URL
    const urlParams = new URLSearchParams(window.location.search);
    const hashParams = new URLSearchParams(window.location.hash.substring(1));
    console.log("URL query params:", Array.from(urlParams.entries()));
    console.log("URL hash params:", Array.from(hashParams.entries()));

    // Check sessionStorage for redirect state (Firebase stores OAuth state here)
    try {
      const sessionKeys = Object.keys(sessionStorage);
      const firebaseKeys = sessionKeys.filter(key => key.includes('firebase') || key.includes('auth'));
      console.log("SessionStorage firebase/auth keys:", firebaseKeys);

      // Try to read and parse the Firebase redirect event
      for (const key of firebaseKeys) {
        try {
          const value = sessionStorage.getItem(key);
          console.log(`SessionStorage ${key}:`, value ? JSON.parse(value) : value);
        } catch (parseError) {
          console.log(`SessionStorage ${key} (raw):`, sessionStorage.getItem(key));
        }
      }
    } catch (e) {
      console.warn("Could not access sessionStorage:", e);
    }

    // Handle Google Redirect Result
    console.log("Calling getRedirectResult...");
    getRedirectResult(auth)
      .then((result) => {
        console.log("getRedirectResult resolved:", result);
        if (result) {
          console.log("Redirect login successful:", result.user?.email, "Provider:", result.providerId);
          // Auth state will be updated by onAuthStateChanged listener
        } else {
          console.log("No redirect result to process");
          // Fallback: check if user is already signed in (might have happened automatically)
          console.log("Auth currentUser:", auth.currentUser?.email);

          // Additional debugging: check if sessionStorage should be cleared
          console.log("Checking if sessionStorage should be cleared...");
          const shouldClearStorage = window.location.search.includes('error') ||
                                    window.location.search.includes('state');
          console.log("Should clear sessionStorage?", shouldClearStorage);
        }
      })
      .catch((error) => {
        console.error("Redirect login error:", error);
        console.error("Error details:", {
          code: error.code,
          message: error.message,
          stack: error.stack
        });
        const msg = getErrorMessage(error as AuthError);
        toast.error(msg);
      });

    // Handle Email Link Sign-in
    if (isSignInWithEmailLink(auth, window.location.href)) {
      let emailForSignIn = window.localStorage.getItem('emailForSignIn');
      if (!emailForSignIn) {
        emailForSignIn = window.prompt('Please provide your email for confirmation');
      }

      if (emailForSignIn) {
        setIsLoading(true);
        signInWithEmailLink(auth, emailForSignIn, window.location.href)
          .then(() => {
            window.localStorage.removeItem('emailForSignIn');
            toast.success("Successfully signed in!");
            // Navigation handled by auth listener in App or router
          })
          .catch((error) => {
            console.error("Email link sign in error:", error);
            const msg = getErrorMessage(error as AuthError);
            toast.error(msg);
            setIsLoading(false);
          });
      }
    }
  }, []);

  const getErrorMessage = (error: AuthError) => {
    switch (error.code) {
      case "auth/invalid-email":
        return "Invalid email address.";
      case "auth/user-disabled":
        return "This user account has been disabled.";
      case "auth/user-not-found":
        return "No account found with this email.";
      case "auth/wrong-password":
        return "Incorrect password.";
      case "auth/email-already-in-use":
        return "An account already exists with this email.";
      case "auth/weak-password":
        return "Password should be at least 6 characters.";
      case "auth/popup-closed-by-user":
        return "Sign in was cancelled.";
      case "auth/popup-blocked":
        return "Popup was blocked by the browser. Redirecting to sign in page...";
      case "auth/unauthorized-domain":
        return "Domain not authorized. Check Firebase Console.";
      case "auth/invalid-action-code":
        return "The login link has expired or has already been used.";
      default:
        return error.message || "An error occurred during authentication.";
    }
  };

  const handleGoogleLogin = async (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    console.log("handleGoogleLogin called.", { auth: !!auth });

    if (!auth) {
      toast.error("Authentication not initialized.");
      return;
    }
    setIsLoading(true);
    const provider = new GoogleAuthProvider();
    console.log("GoogleAuthProvider created");

    // Use redirect flow by default - popup is unreliable with modern browser security policies
    console.log("Using redirect flow for Google Sign-In");
    try {
      console.log("Calling signInWithRedirect...");
      await signInWithRedirect(auth, provider);
      console.log("signInWithRedirect completed, redirect should happen");
      // User will be redirected to Google and back
      return; // Redirect will happen, no further processing needed
    } catch (error) {
      const authError = error as AuthError;
      console.error("Google Sign In (Redirect) failed:", authError);
      toast.error(getErrorMessage(authError));
      setIsLoading(false);
    }
  };

  const handleEmailLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!auth) return;
    setIsLoading(true);
    try {
      await signInWithEmailAndPassword(auth, email, password);
    } catch (e) {
      const msg = getErrorMessage(e as AuthError);
      toast.error(msg);
    } finally {
      setIsLoading(false);
    }
  };

  const handleEmailSignup = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!auth) return;
    setIsLoading(true);
    try {
      await createUserWithEmailAndPassword(auth, email, password);
      toast.success("Account created successfully!");
    } catch (e) {
      const msg = getErrorMessage(e as AuthError);
      toast.error(msg);
    } finally {
      setIsLoading(false);
    }
  };

  const handleMagicLinkSend = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!auth) return;
    if (!magicLinkEmail) {
        toast.error("Please enter your email address.");
        return;
    }

    setIsLoading(true);
    const actionCodeSettings = {
      // URL must be whitelisted in Firebase Console.
      url: window.location.origin + '/login', // Redirect back to login page to handle verification
      handleCodeInApp: true,
    };

    try {
      await sendSignInLinkToEmail(auth, magicLinkEmail, actionCodeSettings);
      window.localStorage.setItem('emailForSignIn', magicLinkEmail);
      toast.success("Login link sent! Please check your email.");
      setMagicLinkEmail(""); // Clear input
    } catch (e) {
      console.error("Send Magic Link Error:", e);
      const msg = getErrorMessage(e as AuthError);
      toast.error(msg);
    } finally {
      setIsLoading(false);
    }
  };

  if (!auth) {
    return (
      <div className="flex h-screen items-center justify-center bg-background">
        <div className="text-center space-y-4 max-w-md p-6">
          <h1 className="text-3xl font-bold text-destructive">Configuration Error</h1>
          <p className="text-muted-foreground">
            Firebase authentication is not initialized. Please check your deployment configuration and environment variables.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen w-full">
      {/* Left Column - Hero/Branding */}
      <div className="hidden lg:flex w-1/2 bg-slate-900 text-white p-12 flex-col justify-between relative overflow-hidden">
        <div className="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1507692049790-de58293a469d?q=80&w=2070&auto=format&fit=crop')] opacity-10 bg-cover bg-center" />

        <div className="relative z-10">
          <div className="flex items-center space-x-2 text-xl font-bold mb-8">
            <Book className="h-6 w-6" />
            <span>Discipleship Journal</span>
          </div>

          <h1 className="text-5xl font-extrabold tracking-tight mb-6 leading-tight">
            Capture your thoughts, <span className="text-primary">prayers</span>, and growth.
          </h1>
          <p className="text-lg text-slate-300 max-w-md mb-8">
            A secure and private space for your spiritual journey. Document your walk, share with small groups, and reflect on your progress.
          </p>
        </div>

        <div className="relative z-10 grid gap-6">
          <div className="flex items-start space-x-4">
            <Shield className="h-6 w-6 text-primary mt-1" />
            <div>
              <h3 className="font-semibold text-lg">Private & Secure</h3>
              <p className="text-slate-400">Your journal entries are private by default. You control what you share.</p>
            </div>
          </div>
          <div className="flex items-start space-x-4">
            <Users className="h-6 w-6 text-primary mt-1" />
            <div>
              <h3 className="font-semibold text-lg">Group Sharing</h3>
              <p className="text-slate-400">Connect with your small group. Share prayer requests and insights effortlessly.</p>
            </div>
          </div>
        </div>

        <div className="relative z-10 text-sm text-slate-500">
          &copy; {new Date().getFullYear()} Discipleship Journal. All rights reserved.
        </div>
      </div>

      {/* Right Column - Auth Form */}
      <div className="flex-1 flex items-center justify-center p-8 bg-background">
        <Card className="w-full max-w-md shadow-lg border-border">
          <CardHeader className="space-y-1 text-center">
            <div className="flex justify-center mb-4 lg:hidden">
              <Book className="h-10 w-10 text-primary" />
            </div>
            <CardTitle className="text-2xl font-bold text-foreground">Welcome back</CardTitle>
            <CardDescription className="text-muted-foreground">
              {activeTab === 'magic-link' ? 'Enter your email to receive a magic link' : 'Enter your email to sign in to your account'}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
              <TabsList className="grid w-full grid-cols-3 mb-4 bg-muted">
                <TabsTrigger value="login" className="data-[state=active]:bg-background data-[state=active]:text-foreground">Sign In</TabsTrigger>
                <TabsTrigger value="magic-link" className="data-[state=active]:bg-background data-[state=active]:text-foreground">Magic Link</TabsTrigger>
                <TabsTrigger value="register" className="data-[state=active]:bg-background data-[state=active]:text-foreground">Register</TabsTrigger>
              </TabsList>

              <TabsContent value="login">
                <form onSubmit={handleEmailLogin} className="space-y-4">
                  <div className="space-y-2">
                    <Input
                      type="email"
                      placeholder="name@example.com"
                      required
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      className="bg-background text-foreground border-input"
                    />
                  </div>
                  <div className="space-y-2">
                    <Input
                      type="password"
                      placeholder="Password"
                      required
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      className="bg-background text-foreground border-input"
                    />
                  </div>
                  <Button type="submit" className="w-full">
                    {isLoading ? "Signing in..." : "Sign In with Email"}
                  </Button>
                </form>
              </TabsContent>

              <TabsContent value="magic-link">
                <form onSubmit={handleMagicLinkSend} className="space-y-4">
                  <div className="space-y-2">
                    <Input
                      type="email"
                      placeholder="name@example.com"
                      required
                      value={magicLinkEmail}
                      onChange={(e) => setMagicLinkEmail(e.target.value)}
                      className="bg-background text-foreground border-input"
                    />
                  </div>
                  <Button type="submit" className="w-full">
                    {isLoading ? "Sending Link..." : "Send Magic Link"}
                  </Button>
                </form>
              </TabsContent>

              <TabsContent value="register">
                <form onSubmit={handleEmailSignup} className="space-y-4">
                  <div className="space-y-2">
                    <Input
                      type="email"
                      placeholder="name@example.com"
                      required
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      className="bg-background text-foreground border-input"
                    />
                  </div>
                  <div className="space-y-2">
                    <Input
                      type="password"
                      placeholder="Create a password"
                      required
                      minLength={6}
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      className="bg-background text-foreground border-input"
                    />
                  </div>

                  <div className="space-y-2 text-sm text-muted-foreground">
                    <div className="flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4 text-green-500" />
                      <span>At least 6 characters</span>
                    </div>
                  </div>

                  <Button type="submit" className="w-full">
                    {isLoading ? "Creating account..." : "Create Account"}
                  </Button>
                </form>
              </TabsContent>
            </Tabs>

            <div className="relative my-6">
              <div className="absolute inset-0 flex items-center">
                <span className="w-full border-t border-border" />
              </div>
              <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-background px-2 text-muted-foreground">
                  Or continue with
                </span>
              </div>
            </div>

            <Button variant="outline" className="w-full border-input text-foreground hover:bg-muted" onClick={handleGoogleLogin} disabled={isLoading}>
              <svg className="mr-2 h-4 w-4" aria-hidden="true" focusable="false" data-prefix="fab" data-icon="google" role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 488 512">
                <path fill="currentColor" d="M488 261.8C488 403.3 391.1 504 248 504 110.8 504 0 393.2 0 256S110.8 8 248 8c66.8 0 123 24.5 166.3 64.9l-67.5 64.9C258.5 52.6 94.3 116.6 94.3 256c0 86.5 69.1 156.6 153.7 156.6 98.2 0 135-70.4 140.8-106.9H248v-85.3h236.1c2.3 12.7 3.9 24.9 3.9 41.4z"></path>
              </svg>
              Google
            </Button>
          </CardContent>
          <CardFooter className="flex justify-center">
            <p className="text-xs text-center text-muted-foreground">
              By clicking continue, you agree to our <Link to="/terms-of-service" className="underline hover:text-primary">Terms of Service</Link> and <Link to="/privacy-policy" className="underline hover:text-primary">Privacy Policy</Link>.
            </p>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
