import { defineConfig } from "vite";
import solid from "vite-plugin-solid";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig(({ mode }) => {
  const base = mode === "production" ? "/eefu/" : "/";
  return {
    plugins: [solid(), tailwindcss()],
    base,
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
  };
});
