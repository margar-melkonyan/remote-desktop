import { createI18n } from "vue-i18n";
import ru_translation from "@/lang/ru"

const i18n = createI18n({
  locale: 'ru',
  fallbackLocale: 'ru',
  messages: {
    ru: {
      ...ru_translation
    },
  }
})

export { i18n }
