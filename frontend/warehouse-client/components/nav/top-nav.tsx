import React from 'react';
import ModeToggle from './mode-toggle';
import {
  Menubar,
  MenubarContent,
  MenubarItem,
  MenubarMenu,
  MenubarSeparator,
  MenubarShortcut,
  MenubarTrigger,
} from '@/components/ui/menubar';

export default function TopNav() {
  return (
    <Menubar>
      <div className="flex-none">
        <MenubarMenu>Logo</MenubarMenu>
      </div>
      <div className="flex flex-grow items-center justify-end gap-1">
        <MenubarMenu>
          <MenubarTrigger className="text-base font-normal">
            Dashboard
          </MenubarTrigger>
          <MenubarContent>
            <MenubarItem>Task 1</MenubarItem>
            <MenubarSeparator />
            <MenubarItem>Task 1</MenubarItem>
          </MenubarContent>
        </MenubarMenu>
        <MenubarMenu>
          <ModeToggle />
        </MenubarMenu>
      </div>
    </Menubar>
  );
}
