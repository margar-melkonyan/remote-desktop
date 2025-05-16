<template>
  <v-card
    class="pa-4"
  >
    <v-card-title
      class="d-flex justify-center"
    >
      {{ $t('splashes.room.wait_opponent') }}
    </v-card-title>
    <v-card-text
      class="d-flex justify-center"
      style="overflow: hidden"
    >
      <v-progress-circular
        indeterminate
        :size="50"
        :width="3"
      />
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">

</script>

