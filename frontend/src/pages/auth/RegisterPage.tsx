/**
 * Register Page Component
 */

import { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useForm, useWatch } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useMutation, useQuery } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { userEndpoints } from '@/entities/user/api/user.endpoints';
import { useAuthStore } from '@/shared/store/useAuthStore';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Check, X, Loader2, AlertCircle } from 'lucide-react';
import { useDebounce } from '@/hooks/useDebounce';
import type { ApiError } from '@/shared/api/http-client';

const registerSchema = z.object({
  display_name: z.string().max(100).optional().or(z.literal('')),
  username: z.string().min(3, 'Username phải có ít nhất 3 ký tự').optional().or(z.literal('')),
  email: z.string().email('Email không hợp lệ').optional().or(z.literal('')),
  password: z.string().min(6, 'Mật khẩu phải có ít nhất 6 ký tự'),
  confirmPassword: z.string(),
}).refine((data) => data.email || data.username, {
  message: 'Vui lòng nhập email hoặc username',
  path: ['email'],
}).refine((data) => data.password === data.confirmPassword, {
  message: 'Mật khẩu xác nhận không khớp',
  path: ['confirmPassword'],
});

type RegisterFormData = z.infer<typeof registerSchema>;

type ValidationState = 'idle' | 'checking' | 'valid' | 'invalid';

