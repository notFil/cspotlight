// import { useState, useEffect } from 'react'; // Removed unused imports
import { useAppContext } from '../../context/AppContext';
import { Moon, Sun } from 'lucide-react';

const ThemeToggle = () => {
  const { state, toggleTheme } = useAppContext();
  const isDark = state.theme === 'dark';

  return (
    <button
      onClick={toggleTheme}
      className="fixed top-4 right-4 z-50 w-12 h-12 rounded-full flex items-center justify-center transition-all duration-300 hover:scale-110 shadow-lg"
      style={{
        backgroundColor: 'var(--color-surface)',
        border: '1px solid var(--color-border)',
        color: 'var(--color-text)',
      }}
      aria-label="Toggle theme"
      title={isDark ? 'Switch to light mode' : 'Switch to dark mode'}
    >
      {isDark ? (
        <Sun className="w-5 h-5" />
      ) : (
        <Moon className="w-5 h-5" />
      )}
    </button>
  );
};

export default ThemeToggle;

