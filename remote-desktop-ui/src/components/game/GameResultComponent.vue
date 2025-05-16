<template>
  <v-row
    v-if="props.wonFlag !== 0"
    class="d-flex justify-center mt-6"
  >
    <v-col
      cols="12"
      class="text-center"
      style="max-width: 500px"
    >
      <v-alert
        :type="currentStatus?.type"
        colored-border
        elevation="2"
      >
        {{ currentStatus?.title }}
      </v-alert>
      <v-btn
        class="mt-4"
        color="primary"
        @click="props.doResetGame"
      >
        🔄 Начать заново
      </v-btn>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
const props = defineProps({
  wonFlag: Number,
  mySymbol: String,
  doResetGame: Function,
});
const currentStatus = ref({})
const gameStatus = () => {
  if (
    (props.wonFlag === 1 && props.mySymbol === 'X') ||
    (props.wonFlag === -1 && props.mySymbol === 'O')
  ) {
    return {
      type: "success",
      title: "Выигрыш!"
    }
  }

  if (
    props.wonFlag === -2
  ) {
    return {
      type: "warning",
      title: "Ничья!"
    }
  }

  return {
    type: "error",
    title: "Проигрыш!"
  }
}

onUpdated(() => {
  currentStatus.value = gameStatus()
})
</script>
