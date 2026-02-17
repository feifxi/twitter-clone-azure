'use client';

import { useState } from 'react';
import { Search, Settings, MailPlus, MoreHorizontal, Image, Smile, Send } from 'lucide-react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { toast } from 'sonner';

// Mock Data
const MOCK_CONVERSATIONS = [
  {
    id: 1,
    user: {
      name: 'Elon Musk',
      handle: 'elonmusk',
      avatar: 'https://pbs.twimg.com/profile_images/1780044485541699584/p78MCn3B_400x400.jpg',
    },
    lastMessage: 'โอนให้ผม 20 wallet เพื่อปลดล็อค premium   ',
    timestamp: '2h',
  },
  {
    id: 2,
    user: {
      name: 'พี่เต้ พระราม7',
      handle: 'tae77',
      avatar: 'https://static.amarintv.com/images/upload/editor/source/wonder/022566/466/346800379_621542926678487_574.jpg',
    },
    lastMessage: 'คุณสนใจเข้าร่วมสภาเจได และเป็นครูฝึกแร็พเตอร์หรือไม่',
    timestamp: '1d',
  },

    {
    id: 4,
    user: {
      name: 'Tim Cook',
      handle: 'tim_cook',
      avatar: 'https://pbs.twimg.com/profile_images/1535420431766671360/Pwq-1eJc_400x400.jpg',
    },
    lastMessage: 'we detect you watch ferry porn in our iphone',
    timestamp: '5d',
  },
];

