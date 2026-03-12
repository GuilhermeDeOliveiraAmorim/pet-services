import { router } from "expo-router";
import { Pressable, SafeAreaView, StyleSheet, Text, View } from "react-native";

import { useAuthStore } from "../../src/features/auth/store";

export default function OwnerHomeScreen() {
  const session = useAuthStore((state) => state.session);
  const clearSession = useAuthStore((state) => state.clearSession);

  const handleLogout = () => {
    clearSession();
    router.replace("/(auth)/login");
  };

  return (
    <SafeAreaView style={styles.safeArea}>
      <View style={styles.container}>
        <Text style={styles.title}>Area do Tutor</Text>
        <Text style={styles.subtitle}>
          Bem-vindo, {session?.user.name ?? "Usuario"}
        </Text>

        <View style={styles.card}>
          <Text style={styles.cardTitle}>Proximos passos</Text>
          <Text style={styles.cardText}>- Listar servicos por localizacao</Text>
          <Text style={styles.cardText}>- Solicitar atendimento para pets</Text>
          <Text style={styles.cardText}>
            - Avaliar parceiros apos o servico
          </Text>
        </View>

        <Pressable
          style={styles.secondaryButton}
          onPress={() => router.push("/(tabs)/partner")}
        >
          <Text style={styles.secondaryButtonLabel}>Ver area do Parceiro</Text>
        </Pressable>

        <Pressable style={styles.primaryButton} onPress={handleLogout}>
          <Text style={styles.primaryButtonLabel}>Sair</Text>
        </Pressable>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: {
    flex: 1,
    backgroundColor: "#F8FAFC",
  },
  container: {
    flex: 1,
    padding: 24,
    gap: 16,
  },
  title: {
    fontSize: 28,
    fontWeight: "700",
    color: "#0F172A",
  },
  subtitle: {
    fontSize: 16,
    color: "#334155",
  },
  card: {
    marginTop: 8,
    padding: 16,
    borderRadius: 14,
    backgroundColor: "#FFFFFF",
    borderWidth: 1,
    borderColor: "#E2E8F0",
    gap: 6,
  },
  cardTitle: {
    fontSize: 16,
    fontWeight: "700",
    color: "#0F172A",
  },
  cardText: {
    fontSize: 14,
    color: "#334155",
  },
  secondaryButton: {
    marginTop: "auto",
    borderRadius: 12,
    paddingVertical: 14,
    alignItems: "center",
    borderWidth: 1,
    borderColor: "#0F766E",
  },
  secondaryButtonLabel: {
    color: "#0F766E",
    fontWeight: "700",
  },
  primaryButton: {
    borderRadius: 12,
    paddingVertical: 14,
    alignItems: "center",
    backgroundColor: "#0F766E",
  },
  primaryButtonLabel: {
    color: "#FFFFFF",
    fontWeight: "700",
  },
});
