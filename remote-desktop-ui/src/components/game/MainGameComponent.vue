<template>
  <v-col>
    <v-dialog
      v-model="chooseSymbolDialog"
      width="500"
      persistent
    >
      <ChooseSplashScreen
        :send-choose-symbol="sendChooseSymbol"
      />
    </v-dialog>
    <v-dialog
      v-model="waitSymbolChoosing"
      width="500"
      persistent
    >
      <WaitSplashScreen />
    </v-dialog>
    <v-dialog
      v-model="isPrivate"
      width="500"
      persistent
    >
      <PrivatePasswordInput
        :room-id="roomInfo?.id"
        :room-name="roomInfo?.name"
        :wss-connect="wssConnect"
        :exit-from-room="exitFromRoom"
      />
    </v-dialog>
    <v-row class="d-flex justify-center">
      <HeaderComponent
        :versus="versus"
        :room-info="roomInfo"
        :exit-from-room="exitFromRoom"
      />
    </v-row>
    <v-row
      v-if="wonFlag === 0 && versus !== null"
      class="d-flex justify-center my-4"
    >
      <span
        v-if="mySymbol === currentPlayer"
        style="color: yellowgreen;"
      >
        Ваш ход
      </span>
      <span
        v-else
        style="color: red;"
      >
        Ходит {{ versus }}
      </span>
    </v-row>
    <v-row
      v-if="wonFlag !== 0 && versus !== null"
      class="d-flex justify-center my-4"
    >
      Игра окончена!
    </v-row>
    <v-row class="d-flex justify-center my-8">
      <GameBoardComponent
        :versus="versus"
        :rows-and-columns="rowsAndColumns"
        :make-step="makeStep"
        :get-cell-style="getCellStyle"
        :get-font-style="getFontStyle"
      />
    </v-row>

    <GameControlsComponent
      v-if="gameStarted === 0 && authStore.user?.id === roomInfo?.creator_id"
      :game-started="gameStarted"
      :rows-and-columns="rowsAndColumns"
      :resize-board="resizeBoard"
      :room-info="roomInfo"
    />

    <GameResultComponent
      :my-symbol="mySymbol"
      :won-flag="wonFlag"
      :do-reset-game="doResetGame"
    />
  </v-col>
</template>

<script setup lang="ts">
import axios from "axios";
import HeaderComponent from './HeaderComponent.vue';
import GameBoardComponent from './GameBoardComponent.vue';
import GameControlsComponent from './GameControlsComponent.vue';
import GameResultComponent from './GameResultComponent.vue';
import { useRouter } from "vue-router"
import { useAuthStore } from "@/stores/auth";
import {
  getCellStyle, getFontStyle, wonState, fetchRoom,
  resetGameBoardCells, resizeCountingArrays,
  mySymbol, currentPlayer, roomInfo, rowsAndColumns,
  wonFlag, gameStarted, versusFetchIntervalId,
  isPrivate, chooseSymbolDialog,
  waitSymbolChoosing, controller, } from "@/plugins/services/utils";
import { connectToRoom } from "@/plugins/services/websocketService";
import { getCurrentInstance, computed, watch } from 'vue';

const { proxy } = getCurrentInstance();
const apiRooms = proxy.$api.rooms;

const authStore = useAuthStore();
const router = useRouter();
const props = defineProps<{
  roomId: string;
}>();

let wss: WebSocket

const wssConnect = (roomID: string, password: string, apiRooms: object, authStore: object) => {
    wss = connectToRoom(roomID, password, apiRooms, authStore, router);
};

if (isPrivate.value === false) {
  wssConnect(props.roomId, "", apiRooms, authStore)
}

const versus = computed(() => {
  return roomInfo.value?.users?.filter((user) => user.name != authStore.user?.name)[0]?.name || null;
});
const sendChooseSymbol = (symbol:string) => {
  mySymbol.value = symbol
  currentPlayer.value = symbol
  chooseSymbolDialog.value = false
  waitSymbolChoosing.value = false
  wss.send(JSON.stringify({
    action: "select symbol",
    symbol: symbol,
  }))
}
function makeStep(i: number, j: number) {
  const cell = document.querySelector(`.grid-index-${i}-${j}>span`);
  if (
    cell && cell.textContent === '' && wonFlag.value === 0 && versus.value !== null &&
    currentPlayer.value == mySymbol.value
  ) {
    wss.send(JSON.stringify({
      data: {
        id: `${i}-${j}`,
        symbol: mySymbol.value,
      },
      action: "step",
    }))
  }
}
function doResetGame() {
  wss.send(JSON.stringify({
    action: 'reset game'
  }));
}
function resizeBoard(size: number) {
  rowsAndColumns.value += size;
  resizeCountingArrays();
  resetGameBoardCells();
  wss.send(JSON.stringify({
    action: "resize",
    size: rowsAndColumns.value
  }))
}
function exitFromRoom() {
  if (versus.value !== null &&
    authStore?.user.id === roomInfo.value.creator_id
  ) {
    wss.send(JSON.stringify({
      action: "close room"
    }))
    axios.delete(apiRooms.urls.room(props.roomId))
  }
  if (authStore?.user.id !== roomInfo.value.creator_id) {
    wss.send(JSON.stringify({
      action: "exit room",
    }))
  }
  router.push({ name: "index" })
}
onMounted(() => {
  versusFetchIntervalId.value = setInterval(() => {
    if (!versus.value) {
      fetchRoom(props.roomId, apiRooms, authStore);
    } else {
      clearInterval(versusFetchIntervalId.value);
    }
  }, 2000);
})
onBeforeUnmount(() => {
  clearInterval(versusFetchIntervalId.value)
  if (controller) controller.abort();
})
watch(wonFlag, function () {
  if (wonFlag.value !== 0) {
    wss.send(JSON.stringify({
      "action": "game end",
      "data": {
        "is_won": wonState(mySymbol.value, wonFlag.value),
        "user_id": authStore?.user?.id,
        "versus_player_nickname": versus.value
        }
    }))
  }
})
</script>
