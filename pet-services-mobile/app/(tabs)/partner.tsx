import { SafeAreaView, StyleSheet, Text, View } from "react-native";

export default function PartnerHomeScreen() {
  return (
    <SafeAreaView style={styles.safeArea}>
      <View style={styles.container}>
        <Text style={styles.title}>Area do Parceiro</Text>
        <Text style={styles.subtitle}>
          Aqui voce pode concluir onboarding, cadastrar servicos e responder
          solicitacoes.
        </Text>

        <View style={styles.card}>
          <Text style={styles.cardTitle}>Checklist inicial</Text>
          <Text style={styles.cardText}>- Completar perfil profissional</Text>
          <Text style={styles.cardText}>- Adicionar local de atendimento</Text>
          <Text style={styles.cardText}>- Publicar servicos e precos</Text>
        </View>
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
    lineHeight: 22,
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
});
