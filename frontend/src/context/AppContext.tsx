import React, { createContext, useContext, useReducer } from 'react';
import type { ReactNode } from 'react';
import { useNavigate } from 'react-router-dom';
import { authService } from '@/services/auth';
import { userService } from '@/services/user';
import type { AppState, User, UserLogin } from '@/types';
import { toast } from 'sonner';

interface AppContextType {
    state: AppState;
    dispatch: React.Dispatch<AppAction>;
    login: (credentials: UserLogin) => Promise<void>;
    logout: () => Promise<void>;
    toggleTheme: () => void;
}

type AppAction =
    | { type: 'SET_LOADING'; payload: boolean }
    | { type: 'SET_USER'; payload: User | null }
    | { type: 'TOGGLE_THEME' }
    | { type: 'SET_ERROR'; payload: string | null };

const getInitialTheme = (): 'light' | 'dark' => {
    if (typeof window !== 'undefined') {
        const savedTheme = localStorage.getItem('theme');
        if (savedTheme === 'dark' || savedTheme === 'light') {
            return savedTheme;
        }
        if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
            return 'dark';
        }
    }
    return 'light'; // Default fallback
};

const initialState: AppState = {
    isLoading: true,
    user: null,
    theme: getInitialTheme(),
    error: null
};

const AppContext = createContext<AppContextType | undefined>(undefined);

const appReducer = (state: AppState, action: AppAction): AppState => {
    switch (action.type) {
        case 'SET_LOADING':
            return { ...state, isLoading: action.payload, error: null }; // Clear error on loading
        case 'SET_USER':
            return { ...state, user: action.payload };
        case 'TOGGLE_THEME':
            return { ...state, theme: state.theme === 'light' ? 'dark' : 'light' };
        case 'SET_ERROR':
            return { ...state, error: action.payload };
        default:
            return state;
    }
}

export const AppProvider = ({ children }: { children: ReactNode }) => {
    const [state, dispatch] = useReducer(appReducer, initialState);
    const navigate = useNavigate();

    // Effect to apply theme changes
    React.useEffect(() => {
        const root = window.document.documentElement;
        root.classList.remove('light', 'dark');
        root.classList.add(state.theme);
        localStorage.setItem('theme', state.theme);
    }, [state.theme]);

    // Effect to restore session
    React.useEffect(() => {
        const restoreSession = async () => {
            try {
                const user = await userService.getCurrentUser();
                dispatch({ type: 'SET_USER', payload: user });
            } catch (error) {
                // Not logged in or session expired
                console.log('No active session');
            } finally {
                dispatch({ type: 'SET_LOADING', payload: false });
            }
        };

        restoreSession();
    }, []);

    const login = async (credentials: UserLogin) => {
        dispatch({ type: 'SET_LOADING', payload: true });
        try {
            // Login now works via cookies. The response contains user data.
            const user = await authService.login(credentials);
            dispatch({ type: 'SET_USER', payload: user });
            navigate('/');
        } catch (error: any) {
            console.error('Login failed', error);
            const message = error.response?.data?.message || 'Failed to sign in. Please check your credentials.';
            dispatch({ type: 'SET_ERROR', payload: message });
            toast.error(message);
        } finally {
            dispatch({ type: 'SET_LOADING', payload: false });
        }
    };

    const logout = async () => {
        try {
            await authService.logout();
        } catch (error) {
            console.error('Logout failed', error);
        } finally {
            localStorage.clear();
            dispatch({ type: 'SET_USER', payload: null });
            navigate('/login');
        }
    };

    const toggleTheme = () => {
        dispatch({ type: 'TOGGLE_THEME' });
    }

    return (
        <AppContext.Provider value={{ state, dispatch, login, logout, toggleTheme }}>
            {children}
        </AppContext.Provider>
    );
};

export const useAppContext = (): AppContextType => {
    const context = useContext(AppContext);
    if (!context) {
        throw new Error('useAppContext must be used within an AppProvider');
    }
    return context;
};