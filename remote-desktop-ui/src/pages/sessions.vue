<template>
  <div>
    <div ref="displayRef" class="guacamole-display" tabindex="0"></div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref } from 'vue'
import Guacamole from '@cyolosecurity/guacamole-common-js'

const displayRef = ref(null)
let client = null
let tunnel = null

// Обработчики событий
const handleKeyEvent = (e, pressed) => {
  if (!client) return

  // Специальная обработка некоторых клавиш
  const keyMap = {
    'Enter': 0xFF0D,
    'Backspace': 0xFF08,
    'Tab': 0xFF09,
    'Escape': 0xFF1B,
    'ArrowUp': 0xFF52,
    'ArrowDown': 0xFF54,
    'ArrowLeft': 0xFF51,
    'ArrowRight': 0xFF53,
    'Shift': 0xFFE1,       // Левый Shift
    'ShiftRight': 0xFFE2,  // Правый Shift
    'Control': 0xFFE3,
    'Alt': 0xFFE9,
    'Meta': 0xFFE7         // Клавиша Win/Command
  }

  const keysym = keyMap[e.key] || e.key.charCodeAt(0)
  client.sendKeyEvent(pressed, keysym)
}

const handleMouse = (e) => {
  if (!client || !displayRef.value) return

  const rect = displayRef.value.getBoundingClientRect()
  const x = e.clientX - rect.left
  const y = e.clientY - rect.top

  // Определяем нажатые кнопки мыши
  const buttonState =
    (e.buttons & 1 ? 1 : 0) |  // Левая кнопка
    (e.buttons & 2 ? 2 : 0) |   // Правая кнопка
    (e.buttons & 4 ? 4 : 0)     // Средняя кнопка

  client.sendMouseState(x, y, buttonState)
}

onMounted(() => {
  if (!displayRef.value) return

  // guac_id - идентификатор клиента
  // token - user token
  // GUAC_TYPE - rdp, s
  // я каждого типа требуются разные параметры:
  // RDP: GUAC_AUDIO, GUAC_TIMEZONE
  // VNC: GUAC_PASSWORD
  // SSH: GUAC_USERNAME, GUAC_PRIVATE_KEY

  // Инициализация подключения
  tunnel = new Guacamole.WebSocketTunnel(
    'ws://192.168.1.4:8080/guacamole/websocket-tunnel?token=268B056B09575661322ECA34B83F4B112CACA3DD3DEC6CF9160E4AF0F298E1C5&GUAC_DATA_SOURCE=postgresql&GUAC_ID=1&GUAC_TYPE=c&GUAC_WIDTH=2560&GUAC_HEIGHT=708&GUAC_DPI=96&GUAC_TIMEZONE=Europe/Moscow&GUAC_AUDIO=audio/L8&GUAC_AUDIO=audio/L16&GUAC_IMAGE=image/jpeg&GUAC_IMAGE=image/png&GUAC_IMAGE=image/webp'
  )

  client = new Guacamole.Client(tunnel)

  // Обработчики ошибок
  client.onerror = error => console.error('Client error:', error)
  tunnel.onerror = error => console.error('Tunnel error:', error)

  // Подключение
  client.connect()

  // Добавляем экран в DOM
  displayRef.value.appendChild(client.getDisplay().getElement())

  // Фокус на элемент для ввода
  displayRef.value.focus()

  // Обработчики событий
  displayRef.value.addEventListener('keydown', e => handleKeyEvent(e, 1))
  displayRef.value.addEventListener('keyup', e => handleKeyEvent(e, 0))
  displayRef.value.addEventListener('mousedown', handleMouse)
  displayRef.value.addEventListener('mousemove', handleMouse)
  displayRef.value.addEventListener('mouseup', handleMouse)

  // Обработка ввода текста (для IME и сложных символов)
  displayRef.value.addEventListener('compositionend', (e) => {
    if (e.data) {
      Array.from(e.data).forEach(char => {
        const keysym = char.charCodeAt(0)
        client.sendKeyEvent(1, keysym)
        client.sendKeyEvent(0, keysym)
      })
    }
  })
})

onUnmounted(() => {
  // Очистка
  if (client) client.disconnect()
  if (tunnel) tunnel.disconnect()

  // Удаление обработчиков
  if (displayRef.value) {
    displayRef.value.removeEventListener('keydown', handleKeyEvent)
    displayRef.value.removeEventListener('keyup', handleKeyEvent)
    displayRef.value.removeEventListener('mousedown', handleMouse)
    displayRef.value.removeEventListener('mousemove', handleMouse)
    displayRef.value.removeEventListener('mouseup', handleMouse)
  }
})
</script>

<style>
.guacamole-display {
  width: 100%;
  height: 80vh;
  border: 1px solid #ccc 5px;
  outline: none; /* Убираем стандартный фокус */
  cursor: text; /* Курсор как в терминале */
  background-color: #000; /* Черный фон как в консоли */
  color: #fff; /* Белый текст */
  font-family: monospace; /* Моноширинный шрифт */
  overflow: hidden;
}
</style>
