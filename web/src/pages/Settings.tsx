import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { syncUser, updateUser } from "@/services/api";
import { useTheme } from "@/components/ThemeProvider";
import { Moon, Sun, Laptop, Palette } from "lucide-react";

export default function Settings() {
  const [username, setUsername] = useState("");
  const [bibleVersion, setBibleVersion] = useState("ESV");
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const { theme, setTheme, mode, setMode } = useTheme();

  useEffect(() => {
    syncUser()
      .then((user) => {
        if (user.username) setUsername(user.username);
        if (user.settings?.bible_version) setBibleVersion(user.settings.bible_version);
      })
      .catch((error) => {
        console.error("Failed to sync user", error);
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  const handleSave = async () => {
    try {
      await updateUser({ username, bible_version: bibleVersion });
      navigate("/");
    } catch (error) {
      console.error("Failed to update settings", error);
      alert("Failed to update settings");
    }
  };

  if (loading) return <div className="p-8">Loading...</div>;

  return (
    <div className="p-4 md:p-8 max-w-2xl mx-auto space-y-8">
      <div>
        <h1 className="text-3xl font-bold">Settings</h1>
        <p className="text-muted-foreground">Manage your account settings and preferences.</p>
      </div>

      <div className="space-y-6">
        {/* Appearance Section */}
        <div className="border rounded-lg p-6 space-y-6">
            <h2 className="text-xl font-semibold flex items-center gap-2">
                <Palette className="h-5 w-5" />
                Appearance
            </h2>

            <div className="space-y-4">
                <div>
                    <label className="text-sm font-medium mb-2 block">Theme</label>
                    <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                        <Button
                            variant={theme === 'default' ? 'default' : 'outline'}
                            className="justify-start h-auto py-3 px-4"
                            onClick={() => setTheme('default')}
                        >
                            <div className="flex flex-col items-start gap-1">
                                <span className="font-semibold">Default</span>
                                <span className="text-xs opacity-70">Slate & Blue</span>
                            </div>
                        </Button>
                        <Button
                            variant={theme === 'serene' ? 'default' : 'outline'}
                            className="justify-start h-auto py-3 px-4"
                            onClick={() => setTheme('serene')}
                        >
                            <div className="flex flex-col items-start gap-1">
                                <span className="font-semibold">Serene</span>
                                <span className="text-xs opacity-70">Stone & Green</span>
                            </div>
                        </Button>
                        <Button
                            variant={theme === 'elegant' ? 'default' : 'outline'}
                            className="justify-start h-auto py-3 px-4"
                            onClick={() => setTheme('elegant')}
                        >
                             <div className="flex flex-col items-start gap-1">
                                <span className="font-semibold">Elegant</span>
                                <span className="text-xs opacity-70">Zinc & Violet</span>
                            </div>
                        </Button>
                    </div>
                </div>

                <div>
                    <label className="text-sm font-medium mb-2 block">Mode</label>
                    <div className="flex gap-2">
                        <Button
                            variant={mode === 'light' ? 'default' : 'outline'}
                            size="sm"
                            onClick={() => setMode('light')}
                        >
                            <Sun className="h-4 w-4 mr-2" />
                            Light
                        </Button>
                        <Button
                            variant={mode === 'dark' ? 'default' : 'outline'}
                            size="sm"
                            onClick={() => setMode('dark')}
                        >
                            <Moon className="h-4 w-4 mr-2" />
                            Dark
                        </Button>
                        <Button
                            variant={mode === 'system' ? 'default' : 'outline'}
                            size="sm"
                            onClick={() => setMode('system')}
                        >
                            <Laptop className="h-4 w-4 mr-2" />
                            System
                        </Button>
                    </div>
                </div>
            </div>
        </div>

        {/* Profile Section */}
        <div className="border rounded-lg p-6 space-y-4">
            <h2 className="text-xl font-semibold">Profile</h2>
            <div>
            <label className="block text-sm font-medium mb-1">Username</label>
            <Input
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="Enter unique username"
                className="max-w-md"
            />
            </div>

            <div>
            <label className="block text-sm font-medium mb-1">Bible Version</label>
            <Input
                value={bibleVersion}
                onChange={(e) => setBibleVersion(e.target.value)}
                placeholder="e.g. ESV, KJV, NIV"
                 className="max-w-md"
            />
            </div>
        </div>

        <div className="flex gap-4 pt-4">
          <Button onClick={handleSave}>Save Changes</Button>
          <Button variant="outline" onClick={() => navigate("/")}>Cancel</Button>
        </div>
      </div>
    </div>
  );
}
