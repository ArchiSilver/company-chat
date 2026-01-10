import AsyncStorage from '@react-native-async-storage/async-storage';
import React, { createContext, ReactNode, useContext, useEffect, useReducer } from 'react';
import apiService from '../services/api';
import { AuthState, LoginCredentials, RegisterData, User } from '../types';

interface AuthContextType extends AuthState {
    login: (credentials: LoginCredentials) => Promise<boolean>;
    register: (data: RegisterData) => Promise<boolean>;
    logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

type AuthAction =
    | { type: 'SET_LOADING'; payload: boolean }
    | { type: 'SET_USER'; payload: { user: User; token: string } }
    | { type: 'LOGOUT' };

const authReducer = (state: AuthState, action: AuthAction): AuthState => {
    switch (action.type) {
        case 'SET_LOADING':
            return { ...state, isLoading: action.payload };
        case 'SET_USER':
            return {
                ...state,
                user: action.payload.user,
                token: action.payload.token,
                isAuthenticated: true,
                isLoading: false,
            };
        case 'LOGOUT':
            return {
                user: null,
                token: null,
                isAuthenticated: false,
                isLoading: false,
            };
        default:
            return state;
    }
};

const initialState: AuthState = {
    user: null,
    token: null,
    isLoading: true,
    isAuthenticated: false,
};

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
    const [state, dispatch] = useReducer(authReducer, initialState);

    useEffect(() => {
        // Check for stored auth data on app start
        const initializeAuth = async () => {
            try {
                const token = await AsyncStorage.getItem('auth_token');
                const userData = await AsyncStorage.getItem('user_data');

                if (token && userData) {
                    const user = JSON.parse(userData);
                    dispatch({ type: 'SET_USER', payload: { user, token } });
                } else {
                    dispatch({ type: 'SET_LOADING', payload: false });
                }
            } catch (error) {
                console.error('Error initializing auth:', error);
                dispatch({ type: 'SET_LOADING', payload: false });
            }
        };

        initializeAuth();
    }, []);

    const login = async (credentials: LoginCredentials): Promise<boolean> => {
        dispatch({ type: 'SET_LOADING', payload: true });

        try {
            const result = await apiService.login(credentials);
            if (result.success && result.data) {
                dispatch({ type: 'SET_USER', payload: result.data });
                return true;
            } else {
                dispatch({ type: 'SET_LOADING', payload: false });
                return false;
            }
        } catch {
            dispatch({ type: 'SET_LOADING', payload: false });
            return false;
        }
    };

    const register = async (data: RegisterData): Promise<boolean> => {
        dispatch({ type: 'SET_LOADING', payload: true });

        try {
            const result = await apiService.register(data);
            if (result.success && result.data) {
                dispatch({ type: 'SET_USER', payload: result.data });
                return true;
            } else {
                dispatch({ type: 'SET_LOADING', payload: false });
                return false;
            }
        } catch {
            dispatch({ type: 'SET_LOADING', payload: false });
            return false;
        }
    };

    const logout = async (): Promise<void> => {
        await apiService.logout();
        dispatch({ type: 'LOGOUT' });
    };

    const value: AuthContextType = {
        ...state,
        login,
        register,
        logout,
    };

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = (): AuthContextType => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};