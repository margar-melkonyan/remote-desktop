<template>
  <div class="d-flex justify-center wrapper ma-4 rounded-4">
    <div
      ref="displayRef"
      class="guacamole-display"
      tabindex="0"
      @contextmenu.prevent
    />
    <div class="d-flex align-self-center">
      <v-card
        v-if="connectionState !== 1"
        max-height="175"
        width="500"
        rounded="xl"
      >
        <v-row v-if="connectionState === 0" class="d-flex justify-center my-6">
          <v-progress-circular
            color="primary"
            indeterminate
            :size="56"
            :width="7"
          />
        </v-row>
        <span v-if="errorText != ''">
          <v-card-title>
            <span style="color: red"> {{ $t("connections.caution") }} </span>
          </v-card-title>
          <v-card-text style="color: grey">
            <span v-html="errorText" />
          </v-card-text>
        </span>
        <v-card-actions>
          <v-btn
            color="blue"
            variant="flat"
            block
            rounded="xl"
            @click="router.push({ name: 'connections' })"
          >
            {{ $t("back") }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from "vue";
import Guacamole from "@cyolosecurity/guacamole-common-js";
import { useRouter } from "vue-router";

const router = useRouter();

const displayRef = ref(null);
const errorText = ref("");
const errorCode = ref("");
const connectionState = ref(0);
let client: Guacamole.Client = null;
let tunnel: Guacamole.Tunnel = null;

const keyMap = {
  Enter: 0xff0d,
  Backspace: 0xff08,
  Tab: 0xff09,
  Escape: 0xff1b,
  ArrowUp: 0xff52,
  ArrowDown: 0xff54,
  ArrowLeft: 0xff51,
  ArrowRight: 0xff53,
  Shift: 0xffe1,
  ShiftRight: 0xffe2,
  Control: 0xffe3,
  Alt: 0xffe9,
  Meta: 0xffe7,
  " ": 0x0020, // Пробел
};

// Блокируемые браузерные комбинации
const blockedCombinations = [
  "Control+t",
  "Control+n",
  "Control+w",
  "Control+Tab",
  "Alt+ArrowLeft",
  "Alt+ArrowRight",
  "F1",
  "F5",
  "Control+r",
  "Meta+Space",
];

const handleKeyEvent = (e, pressed) => {
  if (!client) return;

  if (e.key === " " && (e.metaKey || e.getModifierState("Meta"))) {
    e.preventDefault();
    e.stopPropagation();
    return;
  }

  // Проверка на заблокированные комбинации
  const modifiers = [
    e.ctrlKey ? "Control" : null,
    e.altKey ? "Alt" : null,
    e.shiftKey ? "Shift" : null,
  ].filter(Boolean);

  const currentCombo = [...modifiers, e.key].join("+");

  if (blockedCombinations.includes(currentCombo)) {
    e.preventDefault();
    e.stopPropagation();
    return;
  }

  if (e.key in keyMap) {
    client.sendKeyEvent(pressed, keyMap[e.key]);
    e.preventDefault();
    return;
  }

  if (e.key.length === 1) {
    const keysym = e.key.charCodeAt(0);
    client.sendKeyEvent(pressed, keysym);
    e.preventDefault();
  }
};

let guacMouse = null;

const initMouse = () => {
  if (!client || !displayRef.value) return;

  const displayElement = client.getDisplay().getElement();
  guacMouse = new Guacamole.Mouse(displayElement);

  guacMouse.onmousedown =
    guacMouse.onmousemove =
    guacMouse.onmouseup =
      (state) => {
        client.sendMouseState(state, true);
      };
};

const getDisplayDimensions = () => {
  const rect = document.querySelector(".wrapper").getBoundingClientRect();
  return {
    width: Math.floor(rect.width * 0.9),
    height: Math.floor(rect.height * 0.9),
  };
};

onMounted(() => {
  if (!displayRef.value) return;
  const guacID = router.currentRoute.value.params.id;
  const token = window.localStorage.getItem("guac_token");
  const { width, height } = getDisplayDimensions();

  const websocketUrl =
    `ws://192.168.1.4:8080/guacamole/websocket-tunnel` +
    `?token=${token}` +
    `&GUAC_DATA_SOURCE=postgresql` +
    `&GUAC_ID=${guacID}` +
    `&GUAC_TYPE=c` +
    `&GUAC_WIDTH=${width}` +
    `&GUAC_HEIGHT=${height}` +
    `&GUAC_DPI=96` +
    `&GUAC_TIMEZONE=Europe/Moscow` +
    `&GUAC_AUDIO=audio/L8` +
    `&GUAC_AUDIO=audio/L16` +
    `&GUAC_IMAGE=image/jpeg` +
    `&GUAC_IMAGE=image/png` +
    `&GUAC_IMAGE=image/webp`;

  tunnel = new Guacamole.WebSocketTunnel(websocketUrl);

  client = new Guacamole.Client(tunnel);

  client.onerror = (error) => {
    if (error.code == 519) {
      errorText.value = "«DNS-поиск не удался (неверное имя хоста?)»";
      errorCode.value = error.code;
    }
  };
  tunnel.onerror = (error) => {
    console.log(error.code);
  };

  client.connect();
  displayRef.value.appendChild(client.getDisplay().getElement());

  // Отслеживаем состояние туннеля (WebSocket)
  tunnel.onstatechange = (tunnelState) => {
    connectionState.value = tunnelState;
    switch (tunnelState) {
      case Guacamole.Tunnel.State.CONNECTING:
        console.log("Туннель: подключение...");
        errorText.value = "Установка соединения...";
        break;

      case Guacamole.Tunnel.State.OPEN:
        console.log("Туннель: подключён!");
        errorText.value = "";
        break;

      case Guacamole.Tunnel.State.CLOSED:
        console.log("Туннель: закрыт.");
        errorText.value =
          "Соединение закрыто.<br>Не удалось подключиться!<br>Проверьте введенные данные для подключения";
        break;

      case Guacamole.Tunnel.State.UNSTABLE:
        console.error("Туннель: разрыв соединения!");
        errorText.value = "Ошибка сети. Переподключитесь.";
        break;
    }
  };

  const options = { passive: false };
  displayRef.value.addEventListener(
    "keydown",
    (e) => handleKeyEvent(e, 1),
    options,
  );
  displayRef.value.addEventListener(
    "keyup",
    (e) => handleKeyEvent(e, 0),
    options,
  );

  initMouse();
  displayRef.value.focus();
});

onUnmounted(() => {
  if (client) client.disconnect();
  if (tunnel) tunnel.disconnect();
});
</script>

<style>
.guacamole-display {
  height: 100vh;
  /* font-family: monospace; */
  user-select: none;
  -webkit-user-select: none;
  outline: none;
  cursor: none;
}
.guacamole-display:nth-child(1) {
  display: flex;
  justify-content: center;
  border-color: aqua;
}
</style>
