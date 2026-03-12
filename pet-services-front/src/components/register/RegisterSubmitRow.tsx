import { Button, Flex } from "@chakra-ui/react";

type RegisterSubmitRowProps = {
  isPending: boolean;
  isRedirecting?: boolean;
};

export default function RegisterSubmitRow({
  isPending,
  isRedirecting = false,
}: RegisterSubmitRowProps) {
  const isDisabled = isPending || isRedirecting;

  return (
    <Flex justify="center">
      <Button
        type="submit"
        disabled={isDisabled}
        h="11"
        w="full"
        maxW="xs"
        borderRadius="full"
        bg="green.400"
        color="white"
        _hover={{ bg: "green.500" }}
        _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
      >
        {isPending
          ? "Criando..."
          : isRedirecting
            ? "Redirecionando..."
            : "Criar conta"}
      </Button>
    </Flex>
  );
}
