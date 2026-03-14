'use client';

import { Moon, Sun } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useTheme } from 'next-themes';
import { useEffect, useState } from 'react';

import { DropdownMenuItem } from '@/components/ui/dropdown-menu';

interface ThemeToggleProps {
  className?: string;
  variant?: 'button' | 'menu-item';
}

export function ThemeToggle({ className, variant = 'button' }: ThemeToggleProps) {
  const { theme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setMounted(true);
  }, []);

  const toggleTheme = () => setTheme(theme === 'light' ? 'dark' : 'light');

  if (!mounted) {
    if (variant === 'menu-item') {
      return (
        <DropdownMenuItem disabled className="flex items-center gap-4 p-3 opacity-50">
            <div className="w-5 h-5 border-2 border-muted border-t-foreground rounded-full animate-spin" />
            <span className="text-[15px]">Theme</span>
        </DropdownMenuItem>
      );
    }
    return (
      <Button variant="ghost" size="icon" className={`rounded-full hover:bg-card text-foreground transition-colors ${className || ''}`} disabled>
        <span className="w-6 h-6 border-2 border-muted border-t-foreground rounded-full animate-spin"></span>
      </Button>
    )
  }

  const Icon = theme === 'light' ? Moon : Sun;

  if (variant === 'menu-item') {
    return (
      <DropdownMenuItem 
        onClick={(e) => {
            e.preventDefault(); // Prevent closing on toggle if we want to see change immediately
            toggleTheme();
        }}
        className="flex items-center gap-4 p-3 text-[15px] cursor-pointer focus:bg-accent"
      >
        <Icon className="w-5 h-5" />
        <span>Theme</span>
      </DropdownMenuItem>
    );
  }

  return (
    <Button 
        variant="ghost" 
        size="icon" 
        onClick={toggleTheme}
        className={`rounded-full hover:bg-card text-foreground transition-colors ${className || ''}`}
        aria-label="Toggle theme"
    >
      <Icon className="w-6 h-6" />
    </Button>
  );
}
