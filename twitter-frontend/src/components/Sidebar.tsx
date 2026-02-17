'use client';

import Link from 'next/link';
import { useRouter, usePathname } from 'next/navigation'; // Added usePathname
import { useState, useRef, useEffect } from 'react';
import { Search } from 'lucide-react';
import { useTrendingHashtags, useSuggestedUsers, useSearchUsers, useSearchHashtags } from '@/hooks/useDiscovery';
import { useAuth } from '@/hooks/useAuth';
import { useFollowUser, useUnfollowUser } from '@/hooks/useProfile';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { useDebounce } from '@/hooks/useDebounce';
import { FollowButton } from '@/components/FollowButton';
import type { TrendingHashtagDTO } from '@/types';

export function Sidebar() {
  const router = useRouter();
  const pathname = usePathname();
  const [query, setQuery] = useState('');
  const [isFocused, setIsFocused] = useState(false);
  const debouncedQuery = useDebounce(query, 500);
  const searchRef = useRef<HTMLDivElement>(null);

  const { data: searchResults, isLoading: isSearchLoading } = useSearchUsers(debouncedQuery);
  const { data: hashtagResults, isLoading: isHashtagLoading } = useSearchHashtags(debouncedQuery);
  const { data: trending, isLoading: trendingLoading } = useTrendingHashtags(4);
  const { data: suggested, isLoading: suggestedLoading } = useSuggestedUsers(3);
  const { user: currentUser } = useAuth();

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (searchRef.current && !searchRef.current.contains(event.target as Node)) {
        setIsFocused(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [searchRef]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (query.trim()) {
      router.push(`/search?q=${encodeURIComponent(query)}`);
      setIsFocused(false);
      setQuery(''); // Clear input after search
    }
  };

  const subscribeToPremium = () => {
    router.push('/premium');
  };

  // Helper component for suggested users list
  const SuggestedUsersList = () => (
    <div className="rounded-2xl border border-border bg-card overflow-hidden pt-3 mb-20">
      <h2 className="px-4 pb-2 font-extrabold text-[20px]">
        Who to follow
      </h2>
      {suggestedLoading ? (
        <div className="px-4 pb-3 space-y-4 pt-2">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="h-12 rounded-lg bg-secondary/50 animate-pulse"
            />
          ))}
        </div>
      ) : (
        <div className="px-0">
          {suggested?.content
            ?.filter((u) => u.id !== currentUser?.id)
            .slice(0, 3)
            .map((u) => (
              <div
                key={u.id}
                className="flex items-center gap-3 px-4 py-3 hover:bg-secondary/30 transition-colors cursor-pointer"
              >
                <Link href={`/${u.username}`} className="shrink-0">
                  <Avatar className="w-10 h-10">
                    <AvatarImage src={u.avatarUrl ?? undefined} />
                    <AvatarFallback>{u.displayName[0]}</AvatarFallback>
                  </Avatar>
                </Link>
                <div className="flex-1 min-w-0">
                  <Link
                    href={`/${u.username}`}
                    className="font-bold text-foreground hover:underline block truncate text-[15px] leading-5"
                  >
                    {u.displayName}
                  </Link>
                  <p className="text-muted-foreground text-[15px] truncate leading-5">
                    @{u.username}
                  </p>
                </div>
                <FollowButton userId={u.id} isFollowing={u.followedByMe} />
              </div>
            ))}
          <Link href="/connect_people" className="block px-4 py-3 text-primary text-[15px] hover:bg-secondary/30 transition-colors rounded-b-2xl">
            Show more
          </Link>
        </div>
      )}
    </div>
  );

  // Helper component for Premium subscription
  const PremiumCard = () => (
    <div className="rounded-2xl border border-border bg-card mb-4 p-4 mt-2">
      <h2 className="font-bold text-xl mb-2">Subscribe to Premium</h2>
      <p className="text-[15px] text-muted-foreground mb-3 leading-5">Subscribe to unlock new features and if eligible, receive a share of ads revenue.</p>
      <Button className="rounded-full font-bold px-4 cursor-pointer" size="sm" onClick={subscribeToPremium}>Subscribe</Button>
    </div>
  );

  // Helper component for Trending
  const TrendingList = () => (
    <div className="rounded-2xl border border-border bg-card mb-4 overflow-hidden pt-3">
      <h2 className="px-4 pb-2 font-extrabold text-[20px]">
        What&apos;s happening
      </h2>
      {trendingLoading ? (
        <div className="px-4 pb-3 space-y-4 pt-2">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="h-10 rounded-lg bg-secondary/50 animate-pulse"
            />
          ))}
        </div>
      ) : (
        <div className="px-0">
          {trending?.slice(0, 5).map((item) => (
            <Link
              key={item.hashtag}
              href={`/search?q=${encodeURIComponent('#' + item.hashtag)}`}
              className="block px-4 py-3 hover:bg-secondary/30 transition-colors"
            >
              <div className="flex justify-between items-start">
                <div>
                  <p className="text-muted-foreground text-[13px] leading-4">Trending in generic</p>
                  <p className="text-foreground font-bold text-[15px] pt-0.5">
                    #{item.hashtag}
                  </p>
                  <p className="text-muted-foreground text-[13px] leading-4 pt-0.5">
                    {item.recentCount} posts
                  </p>
                </div>
                <MoreHorizontal className="w-4 h-4 text-muted-foreground" />
              </div>
            </Link>
          ))}
          <Link href="/explore" className="block px-4 py-3 text-primary text-[15px] hover:bg-secondary/30 transition-colors rounded-b-2xl">
            Show more
          </Link>
        </div>
      )}
    </div>
  );

  // If on search page, hide search bar but show other widgets
  if (pathname === '/search' || pathname?.startsWith('/search')) {
    return (
      <div className="sticky top-0 p-4 h-screen overflow-y-auto no-scrollbar pb-10">
        <PremiumCard />
        <TrendingList />
        <SuggestedUsersList />
        <Footer />
      </div>
    );
  }

  return (
    <div className="sticky top-0 p-4 py-2 h-screen overflow-y-auto no-scrollbar pb-10">
      {/* Search Bar */}
      <div className="mb-4 sticky top-0 bg-background z-20 py-1" ref={searchRef}>
        <form onSubmit={handleSearch} className="relative">
          <div className={`flex items-center gap-3 rounded-full bg-secondary/50 border border-transparent px-4 py-2.5 transition-colors ${
            isFocused ? 'bg-background border-primary ring-1 ring-primary' : ''
          }`}>
            <Search className={`w-5 h-5 shrink-0 ${isFocused ? 'text-primary' : 'text-muted-foreground'}`} />
            <input
              type="search"
              placeholder="Search"
              className="flex-1 min-w-0 bg-transparent text-foreground text-[15px] placeholder:text-muted-foreground outline-none border-none leading-5"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              onFocus={() => setIsFocused(true)}
            />
          </div>

          {/* Search Dropdown */}
          {isFocused && query && (
            <div className="absolute top-full left-0 right-0 mt-1 bg-black border border-[#2f3336] rounded-lg shadow-xl overflow-hidden min-h-[100px] z-50 max-h-[400px] overflow-y-auto">
              {isSearchLoading || isHashtagLoading ? (
                <div className="p-4 text-center text-muted-foreground">Loading...</div>
              ) : (
                <ul>
                  {searchResults?.content.length === 0 && hashtagResults?.length === 0 && (
                    <li className="p-4 text-muted-foreground text-center">No results found</li>
                  )}

                  {/* People Results */}
                  {searchResults && searchResults.content.length > 0 && (
                    <li className="px-4 py-2 text-sm text-muted-foreground font-bold bg-secondary/30">People</li>
                  )}
                  {searchResults?.content.map(u => (
                    <li key={u.id}>
                      <Link
                        href={`/${u.username}`}
                        className="flex items-center gap-3 px-4 py-3 hover:bg-[#eff3f41a] transition-colors"
                        onClick={() => setIsFocused(false)}
                      >
                        <Avatar className="w-10 h-10">
                          <AvatarImage src={u.avatarUrl ?? undefined} />
                          <AvatarFallback>{u.displayName[0]}</AvatarFallback>
                        </Avatar>
                        <div className="flex-1 min-w-0">
                          <div className="font-bold text-[#e7e9ea] truncate">{u.displayName}</div>
                          <div className="text-[#71767b] truncate">@{u.username}</div>
                        </div>
                      </Link>
                    </li>
                  ))}

                  {/* Hashtag Results */}
                  {hashtagResults && hashtagResults.length > 0 && (
                    <li className="px-4 py-2 text-sm text-muted-foreground font-bold bg-secondary/30 border-t border-[#2f3336]">Hashtags</li>
                  )}
                  {hashtagResults?.map((h: TrendingHashtagDTO) => (
                    <li key={h.hashtag}>
                      <Link
                        href={`/search?q=${encodeURIComponent('#' + h.hashtag)}`}
                        className="flex items-center gap-3 px-4 py-3 hover:bg-[#eff3f41a] transition-colors"
                        onClick={() => setIsFocused(false)}
                      >
                        <div className="w-10 h-10 rounded-full bg-[#1d9bf0]/10 flex items-center justify-center">
                          <Search className="w-5 h-5 text-[#1d9bf0]" />
                        </div>
                        <div className="flex-1 min-w-0">
                          <div className="font-bold text-[#e7e9ea] truncate">#{h.hashtag}</div>
                          <div className="text-[#71767b] truncate">{h.recentCount} posts</div>
                        </div>
                      </Link>
                    </li>
                  ))}

                  <li className="border-t border-[#2f3336]">
                    <Link
                      href={`/search?q=${encodeURIComponent(query)}`}
                      className="block px-4 py-3 text-[#1d9bf0] hover:bg-[#eff3f41a] transition-colors"
                      onClick={() => setIsFocused(false)}
                    >
                      Search for "{query}"
                    </Link>
                  </li>
                </ul>
              )}
            </div>
          )}
        </form>
      </div>

      <PremiumCard />
      <TrendingList />
      <SuggestedUsersList />
      <Footer />
    </div>
  );
}

// Helper Components
function Footer() {
  return (
    <div className="px-4 text-[13px] text-muted-foreground space-x-2 leading-4">
      <span>Terms of Service</span>
      <span>Privacy Policy</span>
      <span>Cookie Policy</span>
      <span>Accessibility</span>
      <span>Ads info</span>
      <span>More ...</span>
      <p className="mt-1">© 2026 Chanom Corp.</p>
    </div>
  );
}

function MoreHorizontal(props: any) {
  return <svg {...props} xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="1"/><circle cx="19" cy="12" r="1"/><circle cx="5" cy="12" r="1"/></svg>
}
