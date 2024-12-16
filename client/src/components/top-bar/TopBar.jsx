import React from 'react';
import ModeToggle from '@/components/top-bar/mode-toggle';
import {
  Menubar,
  MenubarMenu,
  MenubarTrigger,
  MenubarContent,
  MenubarItem,
  MenubarSeparator,
} from '@/components/ui/menubar';
import {
  SignedIn,
  SignedOut,
  SignInButton,
  UserButton,
} from '@clerk/clerk-react';
import { LogIn, Search } from 'lucide-react';
import notification from '../../assets/notification.svg';
import logo from '../../assets/logo.svg';

export function TopBar() {
  return (
    <Menubar className="h-16 flex items-center px-4 bg-white dark:bg-gray-900 shadow-lg">
      {/* Logo Section */}
      <div className="flex items-center flex-shrink-0">
        <a href="/">
          <img src={logo} alt="logo" className="w-12 h-12" />
        </a>
      </div>

      {/* Search Section */}
      <div className="hidden md:flex flex-grow justify-center">
        <div className="relative w-full max-w-lg">
          <input
            type="text"
            placeholder="Search..."
            className="w-full px-4 h-10 pr-10 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button className="absolute inset-y-0 right-2 flex items-center text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white">
            <Search size={18} />
          </button>
        </div>
      </div>

      {/* Navigation Section */}
      <div className="flex items-center gap-4">
        {/* Notification Icon */}
        <div className="flex items-center justify-center h-10 w-10 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700 transition">
          <img
            src={notification}
            alt="notification"
            className="h-6 w-6 dark:invert"
          />
        </div>

        {/* Theme Toggle */}
        <div className="flex items-center justify-center h-10 w-10 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700 transition">
          <ModeToggle />
        </div>

        {/* UserButton */}
        <SignedOut>
          <div className="flex items-center justify-center h-10 w-10 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700 transition">
            <LogIn size={18} className="text-gray-600 dark:text-gray-300" />
            <SignInButton />
          </div>
        </SignedOut>
        <SignedIn>
          <div className="flex items-center justify-center h-10 w-10 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700 transition">
            <UserButton />
          </div>
        </SignedIn>
      </div>
    </Menubar>
  );
}

export default TopBar;
