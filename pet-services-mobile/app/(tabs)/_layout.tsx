import { Tabs } from "expo-router";

export default function TabsLayout() {
  return (
    <Tabs
      screenOptions={{
        headerShown: false,
        tabBarActiveTintColor: "#0F766E",
      }}
    >
      <Tabs.Screen
        name="owner"
        options={{
          title: "Tutor",
        }}
      />
      <Tabs.Screen
        name="partner"
        options={{
          title: "Parceiro",
        }}
      />
    </Tabs>
  );
}
