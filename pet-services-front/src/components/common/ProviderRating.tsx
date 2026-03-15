import { Badge, type BadgeProps } from "@chakra-ui/react";

type ProviderRatingProps = Omit<BadgeProps, "children"> & {
  rating?: number | null;
  totalReviews?: number | null;
  showCount?: boolean;
  hideWhenZero?: boolean;
  labelPrefix?: string;
};

const formatRating = (value: number): string => {
  if (!Number.isFinite(value)) {
    return "0.0";
  }

  return value.toFixed(1);
};

export default function ProviderRating({
  rating,
  totalReviews,
  showCount = false,
  hideWhenZero = false,
  labelPrefix,
  borderRadius = "full",
  px = 3,
  py = 1,
  colorPalette = "cyan",
  ...badgeProps
}: ProviderRatingProps) {
  const numericRating = Number(rating ?? 0);

  if (hideWhenZero && numericRating <= 0) {
    return null;
  }

  const prefix = labelPrefix ? `${labelPrefix} ` : "";
  const count =
    showCount && typeof totalReviews === "number"
      ? ` · ${totalReviews} avaliações`
      : "";

  return (
    <Badge
      borderRadius={borderRadius}
      px={px}
      py={py}
      colorPalette={colorPalette}
      {...badgeProps}
    >
      {prefix}★ {formatRating(numericRating)}
      {count}
    </Badge>
  );
}
