<template>
  <div class="wrapper ma-4 rounded-4">
    <!-- Основной контейнер с относительным позиционированием -->
    <div class="position-relative" style="height: 100vh">
      <!-- Guacamole дисплей -->
      <div
        ref="displayRef"
        class="guacamole-display"
        tabindex="0"
        @contextmenu.prevent
      />

      <!-- Оверлей с v-card -->
      <div
        v-if="connectionState !== 1"
        class="position-absolute"
        style="
          top: 50%;
          left: 50%;
          transform: translate(-50%, -50%);
          z-index: 1000;
        "
      >
        <v-card max-height="175" width="500" rounded="xl">
          <v-row
            v-if="connectionState === 0"
            class="d-flex justify-center my-6"
          >
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
            <v-row>
              <v-col>
                <v-btn
                  color="blue"
                  variant="flat"
                  block
                  rounded="xl"
                  @click="router.push({ name: 'connections' })"
                  prepend-icon="mdi-arrow-left"
                >
                  {{ $t("back") }}
                </v-btn>
              </v-col>
              <v-col>
                <v-btn
                  color="blue"
                  variant="flat"
                  block
                  rounded="xl"
                  @click="runConnection"
                  prepend-icon="mdi-refresh"
                >
                  {{ $t("connections.reconnect") }}
                </v-btn>
              </v-col>
            </v-row>
          </v-card-actions>
        </v-card>
      </div>
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

const eventHandlers = {
  keydown: null as ((e: KeyboardEvent) => void) | null,
  keyup: null as ((e: KeyboardEvent) => void) | null,
};

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
  if (!displayRef.value) return { width: 800, height: 600 }; // значения по умолчанию
  const rect = document.querySelector(".wrapper").getBoundingClientRect();
  console.log();
  return {
    width: Math.floor(rect.width * 0.9),
    height: Math.floor(rect.height * 0.9),
  };
};

function cleanupConnection() {
  // Удаляем обработчики клавиатуры
  if (displayRef.value) {
    if (eventHandlers.keydown) {
      displayRef.value.removeEventListener("keydown", eventHandlers.keydown);
    }
    if (eventHandlers.keyup) {
      displayRef.value.removeEventListener("keyup", eventHandlers.keyup);
    }
  }

  // Очищаем мышь
  if (guacMouse) {
    guacMouse.onmousedown = guacMouse.onmousemove = guacMouse.onmouseup = null;
    guacMouse = null;
  }

  // Удаляем дисплей
  if (displayRef.value?.firstChild) {
    displayRef.value.removeChild(displayRef.value.firstChild);
  }

  // Отключаем клиент и туннель
  client?.disconnect();
  tunnel?.disconnect();
  client = null;
  tunnel = null;
}

function runConnection() {
  if (!displayRef.value) return;
  cleanupConnection();
  const guacID = router.currentRoute.value.params.id;
  const token = window.localStorage.getItem("guac_token");
  const { width, height } = getDisplayDimensions();
  console.log("Display dimensions:", width, height); // Для отладки
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
        errorText.value = "Соединение закрыто.";
        break;

      case Guacamole.Tunnel.State.UNSTABLE:
        console.error("Туннель: разрыв соединения!");
        errorText.value = "Ошибка сети. Переподключитесь.";
        break;
    }
  };

  eventHandlers.keydown = (e) => handleKeyEvent(e, 1);
  eventHandlers.keyup = (e) => handleKeyEvent(e, 0);
  const options = { passive: false };
  displayRef.value.addEventListener("keydown", eventHandlers.keydown, options);
  displayRef.value.addEventListener("keyup", eventHandlers.keyup, options);

  initMouse();

  setTimeout(() => {
    if (client) {
      client.sendSize(width, height);
    }
  }, 300);
}

const handleResize = () => {
  if (client && displayRef.value) {
    const { width, height } = getDisplayDimensions();
    client.sendSize(width, height);
  }
};

onMounted(() => {
  runConnection();
  displayRef.value.focus();
  window.addEventListener("resize", handleResize);
});

onUnmounted(() => {
  if (client) client.disconnect();
  if (tunnel) tunnel.disconnect();
  window.removeEventListener("resize", handleResize);
});
</script>

<style>
.wrapper {
  position: relative;
  height: calc(100vh - 32px); /* учитываем margin */
}

.guacamole-display {
  width: 100%;
  height: 100%;
  position: absolute;
  top: 0;
  left: 0;
  background-color: #000;
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
