'use client';

import { useSuggestedUsers } from '@/hooks/useDiscovery';
import { FeedList } from '@/components/FeedList'; // We might need a UserList component, but for now we'll map manually or use a new component.
// Reusing logic from search/page.tsx for user list is better.
import Link from 'next/link';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/hooks/useAuth';
import { FollowButton } from '@/components/FollowButton';

import { useRouter } from 'next/navigation';
import { ArrowLeft } from 'lucide-react';

export default function ConnectPeoplePage() {
  const router = useRouter();
  const { data, isLoading } = useSuggestedUsers(20); // Fetch more users
  const users = data?.content ?? [];
  const { user: currentUser } = useAuth();

  return (
    <div className="min-h-screen">
      <div className="sticky top-0 z-30 bg-black/60 backdrop-blur-md border-b border-[#2f3336] px-4 py-3 flex items-center gap-4">
        <div 
            onClick={() => router.back()} 
            className="p-2 rounded-full hover:bg-white/10 cursor-pointer transition-colors -ml-2"
        >
            <ArrowLeft className="w-5 h-5 text-white" />
        </div>
        <h1 className="text-xl font-bold">Connect</h1>
      </div>

      <div className="flex flex-col">
        {isLoading ? (
             <div className="p-4 space-y-4">
                {[1, 2, 3, 4, 5].map((i) => (
                    <div key={i} className="flex items-center gap-3">
                         <div className="w-10 h-10 rounded-full bg-[#2f3336] animate-pulse" />
                         <div className="flex-1 h-10 bg-[#2f3336]/50 rounded animate-pulse" />
                    </div>
                ))}
            </div>
        ) : users.length === 0 ? (
            <div className="p-4 text-center text-[#71767b]">No suggestions available.</div>
        ) : (
            users.map((u) => (
                <div key={u.id} className="flex items-center gap-3 px-4 py-3 hover:bg-white/5 transition-colors cursor-pointer border-b border-[#2f3336]">
                        <Link href={`/${u.username}`} className="shrink-0">
                        <Avatar className="w-10 h-10">
                            <AvatarImage src={u.avatarUrl ?? undefined} />
                            <AvatarFallback>{u.displayName[0]}</AvatarFallback>
                        </Avatar>
                    </Link>
                    <div className="flex-1 min-w-0">
                        <Link href={`/${u.username}`} className="font-bold text-[#e7e9ea] hover:underline block truncate">
                            {u.displayName}
                        </Link>
                        <div className="text-[#71767b] truncate">@{u.username}</div>
                        {u.bio && <p className="text-[#e7e9ea] text-[14px] mt-1 line-clamp-1">{u.bio}</p>}
                    </div>
                    {currentUser?.id !== u.id && (
                        <FollowButton userId={u.id} isFollowing={u.followedByMe} />
                    )}
                </div>
            ))
        )}
      </div>
    </div>
  );
}
