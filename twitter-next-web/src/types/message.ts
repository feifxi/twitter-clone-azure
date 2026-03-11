import type { UserResponse } from './user';

export interface MessageResponse {
  id: number;
  conversationId: number;
  sender: UserResponse;
  content: string;
  createdAt: string;
}

export interface ConversationResponse {
  id: number;
  peer: UserResponse;
  lastMessage: MessageResponse;
  updatedAt: string;
}

