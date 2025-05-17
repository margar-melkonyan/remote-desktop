import { onMounted, onUnmounted, ref } from 'vue'
import Guacamole from '@cyolosecurity/guacamole-common-js'

export function useGuacamole() {
  const displayRef = ref(null)
  let client = null
  let tunnel = null

  // Функция для обработки изменения размеров
  const handleResize = () => {
    if (!client || !displayRef.value) return

    const display = client.getDisplay()
    const scale = Math.min(
      displayRef.value.clientWidth / display.getWidth(),
      displayRef.value.clientHeight / display.getHeight()
    )

    display.scale(scale)
  }

  // Функция инициализации Guacamole
  const initGuacamole = () => {
    if (!displayRef.value) return

    // 1. Создание туннеля
    tunnel = new Guacamole.WebSocketTunnel(
      'ws://192.168.1.4:8080/guacamole/websocket-tunnel?token=FEB0D6502DDA76719DEB54B552CF071D12AC7E52DF543C3DC41F0CD7531F092E&GUAC_DATA_SOURCE=postgresql&GUAC_ID=1&GUAC_TYPE=c&GUAC_WIDTH=2560&GUAC_HEIGHT=708&GUAC_DPI=96&GUAC_TIMEZONE=Europe/Moscow&GUAC_AUDIO=audio/L8&GUAC_AUDIO=audio/L16&GUAC_IMAGE=image/jpeg&GUAC_IMAGE=image/png&GUAC_IMAGE=image/webp'
    )

    // 2. Инициализация клиента
    client = new Guacamole.Client(tunnel)

    // 3. Обработчики событий
    client.onerror = (error) => {
      console.error('Ошибка клиента:', error)
    }

    tunnel.onerror = (error) => {
      console.error('Ошибка туннеля:', error)
    }

    // 4. Подключение
    client.connect()

    // 5. Встраивание экрана
    displayRef.value.innerHTML = '' // Очищаем контейнер
    displayRef.value.appendChild(client.getDisplay().getElement())

    // 6. Первоначальное масштабирование
    handleResize()

    // 7. Добавляем обработчик ресайза
    window.addEventListener('resize', handleResize)
  }

  onMounted(() => {
    initGuacamole()
  })

  onUnmounted(() => {
    // Очистка при размонтировании компонента
    if (client) {
      client.disconnect()
      client = null
    }

    if (tunnel) {
      tunnel.disconnect()
      tunnel = null
    }

    window.removeEventListener('resize', handleResize)
  })

  return {
    displayRef
  }
}
