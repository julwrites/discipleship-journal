import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { syncUser, updateUser } from "@/services/api";

export default function Settings() {
  const [username, setUsername] = useState("");
  const [bibleVersion, setBibleVersion] = useState("ESV");
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    syncUser()
      .then((user) => {
        if (user.username) setUsername(user.username);
        if (user.settings?.bible_version) setBibleVersion(user.settings.bible_version);
        setLoading(false);
      })
      .catch(console.error);
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
    <div className="p-4 md:p-8 max-w-2xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Settings</h1>

      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-1">Username</label>
          <Input
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="Enter unique username"
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">Bible Version</label>
          <Input
            value={bibleVersion}
            onChange={(e) => setBibleVersion(e.target.value)}
            placeholder="e.g. ESV, KJV, NIV"
          />
        </div>

        <div className="flex gap-4 mt-8">
          <Button onClick={handleSave}>Save</Button>
          <Button variant="outline" onClick={() => navigate("/")}>Cancel</Button>
        </div>
      </div>
    </div>
  );
}
