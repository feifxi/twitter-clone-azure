'use client';

import { XLogo } from '@/components/XLogo';
import { Button } from '@/components/ui/button';
import { Check } from 'lucide-react';

export default function PremiumPage() {
  const features = [
    'เป็นที่รักของแก๊ง call center',
    'โพสต์ยาวได้ แต่ไม่มีใครเห็นโพสต์',
    'ได้ blue check "this mf paid for Twitter"',
    'เพื่อนไม่คบ',
    'พ่อแม่ย้ายร่าง',
    'อย่ากดถ้าไม่อยากโดนย่ำกระหม่อม',
  ];

  const handleSubscribe = () => {
    alert('ผมรู้ IP คุณแล้ว เตรียมตัวโดนเช็คอิน');
  };

  return (
    <div className="min-h-screen bg-black text-white p-4 overflow-y-auto no-scrollbar pb-20">
      <div className="max-w-2xl mx-auto flex flex-col items-center text-center mt-10">
        <h1 className="text-4xl md:text-5xl font-extrabold mb-4">Who are you?</h1>
        <div className="relative mb-8">
            <XLogo className="w-16 h-16 fill-white" />
            <div className="absolute -top-2 -right-2 bg-[#1d9bf0] rounded-full p-1">
                <Check className="w-4 h-4 text-white" strokeWidth={4} />
            </div>
        </div>
        
        <h2 className="text-2xl font-bold mb-2">Premium</h2>
        <p className="text-[#71767b] mb-8">Detailed verification. Exclusive features.</p>
        
        <div className="grid grid-cols-1 gap-4 w-full mb-8 max-w-md">
            {/* Life Time Plan */}
             <div className="border border-[#1d9bf0] rounded-2xl p-6 relative bg-black flex flex-col items-center cursor-pointer hover:bg-[#eff3f41a] transition-colors">
                 <div className="absolute -top-3 bg-[#1d9bf0] text-white text-xs font-bold px-2 py-1 rounded">SAVE 80%</div>
                 <h3 className="font-bold text-xl mb-1">Life Time</h3>
                 <p className="text-2xl font-bold mb-4">$99<span className="text-sm font-normal text-[#71767b]">/life time</span></p>
                 <Button 
                    className="cursor-pointer w-full rounded-full font-bold bg-white text-black hover:bg-white/90" 
                    onClick={handleSubscribe}
                 >
                    Subscribe
                 </Button>
             </div>
        </div>

        <div className="w-full text-left">
            <h3 className="font-bold text-xl mb-4">Everything in Premium</h3>
            <ul className="space-y-4">
                {features.map((feature) => (
                    <li key={feature} className="flex items-center gap-3">
                        <div className="bg-[#1d9bf0] rounded-full p-1">
                             <Check className="w-3 h-3 text-white" strokeWidth={4} />
                        </div>
                        <span className="font-medium text-lg">{feature}</span>
                    </li>
                ))}
            </ul>
        </div>
        
        <div className="mt-12 text-[#71767b] text-sm">
            <p className="mb-2">Learn more about Premium and Verified Organizations</p>
        </div>
      </div>
    </div>
  );
}
