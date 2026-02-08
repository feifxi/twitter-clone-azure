import GoogleLoginBtn from '@/components/GoogleLoginBtn';

export default function LoginPage() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gray-50">
      <div className="w-full max-w-md space-y-8 rounded-xl bg-white p-10 shadow-lg">
        <div className="text-center">
          <h2 className="mt-6 text-3xl font-bold tracking-tight text-gray-900">
            Sign in to your account
          </h2>
          <p className="mt-2 text-sm text-gray-600">
            Welcome back to the Twitter Clone
          </p>
        </div>

        <div className="mt-8 flex justify-center">
          {/* The Magic Button */}
          <GoogleLoginBtn />
        </div>
      </div>
    </div>
  );
}