'use client';

import {
  useQuery,
  useInfiniteQuery,
  useQueryClient,
  useMutation,
  type QueryKey,
} from '@tanstack/react-query';
import { axiosInstance } from '@/api/axiosInstance';
import type { PageResponse, TweetResponse } from '@/types';
import { toast } from 'sonner';

export const tweetQueryKey = (id: number) => ['tweets', id] as const;

export function useTweet(tweetId: number | null) {
  return useQuery({
    queryKey: tweetQueryKey(tweetId!),
    queryFn: async (): Promise<TweetResponse> => {
      const { data } = await axiosInstance.get<TweetResponse>(`/tweets/${tweetId}`);
      return data;
    },
    enabled: tweetId != null && tweetId > 0,
  });
}

export function useRepliesInfinite(tweetId: number | null, pageSize = 20) {
  return useInfiniteQuery({
    queryKey: ['tweets', tweetId, 'replies'],
    queryFn: async ({
      pageParam,
    }: {
      pageParam: number;
    }): Promise<PageResponse<TweetResponse>> => {
      const { data } = await axiosInstance.get<PageResponse<TweetResponse>>(
        `/tweets/${tweetId}/replies`,
        { params: { page: pageParam, size: pageSize } }
      );
      return data;
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) =>
      lastPage.last ? undefined : lastPage.page + 1,
    enabled: tweetId != null && tweetId > 0,
  });
}

export function useInvalidateTweet() {
  const queryClient = useQueryClient();
  return (tweetId: number) =>
    queryClient.invalidateQueries({ queryKey: tweetQueryKey(tweetId) });
}

import { useAuthStore } from '@/store/useAuthStore';
import { feedQueryKey } from './useFeed';
import type { InfiniteData } from '@tanstack/react-query';

export function useCreateTweet() {
  const queryClient = useQueryClient();
  const { user } = useAuthStore();

  return useMutation({
    mutationFn: async ({ content, media, parentId }: { content: string; media?: File; parentId?: number }) => {
      const formData = new FormData();
      formData.append(
        'data',
        new Blob([JSON.stringify({ content, parentId })], { type: 'application/json' })
      );
      if (media) {
        formData.append('media', media);
      }
      const { data } = await axiosInstance.post<TweetResponse>('/tweets', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
      return data;
    },
    onMutate: async (newTweet) => {
      // Stop outgoing refetches
      await queryClient.cancelQueries({ queryKey: ['feeds'] });

      // Snapshot previous value
      const previousGlobal = queryClient.getQueryData<InfiniteData<PageResponse<TweetResponse>>>(feedQueryKey('global'));

      // Optimistic update
      if (user) {
        const tempId = Date.now();
        const optimisticTweet: TweetResponse = {
          id: tempId,
          content: newTweet.content,
          createdAt: new Date().toISOString(),
          user: user,
          likeCount: 0,
          retweetCount: 0,
          replyCount: 0,
          likedByMe: false,
          retweetedByMe: false,
          mediaUrl: newTweet.media ? URL.createObjectURL(newTweet.media) : null,
          mediaType: newTweet.media ? 'IMAGE' : null,
          replyToTweetId: newTweet.parentId ?? null,
          originalTweet: null, // Not a retweet
        };





        // Helper to prepend tweet to feed

        const updateFeed = (key: QueryKey) => {
          queryClient.setQueryData<InfiniteData<PageResponse<TweetResponse>>>(key, (old) => {
            if (!old) return old;
            const newPages = [...old.pages];
            if (newPages.length > 0) {
              newPages[0] = {
                ...newPages[0],
                content: [optimisticTweet, ...newPages[0].content],
              };
            }
            return { ...old, pages: newPages };
          });
        }

        updateFeed(feedQueryKey('global'));
        updateFeed(feedQueryKey('following'));

        // Optimistic update for replies if on tweet detail page
        if (newTweet.parentId) {
          queryClient.setQueryData<InfiniteData<PageResponse<TweetResponse>>>(
            ['tweets', newTweet.parentId, 'replies'],
            (old) => {
              if (!old) return old;
              const newPages = [...old.pages];
              // Add to first page or create a new page if empty
              if (newPages.length > 0) {
                newPages[0] = {
                  ...newPages[0],
                  content: [optimisticTweet, ...newPages[0].content],
                };
              }
              return { ...old, pages: newPages };
            }
          );
        }
      }

      return { previousGlobal };
    },
    onError: (_err, _newTweet, context) => {
      if (context?.previousGlobal) {
        queryClient.setQueryData(feedQueryKey('global'), context.previousGlobal);
        // We could also rollback following feed but simplistic approach is enough, usually just invalidate on error
        queryClient.invalidateQueries({ queryKey: ['feeds'] });
      }
      // Invalidate replies query too
      if (_newTweet.parentId) {
        queryClient.invalidateQueries({ queryKey: ['tweets', _newTweet.parentId, 'replies'] });
      }
    },
    onSettled: () => {

    },
    onSuccess: (realTweet, variables) => {
      // Helper to swap optimistic tweet with real tweet in a feed
      const swapOptimisticTweet = (key: QueryKey) => {
        queryClient.setQueryData<InfiniteData<PageResponse<TweetResponse>>>(key, (old) => {
          if (!old) return old;
          return {
            ...old,
            pages: old.pages.map(page => ({
              ...page,
              content: page.content.map(t => {
                if (t.content === variables.content && t.user.id === user?.id && t.id > 1700000000000) {
                  return realTweet;
                }
                return t;
              })
            }))
          };
        });
      };

      swapOptimisticTweet(feedQueryKey('global'));
      swapOptimisticTweet(feedQueryKey('following'));

      // Also update replies if it was a reply
      if (variables.parentId) {
        swapOptimisticTweet(['tweets', variables.parentId, 'replies']);
      }
    }
  });
}

export function useDeleteTweet() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (tweetId: number) => {
      await axiosInstance.delete(`/tweets/${tweetId}`);
    },
    onMutate: async (tweetId) => {
      // Cancel any outgoing refetches so they don't overwrite our optimistic update
      await queryClient.cancelQueries({ queryKey: ['feeds'] });
      await queryClient.cancelQueries({ queryKey: ['tweets', tweetId] });

      // Helper to remove tweet from a feed
      const removeTweetFromFeed = (feedKey: QueryKey) => {
        const previousData = queryClient.getQueryData<InfiniteData<PageResponse<TweetResponse>>>(feedKey);

        if (previousData) {
          queryClient.setQueryData<InfiniteData<PageResponse<TweetResponse>>>(feedKey, (old) => {
            if (!old) return old;
            return {
              ...old,
              pages: old.pages.map((page) => ({
                ...page,
                content: page.content.filter((t) => t.id !== tweetId),
              })),
            };
          });
        }
        return previousData;
      };

      const previousGlobal = removeTweetFromFeed(feedQueryKey('global'));
      const previousFollowing = removeTweetFromFeed(feedQueryKey('following'));

      return { previousGlobal, previousFollowing };
    },
    onError: (_err, _tweetId, context) => {
      if (context?.previousGlobal) {
        queryClient.setQueryData(feedQueryKey('global'), context.previousGlobal);
      }
      if (context?.previousFollowing) {
        queryClient.setQueryData(feedQueryKey('following'), context.previousFollowing);
      }
      toast.error('Failed to delete tweet');
    },
    onSuccess: () => {
      toast.success('Tweet deleted');
      // Invalidate to ensure consistency, though optimistic update handles immediate UI
      queryClient.invalidateQueries({ queryKey: ['feeds'] });
      queryClient.invalidateQueries({ queryKey: ['profile'] });
    },
  });
}
