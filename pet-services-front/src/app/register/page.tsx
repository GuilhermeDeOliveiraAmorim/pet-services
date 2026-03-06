import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import RegisterAside from "@/components/register/RegisterAside";
import RegisterForm from "@/components/register/RegisterForm";

export default function RegisterPage() {
  return (
    <PageWrapper gap={16}>
      <MainNav showLinks={false} showActions={false} />
      <div className="grid w-full gap-10 lg:grid-cols-[0.9fr_1.1fr]">
        <RegisterAside />
        <RegisterForm />
      </div>
    </PageWrapper>
  );
}
