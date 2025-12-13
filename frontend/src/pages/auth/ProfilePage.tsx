/**
 * Profile Page Component
 */

import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { userEndpoints } from '@/entities/user/api/user.endpoints';
import { useAuthStore } from '@/shared/store/useAuthStore';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { LogOut, AlertCircle } from 'lucide-react';

const profileSchema = z.object({
  display_name: z.string().max(100).optional().or(z.literal('')),
  avatar_url: z.string().url().max(500).optional().or(z.literal('')),
  birth_day: z.string().regex(/^\d{4}-\d{2}-\d{2}$/).optional().or(z.literal('')),
  bio: z.string().optional(),
});

type ProfileFormData = z.infer<typeof profileSchema>;

export default function ProfilePage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { user, logout } = useAuthStore();

  const { data: profile, isLoading, error } = useQuery({
    queryKey: ['userProfile'],
    queryFn: () => userEndpoints.getProfile(),
  });

  const form = useForm<ProfileFormData>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      display_name: '',
      avatar_url: '',
      birth_day: '',
      bio: '',
    },
  });

  useEffect(() => {
    if (profile) {
      form.reset({
        display_name: profile.display_name || '',
        avatar_url: profile.avatar_url || '',
        birth_day: profile.birth_day || '',
        bio: profile.bio || '',
      });
    }
  }, [profile, form]);

  const updateMutation = useMutation({
    mutationFn: (data: ProfileFormData) => {
      const updateData: any = {};
      if (data.display_name) updateData.display_name = data.display_name;
      if (data.avatar_url) updateData.avatar_url = data.avatar_url;
      if (data.birth_day) updateData.birth_day = data.birth_day;
      if (data.bio !== undefined) updateData.bio = data.bio;
      return userEndpoints.updateProfile(updateData);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['userProfile'] });
      alert('Cập nhật profile thành công!');
    },
    onError: (err: any) => {
      alert(err.message || 'Cập nhật profile thất bại');
    },
  });

  const onSubmit = (data: ProfileFormData) => {
    updateMutation.mutate(data);
  };

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div>Đang tải...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen p-4 bg-gradient-to-br from-background to-muted/20">
      <div className="max-w-2xl mx-auto space-y-6">
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-2xl">Quản Lý Profile</CardTitle>
                <CardDescription>
                  Cập nhật thông tin cá nhân của bạn
                </CardDescription>
              </div>
              <Button variant="outline" onClick={handleLogout}>
                <LogOut className="h-4 w-4 mr-2" />
                Đăng Xuất
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            {error && (
              <Alert variant="destructive" className="mb-4">
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>Không thể tải thông tin profile</AlertDescription>
              </Alert>
            )}

            <div className="mb-6 p-4 bg-muted rounded-lg">
              <div className="text-sm text-muted-foreground mb-2">Thông tin tài khoản</div>
              <div className="space-y-1">
                <div><strong>ID:</strong> {user?.id}</div>
                {user?.email && <div><strong>Email:</strong> {user.email}</div>}
                {user?.username && <div><strong>Username:</strong> {user.username}</div>}
              </div>
            </div>

            <Form {...form}>
              <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                <FormField
                  control={form.control}
                  name="display_name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Tên hiển thị</FormLabel>
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
                  name="avatar_url"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>URL Avatar</FormLabel>
                      <FormControl>
                        <input
                          type="url"
                          className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                          placeholder="https://example.com/avatar.jpg"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="birth_day"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Ngày sinh</FormLabel>
                      <FormControl>
                        <input
                          type="date"
                          className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="bio"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Giới thiệu</FormLabel>
                      <FormControl>
                        <textarea
                          className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                          placeholder="Giới thiệu về bản thân..."
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
                  disabled={updateMutation.isPending}
                >
                  {updateMutation.isPending ? 'Đang cập nhật...' : 'Cập Nhật Profile'}
                </Button>
              </form>
            </Form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
