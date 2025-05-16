<template>
  <v-container
    max-width="800px"
    class="pa-0"
  >
    <v-divider class="my-8" />
    <v-row
      class="d-flex justify-center"
      style="font-family: 'Roboto', sans-serif"
    >
      {{ $t('room.board_size', [rowsAndColumns]) }}
    </v-row>
    <v-divider class="mb-4 mt-8" />
    <v-row class="d-flex justify-center">
      <v-col
        cols="4"
      >
        <v-btn
          color="red"
          block
          :disabled="props.rowsAndColumns === 3"
          @click="props.resizeBoard(-1)"
        >
          <v-icon>mdi-minus</v-icon>
        </v-btn>
      </v-col>
      <v-col
        cols="4"
      >
        <v-btn
          color="green"
          block
          :disabled="props.rowsAndColumns === 6"
          @click="props.resizeBoard(1)"
        >
          <v-icon>mdi-plus</v-icon>
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";

const props = defineProps({
  gameStarted: Number,
  rowsAndColumns: Number,
  roomInfo: Object,
  resizeBoard: Function,
});
</script>
