<template>
  <div class="main-board">
    <div
      v-for="i in props.rowsAndColumns"
      :key="i"
      style="display: flex; width: 100%;"
    >
      <div
        v-for="j in props.rowsAndColumns"
        :key="j"
        :class="`grid-index-${i}-${j}`"
        :style="props.getCellStyle()"
        @click="props.makeStep(i, j)"
      >
        <span
          class="no-select coming-soon-regular"
          :style="props.getFontStyle(rowsAndColumns)"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue';

const props = defineProps({
  rowsAndColumns: Number,
  makeStep: Function,
  getCellStyle: Function,
  getFontStyle: Function,
});
</script>

<style>
.main-board {
  display: flex;
  flex-direction: column;
  width: 100%;
  max-width: 600px;
}
.coming-soon-regular {
  font-family: "Coming Soon", cursive;
  font-weight: 400;
  font-style: normal;
}

</style>
