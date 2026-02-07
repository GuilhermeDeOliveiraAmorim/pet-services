import type { ReactNode } from "react";
import RegisterFormHeader from "./RegisterFormHeader";

type RegisterFormCardProps = {
  children: ReactNode;
};

export default function RegisterFormCard({ children }: RegisterFormCardProps) {
  return (
    <div className="flex items-center justify-center">
      <div className="w-full max-w-2xl rounded-4xl bg-white p-10 shadow-[0_30px_80px_rgba(124,139,255,0.15)]">
        <RegisterFormHeader />
        {children}
      </div>
    </div>
  );
}
