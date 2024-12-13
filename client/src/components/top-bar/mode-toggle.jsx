import { Moon, Sun } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useTheme } from '@/components/theme/theme-provider';
import React from 'react';

export default function ModeToggle() {
  const { theme, setTheme } = useTheme();

  return (
    <Button
      variant="link"
      size="icon"
      onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
    >
      {theme === 'dark' ? (
        <Sun width={24} height={24} />
      ) : (
        <Moon width={24} height={24} />
      )}
    </Button>
  );
}
