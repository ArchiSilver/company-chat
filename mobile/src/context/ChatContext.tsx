import React, { createContext, ReactNode, useCallback, useContext, useEffect, useReducer } from 'react';
import apiService from '../services/api';
import { Chat, Message, WebSocketMessage } from '../types';
import { useAuth } from './AuthContext';

interface ChatState {
    chats: Chat[];
    currentChat: Chat | null;
    messages: Message[];
    isLoading: boolean;
    error: string | null;
    isConnected: boolean;
}

interface ChatContextType extends ChatState {
    loadChats: () => Promise<void>;
    selectChat: (chatId: string) => Promise<void>;
    sendMessage: (content: string) => Promise<void>;
    connectWebSocket: (chatId?: string) => void;
    disconnectWebSocket: () => void;
}

const ChatContext = createContext<ChatContextType | undefined>(undefined);

type ChatAction =
    | { type: 'SET_LOADING'; payload: boolean }
    | { type: 'SET_CHATS'; payload: Chat[] }
    | { type: 'SET_CURRENT_CHAT'; payload: Chat }
    | { type: 'SET_MESSAGES'; payload: Message[] }
    | { type: 'ADD_MESSAGE'; payload: Message }
    | { type: 'SET_ERROR'; payload: string }
    | { type: 'CLEAR_ERROR' }
    | { type: 'SET_CONNECTED'; payload: boolean };

const chatReducer = (state: ChatState, action: ChatAction): ChatState => {
    switch (action.type) {
        case 'SET_LOADING':
            return { ...state, isLoading: action.payload };
        case 'SET_CHATS':
            return { ...state, chats: action.payload, isLoading: false };
        case 'SET_CURRENT_CHAT':
            return { ...state, currentChat: action.payload };
        case 'SET_MESSAGES':
            return { ...state, messages: action.payload };
        case 'ADD_MESSAGE':
            return {
                ...state,
                messages: [...state.messages, action.payload],
            };
        case 'SET_ERROR':
            return { ...state, error: action.payload, isLoading: false };
        case 'CLEAR_ERROR':
            return { ...state, error: null };
        case 'SET_CONNECTED':
            return { ...state, isConnected: action.payload };
        default:
            return state;
    }
};

const initialState: ChatState = {
    chats: [],
    currentChat: null,
    messages: [],
    isLoading: false,
    error: null,
    isConnected: false,
};

export const ChatProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
    const [state, dispatch] = useReducer(chatReducer, initialState);
    const { token, isAuthenticated } = useAuth();

    // WebSocket connection
    const wsRef = React.useRef<WebSocket | null>(null);

    const connectWebSocket = (chatId?: string) => {
        if (!token || wsRef.current) return;

        const wsUrl = apiService.getWebSocketURL(token, chatId);
        const ws = new WebSocket(wsUrl);

        ws.onopen = () => {
            dispatch({ type: 'SET_CONNECTED', payload: true });
            console.log('WebSocket connected');
        };

        ws.onmessage = event => {
            try {
                const message: WebSocketMessage = JSON.parse(event.data);
                handleWebSocketMessage(message);
            } catch (error) {
                console.error('Error parsing WebSocket message:', error);
            }
        };

        ws.onclose = () => {
            dispatch({ type: 'SET_CONNECTED', payload: false });
            console.log('WebSocket disconnected');
            wsRef.current = null;
        };

        ws.onerror = error => {
            console.error('WebSocket error:', error);
            dispatch({ type: 'SET_CONNECTED', payload: false });
        };

        wsRef.current = ws;
    };

    const disconnectWebSocket = () => {
        if (wsRef.current) {
            wsRef.current.close();
            wsRef.current = null;
        }
    };

    const handleWebSocketMessage = (message: WebSocketMessage) => {
        switch (message.type) {
            case 'message':
                if (message.chat_id && message.user_id) {
                    // This would be a new message from WebSocket
                    // For now, we'll reload messages when selecting chat
                }
                break;
            default:
                break;
        }
    };

    const loadChats = useCallback(async (): Promise<void> => {
        if (!isAuthenticated) return;

        dispatch({ type: 'SET_LOADING', payload: true });
        dispatch({ type: 'CLEAR_ERROR' });

        try {
            const result = await apiService.getChats();
            if (result.success && result.data) {
                dispatch({ type: 'SET_CHATS', payload: result.data });
            } else {
                dispatch({ type: 'SET_ERROR', payload: result.error || 'Failed to load chats' });
            }
        } catch {
            dispatch({ type: 'SET_ERROR', payload: 'Network error' });
        }
    }, [isAuthenticated]);

    const selectChat = async (chatId: string): Promise<void> => {
        dispatch({ type: 'SET_LOADING', payload: true });

        try {
            const [chatResult, messagesResult] = await Promise.all([
                apiService.getChat(chatId),
                apiService.getMessages(chatId),
            ]);

            if (chatResult.success && chatResult.data) {
                dispatch({ type: 'SET_CURRENT_CHAT', payload: chatResult.data });
            }

            if (messagesResult.success && messagesResult.data) {
                dispatch({ type: 'SET_MESSAGES', payload: messagesResult.data });
            }

            dispatch({ type: 'SET_LOADING', payload: false });
        } catch {
            dispatch({ type: 'SET_ERROR', payload: 'Failed to load chat' });
        }
    };

    const sendMessage = async (content: string): Promise<void> => {
        if (!state.currentChat) return;

        try {
            const result = await apiService.sendMessage(state.currentChat.id, content);
            if (result.success && result.data) {
                dispatch({ type: 'ADD_MESSAGE', payload: result.data });
            } else {
                dispatch({ type: 'SET_ERROR', payload: result.error || 'Failed to send message' });
            }
        } catch {
            dispatch({ type: 'SET_ERROR', payload: 'Network error' });
        }
    };

    // Load chats when user is authenticated
    useEffect(() => {
        if (isAuthenticated) {
            loadChats();
        }
    }, [isAuthenticated, loadChats]);

    // Cleanup WebSocket on unmount
    useEffect(() => {
        return () => {
            disconnectWebSocket();
        };
    }, []);

    const value: ChatContextType = {
        ...state,
        loadChats,
        selectChat,
        sendMessage,
        connectWebSocket,
        disconnectWebSocket,
    };

    return <ChatContext.Provider value={value}>{children}</ChatContext.Provider>;
};

export const useChat = (): ChatContextType => {
    const context = useContext(ChatContext);
    if (context === undefined) {
        throw new Error('useChat must be used within a ChatProvider');
    }
    return context;
};