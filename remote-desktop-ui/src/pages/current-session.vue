<template>
  <div class="d-flex justify-center wrapper ma-4 rounded-4">
    <div
      ref="displayRef"
      class="guacamole-display"
      tabindex="0"
      @contextmenu.prevent
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import Guacamole from '@cyolosecurity/guacamole-common-js'
import { useRouter } from "vue-router";

const router = useRouter();

const displayRef = ref(null)
let client = null
let tunnel = null

const keyMap = {
  'Enter': 0xFF0D,
  'Backspace': 0xFF08,
  'Tab': 0xFF09,
  'Escape': 0xFF1B,
  'ArrowUp': 0xFF52,
  'ArrowDown': 0xFF54,
  'ArrowLeft': 0xFF51,
  'ArrowRight': 0xFF53,
  'Shift': 0xFFE1,
  'ShiftRight': 0xFFE2,
  'Control': 0xFFE3,
  'Alt': 0xFFE9,
  'Meta': 0xFFE7,
  ' ': 0x0020 // Пробел
}

// Блокируемые браузерные комбинации
const blockedCombinations = [
  'Control+t', 'Control+n', 'Control+w',
  'Control+Tab', 'Alt+ArrowLeft', 'Alt+ArrowRight',
  'F1', 'F5', 'Control+r', 'Meta+Space'
]

const handleKeyEvent = (e, pressed) => {
  if (!client) return

  if (e.key === ' ' && (e.metaKey || e.getModifierState('Meta'))) {
    e.preventDefault();
    e.stopPropagation();
    return;
  }

  // Проверка на заблокированные комбинации
  const modifiers = [
    e.ctrlKey ? 'Control' : null,
    e.altKey ? 'Alt' : null,
    e.shiftKey ? 'Shift' : null
  ].filter(Boolean)

  const currentCombo = [...modifiers, e.key].join('+')

  if (blockedCombinations.includes(currentCombo)) {
    e.preventDefault()
    e.stopPropagation()
    return
  }

  if (e.key in keyMap) {
    client.sendKeyEvent(pressed, keyMap[e.key])
    e.preventDefault()
    return
  }

  if (e.key.length === 1) {
    const keysym = e.key.charCodeAt(0)
    client.sendKeyEvent(pressed, keysym)
    e.preventDefault()
  }
}

let guacMouse = null;

const initMouse = () => {
  if (!client || !displayRef.value) return;

  const displayElement = client.getDisplay().getElement();
  guacMouse = new Guacamole.Mouse(displayElement);

  guacMouse.onmousedown =
  guacMouse.onmousemove =
  guacMouse.onmouseup = (state) => {
    client.sendMouseState(state, true)
  };
};

const getDisplayDimensions = () => {
  const rect = document.querySelector('.wrapper').getBoundingClientRect();
  return {
    width: Math.floor(rect.width * 0.9),
    height: Math.floor(rect.height * 0.9)
  };
};

onMounted(() => {
  if (!displayRef.value) return
  const guacID = router.currentRoute.value.params.id;
  const token = 'C60A97B43D92B7AE1AB0E36F9C2D4F395EF71C413062360498E55F34AF714933';
  const { width, height } = getDisplayDimensions();
  tunnel = new Guacamole.WebSocketTunnel(
    `ws://192.168.1.4:8080/guacamole/websocket-tunnel?token=${token}&GUAC_DATA_SOURCE=postgresql&GUAC_ID=${guacID}&GUAC_TYPE=c&GUAC_WIDTH=${width}&GUAC_HEIGHT=${height}&GUAC_DPI=96&GUAC_TIMEZONE=Europe/Moscow&GUAC_AUDIO=audio/L8&GUAC_AUDIO=audio/L16&GUAC_IMAGE=image/jpeg&GUAC_IMAGE=image/png&GUAC_IMAGE=image/webp`
  )

  client = new Guacamole.Client(tunnel)

  client.onerror = (error) => {
    if (error.code === 519) {
      alert(`Не удалось подключить проверьте указанные данные в подключении или удостоверьтесь что устройство к которому подключаетесь работает ${error.message}`)
    }
  }
  tunnel.onerror = error => console.error('Tunnel error:', error)
  client.onmessage = (data) => console.log(data)
  tunnel.onmessage = (data) => console.log(data)

  client.connect()
  displayRef.value.appendChild(client.getDisplay().getElement())

  const options = { passive: false }
  displayRef.value.addEventListener('keydown', e => handleKeyEvent(e, 1), options)
  displayRef.value.addEventListener('keyup', e => handleKeyEvent(e, 0), options)

  initMouse();
  displayRef.value.focus()
})

onUnmounted(() => {
  if (client) client.disconnect()
  if (tunnel) tunnel.disconnect()
})
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
