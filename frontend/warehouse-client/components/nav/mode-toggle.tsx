'use client';
import React from 'react';
import { Moon, Sun } from 'lucide-react';
import { Button } from '../ui/button';
import { useTheme } from 'next-themes';
import { useEffect, useState } from 'react';

export default function ModeToggle() {
  const [mounted, setMounted] = useState(false);
  useEffect(() => {
    setMounted(true);
  }, []);
  if (!mounted) {
    return null;
  }
  const { theme, setTheme } = useTheme();
  return (
    <Button
      variant="link"
      size="icon"
      onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
    >
      {theme === 'dark' ? <Sun size={18} /> : <Moon size={18} />}
    </Button>
  );
}
