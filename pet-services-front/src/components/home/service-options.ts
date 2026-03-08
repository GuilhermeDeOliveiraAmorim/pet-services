import type { LucideIcon } from "lucide-react";
import {
  Bath,
  Dog,
  GraduationCap,
  HeartHandshake,
  Hospital,
  ShoppingBag,
  Trees,
} from "lucide-react";

export type ServiceOption = {
  label: string;
  icon: LucideIcon;
  href: string;
  iconColor: string;
  iconBg: string;
};

export const serviceOptions: ServiceOption[] = [
  {
    label: "Clínica Veterinária",
    icon: Hospital,
    href: "#servico-clinica-veterinaria",
    iconColor: "red.500",
    iconBg: "red.50",
  },
  {
    label: "Pet Shop",
    icon: ShoppingBag,
    href: "#servico-pet-shop",
    iconColor: "blue.500",
    iconBg: "blue.50",
  },
  {
    label: "Banho e Tosa",
    icon: Bath,
    href: "#servico-banho-e-tosa",
    iconColor: "teal.500",
    iconBg: "teal.50",
  },
  {
    label: "Hotelzinho e Creche",
    icon: Trees,
    href: "#servico-hotelzinho-e-creche",
    iconColor: "green.500",
    iconBg: "green.50",
  },
  {
    label: "Passeador(a)",
    icon: Dog,
    href: "#servico-passeador",
    iconColor: "orange.500",
    iconBg: "orange.50",
  },
  {
    label: "Pet Sitter",
    icon: HeartHandshake,
    href: "#servico-pet-sitter",
    iconColor: "pink.500",
    iconBg: "pink.50",
  },
  {
    label: "Adestrador",
    icon: GraduationCap,
    href: "#servico-adestrador",
    iconColor: "purple.500",
    iconBg: "purple.50",
  },
];