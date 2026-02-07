import type { ReactNode } from "react";

type PageWrapperProps = {
  children: ReactNode;
  className?: string;
};

export default function PageWrapper({
  children,
  className = "",
}: PageWrapperProps) {
  return (
    <div className="min-h-screen bg-[#f7f9ff] text-slate-900">
      <div
        className={`mx-auto flex min-h-screen w-full max-w-6xl flex-col gap-10 px-6 py-10 lg:px-8 ${className}`.trim()}
      >
        {children}
      </div>
    </div>
  );
}
