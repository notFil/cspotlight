import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";

interface UserInfoCardProps {
  name: string;
  email: string;
  username: string;
  role: string;
}

export function UserInfoCard({ name, email, username, role }: UserInfoCardProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>User Information</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="name">Name</Label>
          <Input id="name" value={name} readOnly className="bg-muted" />
        </div>
        <div className="space-y-2">
          <Label htmlFor="email">Email</Label>
          <Input id="email" value={email} readOnly className="bg-muted" />
        </div>
        <div className="space-y-2">
          <Label htmlFor="username">Username</Label>
          <Input id="username" value={username} readOnly className="bg-muted" />
        </div>
        <div className="space-y-2">
          <Label htmlFor="role">Role</Label>
          <Input id="role" value={role} readOnly className="bg-muted" />
        </div>
      </CardContent>
    </Card>
  );
}
