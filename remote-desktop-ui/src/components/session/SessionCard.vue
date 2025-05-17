<template>
  <v-card rounded="xl">
    <v-card-title class="mt-2">
      {{ props.session.connection }} ({{ props.session.hostname }}:{{ props.session.port }})
    </v-card-title>
    <v-divider class="my-2" />
    <v-card-text>
      ({{ props.session.description }})
    </v-card-text>
    <v-divider class="my-2" />
    <v-card-actions>
      <v-row class="d-flex justify-space-between px-5 my-2">
        <div class="align-self-center">
          ({{ props.session.protocol }})
        </div>
        <v-btn
          color="blue"
          variant="tonal"
          rounded="lg"
          @click="openConnection"
        >
          {{ $t('sessions.connect') }}
        </v-btn>
      </v-row>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
const props = defineProps({
  session: {
    connection: String,
    hostname: String,
    port: Number,
    description: String,
    protocol: String,
  }
})
</script>
