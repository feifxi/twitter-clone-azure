import type { UserResponse } from './user';

export interface GoogleAuthRequest {
    token: string;
}

export interface AuthResponse {
    accessToken: string;
    user: UserResponse;
}
