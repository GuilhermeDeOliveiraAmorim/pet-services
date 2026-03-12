import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import ContactSection from "@/components/home/ContactSection";
import HomeFooter from "@/components/home/HomeFooter";
import PartnerBenefitsSection from "@/components/partner/PartnerBenefitsSection";
import PartnerHeroSection from "@/components/partner/PartnerHeroSection";
import PartnerOnboardingSection from "@/components/partner/PartnerOnboardingSection";

export default function PartnerPage() {
  return (
    <PageWrapper gap={16}>
      <MainNav />
      <PartnerHeroSection />
      <PartnerBenefitsSection />
      <PartnerOnboardingSection />
      <ContactSection />
      <HomeFooter />
    </PageWrapper>
  );
}
