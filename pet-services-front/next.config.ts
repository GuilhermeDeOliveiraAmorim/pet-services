import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  outputFileTracingRoot: __dirname,
  experimental: {
    optimizePackageImports: ["@chakra-ui/react"],
  },
  turbopack: {
    root: __dirname,
  },
};

export default nextConfig;
