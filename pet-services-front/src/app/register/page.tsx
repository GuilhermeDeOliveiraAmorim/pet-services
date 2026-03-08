import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import RegisterAside from "@/components/register/RegisterAside";
import RegisterForm from "@/components/register/RegisterForm";
import { Grid } from "@chakra-ui/react";

export default function RegisterPage() {
  return (
    <PageWrapper gap={16}>
      <MainNav showLinks={false} showActions={false} />
      <Grid
        w="full"
        gap={10}
        templateColumns={{ base: "1fr", lg: "0.9fr 1.1fr" }}
      >
        <RegisterAside />
        <RegisterForm />
      </Grid>
    </PageWrapper>
  );
}
