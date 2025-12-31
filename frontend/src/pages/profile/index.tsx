

import { UserInfoCard } from "./components/user-info-card";
import { PasswordUpdateForm } from "./components/password-update-form";
import { AvatarUpload } from "./components/avatar-upload";
import { useAuth } from "@/hooks/use-auth"

const Profile = () => {
  const { user } = useAuth()

  if (!user) {
    return <div>Loading...</div>
  }

  return (
    <div className="p-6 space-y-6">
      <div className="max-w-7xl mx-auto">
        <h1 className="text-3xl font-bold mb-8">Profile</h1>

        <div className="grid gap-6 md:grid-cols-12">
          {/* Left Column: Avatar and User Info */}
          <div className="md:col-span-4 space-y-6">
            <AvatarUpload
              currentAvatarUrl={user.image || ""}
              username={user.username}
              userId={user.id || ""}
            />
            <UserInfoCard
              name={`${user.firstName} ${user.lastName}`}
              email={user.email}
              username={user.username}
              role={user.role}
            />
          </div>

          {/* Right Column: Password Update and other settings */}
          <div className="md:col-span-8 space-y-6">
            <PasswordUpdateForm userId={user.id || ""} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default Profile;
