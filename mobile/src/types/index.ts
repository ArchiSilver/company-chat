// Types for Company Chat Mobile App

export interface User {
    id: string;
    email: string;
    username: string;
    role: string;
    created_at: string;
}

export interface Chat {
    id: string;
    name: string;
    type: 'direct' | 'group';
    participants: User[];
    last_message?: Message;
    created_at: string;
    updated_at: string;
}

export interface Message {
    id: string;
    chat_id: string;
    user_id: string;
    content: string;
    type: 'text' | 'image' | 'file';
    created_at: string;
    user?: User;
}

export interface AuthState {
    user: User | null;
    token: string | null;
    isLoading: boolean;
    isAuthenticated: boolean;
}

export interface LoginCredentials {
    email: string;
    password: string;
}

export interface RegisterData {
    email: string;
    username: string;
    password: string;
}

export interface WebSocketMessage {
    type: 'message' | 'join' | 'leave' | 'typing';
    chat_id?: string;
    content?: string;
    user_id?: string;
}

export interface ApiResponse<T> {
    success: boolean;
    data?: T;
    error?: string;
}