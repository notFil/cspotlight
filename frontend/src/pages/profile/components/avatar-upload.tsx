import { useState, useRef } from "react";
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { toast } from "sonner";
import { Upload } from "lucide-react";
import { useChangeImage } from "@/hooks/use-users";

interface AvatarUploadProps {
  currentAvatarUrl?: string;
  username: string;
  userId: string;
}

export function AvatarUpload({ currentAvatarUrl, username, userId }: AvatarUploadProps) {
  const [previewUrl, setPreviewUrl] = useState<string | null>(currentAvatarUrl || null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const changeImageMutation = useChangeImage();

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      if (file.size > 2 * 1024 * 1024) {
        toast.error("File size must be less than 2MB");
        return;
      }

      setSelectedFile(file);
      const reader = new FileReader();
      reader.onloadend = () => {
        setPreviewUrl(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleUploadClick = () => {
    fileInputRef.current?.click();
  };

  const handleSave = () => {
    if (!selectedFile) return;

    changeImageMutation.mutate(
      { id: userId, data: selectedFile },
      {
        onSuccess: () => {
          toast.success("Avatar updated successfully");
          setSelectedFile(null);
        },
        onError: (error: any) => {
          console.error("Failed to update avatar", error);
          toast.error(error.response?.data?.message || "Failed to update avatar");
        },
      }
    );
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Avatar</CardTitle>
        <CardDescription>
          This is your avatar. Click on the avatar to upload a custom one from your files.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col items-center justify-center space-y-4">
        <div className="relative group cursor-pointer" onClick={handleUploadClick}>
          <Avatar className="h-32 w-32">
            <AvatarImage src={previewUrl || ""} alt={username} />
            <AvatarFallback className="text-4xl">{username.slice(0, 2).toUpperCase()}</AvatarFallback>
          </Avatar>
          <div className="absolute inset-0 flex items-center justify-center bg-black/60 rounded-full opacity-0 group-hover:opacity-100 transition-opacity">
            <Upload className="h-8 w-8 text-white" />
          </div>
        </div>
        <input
          type="file"
          ref={fileInputRef}
          className="hidden"
          accept="image/*"
          onChange={handleFileChange}
        />
        <div className="text-sm text-muted-foreground">
          Click to upload. JPG, GIF or PNG. 2MB max.
        </div>
      </CardContent>
      <CardFooter className="border-t px-6 py-4 flex justify-between">
        <Button
          variant="outline"
          onClick={() => {
            setPreviewUrl(currentAvatarUrl || null);
            setSelectedFile(null);
          }}
        >
          Reset
        </Button>
        <Button
          onClick={handleSave}
          disabled={changeImageMutation.isPending || !selectedFile}
        >
          {changeImageMutation.isPending ? "Saving..." : "Save Avatar"}
        </Button>
      </CardFooter>
    </Card>
  );
}
