import { z } from "zod";

export const PASSWORD_MIN_LENGTH = 8;
export const PASSWORD_MAX_LENGTH = 64;

export const PasswordRegex = {
  lowercase: /[a-z]/,
  uppercase: /[A-Z]/,
  digit: /[0-9]/,
  special: /[!@#$%^&*()_+\-=[\]{}|;:'",.<>?/]/,
};

export const passwordSchema = z.string()
  .min(PASSWORD_MIN_LENGTH, `Password must be at least ${PASSWORD_MIN_LENGTH} characters long`)
  .max(PASSWORD_MAX_LENGTH, `Password must be at most ${PASSWORD_MAX_LENGTH} characters long`)
  .regex(PasswordRegex.lowercase, "Password must contain at least one lowercase letter")
  .regex(PasswordRegex.uppercase, "Password must contain at least one uppercase letter")
  .regex(PasswordRegex.digit, "Password must contain at least one digit")
  .regex(PasswordRegex.special, "Password must contain at least one special character");
