/**
 * plugins/index.ts
 *
 * Automatically included in `./src/main.ts`
 */

// Plugins
import vuetify from './vuetify'
import pinia from '../stores'
import router from '../router'
import './axios'

// Types
import type { App } from 'vue'
import { i18n } from './i18n'
import Vue3Toastify, { type ToastContainerOptions } from 'vue3-toastify';
import { toast } from 'vue3-toastify'
import API from '@/api'

export function registerPlugins (app: App) {
  app
    .use(vuetify)
    .use(router)
    .use(pinia)
    .use(i18n)
    .use(Vue3Toastify, {
      autoClose: 5000,
      position: toast.POSITION.BOTTOM_RIGHT,
    } as ToastContainerOptions);

    app.config.globalProperties.$api = new API().api
}
