import { describe, it, expect } from 'vitest';
import { passwordSchema } from './validation';

describe('passwordSchema', () => {
  it('should pass for a valid complex password', () => {
    const validPassword = 'ComplexPassword123!';
    const result = passwordSchema.safeParse(validPassword);
    expect(result.success).toBe(true);
  });

  it('should fail for passwords strictly shorter than 8 characters', () => {
    const shortPassword = 'Short1!';
    const result = passwordSchema.safeParse(shortPassword);
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toContain('at least 8 characters');
    }
  });

  it('should fail for passwords longer than 64 characters', () => {
    const longPassword = 'A'.repeat(65) + '1!';
    const result = passwordSchema.safeParse(longPassword);
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toContain('at most 64 characters');
    }
  });

  it('should fail if missing lowercase letter', () => {
    const noLower = 'UPPERCASE123!';
    const result = passwordSchema.safeParse(noLower);
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toContain('lowercase letter');
    }
  });

  it('should fail if missing uppercase letter', () => {
    const noUpper = 'lowercase123!';
    const result = passwordSchema.safeParse(noUpper);
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toContain('uppercase letter');
    }
  });

  it('should fail if missing digit', () => {
    const noDigit = 'NoDigitHere!';
    const result = passwordSchema.safeParse(noDigit);
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toContain('digit');
    }
  });

  it('should fail if missing special character', () => {
    const noSpecial = 'NoSpecial123';
    const result = passwordSchema.safeParse(noSpecial);
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toContain('special character');
    }
  });
});
