export type Role = 'USER' | 'ADMIN';

export interface UserResponse {
    id: number;
    username: string;
    email: string;
    displayName: string;
    avatarUrl: string | null;
    bio: string | null;
    role: Role;
    followersCount: number;
    followingCount: number;
    followedByMe: boolean;
}

export interface UpdateProfileRequest {
    displayName?: string;
    bio?: string;
}
