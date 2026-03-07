import { type ChangeEvent } from "react";
import { Box, Grid, Input, Text } from "@chakra-ui/react";

type RegisterAccountFieldsProps = {
  name: string;
  onNameChange: (value: string) => void;
  email: string;
  onEmailChange: (value: string) => void;
  password: string;
  onPasswordChange: (value: string) => void;
};

export default function RegisterAccountFields({
  name,
  onNameChange,
  email,
  onEmailChange,
  password,
  onPasswordChange,
}: RegisterAccountFieldsProps) {
  return (
    <>
      <Grid gap={4} templateColumns={{ base: "1fr" }}>
        <Box minW={0}>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Nome completo
          </Text>
          <Input
            id="name"
            name="name"
            value={name}
            onChange={(event: ChangeEvent<HTMLInputElement>) =>
              onNameChange(event.target.value)
            }
            h="11"
            borderRadius="xl"
            bg="gray.50"
            borderColor="gray.200"
            focusRingColor="teal.200"
            placeholder="Seu nome"
            required
            w="full"
          />
        </Box>
      </Grid>

      <Grid gap={4} templateColumns={{ base: "1fr", sm: "1fr 1fr" }}>
        <Box minW={0}>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Email
          </Text>
          <Input
            id="email"
            name="email"
            type="email"
            value={email}
            onChange={(event: ChangeEvent<HTMLInputElement>) =>
              onEmailChange(event.target.value)
            }
            h="11"
            borderRadius="xl"
            bg="gray.50"
            borderColor="gray.200"
            focusRingColor="teal.200"
            placeholder="voce@email.com"
            required
            w="full"
          />
        </Box>

        <Box minW={0}>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Senha
          </Text>
          <Input
            id="password"
            name="password"
            type="password"
            value={password}
            onChange={(event: ChangeEvent<HTMLInputElement>) =>
              onPasswordChange(event.target.value)
            }
            h="11"
            borderRadius="xl"
            bg="gray.50"
            borderColor="gray.200"
            focusRingColor="teal.200"
            placeholder="********"
            required
            autoComplete="new-password"
            w="full"
          />
        </Box>
      </Grid>
    </>
  );
}
