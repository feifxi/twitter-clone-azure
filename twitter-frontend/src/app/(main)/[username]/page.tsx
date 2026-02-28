'use client';

import { use, Suspense } from 'react';
import { FeedList } from '@/components/FeedList';
import { useUserProfileByUsername, useUserFeed } from '@/hooks/useProfile';
import { FollowButton } from '@/components/FollowButton';
import { useAuth } from '@/hooks/useAuth';
import { EditProfileModal } from '@/components/EditProfileModal';
import { UserListModal } from '@/components/UserListModal';
import { ArrowLeft } from 'lucide-react';
import { useRouter, useSearchParams, usePathname } from 'next/navigation';
import { useState } from 'react';

function ProfileContent({ params }: { params: Promise<{ username: string }> }) {
  const { username } = use(params);
  const { data: user, isLoading: userLoading } = useUserProfileByUsername(username);
  const userId = user?.id ?? null;
  const feed = useUserFeed(userId);
  const { user: currentUser, isLoggedIn } = useAuth();
  const [showEditProfile, setShowEditProfile] = useState(false);
  const [listType, setListType] = useState<'followers' | 'following' | null>(null);
  
  const router = useRouter();
  const searchParams = useSearchParams();
  const pathname = usePathname();

  const currentTab = searchParams.get('tab') === 'media' ? 'media' : 'tweets';

  const tweets = feed.data?.pages.flatMap((p) => p.content) ?? [];
  const tweetsFiltered =
    currentTab === 'media'
      ? tweets.filter((t) => (t.originalTweet ?? t).mediaUrl)
      : tweets;

  const handleTabChange = (newTab: 'tweets' | 'media') => {
    if (newTab !== currentTab) {
        const params = new URLSearchParams(searchParams);
        params.set('tab', newTab);
        router.replace(`${pathname}?${params.toString()}`);
    }
  };

  if (userLoading || !user) {
    return (
      <div className="p-4">
        <div className="h-32 w-32 rounded-full bg-card animate-pulse" />
        <div className="h-6 w-48 bg-card rounded mt-4 animate-pulse" />
        <div className="h-4 w-32 bg-card/60 rounded mt-2 animate-pulse" />
      </div>
    );
  }

  const isOwnProfile = currentUser?.id === user.id;

  return (
    <div className="min-h-screen">
      <div className="sticky top-0 z-30 bg-background/60 backdrop-blur-md px-4 py-1 flex items-center gap-6 border-b border-border">
          <div 
              onClick={() => router.back()} 
              className="p-2 rounded-full hover:bg-card cursor-pointer transition-colors -ml-2"
          >
              <ArrowLeft className="w-5 h-5 text-foreground" />
          </div>
          <div className="flex flex-col">
              <h2 className="font-bold text-[20px] leading-6 text-foreground">{user.displayName}</h2>
              <span className="text-[13px] text-muted-foreground">@{user.username}</span> 
          </div>
      </div>
      <div className="p-4 flex flex-col gap-4">
        <div className="flex flex-col sm:flex-row sm:items-end gap-4">
          <img
            src={user.avatarUrl ?? undefined}
            alt={user.displayName}
            className="w-20 h-20 rounded-full object-cover bg-card"
          />
          <div className="flex-1">
            <h1 className="text-foreground font-bold text-[20px]">
              {user.displayName}
            </h1>
            <p className="text-muted-foreground text-[15px]">@{user.username}</p>
            {user.bio && (
              <p className="text-foreground text-[15px] mt-1">{user.bio}</p>
            )}
            <div className="flex gap-4 mt-2 text-muted-foreground text-[15px]">
              <button 
                className="hover:underline cursor-pointer"
                onClick={() => setListType('following')}
              >
                <strong className="text-foreground">{user.followingCount}</strong>{' '}
                Following
              </button>
              <button 
                className="hover:underline cursor-pointer"
                onClick={() => setListType('followers')}
              >
                <strong className="text-foreground">{user.followersCount}</strong>{' '}
                Followers
              </button>
            </div>
            {isLoggedIn && !isOwnProfile && (
              <div className="mt-2">
                <FollowButton userId={user.id} isFollowing={user.followedByMe} />
              </div>
            )}
            {isLoggedIn && isOwnProfile && (
                <button
                    type="button"
                    className="mt-2 py-1.5 px-4 rounded-full border border-border text-foreground font-bold text-[14px] hover:bg-card transition-colors"
                    onClick={() => setShowEditProfile(true)}
                >
                    Edit profile
                </button>
            )}
          </div>
        </div>
        
        <EditProfileModal user={user} isOpen={showEditProfile} onClose={() => setShowEditProfile(false)} />
        
        {user.id && (
          <UserListModal 
            userId={user.id} 
            type={listType} 
            isOpen={!!listType} 
            onClose={() => setListType(null)} 
          />
        )}

        <div className="flex border-b border-border">
          {(['tweets', 'media'] as const).map((t) => (
            <button
              key={t}
              type="button"
              onClick={() => handleTabChange(t)}
              className={`flex-1 py-4 text-[15px] font-bold capitalize transition-colors hover:bg-card cursor-pointer ${
                currentTab === t
                  ? 'text-foreground border-b-2 border-primary'
                  : 'text-muted-foreground'
              }`}
            >
              {t}
            </button>
          ))}
        </div>

        <FeedList
            tweets={tweetsFiltered}
            isLoading={feed.isLoading}
            isFetchingNextPage={feed.isFetchingNextPage}
            hasNextPage={feed.hasNextPage}
            fetchNextPage={feed.fetchNextPage}
        />
      </div>
    </div>
  );
}

export default function ProfilePage({ params }: { params: Promise<{ username: string }> }) {
    return (
        <Suspense fallback={<div className="p-4 text-center text-muted-foreground">Loading...</div>}>
            <ProfileContent params={params} />
        </Suspense>
    );
}
