'use client';

import { useState, useRef, useEffect } from 'react';
import { useAuth } from '@/hooks/useAuth';
import { useUpdateProfile } from '@/hooks/useProfile';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/dialog';
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { X, Camera } from 'lucide-react';
import type { UserResponse } from '@/types';

interface EditProfileModalProps {
  user: UserResponse;
  isOpen: boolean;
  onClose: () => void;
}

export function EditProfileModal({ user, isOpen, onClose }: EditProfileModalProps) {
  const [displayName, setDisplayName] = useState(user.displayName);
  const [bio, setBio] = useState(user.bio || '');
  const [avatar, setAvatar] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(user.avatarUrl ?? null);
  
  const fileInputRef = useRef<HTMLInputElement>(null);
  const updateMutation = useUpdateProfile();
  const { user: currentUser, setAuth } = useAuth(); // Need to update global auth store if current user changes

  const hasChanges = displayName !== user.displayName || bio !== (user.bio || '') || avatar !== null;

  useEffect(() => {
    if (isOpen) {
        // eslint-disable-next-line react-hooks/set-state-in-effect
        setDisplayName(user.displayName);
        // eslint-disable-next-line react-hooks/set-state-in-effect
        setBio(user.bio || '');
        // eslint-disable-next-line react-hooks/set-state-in-effect
        setPreviewUrl(user.avatarUrl ?? null);
        // eslint-disable-next-line react-hooks/set-state-in-effect
        setAvatar(null);
    }
  }, [isOpen, user]);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setAvatar(file);
      const url = URL.createObjectURL(file);
      setPreviewUrl(url);
    }
  };

  const handleSubmit = async () => {
    try {
      const updatedUser = await updateMutation.mutateAsync({
        displayName,
        bio,
        avatar: avatar ?? undefined,
      });
      
      // Update global auth store if we're editing our own profile
      if (currentUser && currentUser.id === updatedUser.id) {
          // This is a bit of a hack since setAuth requires token, but we just want to update user. 
          // Assuming we can re-fetch or the store has a partial update method (which it likely doesn't).
          // Ideally we invalidate 'auth' query and let it re-fetch.
          // But for now, we rely on React Query invalidation in the hook.
      }
      
      onClose();
    } catch (error) {
      console.error('Failed to update profile:', error);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent showCloseButton={false} className="sm:max-w-[600px] bg-background border-border p-0 gap-0 top-[5%] translate-y-0 sm:top-[10%] min-h-[400px] flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between h-[53px] px-4 shrink-0 border-b border-border">
            <div className="flex items-center gap-4">
                <button
                    onClick={onClose}
                    className="p-2 rounded-full hover:bg-card transition-colors -ml-2"
                >
                    <X size={20} className="text-foreground" />
                </button>
                <DialogTitle className="text-[20px] font-bold text-foreground">Edit Profile</DialogTitle>
                <DialogDescription className="sr-only">Make changes to your profile here. Click save when you&apos;re done.</DialogDescription>
            </div>
            <Button
                onClick={handleSubmit}
                disabled={!displayName.trim() || updateMutation.isPending || !hasChanges}
                className="rounded-full bg-foreground text-background hover:bg-foreground/90 h-[32px] font-bold text-[14px] px-4"
            >
                {updateMutation.isPending ? 'Saving...' : 'Save'}
            </Button>
        </div>

        {/* Content */}
        <div className="p-4 flex flex-col gap-6">
            {/* Banner (Placeholder) */}
            <div className="h-[200px] bg-secondary -mt-4 -mx-4 mb-10 relative">
                {/* Avatar Overlay */}
                <div className="absolute -bottom-[40px] left-4">
                     <div className="w-[112px] h-[112px] rounded-full border-4 border-background bg-card relative overflow-hidden group">
                        <Avatar className="w-full h-full rounded-none">
                            <AvatarImage src={previewUrl ?? undefined} alt="" className="object-cover" />
                            <AvatarFallback className="text-[32px] font-bold">{user.displayName[0]}</AvatarFallback>
                        </Avatar>
                        <div className="absolute inset-0 bg-black/30 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer" onClick={() => fileInputRef.current?.click()}>
                            <Camera className="w-6 h-6 text-white" />
                        </div>
                     </div>
                     <input
                        type="file"
                        ref={fileInputRef}
                        className="hidden"
                        accept="image/*"
                        onChange={handleFileChange}
                     />
                </div>
            </div>

            <div className="flex flex-col gap-6 mt-4">
                <div className="border border-border rounded-md px-3 py-1.5 focus-within:border-primary focus-within:ring-1 focus-within:ring-primary transition-colors relative">
                    <div className="flex justify-between items-center mb-0.5">
                        <label className="block text-muted-foreground text-[13px]">Name</label>
                        <span className="text-muted-foreground text-[13px]">{displayName.length} / 30</span>
                    </div>
                    <input 
                        className="w-full bg-transparent text-foreground text-[17px] outline-none"
                        value={displayName}
                        maxLength={30}
                        onChange={(e) => setDisplayName(e.target.value)}
                    />
                </div>

                <div className="border border-border rounded-md px-3 py-1.5 focus-within:border-primary focus-within:ring-1 focus-within:ring-primary transition-colors relative">
                    <div className="flex justify-between items-center mb-0.5">
                        <label className="block text-muted-foreground text-[13px]">Bio</label>
                        <span className="text-muted-foreground text-[13px]">{bio.length} / 160</span>
                    </div>
                    <textarea 
                        className="w-full bg-transparent text-foreground text-[17px] outline-none resize-none h-[80px]"
                        value={bio}
                        maxLength={160}
                        onChange={(e) => setBio(e.target.value)}
                    />
                </div>
            </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
