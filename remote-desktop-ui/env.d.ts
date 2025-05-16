/// <reference types="vite/client" />
/// <reference types="unplugin-vue-router/client" />
/// <reference types="vite-plugin-vue-layouts/client" />

interface ImportMetaEnv {
  readonly VITE_API_URL: string
  readonly VITE_PORT: string
  readonly VUE_APP_TITLE?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
