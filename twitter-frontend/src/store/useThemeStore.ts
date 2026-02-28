import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { setTheme as applyCssTheme } from '@/theme/applyTheme';

type ThemeMode = 'light' | 'dark';

interface ThemeState {
    theme: ThemeMode;
    setTheme: (theme: ThemeMode) => void;
    toggleTheme: () => void;
}

export const useThemeStore = create<ThemeState>()(
    persist(
        (set) => ({
            theme: 'dark', // Default to dark mode
            setTheme: (theme) => {
                applyCssTheme(theme);
                set({ theme });
            },
            toggleTheme: () => set((state) => {
                const newTheme = state.theme === 'light' ? 'dark' : 'light';
                applyCssTheme(newTheme);
                return { theme: newTheme };
            }),
        }),
        {
            name: 'twitter-theme-storage',
        }
    )
);
