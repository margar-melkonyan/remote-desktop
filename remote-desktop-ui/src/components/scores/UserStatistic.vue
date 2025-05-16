<template>
  <v-card class="pa-4">
    <v-card-title class="d-flex justify-center">
      {{ $t('statistics.title') }}
    </v-card-title>
    <v-divider class="my-4" />
    <div style="max-height: 350px; overflow-y: auto;">
      <template v-if="scores.length > 0">
        <v-card-text
          v-for="(score, key) in scores"
          :key="`score-${score.id}-${key}`"
          class="py-2"
        >
          <div
            class="py-4 d-flex justify-center"
            :style="`background-color: ${labelWithColor(score).color}; border-radius: 0.5rem;`"
          >
            {{ labelWithColor(score).title }} {{ score.nickname }} {{ moment(score.created_at).utc().format('DD.MM.YYYY HH:mm:ss') }}
          </div>
        </v-card-text>
      </template>
      <v-card-text
        v-else
        class="d-flex justify-center"
      >
        {{ $t('statistics.unavailable') }}
      </v-card-text>
    </div>
    <v-divider class="my-4" />
    <v-card-actions>
      <v-btn
        variant="tonal"
        block
        @click="close"
      >
        {{ $t('close') }}
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script lang="ts" setup>
import { ref, getCurrentInstance, defineEmits } from 'vue';
import moment from "moment";
import axios from "axios";
const { proxy } = getCurrentInstance();
const apiScores = proxy.$api.scores;
const scores = ref([])
const emit = defineEmits([
  "closeDialog",
])
const labelWithColor = (score) => {
  if (score.is_won === 1) {
    return {title:"Выиграл у игрока", color: "green"}
  }
  if (score.is_won === 0) {
    return {title:"Проиграл игроку", color: "red"}
  }
  if (score.is_won === -1) {
    return {title:"Ничья c", color: "orange"}
  }
}
const close = () => {
  emit("closeDialog")
}
axios.get(apiScores.urls.scores())
  .then(({data}) => {
    scores.value = data.data
  })
</script>
