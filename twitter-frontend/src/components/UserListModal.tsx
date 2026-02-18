import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { useUserFollowers, useUserFollowing } from '@/hooks/useProfile';
import { FollowButton } from '@/components/FollowButton';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import Link from 'next/link';
import { useInView } from 'react-intersection-observer';
import { useEffect } from 'react';
import { useAuth } from '@/hooks/useAuth';

interface UserListModalProps {
  userId: number | null;
  type: 'followers' | 'following' | null;
  isOpen: boolean;
  onClose: () => void;
}

export function UserListModal({ userId, type, isOpen, onClose }: UserListModalProps) {
  const { ref, inView } = useInView();
  const { user: currentUser } = useAuth();

  const isFollowers = type === 'followers';
  const queryHook = isFollowers ? useUserFollowers : useUserFollowing;
  
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, isLoading } = queryHook(userId);

  useEffect(() => {
    if (inView && hasNextPage) {
      fetchNextPage();
    }
  }, [inView, hasNextPage, fetchNextPage]);

  const users = data?.pages.flatMap((page) => page.content) ?? [];

  if (!type) return null;

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[480px] p-0 gap-0 bg-black border-[#2f3336] text-[#e7e9ea] h-[600px] flex flex-col">
        <DialogHeader className="px-4 py-3 border-b border-[#2f3336]">
          <DialogTitle className="text-xl font-bold">
            {isFollowers ? 'Followers' : 'Following'}
          </DialogTitle>
        </DialogHeader>
        
        <div className="flex-1 overflow-y-auto min-h-0">
          {isLoading ? (
             <div className="p-4 text-center text-[#71767b]">Loading...</div>
          ) : users.length === 0 ? (
            <div className="p-8 text-center text-[#71767b]">
              <p className="font-bold text-lg mb-2">
                {isFollowers ? 'Looking empty here' : 'Be the first to follow'}
              </p>
              <p className="text-sm">
                {isFollowers 
                  ? 'This user has no followers yet.' 
                  : 'This user is not following anyone yet.'}
              </p>
            </div>
          ) : (
            <div className="flex flex-col">
              {users.map((user) => (
                <div key={user.id} className="flex items-center gap-3 px-4 py-3 hover:bg-[#eff3f41a] transition-colors">
                  <Link href={`/${user.username}`} onClick={onClose} className="shrink-0">
                    <Avatar className="w-10 h-10">
                      <AvatarImage src={user.avatarUrl ?? undefined} />
                      <AvatarFallback>{user.displayName[0]}</AvatarFallback>
                    </Avatar>
                  </Link>
                  
                  <div className="flex-1 min-w-0">
                    <Link href={`/${user.username}`} onClick={onClose} className="block group">
                      <div className="font-bold text-[#e7e9ea] truncate group-hover:underline">
                        {user.displayName}
                      </div>
                      <div className="text-[#71767b] truncate">@{user.username}</div>
                    </Link>
                    {user.bio && <p className="text-[#e7e9ea] text-[14px] truncate mt-0.5">{user.bio}</p>}
                  </div>

                  {currentUser?.id !== user.id && (
                    <FollowButton userId={user.id} isFollowing={user.followedByMe} />
                  )}
                </div>
              ))}
              
              {/* Infinite scroll loader */}
              <div ref={ref} className="h-10 flex items-center justify-center p-4">
                {isFetchingNextPage && <div className="w-6 h-6 border-2 border-[#1d9bf0] border-t-transparent rounded-full animate-spin" />}
              </div>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
