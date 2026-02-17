import { useMutation, useQueryClient } from '@tanstack/react-query';
import { axiosInstance } from '@/api/axiosInstance';
import type { TweetResponse } from '@/types';
import { tweetQueryKey } from './useTweet';
import { useUpdateTweetCache } from './useOptimisticTweetUpdate';

import { useAuth } from '@/hooks/useAuth';

function toggleRetweetInTweet(
  t: TweetResponse,
  tweetId: number,
  retweeted: boolean,
  currentUserId?: number
): TweetResponse | null {
  // If un-retweeting (retweeted === false)
  if (!retweeted && currentUserId) {
    // Logic: If this tweet is a Retweet Wrapper (t.originalTweet exists) AND
    // the wrapper's author is ME (t.user.id === currentUserId) AND
    // the original tweet ID matches the target tweetId (t.originalTweet.id === tweetId)
    // THEN: Delete this wrapper from my feed view.
    if (t.originalTweet?.id === tweetId && t.user.id === currentUserId) {
      return null; // Delete!
    }
  }

  const delta = retweeted ? 1 : -1;
  const update = (x: TweetResponse) =>
    x.id === tweetId
      ? { ...x, retweetedByMe: retweeted, retweetCount: Math.max(0, x.retweetCount + delta) }
      : x;

  if (t.id === tweetId) return update(t) as TweetResponse;

  // If it's a wrapper but NOT my wrapper (or I'm not unretweeting), just update the inner content
  if (t.originalTweet?.id === tweetId)
    return { ...t, originalTweet: update(t.originalTweet) as TweetResponse };

  return t;
}

export function useRetweet() {
  const queryClient = useQueryClient();
  const updateTweetCache = useUpdateTweetCache();
  const { user } = useAuth();
  const currentUserId = user?.id;

  return useMutation({
    mutationFn: async (tweetId: number) => {
      await axiosInstance.post(`/tweets/${tweetId}/retweet`);
    },
    onMutate: async (tweetId) => {
      await queryClient.cancelQueries({ queryKey: ['feeds'] });
      await queryClient.cancelQueries({ queryKey: ['tweets', tweetId] });
      await queryClient.cancelQueries({ queryKey: ['tweets'] });
      await queryClient.cancelQueries({ queryKey: ['search'] });

      updateTweetCache(tweetId, (t) => toggleRetweetInTweet(t, tweetId, true, currentUserId));
    },
    onError: (_err, tweetId) => {
      queryClient.invalidateQueries({ queryKey: ['feeds'] });
      queryClient.invalidateQueries({ queryKey: ['tweets'] });
      queryClient.invalidateQueries({ queryKey: ['search'] });
    },
    onSettled: (_data, _err, tweetId) => {
      queryClient.invalidateQueries({ queryKey: tweetQueryKey(tweetId) });
    },
  });
}

export function useUnretweet() {
  const queryClient = useQueryClient();
  const updateTweetCache = useUpdateTweetCache();
  const { user } = useAuth();
  const currentUserId = user?.id;

  return useMutation({
    mutationFn: async (tweetId: number) => {
      await axiosInstance.delete(`/tweets/${tweetId}/retweet`);
    },
    onMutate: async (tweetId) => {
      await queryClient.cancelQueries({ queryKey: ['feeds'] });
      await queryClient.cancelQueries({ queryKey: ['tweets', tweetId] });
      await queryClient.cancelQueries({ queryKey: ['tweets'] });
      await queryClient.cancelQueries({ queryKey: ['search'] });

      updateTweetCache(tweetId, (t) => toggleRetweetInTweet(t, tweetId, false, currentUserId));
    },
    onError: (_err, tweetId) => {
      queryClient.invalidateQueries({ queryKey: ['feeds'] });
      queryClient.invalidateQueries({ queryKey: ['tweets'] });
      queryClient.invalidateQueries({ queryKey: ['search'] });
    },
    onSettled: (_data, _err, tweetId) => {
      queryClient.invalidateQueries({ queryKey: tweetQueryKey(tweetId) });
    },
  });
}
