<template>
  <v-card rounded="xl">
    <v-card-title class="mt-2">
      {{ props.session.name }}
    </v-card-title>
    <v-divider class="my-2" />
    <v-card-actions>
      <v-row class="d-flex justify-space-between px-5 my-2">
        <div class="align-self-center">
          {{ $t('connections.protocol', [props.session.protocol]) }}
        </div>
        <v-btn
          color="blue"
          variant="tonal"
          rounded="lg"
          @click="openConnection"
        >
          {{ $t('connections.connect') }}
        </v-btn>
      </v-row>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
const router = useRouter();

const props = defineProps({
  session: {
    identifier: String,
    connection: String,
    protocol: String,
  }
})

const openConnection = () => {
  router.push({ name:  'current-session', params: {id: props.session.identifier} })
}
</script>
