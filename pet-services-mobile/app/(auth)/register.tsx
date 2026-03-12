import { Link, router, useLocalSearchParams } from "expo-router";
import { useMemo, useState } from "react";
import {
  Pressable,
  SafeAreaView,
  StyleSheet,
  Text,
  TextInput,
  View,
} from "react-native";

export default function RegisterScreen() {
  const { user_type: userTypeParam } = useLocalSearchParams<{ user_type?: string }>();
  const defaultUserType = useMemo(() => {
    return userTypeParam === "provider" ? "provider" : "owner";
  }, [userTypeParam]);

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [userType, setUserType] = useState<"owner" | "provider">(defaultUserType);

  const submit = async () => {
    // Placeholder para integrar com endpoint real de cadastro.
    console.log("register", { name, email, password, userType });
    router.replace("/(auth)/login");
  };

  return (
    <SafeAreaView style={styles.safeArea}>
      <View style={styles.container}>
        <Text style={styles.title}>Criar conta</Text>

        <TextInput
          value={name}
          onChangeText={setName}
          placeholder="Seu nome"
          style={styles.input}
        />

        <TextInput
          value={email}
          onChangeText={setEmail}
          autoCapitalize="none"
          keyboardType="email-address"
          placeholder="Seu email"
          style={styles.input}
        />

        <TextInput
          value={password}
          onChangeText={setPassword}
          secureTextEntry
          placeholder="Senha"
          style={styles.input}
        />

        <View style={styles.toggleRow}>
          <Pressable
            style={[
              styles.chip,
              userType === "owner" ? styles.chipActive : styles.chipInactive,
            ]}
            onPress={() => setUserType("owner")}
          >
            <Text style={userType === "owner" ? styles.chipLabelActive : styles.chipLabelInactive}>
              Tutor
            </Text>
          </Pressable>
          <Pressable
            style={[
              styles.chip,
              userType === "provider" ? styles.chipActive : styles.chipInactive,
            ]}
            onPress={() => setUserType("provider")}
          >
            <Text
              style={
                userType === "provider" ? styles.chipLabelActive : styles.chipLabelInactive
              }
            >
              Parceiro
            </Text>
          </Pressable>
        </View>

        <Pressable style={styles.primaryButton} onPress={submit}>
          <Text style={styles.primaryButtonLabel}>Cadastrar</Text>
        </Pressable>

        <Link href="/(auth)/login" asChild>
          <Pressable style={styles.linkButton}>
            <Text style={styles.linkButtonLabel}>Ja tenho conta</Text>
          </Pressable>
        </Link>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: {
    flex: 1,
    backgroundColor: "#F4F8F7",
  },
  container: {
    flex: 1,
    justifyContent: "center",
    gap: 12,
    paddingHorizontal: 24,
  },
  title: {
    fontSize: 32,
    fontWeight: "700",
    color: "#0F172A",
    marginBottom: 8,
  },
  input: {
    borderWidth: 1,
    borderColor: "#D1D5DB",
    borderRadius: 12,
    backgroundColor: "#FFFFFF",
    paddingHorizontal: 14,
    paddingVertical: 12,
    fontSize: 16,
  },
  toggleRow: {
    flexDirection: "row",
    gap: 10,
    marginTop: 4,
  },
  chip: {
    borderRadius: 999,
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderWidth: 1,
  },
  chipActive: {
    backgroundColor: "#0F766E",
    borderColor: "#0F766E",
  },
  chipInactive: {
    backgroundColor: "#FFFFFF",
    borderColor: "#D1D5DB",
  },
  chipLabelActive: {
    color: "#FFFFFF",
    fontWeight: "700",
  },
  chipLabelInactive: {
    color: "#334155",
    fontWeight: "600",
  },
  primaryButton: {
    marginTop: 8,
    backgroundColor: "#0F766E",
    borderRadius: 12,
    paddingVertical: 14,
    alignItems: "center",
  },
  primaryButtonLabel: {
    color: "#FFFFFF",
    fontWeight: "600",
    fontSize: 16,
  },
  linkButton: {
    alignItems: "center",
    paddingVertical: 8,
  },
  linkButtonLabel: {
    color: "#0F766E",
    fontWeight: "600",
  },
});
