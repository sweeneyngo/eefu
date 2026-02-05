/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly MODE: "development" | "production" | string;
  readonly VITE_API_URL?: string; // optional for your custom env vars
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
