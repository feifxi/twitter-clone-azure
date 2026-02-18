'use client';

import { useState } from 'react';
import { XLogo } from '../XLogo';
import GoogleLoginBtn from '../GoogleLoginBtn';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { X } from 'lucide-react';

import { useUIStore } from '@/store/useUIStore';

// Removed Props interface as it's now controlled globally
export function SignInModal() {
  const isOpen = useUIStore((s) => s.isSignInModalOpen);
  const onClose = useUIStore((s) => s.closeSignInModal);
  const openSignUpModal = useUIStore((s) => s.openSignUpModal);

  const [identifier, setIdentifier] = useState('');
  const [step, setStep] = useState(1);

  // Reset state when closing
  const handleClose = () => {
    setStep(1);
    setIdentifier('');
    onClose();
  };

  const handleNext = () => {
    alert('สมัครไม่ได้หรอกโว้ย บังคับ Google login เท่านั้น 5555');
  }

  const handleForgot = () => {
    alert('ลืมรหัสผ่านหรอ ว้ายแย่จุง 5555');
  }

  const handleSignUpClick = () => {
    onClose();
    openSignUpModal();
  };

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent showCloseButton={false} className="sm:max-w-[600px] bg-black text-[#e7e9ea] border-[#2f3336] p-0 gap-0 overflow-hidden h-[650px] flex flex-col">
        {/* Header with Logo */}
        <div className="flex items-center h-[53px] px-4 shrink-0 relative">
          <button 
            onClick={handleClose}
            className="p-2 rounded-full hover:bg-[#eff3f41a] transition-colors -ml-2"
          >
            <X size={20} />
          </button>
          <div className="absolute left-1/2 -translate-x-1/2">
            <XLogo className="w-8 h-8 text-[#e7e9ea]" />
          </div>
        </div>

        {/* Content */}
        <div className="flex-1 flex flex-col px-20 pb-12 pt-4">
          <div className="w-full max-w-[364px] mx-auto flex flex-col flex-1">
            <DialogHeader className="mb-8">
              <DialogTitle className="text-3xl font-bold text-left mb-8">
                Sign in to X
              </DialogTitle>
              <DialogDescription className="sr-only">
                Sign in to your X account using Google or your credentials.
              </DialogDescription>
            </DialogHeader>

            <div className="flex flex-col gap-4 flex-1">
              <GoogleLoginBtn />

              <div className="flex items-center gap-2 my-4">
                <span className="flex-1 h-px bg-[#2f3336]" />
                <span className="text-[#e7e9ea] text-[15px]">or</span>
                <span className="flex-1 h-px bg-[#2f3336]" />
              </div>

              <div className="space-y-4">
                 <Input
                  type="text"
                  placeholder="Phone, email, or username"
                  className="bg-black border-[#333639] focus-visible:ring-1 focus-visible:ring-[#1d9bf0] h-[56px] text-lg"
                  value={identifier}
                  onChange={(e) => setIdentifier(e.target.value)}
                />
                <Button 
                    variant="default"
                    className="w-full rounded-full bg-[#eff3f4] text-black hover:bg-[#d7dbdc] h-[36px] font-bold text-[15px] cursor-pointer"
                    onClick={handleNext}
                >
                    Next
                </Button>
                 <Button 
                    variant="outline"
                    className="w-full rounded-full border-[#536471] text-white hover:bg-[#eff3f41a] h-[36px] font-bold text-[15px] cursor-pointer"
                    onClick={handleForgot}
                >
                    Forgot password?
                </Button>
              </div>
            </div>
            
             <p className="text-[#71767b] text-[15px] mt-12">
              Don't have an account? <span className="text-[#1d9bf0] cursor-pointer hover:underline" onClick={handleSignUpClick}>Sign up</span>
            </p>

          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
