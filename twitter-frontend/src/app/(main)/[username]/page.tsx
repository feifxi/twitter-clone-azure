'use client';

import { use, Suspense } from 'react';
import { FeedList } from '@/components/FeedList';
import { useUserProfileByUsername, useUserFeed } from '@/hooks/useProfile';
import { FollowButton } from '@/components/FollowButton';
import { useAuth } from '@/hooks/useAuth';
import { EditProfileModal } from '@/components/EditProfileModal';
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
        <div className="h-32 w-32 rounded-full bg-[#2f3336] animate-pulse" />
        <div className="h-6 w-48 bg-[#2f3336] rounded mt-4 animate-pulse" />
        <div className="h-4 w-32 bg-[#2f3336]/60 rounded mt-2 animate-pulse" />
      </div>
    );
  }

  const isOwnProfile = currentUser?.id === user.id;

  return (
    <div className="min-h-screen">
      <div className="sticky top-0 z-30 bg-black/60 backdrop-blur-md px-4 py-1 flex items-center gap-6 border-b border-[#2f3336]">
          <div 
              onClick={() => router.back()} 
              className="p-2 rounded-full hover:bg-white/10 cursor-pointer transition-colors -ml-2"
          >
              <ArrowLeft className="w-5 h-5 text-white" />
          </div>
          <div className="flex flex-col">
              <h2 className="font-bold text-[20px] leading-6 text-[#e7e9ea]">{user.displayName}</h2>
              <span className="text-[13px] text-[#71767b]">@{user.username}</span> 
          </div>
      </div>
      <div className="p-4 flex flex-col gap-4">
        <div className="flex flex-col sm:flex-row sm:items-end gap-4">
          <img
            src={user.avatarUrl ?? undefined}
            alt={user.displayName}
            className="w-20 h-20 rounded-full object-cover bg-[#2f3336]"
          />
          <div className="flex-1">
            <h1 className="text-[#e7e9ea] font-bold text-[20px]">
              {user.displayName}
            </h1>
            <p className="text-[#71767b] text-[15px]">@{user.username}</p>
            {user.bio && (
              <p className="text-[#e7e9ea] text-[15px] mt-1">{user.bio}</p>
            )}
            <div className="flex gap-4 mt-2 text-[#71767b] text-[15px]">
              <span>
                <strong className="text-[#e7e9ea]">{user.followingCount}</strong>{' '}
                Following
              </span>
              <span>
                <strong className="text-[#e7e9ea]">{user.followersCount}</strong>{' '}
                Followers
              </span>
            </div>
            {isLoggedIn && !isOwnProfile && (
              <div className="mt-2">
                <FollowButton userId={user.id} isFollowing={user.followedByMe} />
              </div>
            )}
            {isLoggedIn && isOwnProfile && (
                <button
                    type="button"
                    className="mt-2 py-1.5 px-4 rounded-full border border-[#536471] text-[#e7e9ea] font-bold text-[14px] hover:bg-[#eff3f41a] transition-colors"
                    onClick={() => setShowEditProfile(true)}
                >
                    Edit profile
                </button>
            )}
          </div>
        </div>
        
        <EditProfileModal user={user} isOpen={showEditProfile} onClose={() => setShowEditProfile(false)} />

        <div className="flex border-b border-[#2f3336]">
          {(['tweets', 'media'] as const).map((t) => (
            <button
              key={t}
              type="button"
              onClick={() => handleTabChange(t)}
              className={`flex-1 py-4 text-[15px] font-bold capitalize transition-colors hover:bg-[#080808] cursor-pointer ${
                currentTab === t
                  ? 'text-[#e7e9ea] border-b-2 border-[#1d9bf0]'
                  : 'text-[#71767b]'
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
        <Suspense fallback={<div className="p-4 text-center text-[#71767b]">Loading...</div>}>
            <ProfileContent params={params} />
        </Suspense>
    );
}
