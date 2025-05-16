<template>
  <v-card
    class="py-4"
  >
    <v-card-text>
      <v-row class="d-flex align-center">
        <v-col cols="6">
          <v-card
            color="red"
            :style="selectedSymbol === 'X' ? 'border: #ff7fea solid 0.3rem' : ''"
            @click="choseSymbol('X')"
          >
            <v-card-text
              class="d-flex justify-center"
            >
              <span class="symbol-size">
                X
              </span>
            </v-card-text>
          </v-card>
        </v-col>
        <v-col cols="6">
          <v-card
            color="green"
            :style="selectedSymbol === 'O' ? 'border: #ff7fea solid 0.3rem' : ''"
            @click="choseSymbol('O')"
          >
            <v-card-text
              class="d-flex justify-center"
            >
              <span class="symbol-size">
                O
              </span>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-card-text>
    <v-card-actions>
      <v-row
        class="px-4"
        no-gutters
      >
        <v-btn
          block
          variant="tonal"
          :disabled="selectedSymbol === ''"
          @click="selectSymbol"
        >
          Выбрать
        </v-btn>
      </v-row>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import { ref, defineProps } from "vue";
const props = defineProps<{
  sendChooseSymbol: Function;
}>();
const selectedSymbol = ref<string>("")
const choseSymbol = (symbol:string) => {
  selectedSymbol.value = symbol
}
const selectSymbol = () => {
  props.sendChooseSymbol(selectedSymbol.value)
}
</script>

<style scoped>
.symbol-size {
  font-size: 90px;
}
</style>
