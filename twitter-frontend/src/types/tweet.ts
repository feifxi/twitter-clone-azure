import type { UserResponse } from './user';

export interface TweetRequest {
    content?: string;
    parentId?: number | null;
}

export interface TweetResponse {
    id: number;
    content: string | null;
    mediaType: string | null;
    mediaUrl: string | null;
    user: UserResponse;
    replyCount: number;
    likeCount: number;
    retweetCount: number;
    likedByMe: boolean;
    retweetedByMe: boolean;
    originalTweet: TweetResponse | null;
    replyToTweetId?: number | null;
    replyToUserHandle?: string | null;
    createdAt: string;
}
