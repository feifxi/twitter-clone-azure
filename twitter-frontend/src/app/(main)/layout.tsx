'use client';

import { AppNav } from '@/components/AppNav';
import { Sidebar } from '@/components/Sidebar';

import { usePathname } from 'next/navigation';

export default function MainLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const isMessagesPage = pathname?.startsWith('/messages');

  return (
    <div
      className="min-h-screen bg-background text-foreground flex justify-center"
      style={{ fontFamily: 'system-ui, -apple-system, BlinkMacSystemFont, sans-serif' }}
    >
      {/* Left gutter: flex grow so content is centered, then fixed columns */}
      <div className={`flex justify-center flex-1 min-w-0 ${isMessagesPage ? 'max-w-[1500px]' : 'max-w-[1280px]'}`}>
        {/* Left Sidebar */}
        <aside className="w-[68px] xl:w-[275px] shrink-0 flex justify-end">
          <AppNav />
        </aside>
        {/* Main Feed */}
        <main className={`w-full border-x border-border min-h-screen ${isMessagesPage ? 'max-w-[1000px] flex-1' : 'max-w-[700px]'}`}>
          {children}
        </main>
        {/* Right Sidebar: hidden on smaller screens, and hidden on messages page */}
        {!isMessagesPage && (
          <aside className="w-[350px] shrink-0 hidden lg:block">
            <Sidebar />
          </aside>
        )}
      </div>
    </div>
  );
}
