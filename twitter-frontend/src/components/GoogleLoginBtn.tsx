'use client';

import { GoogleLogin, CredentialResponse } from '@react-oauth/google';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api'; // Import your configured axios instance

export default function GoogleLoginBtn() {
  const router = useRouter();

  const handleSuccess = async (credentialResponse: CredentialResponse) => {
    try {
      // 1. Get Google Token
      const googleToken = credentialResponse.credential;
    
      console.log(googleToken)
      return

      // 2. Exchange for Your Backend Tokens
      const res = await api.post('/auth/google', {
        token: googleToken,
      });

      const { accessToken, refreshToken, user } = res.data;

      // 3. Save to Storage (For MVP, localStorage is fine)
      localStorage.setItem('accessToken', accessToken);
      localStorage.setItem('refreshToken', refreshToken);
      localStorage.setItem('user', JSON.stringify(user)); // Optional: Store user info

      // 4. Redirect to Feed
      router.push('/feed');
      
    } catch (error) {
      console.error('Login Failed:', error);
      alert('Login failed. Please try again.');
    }
  };

  return (
    <div className="flex flex-col items-center gap-4">
      <GoogleLogin
        onSuccess={handleSuccess}
        onError={() => console.log('Google Login Failed')}
        useOneTap // Optional: Shows the popup automatically
        theme="filled_blue"
        shape="pill"
      />
    </div>
  );
}