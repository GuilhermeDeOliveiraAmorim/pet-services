import MainNav from "@/components/common/MainNav";
import HeroSection from "@/components/home/HeroSection";
import SegmentBlocks from "@/components/home/SegmentBlocks";
import QuickActionsBar from "@/components/home/QuickActionsBar";
import HighlightsSection from "@/components/home/HighlightsSection";
import PromoBannerSection from "@/components/home/PromoBannerSection";
import PartnerCallToAction from "@/components/home/PartnerCallToAction";
import InstitutionalSection from "@/components/home/InstitutionalSection";
import Footer from "@/components/home/Footer";

export default function Home() {
  return (
    <>
      <MainNav />
      <HeroSection />
      <SegmentBlocks />
      <QuickActionsBar />
      <HighlightsSection />
      <PromoBannerSection />
      <PartnerCallToAction />
      <InstitutionalSection />
      <Footer />
    </>
  );
}
