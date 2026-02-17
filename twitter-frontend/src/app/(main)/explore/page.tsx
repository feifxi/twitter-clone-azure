'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Search, ArrowLeft } from 'lucide-react';
// ... imports
import { useTrendingHashtags } from '@/hooks/useDiscovery';
import Link from 'next/link';

export default function ExplorePage() {
  const router = useRouter();
  const [query, setQuery] = useState('');
  const { data: trending, isLoading: trendingLoading } = useTrendingHashtags(20);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (query.trim()) {
      router.push(`/search?q=${encodeURIComponent(query)}`);
    }
  };

  return (
    <div className="min-h-screen">
      {/* Search Bar Header */}
      <div className="sticky top-0 z-30 bg-black/60 backdrop-blur-md border-b border-[#2f3336] px-4 py-3 flex items-center gap-4">
        <div 
            onClick={() => router.back()} 
            className="p-2 rounded-full hover:bg-white/10 cursor-pointer transition-colors -ml-2"
        >
            <ArrowLeft className="w-5 h-5 text-white" />
        </div>
        <form onSubmit={handleSearch} className="relative flex-1">
          <div className="flex items-center gap-2 rounded-full bg-[#202327] border border-transparent focus-within:border-[#1d9bf0] focus-within:bg-transparent px-4 py-2.5 transition-colors">
            <Search className="w-5 h-5 text-[#71767b] shrink-0" />
            <input
              type="search"
              placeholder="Search Twitter"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              className="flex-1 min-w-0 bg-transparent text-[#e7e9ea] text-[15px] placeholder:text-[#71767b] outline-none border-none"
            />
          </div>
        </form>
      </div>

      {/* Trending List */}
      <h1 className="text-xl font-bold p-4 border-b border-[#2f3336]">Trends for you</h1>
      
      {trendingLoading ? (
         <div className="p-4 space-y-4">
            {[1, 2, 3, 4, 5].map((i) => (
                <div key={i} className="h-16 bg-[#2f3336]/30 rounded-lg animate-pulse" />
            ))}
         </div>
      ) : (
          <div>
              {trending?.map((item) => (
                <Link
                    key={item.hashtag}
                    href={`/search?q=${encodeURIComponent('#' + item.hashtag)}`}
                    className="block px-4 py-3 hover:bg-[#eff3f41a] transition-colors border-b border-[#2f3336]"
                >
                    <div className="text-[13px] text-[#71767b]">Trending</div>
                    <div className="font-bold text-[#e7e9ea] text-[16px]">#{item.hashtag}</div>
                    <div className="text-[13px] text-[#71767b]">{item.recentCount} posts</div>
                </Link>
              ))}
          </div>
      )}
    </div>
  );
}
