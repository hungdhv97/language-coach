/**
 * Zod utilities for React Hook Form integration
 */

import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

// Export zodResolver for convenience
export { zodResolver };

// Common validation schemas
export const commonSchemas = {
  email: z.string().email('Email không hợp lệ'),
  password: z.string().min(8, 'Mật khẩu phải có ít nhất 8 ký tự'),
  required: (message?: string) => z.string().min(1, message || 'Trường này là bắt buộc'),
  optionalString: z.string().optional(),
  url: z.string().url('URL không hợp lệ'),
  positiveNumber: z.number().positive('Số phải lớn hơn 0'),
};

// Helper function to create form resolver
export function createFormResolver<T extends z.ZodTypeAny>(schema: T) {
  return zodResolver(schema);
}

