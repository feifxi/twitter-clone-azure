import { XLogo } from '../XLogo';

export function LoadingScreen() {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black">
      <div className="flex flex-col items-center">
        <XLogo className="w-20 h-20 text-[#e7e9ea] animate-pulse" />
      </div>
    </div>
  );
}
