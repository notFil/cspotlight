import React, { createContext, useContext, useReducer } from 'react';
import type { ReactNode } from 'react';
import type { AppState, User } from '@/types';

interface AppContextType {
    state: AppState;
    dispatch: React.Dispatch<AppAction>;
    login: (user: User) => void;
    logout: () => void;
    toggleTheme: () => void;
}

type AppAction =
    | { type: 'SET_LOADING'; payload: boolean }
    | { type: 'SET_USER'; payload: User | null }
    | { type: 'TOGGLE_THEME' };

const initialState: AppState = {
    isLoading: false,
    user: null,
    theme: 'dark'
};

const AppContext = createContext<AppContextType | undefined>(undefined);

const appReducer = (state: AppState, action: AppAction): AppState => {
    switch (action.type) {
        case 'SET_LOADING':
            return { ...state, isLoading: action.payload };
        case 'SET_USER':
            return { ...state, user: action.payload };
        case 'TOGGLE_THEME':
            return { ...state, theme: state.theme === 'light' ? 'dark' : 'light' };
        default:
            return state;
    }
}

export const AppProvider = ({ children }: { children: ReactNode }) => {
    const [state, dispatch] = useReducer(appReducer, initialState);

    const login = (user: User) => {
        dispatch({ type: 'SET_USER', payload: user });
    };

    const logout = () => {
        dispatch({ type: 'SET_USER', payload: null });
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