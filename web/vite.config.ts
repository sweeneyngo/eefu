import { defineConfig } from "vite";
import solid from "vite-plugin-solid";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [solid(), tailwindcss()],
  base: "/eefu/",
  resolve: {
    alias: {
      "@": "/src",
      api: "/src/api",
      components: "/src/components",
      resources: "/src/resources",
      pages: "/src/pages",
      types: "/src/types",
      utils: "/src/utils",
    },
  },
});
