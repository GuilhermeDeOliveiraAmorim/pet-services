import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import ContactSection from "@/components/home/ContactSection";
import HeroSection from "@/components/home/HeroSection";
import HomeFooter from "@/components/home/HomeFooter";
import ServicesSection from "@/components/home/ServicesSection";
import TestimonialsSection from "@/components/home/TestimonialsSection";

export default function Home() {
  return (
    <PageWrapper className="gap-16">
      <MainNav />
      <HeroSection />
      <ServicesSection />
      <TestimonialsSection />
      <ContactSection />
      <HomeFooter />
    </PageWrapper>
  );
}
