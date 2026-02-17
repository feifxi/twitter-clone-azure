'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useNotifications, useMarkAllRead } from '@/hooks/useNotifications';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { ArrowLeft, Heart, Repeat2, User, MessageCircle } from 'lucide-react';
import Link from 'next/link';
import { formatDistanceToNowStrict } from 'date-fns';
import { useInView } from 'react-intersection-observer';
import type { NotificationResponse } from '@/types';

function NotificationItem({ notification }: { notification: NotificationResponse }) {
  const { actor, type, tweetId, tweetContent, createdAt } = notification;

  let Icon = User;
  let iconColor = 'text-[#1d9bf0]';
  let message = '';

  switch (type) {
    case 'LIKE':
        Icon = Heart;
        iconColor = 'text-[#f91880]';
        message = 'liked your tweet';
        break;
    case 'RETWEET':
        Icon = Repeat2;
        iconColor = 'text-[#00ba7c]';
        message = 'reposted your tweet';
        break;
    case 'FOLLOW':
        Icon = User;
        iconColor = 'text-[#1d9bf0]';
        message = 'followed you';
        break;
    case 'REPLY':
        Icon = MessageCircle;
        iconColor = 'text-[#1d9bf0]';
        message = 'replied to your tweet';
        break;
  }

  return (
    <div className={`flex gap-3 px-4 py-3 border-b border-[#2f3336] hover:bg-white/3 transition-colors ${!notification.isRead ? 'bg-[#16181c]' : ''}`}>
        <div className="w-10 shrink-0 flex justify-end">
            <Icon className={`w-7 h-7 ${iconColor} fill-current`} />
        </div>
        <div className="flex-1 flex flex-col gap-2">
             <div className="flex items-center gap-2">
                <Link href={`/${actor.username}`}>
                    <Avatar className="w-8 h-8">
                        <AvatarImage src={actor.avatarUrl ?? undefined} />
                        <AvatarFallback>{actor.displayName[0]}</AvatarFallback>
                    </Avatar>
                </Link>
             </div>
             <div className="text-[15px] text-[#e7e9ea]">
                <Link href={`/${actor.username}`} className="font-bold hover:underline">
                    {actor.displayName}
                </Link>
                {' '}
                <span className="text-[#e7e9ea]">{message}</span>
             </div>
             {(tweetId || tweetContent) && (
                <Link href={`/tweet/${tweetId}`} className="text-[#71767b] text-[15px] whitespace-pre-wrap line-clamp-3 hover:text-[#e7e9ea] transition-colors">
                    {tweetContent}
                </Link>
             )}
        </div>
    </div>
  );
}

export default function NotificationsPage() {
  const router = useRouter();
  const { data, isLoading, fetchNextPage, hasNextPage, isFetchingNextPage } = useNotifications();
  const markReadMutation = useMarkAllRead();
  const { ref, inView } = useInView();

  useEffect(() => {
    if (inView && hasNextPage) {
      fetchNextPage();
    }
  }, [inView, hasNextPage, fetchNextPage]);

  // Mark all as read on mount
  useEffect(() => {
      markReadMutation.mutate();
  }, []);

  const notifications = data?.pages.flatMap((p) => p.content) ?? [];

  return (
    <div className="min-h-screen">
      <div className="sticky top-0 z-30 bg-black/80 backdrop-blur border-b border-[#2f3336] px-4 py-3 flex items-center gap-4">
        <div 
            onClick={() => router.back()} 
            className="p-2 rounded-full hover:bg-white/10 cursor-pointer transition-colors -ml-2"
        >
            <ArrowLeft className="w-5 h-5 text-white" />
        </div>
        <h1 className="text-[20px] font-bold text-[#e7e9ea]">Notifications</h1>
      </div>

      {isLoading ? (
        <div className="p-4 text-[#71767b]">Loading...</div>
      ) : notifications.length === 0 ? (
        <div className="p-8 text-center text-[#71767b]">
            <p className="text-[30px] font-bold text-[#e7e9ea] mb-2">Nothing to see here — yet</p>
            <p>From likes to reposts and a whole lot more, this is where all the action happens.</p>
        </div>
      ) : (
        <div>
           {notifications.map((n) => (
             <NotificationItem key={n.id} notification={n} />
           ))}
           <div ref={ref} className="h-4" />
           {isFetchingNextPage && <div className="p-4 text-center text-[#71767b]">Loading more...</div>}
        </div>
      )}
    </div>
  );
}
