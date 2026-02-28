'use client';

import { Dialog, DialogContent, DialogTitle, DialogDescription } from '@/components/ui/dialog';
import { useRouter } from 'next/navigation';
import { CreateTweet } from './CreateTweet';
import { Tweet } from './Tweet';
import type { TweetResponse } from '@/types';
import { X } from 'lucide-react';

interface ReplyModalProps {
  tweet: TweetResponse | null;
  isOpen: boolean;
  onClose: () => void;
}

export function ReplyModal({ tweet, isOpen, onClose }: ReplyModalProps) {
  const router = useRouter();
  if (!tweet) return null;

  return (
    <div onClick={(e) => e.stopPropagation()}>
      <Dialog open={isOpen} onOpenChange={onClose}>
        <DialogContent showCloseButton={false} className="sm:max-w-[600px] bg-background border-border p-0 gap-0 top-[5%] translate-y-0 sm:top-[10%]">
          <DialogTitle className="sr-only">Reply to Tweet</DialogTitle>
          <DialogDescription className="sr-only">Replying to @{tweet.user.username}</DialogDescription>
          <div className="flex items-center h-[53px] px-4 shrink-0">
            <button
               onClick={onClose}
               className="p-2 rounded-full hover:bg-card transition-colors -ml-2"
            >
               <X size={20} className="text-foreground cursor-pointer" />
            </button>
          </div>
          
          <div className="px-4 pb-4">
              {/* Original Tweet Preview (simplified) */}
              <div className="flex gap-3 relative">
                  <div className="flex flex-col items-center">
                      <div className="w-10 h-10 rounded-full bg-card overflow-hidden shrink-0">
                          <img src={tweet.user.avatarUrl ?? undefined} alt="" className="w-full h-full object-cover"/>
                      </div>
                      <div className="w-0.5 grow bg-card my-2" />
                  </div>
                  <div className="flex-1 pb-4">
                       <div className="flex items-center gap-1 text-[15px]">
                          <span className="font-bold text-foreground">{tweet.user.displayName}</span>
                          <span className="text-muted-foreground">@{tweet.user.username}</span>
                          <span className="text-muted-foreground">· 1h</span>
                      </div>
                      <div className="text-foreground text-[15px] mt-1 whitespace-pre-wrap">{tweet.content}</div>
                       <div className="text-muted-foreground text-[15px] mt-2">
                          Replying to <span className="text-primary">@{tweet.user.username}</span>
                      </div>
                  </div>
              </div>

              {/* Create Reply */}
              <CreateTweet 
                  isReply 
                  replyToId={tweet.id} 
                  onSuccess={() => {
                    onClose();
                    router.push(`/tweet/${tweet.id}`);
                  }} 
                  className="border-none px-0 py-0"
                  placeholder="Post your reply"
              />
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
