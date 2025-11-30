

import { UserInfoCard } from "./components/user-info-card";
import { PasswordUpdateForm } from "./components/password-update-form";
import { AvatarUpload } from "./components/avatar-upload";

const Profile = () => {
  // Mock User Data
  const user = {
    name: "Phil Ogb",
    email: "phil.ogb@example.com",
    username: "philogb",
    avatarUrl: "https://github.com/shadcn.png", // Placeholder image
  };

  return (
    <div className="p-6 space-y-6">
      <div className="max-w-7xl mx-auto">
        <h1 className="text-3xl font-bold mb-8">Profile</h1>

        <div className="grid gap-6 md:grid-cols-12">
          {/* Left Column: Avatar and User Info */}
          <div className="md:col-span-4 space-y-6">
            <AvatarUpload
              currentAvatarUrl={user.avatarUrl}
              username={user.username}
            />
            <UserInfoCard
              name={user.name}
              email={user.email}
              username={user.username}
            />
          </div>

          {/* Right Column: Password Update and other settings */}
          <div className="md:col-span-8 space-y-6">
            <PasswordUpdateForm />
          </div>
        </div>
      </div>
    </div>
  );
};

export default Profile;
