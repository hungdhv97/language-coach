/**
 * Login Page Component
 */

import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useMutation } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { userEndpoints } from '@/entities/user/api/user.endpoints';
import { useAuthStore } from '@/shared/store/useAuthStore';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Checkbox } from '@/components/ui/checkbox';
import { AlertCircle } from 'lucide-react';

const loginSchema = z.object({
  emailOrUsername: z.string().min(1, 'Vui lòng nhập email hoặc username'),
  password: z.string().min(6, 'Mật khẩu phải có ít nhất 6 ký tự'),
  rememberMe: z.boolean(),
});

type LoginFormData = z.infer<typeof loginSchema>;

export default function LoginPage() {
  const navigate = useNavigate();
  const login = useAuthStore((state) => state.login);
  const [error, setError] = useState<string | null>(null);

  const form = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      emailOrUsername: '',
      password: '',
      rememberMe: false,
    },
  });

  const loginMutation = useMutation({
    mutationFn: async (data: LoginFormData) => {
      // Determine if input is email or username
      const isEmail = data.emailOrUsername.includes('@');
      const loginData: { email?: string; username?: string; password: string } = {
        password: data.password,
      };
      
      if (isEmail) {
        loginData.email = data.emailOrUsername;
      } else {
        loginData.username = data.emailOrUsername;
      }
      
      return userEndpoints.login(loginData);
    },
    onSuccess: (response) => {
      login(response.token, {
        id: response.user_id,
        email: response.email,
        username: response.username,
        created_at: '',
        updated_at: '',
        is_active: true,
      });
      navigate('/games');
    },
    onError: (err: any) => {
      // Extract error message from Error instance or ApiError object
      const errorMessage = err?.message || err?.error?.message || 'Đăng nhập thất bại. Vui lòng kiểm tra lại thông tin.';
      setError(errorMessage);
    },
  });

  const onSubmit = (data: LoginFormData) => {
    setError(null);
    loginMutation.mutate(data);
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-4 bg-gradient-to-br from-background to-muted/20">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-2xl">Đăng Nhập</CardTitle>
          <CardDescription>
            Đăng nhập để chơi game và quản lý profile
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
                name="emailOrUsername"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email hoặc Username</FormLabel>
                    <FormControl>
                      <input
                        type="text"
                        className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                        placeholder="email@example.com hoặc username"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
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

              <div className="flex items-center justify-between">
                <FormField
                  control={form.control}
                  name="rememberMe"
                  render={({ field }) => (
                    <FormItem className="flex flex-row items-start space-x-3 space-y-0">
                      <FormControl>
                        <Checkbox
                          checked={field.value}
                          onCheckedChange={field.onChange}
                        />
                      </FormControl>
                      <div className="space-y-1 leading-none">
                        <FormLabel className="text-sm font-normal cursor-pointer">
                          Ghi nhớ đăng nhập
                        </FormLabel>
                      </div>
                    </FormItem>
                  )}
                />
                <Link to="#" className="text-sm text-primary hover:underline">
                  Quên mật khẩu?
                </Link>
              </div>

              <Button
                type="submit"
                className="w-full"
                disabled={loginMutation.isPending}
              >
                {loginMutation.isPending ? 'Đang đăng nhập...' : 'Đăng Nhập'}
              </Button>

              <div className="text-center text-sm">
                <span className="text-muted-foreground">Chưa có tài khoản? </span>
                <Link to="/auth/register" className="text-primary hover:underline">
                  Đăng ký ngay
                </Link>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
