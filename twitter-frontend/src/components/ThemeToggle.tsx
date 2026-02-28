'use client';

import { Moon, Sun } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useThemeStore } from '@/store/useThemeStore';

interface ThemeToggleProps {
  className?: string;
}

export function ThemeToggle({ className }: ThemeToggleProps) {
  const { theme, toggleTheme } = useThemeStore();

  return (
    <Button 
        variant="ghost" 
        size="icon" 
        onClick={toggleTheme}
        className={`rounded-full hover:bg-card text-foreground transition-colors ${className || ''}`}
        aria-label="Toggle theme"
    >
      {theme === 'light' ? (
        <Moon className="w-6 h-6" />
      ) : (
        <Sun className="w-6 h-6" />
      )}
    </Button>
  );
}
