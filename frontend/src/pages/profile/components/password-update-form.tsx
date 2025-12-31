import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { useChangePassword } from "@/hooks/use-users";
import { useAuth } from "@/hooks/use-auth";
import type { ChangePasswordRequest } from "@/types";

const passwordUpdateSchema = z
  .object({
    currentPassword: z.string().min(8, "Current password is required"),
    newPassword: z.string().min(8, "Password must be at least 8 characters long"),
    confirmNewPassword: z.string().min(8, "Please confirm your new password"),
  })
  .refine((data) => data.newPassword === data.confirmNewPassword, {
    message: "New passwords do not match",
    path: ["confirmNewPassword"],
  });

type PasswordUpdateFormValues = z.infer<typeof passwordUpdateSchema>;

interface PasswordUpdateFormProps {
  userId: string;
}

export function PasswordUpdateForm({ userId }: PasswordUpdateFormProps) {
  const { mutateAsync: changePassword } = useChangePassword();
  const { logout } = useAuth();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<PasswordUpdateFormValues>({
    resolver: zodResolver(passwordUpdateSchema),
    defaultValues: {
      currentPassword: "",
      newPassword: "",
      confirmNewPassword: "",
    },
  });

  const onSubmit = async (data: PasswordUpdateFormValues) => {
    if (!userId) {
      toast.error("User ID is missing");
      return;
    }

    const payload: ChangePasswordRequest = {
      currentPassword: data.currentPassword,
      newPassword: data.newPassword,
      confirmNewPassword: data.confirmNewPassword,
    };

    try {
      await changePassword({
        id: userId,
        data: payload,
      });
      toast.success("Password updated successfully");
      reset();
      await logout();
    } catch (error: any) {
      console.error(error);
      toast.error(error.response?.data?.message || "Failed to update password");
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Update Password</CardTitle>
        <CardDescription>
          Ensure your account is using a long, random password to stay secure.
        </CardDescription>
      </CardHeader>
      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-6">
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="current-password">Current Password</Label>
            <Input
              id="current-password"
              type="password"
              {...register("currentPassword")}
            />
            {errors.currentPassword && (
              <p className="text-sm text-red-500">{errors.currentPassword.message}</p>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="new-password">New Password</Label>
            <Input
              id="new-password"
              type="password"
              {...register("newPassword")}
            />
            {errors.newPassword && (
              <p className="text-sm text-red-500">{errors.newPassword.message}</p>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="confirm-password">Confirm New Password</Label>
            <Input
              id="confirm-password"
              type="password"
              {...register("confirmNewPassword")}
            />
            {errors.confirmNewPassword && (
              <p className="text-sm text-red-500">{errors.confirmNewPassword.message}</p>
            )}
          </div>
        </CardContent>
        <CardFooter className="border-t px-6 py-4">
          <Button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "Updating..." : "Save Password"}
          </Button>
        </CardFooter>
      </form>
    </Card>
  );
}