export default function RegisterPage() {
  const navigate = useNavigate();
  const login = useAuthStore((state) => state.login);
  const [error, setError] = useState<string | null>(null);
  const [usernameValidation, setUsernameValidation] = useState<ValidationState>('idle');
  const [emailValidation, setEmailValidation] = useState<ValidationState>('idle');

  const form = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      display_name: '',
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
    },
  });

  const username = useWatch({ control: form.control, name: 'username' });
  const email = useWatch({ control: form.control, name: 'email' });
  const debouncedUsername = useDebounce(username, 500);
  const debouncedEmail = useDebounce(email, 500);

  // Check username availability
  const { data: usernameCheck } = useQuery({
    queryKey: ['checkUsername', debouncedUsername],
    queryFn: () => {
      if (!debouncedUsername) throw new Error('Username is required');
      return userEndpoints.checkUsernameAvailability(debouncedUsername);
    },
    enabled: !!debouncedUsername && debouncedUsername.length >= 3,
    retry: false,
  });

  // Check email availability
  const { data: emailCheck } = useQuery({
    queryKey: ['checkEmail', debouncedEmail],
    queryFn: () => {
      if (!debouncedEmail) throw new Error('Email is required');
      return userEndpoints.checkEmailAvailability(debouncedEmail);
    },
    enabled: !!debouncedEmail && /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(debouncedEmail),
    retry: false,
  });

  useEffect(() => {
    if (!debouncedUsername || debouncedUsername.length < 3) {
      setUsernameValidation('idle');
      return;
    }

    if (usernameCheck) {
      setUsernameValidation(usernameCheck.available ? 'valid' : 'invalid');
    } else {
      setUsernameValidation('checking');
    }
  }, [debouncedUsername, usernameCheck]);

  useEffect(() => {
    if (!debouncedEmail || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(debouncedEmail)) {
      setEmailValidation('idle');
      return;
    }

    if (emailCheck) {
      setEmailValidation(emailCheck.available ? 'valid' : 'invalid');
    } else {
      setEmailValidation('checking');
    }
  }, [debouncedEmail, emailCheck]);

  const registerMutation = useMutation({
    mutationFn: async (data: RegisterFormData) => {
      const registerData: { email?: string; username?: string; password: string; display_name?: string } = {
        password: data.password,
      };
      if (data.email) registerData.email = data.email;
      if (data.username) registerData.username = data.username;
      if (data.display_name) registerData.display_name = data.display_name;
      return userEndpoints.register(registerData);
    },
    onSuccess: async (response) => {
      // Auto login after registration
      const loginData: { email?: string; username?: string; password: string } = {
        password: form.getValues('password'),
      };
      if (response.email) loginData.email = response.email;
      if (response.username) loginData.username = response.username;
      
      try {
        const loginResponse = await userEndpoints.login(loginData);
        login(loginResponse.token, {
          id: loginResponse.user_id,
          email: loginResponse.email,
          username: response.username,
          created_at: '',
          updated_at: '',
          is_active: true,
        });
        navigate('/games');
      } catch {
        setError('Đăng ký thành công nhưng đăng nhập thất bại. Vui lòng đăng nhập thủ công.');
      }
    },
    onError: (err: Error | ApiError) => {
      const errorMessage = err instanceof Error 
        ? err.message 
        : (err as ApiError).message || 'Đăng ký thất bại. Vui lòng thử lại.';
      setError(errorMessage);
    },
  });

  const onSubmit = (data: RegisterFormData) => {
    setError(null);
    registerMutation.mutate(data);
  };

  const getUsernameIcon = () => {
    if (usernameValidation === 'checking') {
      return <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />;
    }
    if (usernameValidation === 'valid') {
      return <Check className="h-4 w-4 text-green-500" />;
    }
    if (usernameValidation === 'invalid') {
      return <X className="h-4 w-4 text-red-500" />;
    }
    return null;
  };

  const getEmailIcon = () => {
    if (emailValidation === 'checking') {
      return <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />;
    }
    if (emailValidation === 'valid') {
      return <Check className="h-4 w-4 text-green-500" />;
    }
    if (emailValidation === 'invalid') {
      return <X className="h-4 w-4 text-red-500" />;
    }
    return null;
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-4 bg-gradient-to-br from-background to-muted/20">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-2xl">Đăng Ký</CardTitle>
          <CardDescription>
            Tạo tài khoản mới để bắt đầu học từ vựng
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              {error && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>{error}</AlertDescription>
                </Alert>
              )}

              <FormField
                control={form.control}
                name="display_name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Tên hiển thị (tùy chọn)</FormLabel>
                    <FormControl>
                      <input
                        type="text"
                        className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                        placeholder="Tên hiển thị"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="username"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Username</FormLabel>
                    <FormControl>
                      <div className="relative">
                        <input
                          type="text"
                          className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 pr-10 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                          placeholder="username"
                          {...field}
                        />
                        {username && (
                          <div className="absolute right-3 top-1/2 -translate-y-1/2">
                            {getUsernameIcon()}
                          </div>
                        )}
                      </div>
                    </FormControl>
                    <FormMessage />
                    {usernameValidation === 'invalid' && (
                      <p className="text-sm text-red-500">Username đã tồn tại</p>
                    )}
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <div className="relative">
                        <input
                          type="email"
                          className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 pr-10 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                          placeholder="your@email.com"
                          {...field}
                        />
                        {email && (
                          <div className="absolute right-3 top-1/2 -translate-y-1/2">
                            {getEmailIcon()}
                          </div>
                        )}
                      </div>
                    </FormControl>
                    <FormMessage />
                    {emailValidation === 'invalid' && (
                      <p className="text-sm text-red-500">Email đã tồn tại</p>
                    )}
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Mật khẩu</FormLabel>
                    <FormControl>
                      <input
                        type="password"
                        className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                        placeholder="••••••••"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="confirmPassword"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Xác nhận mật khẩu</FormLabel>
                    <FormControl>
                      <input
                        type="password"
                        className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                        placeholder="••••••••"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <Button
                type="submit"
                className="w-full"
                disabled={registerMutation.isPending || usernameValidation === 'invalid' || emailValidation === 'invalid'}
              >
                {registerMutation.isPending ? 'Đang đăng ký...' : 'Đăng Ký'}
              </Button>

              <div className="text-center text-sm">
                <span className="text-muted-foreground">Đã có tài khoản? </span>
                <Link to="/auth/login" className="text-primary hover:underline">
                  Đăng nhập ngay
                </Link>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
