import MainNav from "@/components/common/MainNav";
import HeroSection from "@/components/home/HeroSection";
import SegmentBlocks from "@/components/home/SegmentBlocks";
import QuickActionsBar from "@/components/home/QuickActionsBar";
import HighlightsSection from "@/components/home/HighlightsSection";
import PartnerCallToAction from "@/components/home/PartnerCallToAction";
import InstitutionalSection from "@/components/home/InstitutionalSection";
import Footer from "@/components/home/Footer";

export default function Home() {
  return (
    <>
      <MainNav />
      <HeroSection />

      {/* Barra de ações rápidas logo após o hero */}
      <QuickActionsBar />

      {/* Blocos de segmentos/serviços */}
      <SegmentBlocks />

      {/* Destaques visuais (promoções, pets, serviços) */}
      <HighlightsSection />

      {/* Chamada para parceiros/ONGs/voluntários */}
      <PartnerCallToAction />

      {/* Links institucionais e redes sociais */}
      <InstitutionalSection />

      {/* Rodapé informativo */}
      <Footer />
    </>
  );
}
