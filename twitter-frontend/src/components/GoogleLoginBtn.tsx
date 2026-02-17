'use client';

import { GoogleLogin, CredentialResponse } from '@react-oauth/google';
import { useRouter } from 'next/navigation';
import { axiosInstance } from '@/api/axiosInstance';
import { useAuthStore } from '@/store/useAuthStore';
import type { AuthResponse } from '@/types';

import { useUIStore } from '@/store/useUIStore';

import { useQueryClient } from '@tanstack/react-query';

export default function GoogleLoginBtn() {
  const router = useRouter();
  const setAuth = useAuthStore((s) => s.setAuth);
  const closeSignInModal = useUIStore((s) => s.closeSignInModal);
  const queryClient = useQueryClient();

  const handleSuccess = async (credentialResponse: CredentialResponse) => {
    const googleToken = credentialResponse.credential;
    if (!googleToken) return;
    try {
      const { data } = await axiosInstance.post<AuthResponse>('/auth/google', {
        token: googleToken,
      });
      setAuth(data.accessToken, data.user);
      
      // Invalidate queries to refresh data with new auth state
      await queryClient.invalidateQueries({ queryKey: ['feeds'] });
      await queryClient.invalidateQueries({ queryKey: ['discovery'] });
      await queryClient.invalidateQueries({ queryKey: ['users'] });
      await queryClient.invalidateQueries({ queryKey: ['search'] });

      closeSignInModal();
      router.push('/');
    } catch (error) {
      console.error('Login Failed:', error);
      alert('Login failed. Please try again.');
    }
  };

  return (
    <div className="w-full [&_iframe]:min-h-[44px]! [&_iframe]:w-full!">
      <GoogleLogin
        onSuccess={handleSuccess}
        onError={() => console.log('Google Login Failed')}
        theme="filled_black"
        shape="pill"
        size="large"
        width="100%"
        text="signup_with"
      />
    </div>
  );
}
