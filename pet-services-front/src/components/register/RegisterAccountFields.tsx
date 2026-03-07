import { type ChangeEvent } from "react";
import { Box, Grid, Input, NativeSelect, Text } from "@chakra-ui/react";

import { type UserType, UserTypes } from "@/domain";

type RegisterAccountFieldsProps = {
  name: string;
  onNameChange: (value: string) => void;
  userType: UserType;
  onUserTypeChange: (value: UserType) => void;
  email: string;
  onEmailChange: (value: string) => void;
  password: string;
  onPasswordChange: (value: string) => void;
};

export default function RegisterAccountFields({
  name,
  onNameChange,
  userType,
  onUserTypeChange,
  email,
  onEmailChange,
  password,
  onPasswordChange,
}: RegisterAccountFieldsProps) {
  return (
    <>
      <Grid gap={4} templateColumns={{ base: "1fr", sm: "1.3fr 0.7fr" }}>
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

        <Box minW={0}>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Tipo de usuário
          </Text>
          <NativeSelect.Root
            size="md"
            h="11"
            borderRadius="xl"
            bg="gray.50"
            borderColor="gray.200"
            focusRingColor="teal.200"
            w="full"
          >
            <NativeSelect.Field
              id="userType"
              name="userType"
              value={userType}
              onChange={(event: ChangeEvent<HTMLSelectElement>) =>
                onUserTypeChange(event.target.value as UserType)
              }
            >
              <option value={UserTypes.Owner}>Tutor</option>
              <option value={UserTypes.Provider}>Prestador</option>
            </NativeSelect.Field>
            <NativeSelect.Indicator />
          </NativeSelect.Root>
        </Box>
      </Grid>

      <Grid gap={4} templateColumns={{ base: "1fr", sm: "1.3fr 0.7fr" }}>
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
