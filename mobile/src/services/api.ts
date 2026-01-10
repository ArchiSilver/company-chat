import AsyncStorage from '@react-native-async-storage/async-storage';
import axios, { AxiosInstance, AxiosResponse } from 'axios';
import { ApiResponse, Chat, LoginCredentials, Message, RegisterData, User } from '../types';

class ApiService {
    private api: AxiosInstance;
    private baseURL = 'http://localhost:8080'; // Change this to your backend URL

    constructor() {
        this.api = axios.create({
            baseURL: this.baseURL,
            timeout: 10000,
            headers: {
                'Content-Type': 'application/json',
            },
        });

        // Add request interceptor to include auth token
        this.api.interceptors.request.use(
            async config => {
                const token = await AsyncStorage.getItem('auth_token');
                if (token) {
                    config.headers.Authorization = `Bearer ${token}`;
                }
                return config;
            },
            error => Promise.reject(error),
        );

        // Add response interceptor for error handling
        this.api.interceptors.response.use(
            response => response,
            error => {
                if (error.response?.status === 401) {
                    // Token expired, clear storage
                    AsyncStorage.removeItem('auth_token');
                    AsyncStorage.removeItem('user_data');
                }
                return Promise.reject(error);
            },
        );
    }

    // Authentication methods
    async login(credentials: LoginCredentials): Promise<ApiResponse<{ user: User; token: string }>> {
        try {
            const response: AxiosResponse = await this.api.post('/auth/login', credentials);
            const { user, token } = response.data;

            // Store token and user data
            await AsyncStorage.setItem('auth_token', token);
            await AsyncStorage.setItem('user_data', JSON.stringify(user));

            return { success: true, data: { user, token } };
        } catch (error: any) {
            return {
                success: false,
                error: error.response?.data?.error || 'Login failed',
            };
        }
    }

    async register(data: RegisterData): Promise<ApiResponse<{ user: User; token: string }>> {
        try {
            const response: AxiosResponse = await this.api.post('/auth/register', data);
            const { user, token } = response.data;

            await AsyncStorage.setItem('auth_token', token);
            await AsyncStorage.setItem('user_data', JSON.stringify(user));

            return { success: true, data: { user, token } };
        } catch (error: any) {
            return {
                success: false,
                error: error.response?.data?.error || 'Registration failed',
            };
        }
    }

    async logout(): Promise<void> {
        await AsyncStorage.removeItem('auth_token');
        await AsyncStorage.removeItem('user_data');
    }

    async getCurrentUser(): Promise<User | null> {
        try {
            const userData = await AsyncStorage.getItem('user_data');
            return userData ? JSON.parse(userData) : null;
        } catch {
            return null;
        }
    }

    // Chat methods
    async getChats(): Promise<ApiResponse<Chat[]>> {
        try {
            const response: AxiosResponse = await this.api.get('/chats');
            return { success: true, data: response.data };
        } catch (error: any) {
            return {
                success: false,
                error: error.response?.data?.error || 'Failed to load chats',
            };
        }
    }

    async getChat(chatId: string): Promise<ApiResponse<Chat>> {
        try {
            const response: AxiosResponse = await this.api.get(`/chats/${chatId}`);
            return { success: true, data: response.data };
        } catch (error: any) {
            return {
                success: false,
                error: error.response?.data?.error || 'Failed to load chat',
            };
        }
    }

    async getMessages(chatId: string, limit = 50, offset = 0): Promise<ApiResponse<Message[]>> {
        try {
            const response: AxiosResponse = await this.api.get(
                `/chats/${chatId}/messages?limit=${limit}&offset=${offset}`,
            );
            return { success: true, data: response.data };
        } catch (error: any) {
            return {
                success: false,
                error: error.response?.data?.error || 'Failed to load messages',
            };
        }
    }

    async sendMessage(chatId: string, content: string): Promise<ApiResponse<Message>> {
        try {
            const response: AxiosResponse = await this.api.post(`/chats/${chatId}/messages`, {
                content,
                type: 'text',
            });
            return { success: true, data: response.data };
        } catch (error: any) {
            return {
                success: false,
                error: error.response?.data?.error || 'Failed to send message',
            };
        }
    }

    // WebSocket connection URL
    getWebSocketURL(token: string, chatId?: string): string {
        const baseUrl = this.baseURL.replace('http', 'ws');
        let url = `${baseUrl}/ws?token=${encodeURIComponent(token)}`;
        if (chatId) {
            url += `&chat_id=${encodeURIComponent(chatId)}`;
        }
        return url;
    }
}

export default new ApiService();