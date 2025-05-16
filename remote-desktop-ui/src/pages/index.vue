<template>
  <v-container
    max-width="800"
    max-height="100vh"
  >
    <v-row>
      <v-col cols="12">
        <v-card
          elevation="12"
          rounded="xl"
        >
          <v-img
            height="200"
            cover
            gradient="to bottom, rgba(0,0,0,0.1), rgba(0,0,0,0.5)"
          >
            <v-card-title
              class="text-h3 font-weight-bold text-white pt-12"
              style="
                word-break: break-word;
                overflow-wrap: break-word;
                white-space: normal;
                line-height: 1.2;
              "
            >
              Добро пожаловать в RemoteDesktop
            </v-card-title>
          </v-img>

          <v-card-text class="pa-6">
            <v-list lines="two">
              <v-list-item
                v-for="(feature, i) in features"
                :key="i"
                :prepend-icon="feature.icon"
              >
                <template #title>
                  <span class="font-weight-bold">{{ feature.title }}</span>
                </template>
                {{ feature.description }}
              </v-list-item>
            </v-list>

            <v-divider class="my-4" />

            <div class="text-center mt-6">
              <v-btn
                class="px-8"
                color="primary"
                size="large"
                rounded="lg"
                @click="startSession"
              >
                <v-icon start>mdi-remote-desktop</v-icon>
                Начать подключение
              </v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const features = [
  {
    icon: 'mdi-shield-check',
    title: 'Безопасное подключение',
    description: 'Все соединения защищены 256-битным шифрованием'
  },
  {
    icon: 'mdi-devices',
    title: 'Доступ с любых устройств',
    description: 'Работает на компьютерах, планшетах и смартфонах'
  }
]

const startSession = () => {
  router.push('/connect')
}
</script>