export default function MessagesPage() {
  const [selectedConversation, setSelectedConversation] = useState<number | null>(null);
  const [messageInput, setMessageInput] = useState('');

  const currentConversation = MOCK_CONVERSATIONS.find(c => c.id === selectedConversation);

  const handleSendMessage = () => {
    if (!messageInput.trim()) return;
    toast.info("Messaging feature is coming soon! อิอิ");
    setMessageInput('');
  };

  return (
    <div className="flex h-screen max-h-screen overflow-hidden">
      {/* Conversation List (Hidden on mobile when chat is open) */}
      <div className={`flex-1 md:flex-[0.4] border-r border-[#2f3336] flex flex-col ${selectedConversation ? 'hidden md:flex' : 'flex'}`}>
        
        {/* Header */}
        <div className="flex items-center justify-between px-4 py-3 sticky top-0 bg-black/80 backdrop-blur-md z-10">
          <h1 className="text-xl font-bold text-[#e7e9ea]">Messages</h1>
          <div className="flex gap-2">
            <Button variant="ghost" size="icon" className="rounded-full hover:bg-[#eff3f41a]">
              <Settings className="w-5 h-5" />
            </Button>
            <Button variant="ghost" size="icon" className="rounded-full hover:bg-[#eff3f41a]">
              <MailPlus className="w-5 h-5" />
            </Button>
          </div>
        </div>

        {/* Search */}
        <div className="px-4 pb-2">
          <div className="relative group">
            <div className="absolute inset-y-0 left-3 flex items-center pointer-events-none">
              <Search className="w-4 h-4 text-[#71767b] group-focus-within:text-[#1d9bf0]" />
            </div>
            <input 
              type="text" 
              placeholder="Search Direct Messages" 
              className="w-full bg-[#202327] text-[#e7e9ea] rounded-full py-2 pl-10 pr-4 outline-none border border-transparent focus:border-[#1d9bf0] focus:bg-black placeholder-[#71767b] transition-all text-[15px]"
            />
          </div>
        </div>

        {/* List */}
        <div className="flex-1 overflow-y-auto">
          {MOCK_CONVERSATIONS.map((conv) => (
            <div 
              key={conv.id}
              onClick={() => setSelectedConversation(conv.id)}
              className={`flex gap-3 px-4 py-3 cursor-pointer transition-colors hover:bg-[#eff3f41a] ${selectedConversation === conv.id ? 'border-r-2 border-[#1d9bf0] bg-[#eff3f41a]' : ''}`}
            >
              <Avatar className="w-10 h-10">
                <AvatarImage src={conv.user.avatar} />
                <AvatarFallback>{conv.user.name[0]}</AvatarFallback>
              </Avatar>
              <div className="flex-1 min-w-0">
                <div className="flex justify-between items-baseline">
                  <div className="truncate text-[#e7e9ea] font-bold text-[15px]">
                    {conv.user.name} 
                    <span className="font-normal text-[#71767b] ml-1">@{conv.user.handle}</span>
                  </div>
                  <span className="text-[#71767b] text-[13px] whitespace-nowrap ml-1">{conv.timestamp}</span>
                </div>
                <div className="text-[#71767b] text-[15px] truncate">
                  {conv.lastMessage}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Chat Interface (Hidden on mobile when no conversation selected) */}
      <div className={`flex-1 flex flex-col ${!selectedConversation ? 'hidden md:flex' : 'flex'}`}>
        {selectedConversation && currentConversation ? (
          <>
            {/* Chat Header */}
            <div className="flex items-center justify-between px-4 py-3 border-b border-[#2f3336] bg-black/80 backdrop-blur-md sticky top-0 z-10">
               <div className="flex items-center gap-3">
                 {/* Back button on mobile */}
                 <Button 
                    variant="ghost" 
                    size="icon" 
                    className="md:hidden -ml-2 rounded-full"
                    onClick={() => setSelectedConversation(null)}
                 >
                    <svg viewBox="0 0 24 24" aria-hidden="true" className="w-5 h-5 fill-current"><g><path d="M7.414 13l5.043 5.04-1.414 1.42L3.586 12l7.457-7.46 1.414 1.42L7.414 11H21v2H7.414z"></path></g></svg>
                 </Button>
                 <div className="flex flex-col">
                   <span className="font-bold text-[17px] text-[#e7e9ea]">{currentConversation.user.name}</span>
                   <span className="text-[13px] text-[#71767b]">@{currentConversation.user.handle}</span>
                 </div>
               </div>
               <Button variant="ghost" size="icon" className="rounded-full text-[#eff3f4]">
                 <MoreHorizontal className="w-5 h-5" />
               </Button>
            </div>

            {/* Chat History (Mock) */}
            <div className="flex-1 overflow-y-auto p-4 flex flex-col gap-4">
               {/* Welcome Message */}
               <div className="flex flex-col items-center justify-center my-8 text-center">
                  <Avatar className="w-16 h-16 mb-2">
                    <AvatarImage src={currentConversation.user.avatar} />
                    <AvatarFallback>{currentConversation.user.name[0]}</AvatarFallback>
                  </Avatar>
                  <h2 className="font-bold text-lg text-[#e7e9ea]">{currentConversation.user.name}</h2>
                  <p className="text-[#71767b]">@{currentConversation.user.handle}</p>
                  <p className="text-[#71767b] mt-2 text-sm">Joined January 2026 • 20 Followers</p>
               </div>
               
               {/* Inbound Message */}
               <div className="self-start max-w-[70%] bg-[#2f3336] text-white px-4 py-3 rounded-t-2xl rounded-r-2xl rounded-bl-none text-[15px]">
                 {currentConversation.lastMessage}
               </div>

               {/* Timestamp */}
               <div className="self-center text-[#71767b] text-xs">
                 {currentConversation.timestamp} ago
               </div>
            </div>

            {/* Input Area */}
            <div className="p-3 border-t border-[#2f3336]">
              <div className="bg-[#202327] rounded-2xl flex items-center px-2 py-1">
                 <Button variant="ghost" size="icon" className="text-[#1d9bf0] hover:bg-[#1d9bf0]/10 rounded-full w-9 h-9">
                    <Image className="w-5 h-5" />
                 </Button>
                 <Button variant="ghost" size="icon" className="text-[#1d9bf0] hover:bg-[#1d9bf0]/10 rounded-full w-9 h-9">
                    <Smile className="w-5 h-5" />
                 </Button>
                 <Input 
                    placeholder="Start a new message" 
                    className="flex-1 border-none bg-transparent focus-visible:ring-0 text-[#e7e9ea] placeholder-[#71767b] text-[15px]" 
                    value={messageInput}
                    onChange={(e) => setMessageInput(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && handleSendMessage()}
                 />
                 <Button 
                    variant="ghost" 
                    size="icon" 
                    className={`rounded-full w-9 h-9 transition-colors ${messageInput.trim() ? 'text-[#1d9bf0] hover:bg-[#1d9bf0]/10' : 'text-[#71767b] cursor-default hover:bg-transparent'}`}
                    onClick={handleSendMessage}
                    disabled={!messageInput.trim()}
                 >
                    <Send className="w-5 h-5" />
                 </Button>
              </div>
            </div>
          </>
        ) : (
          <div className="flex-1 flex flex-col items-center justify-center p-8 text-center h-full">
            <h2 className="text-3xl font-bold text-[#e7e9ea] mb-2">Select a message</h2>
            <p className="text-[#71767b] max-w-[350px]">Choose from your existing conversations, start a new one, or just keep swimming.</p>
            <Button className="mt-6 rounded-full bg-[#1d9bf0] hover:bg-[#1a8cd8] text-white font-bold px-8 py-6 text-lg">
               New message
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}
