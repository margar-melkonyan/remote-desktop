import {
  chooseSymbolDialog, currentPlayer,
  isPrivate, mySymbol, versusFetchIntervalId, waitSymbolChoosing,
  wssIsSuccess, rowsAndColumns,
  fetchRoom, playerStep, resetGame, resizeCountingArrays,
} from "@/plugins/services/utils";
import { toast } from "vue3-toastify";

/**
 * Подключение к WebSocket серверу для взаимодействия в комнате.
 *
 * Эта функция создает новое соединение WebSocket с сервером, подключается к указанной комнате,
 * а также обрабатывает различные события WebSocket (открытие соединения, ошибки, получение сообщений и закрытие).
 * После успешного подключения она отправляет команду на сервер, чтобы инициировать соединение с указанным паролем.
 *
 * @param {string} roomId - Идентификатор комнаты, к которой нужно подключиться.
 * @param {string} password - Пароль для входа в комнату. Может быть пустым, если комната открытая.
 * @param {object} apiRooms - Объект API для работы с комнатами. Должен содержать URL для подключения.
 * @param {object} authStore - Объект состояния аутентификации (например, Vuex store или Pinia store), который хранит информацию о пользователе.
 * @param {object} router - Объект роутера, используемый для навигации в приложении.
 *
 * @returns {WebSocket} Возвращает объект WebSocket, который используется для связи с сервером.
 *
 * @throws {Error} В случае ошибок при подключении или обмене данными с сервером.
 *
 * Основные события WebSocket:
 * - **onopen**: При успешном открытии соединения отправляется запрос на подключение.
 * - **onerror**: Обрабатываются ошибки при установке соединения.
 * - **onclose**: При закрытии соединения выполняются действия в зависимости от причины (например, перенаправление пользователя).
 * - **onmessage**: Обрабатываются сообщения от сервера, включая различные команды (например, выбор символа, обновление состояния игры).
 */
export function connectToRoom(
  roomId: string, password: string, apiRooms, authStore, router): WebSocket {
  const ws = new WebSocket(`${apiRooms.urls.room(roomId).replace("http", "ws")}?token=${localStorage.getItem("token")}`);
  ws.onopen = () => {
    ws.send(JSON.stringify({
      action: "new connection to room",
      password: password,
    }))
    wssIsSuccess.value = true;
    isPrivate.value = false;
  };
  ws.onerror = (event) => {
    console.error("WebSocket error observed:", event);
    router.push({ name: "index" })
  };
  ws.onclose = (event) => {
    toast.error(event.reason);
    if( event.code === 1013 ||
      event.reason === 'connection is close' ||
      event.reason === 'cannot find room'
    ) {
      clearInterval(versusFetchIntervalId.value)
      router.push({ name: "index" })
    }
    if (event.code === 1008) {
      wssIsSuccess.value = false;
      isPrivate.value = true;
    }
  };
  ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    switch(data.action) {
      case "reset game":
        resetGame();
        break;
      case "new connection to room":
        fetchRoom(roomId, apiRooms, authStore);
        break;
      case "get positions":
        const positions = data.data.positions;
        positions.forEach((position) => {
          const pos = position.id.split("-");
          const i = Number(pos[0]);
          const j = Number(pos[1]);
          playerStep(i, j, position.symbol);
        });
        currentPlayer.value = data.symbol
        if (positions.length === 0) {
          resetGame();
          fetchRoom(roomId, apiRooms, authStore);
          currentPlayer.value = "X"
        }
        break;
      case "resize":
        rowsAndColumns.value = data.size;
        resizeCountingArrays();
        break;
      case "choose symbol":
        mySymbol.value = ""
        if(authStore?.user?.id === data.user_id) {
          chooseSymbolDialog.value = true
        } else {
          waitSymbolChoosing.value = true
        }
        break
      case "selected symbol":
        mySymbol.value = data.symbol
        currentPlayer.value = 'X'
        waitSymbolChoosing.value = false
        break
      case "restart game":
        currentPlayer.value = 'X'
        break
      case "sync symbol":
        mySymbol.value = data.symbol
        waitSymbolChoosing.value = false
        break
    }
  };
  return ws;
}
