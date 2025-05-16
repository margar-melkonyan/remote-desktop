import { ref } from "vue";
import axios from "axios";

/**
 * Содержит символ игрока, который участвует в текущей игре.
 * Значение может быть либо пустой строкой, либо символом ('X' или 'O').
 *
 * @type {Ref<string>}
 * @default ""
 */
export const mySymbol = ref<string>("");

/**
 * Содержит символ текущего игрока ('X' или 'O'), чей ход сейчас.
 *
 * @type {Ref<string>}
 * @default ""
 */
export const currentPlayer = ref<string>("");

/**
 * Содержит информацию о текущей комнате.
 * Используется для хранения данных о комнате, таких как пользователи и статус.
 *
 * @type {Ref<any>}
 * @default null
 */
export const roomInfo = ref(null);

/**
 * Содержит количество строк и столбцов на игровом поле.
 * Используется для адаптивного изменения размеров поля.
 *
 * @type {Ref<number>}
 * @default 3
 */
export const rowsAndColumns = ref<number>(3);

/**
 * Флаг, указывающий на победителя игры.
 * Значение:
 * - 0: нет победителя.
 * - 1: выиграл игрок с символом 'X'.
 * - -1: выиграл игрок с символом 'O'.
 * - -2: ничья.
 *
 * @type {Ref<number>}
 * @default 0
 */
export const wonFlag = ref<number>(0);

/**
 * Флаг, указывающий на начало игры.
 * Значение:
 * - 0: игра не начата.
 * - 1: игра начата.
 *
 * @type {Ref<number>}
 * @default 0
 */
export const gameStarted = ref<number>(0);

/**
 * Флаг, указывающий, закончена ли игра.
 * Значение:
 * - 0: игра не закончена.
 * - 1: игра завершена (победитель найден или ничья).
 *
 * @type {Ref<number>}
 * @default 0
 */
export const gameEnd = ref<number>(0);

/**
 * Интервал, использующийся для периодического обновления состояния игры.
 * Значение — идентификатор интервала для дальнейшего его очистки.
 *
 * @type {Ref<number>}
 * @default 0
 */
export const versusFetchIntervalId = ref<number>(0);

/**
 * Флаг, указывающий, является ли комната частной.
 *
 * @type {Ref<boolean>}
 * @default false
 */
export const isPrivate = ref<boolean>(false);

/**
 * Флаг, указывающий, был ли успешен WebSocket-соединение.
 *
 * @type {Ref<boolean>}
 * @default false
 */
export const wssIsSuccess = ref<boolean>(false);

/**
 * Флаг, показывающий, отображается ли диалог выбора символа.
 * Если значение true, то отображается диалог для выбора символа игроком.
 *
 * @type {Ref<boolean>}
 * @default false
 */
export const chooseSymbolDialog = ref<boolean>(false);

/**
 * Флаг, показывающий, ожидает ли игрок выбора символа.
 * Если значение true, то игрок ожидает, чтобы выбрать символ.
 *
 * @type {Ref<boolean>}
 * @default false
 */
export const waitSymbolChoosing = ref<boolean>(false);

/**
 * Массив для подсчета количества символов 'X' в каждой строке/столбце (или диагонали).
 * Динамически изменяется в зависимости от размера поля.
 *
 * @type {number[]}
 * @default [0, 0, 0] // Для поля 3x3
 */
export const xCount = new Array(rowsAndColumns.value).fill(0);

/**
 * Массив для подсчета количества символов 'O' в каждой строке/столбце (или диагонали).
 * Динамически изменяется в зависимости от размера поля.
 *
 * @type {number[]}
 * @default [0, 0, 0] // Для поля 3x3
 */
export const oCount = new Array(rowsAndColumns.value).fill(0);

/**
 * Объект AbortController для управления отменой HTTP-запросов (например, для отмены при выходе из комнаты или обновлении состояния).
 *
 * @type {AbortController}
 * @default undefined
 */
export let controller: AbortController;



/**
 * Генерирует объект стилей для ячейки игрового поля.
 *
 * Эта функция возвращает объект с CSS стилями для ячейки, который применяется для настройки внешнего вида ячейки игрового поля.
 * Стили включают в себя настройки для размеров, границ, фона, выравнивания содержимого и курсора.
 *
 * @returns {object} Объект стилей CSS, который можно применить к элементу.
 *
 */
export function getCellStyle() {
  return {
    flex: '1',
    aspectRatio: '1',
    border: '0.25rem solid #ff7fea',
    'border-radius': '0.25rem',
    margin: '0.25rem',
    backgroundColor: 'white',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    cursor: 'pointer',
  };
}

/**
 * Генерирует объект стилей для шрифта в зависимости от размеров игрового поля.
 *
 * Эта функция вычисляет размер шрифта для отображения текста в ячейках игрового поля на основе количества строк и колонок.
 * Чем больше поле, тем меньше размер шрифта.
 *
 * @param {number} rowsAndColumns - Количество строк и колонок на игровом поле. Используется для вычисления размера шрифта.
 *
 * @returns {object} Объект с CSS стилями для шрифта, включая цвет и размер шрифта.
 *
 * Пример стилей, которые возвращаются функцией:
 * - `color: 'black'`: Цвет текста — черный.
 * - `fontSize: `${size}px``: Размер шрифта, вычисляемый на основе количества строк и колонок. Чем больше поле, тем меньше шрифт.
 *
 * Пример вычисления размера:
 * Если передано значение `rowsAndColumns = 10`, то размер шрифта будет рассчитан как `300 / 10 = 30px`.
 */
export function getFontStyle(rowsAndColumns: number):object {
  const base = 300;
  const size = Math.floor(base / rowsAndColumns);
  return {
    color: 'black',
    fontSize: `${size}px`,
  };
}

/**
 * Определяет состояние игры на основе символа игрока и флага победы.
 *
 * Эта функция проверяет, выиграл ли игрок в игре, а также учитывает флаг победы, который может указать на проигрыш
 * или ничью. В зависимости от этих параметров функция возвращает:
 * - `1`, если игрок выиграл,
 * - `-1`, если игрок проиграл,
 * - `0`, если игра не завершена или ничья.
 *
 * @param {string} mySymbol - Символ игрока (например, `'X'` или `'O'`).
 * @param {number} wonFlag - Флаг победы:
 *   - `1`: Игрок с символом `'X'` победил.
 *   - `-1`: Игрок с символом `'O'` победил.
 *   - `-2`: Игра завершена ничьей.
 *   - `0`: Игра продолжается.
 *
 * @returns {number} - Состояние игры:
 *   - `1`: Если игрок победил.
 *   - `-1`: Если игрок проиграл.
 *   - `0`: Если игра продолжается или ничья.
 */
export const wonState = (mySymbol: string, wonFlag: number): number => {
  if (mySymbol === 'X' && wonFlag === 1) {
    return 1
  }
  if (mySymbol === 'O' && wonFlag === -1) {
    return 1
  }
  if (wonFlag === -2) {
    return -1
  }
  return 0
}

/**
 * Загружает информацию о комнате и обновляет состояние игры.
 *
 * Эта функция выполняет запрос к API для получения данных о комнате по её идентификатору. После получения данных,
 * она обновляет состояние, связанное с пользователями, символами и конфиденциальностью комнаты. Также она проверяет,
 * выбран ли символ текущим пользователем, и изменяет соответствующие флаги и значения.
 *
 * @param {string} roomId - Идентификатор комнаты, для которой нужно получить информацию.
 * @param {object} apiRooms - Объект API для работы с комнатами.
 * @param {object} authStore - Хранилище с данными об аутентификации пользователя (содержит информацию о текущем пользователе).
 *
 * @returns {Promise<void>} - Возвращает промис, который выполняется после получения данных.
 */
export async function fetchRoom(roomId, apiRooms, authStore:any): Promise<void> {
  if (controller) controller.abort();
  controller = new AbortController();
  try {
    const { data } = await axios.get(apiRooms.urls.roomInfo(roomId));
    roomInfo.value = data.data;
    const currentUser = roomInfo?.value.users.filter((user) => user.id === authStore?.user?.id)[0]
    if (currentUser?.symbol !== undefined && currentUser?.symbol !== "") {
      mySymbol.value = currentUser?.symbol
      chooseSymbolDialog.value = false
      waitSymbolChoosing.value = false
    }
    if (!wssIsSuccess.value) {
      isPrivate.value = roomInfo?.value.is_private
    }
  } catch (err) {
    console.error("Failed to fetch room info:", err);
  }
}

/**
 * Очищает клетки игрового поля.
 *
 * Эта функция находит все элементы `span`, которые являются дочерними элементами элементов с классами, начинающимися на `grid-index-`,
 * и очищает их текстовое содержимое. Это полезно для сброса состояния игры, например, при перезапуске игры или очистке игрового поля.
 *
 * @returns {void} - Функция не возвращает значений.
 */
export function resetGameBoardCells() {
  const gridSpans = document.querySelectorAll('[class^="grid-index-"] > span');
  gridSpans.forEach(span => {
    span.textContent = '';
  });
}

/**
 * Осуществляет шаг игрока на игровом поле.
 *
 * Эта функция выполняет ход игрока, обновляя состояние клетки на игровом поле и проверяя, не завершена ли игра.
 * После выполнения хода функция выполняет проверки на выигрыш, ничью и обновляет состояние игры.
 *
 * @param {number} i - Индекс строки игрового поля, на которой был сделан ход (0-based).
 * @param {number} j - Индекс столбца игрового поля, на котором был сделан ход (0-based).
 * @param {string} symbol - Символ игрока, который выполняет ход ('X' или 'O').
 *
 * @returns {void} - Функция не возвращает значений, но изменяет состояние игрового поля.
 */
export function playerStep(i: number, j: number, symbol: string) {
  gameStarted.value = 1;
  const cell = document.querySelector(`.grid-index-${i}-${j}>span`);
  if (wonFlag.value !== 0 || !cell) {
    return;
  }
  if (cell.textContent === 'O' || cell.textContent === 'X') {
    return;
  }
  cell.textContent = symbol;
  resetCounting();
  verticalCheck();
  if (wonFlag.value === 0) {
    resetCounting();
    horizontalCheck();
  }
  if (wonFlag.value === 0) {
    resetCounting();
    mainDiagonalCheck();
  }
  if (wonFlag.value === 0) {
    resetCounting();
    sideDiagonalCheck();
  }
  checkDraw();
}

/**
 * Сбрасывает счётчики для символов 'X' и 'O'.
 *
 * Эта функция обнуляет значения в массивах счётчиков `xCount` и `oCount`, которые отслеживают количество
 * символов на игровом поле для каждого игрока. Это необходимо делать после каждого хода, чтобы корректно
 * отслеживать состояние игры.
 *
 * @returns {void} - Функция не возвращает значений, но изменяет состояние счётчиков.
 */
function resetCounting() {
  xCount.fill(0);
  oCount.fill(0);
}

/**
 * Изменяет длину массивов счётчиков для символов 'X' и 'O' в зависимости от размера игрового поля.
 *
 * Эта функция изменяет длину массивов `xCount` и `oCount` на основе значения `rowsAndColumns`,
 * которое определяет размер игрового поля. После изменения длины, массивы заполняются нулями,
 * чтобы подготовить их для отслеживания состояния игры на поле с новым размером.
 *
 * @returns {void} - Функция не возвращает значений, но изменяет длину массивов и их содержимое.
 */
export  function resizeCountingArrays() {
  xCount.length = rowsAndColumns.value;
  oCount.length = rowsAndColumns.value;
  xCount.fill(0);
  oCount.fill(0);
}

/**
 * Проверяет главную диагональ игрового поля и обновляет соответствующие счётчики для символов 'X' и 'O'.
 *
 * Эта функция итерирует по клеткам главной диагонали (клетки, где индексы строки и столбца равны),
 * проверяя их содержимое. В зависимости от символа в клетке ('X' или 'O'), она увеличивает соответствующие
 * счётчики (`xCount` и `oCount`) для каждого игрока. После проверки главной диагонали вызывается функция
 * `diagonalChecker`, которая, вероятно, выполняет дополнительную логику для проверки диагоналей.
 *
 * @returns {void} - Функция не возвращает значений, но изменяет состояние счётчиков.
 */
function mainDiagonalCheck() {
  for (let i = 1; i <= rowsAndColumns.value; i++) {
    const cell = document.querySelector(`.grid-index-${i}-${i}>span`);
    if (!cell) continue;
    if (cell.textContent === 'O') oCount[i - 1] += 1;
    if (cell.textContent === 'X') xCount[i - 1] += 1;
  }
  diagonalChecker();
}

/**
 * Проверяет побочную диагональ игрового поля и обновляет соответствующие счётчики для символов 'X' и 'O'.
 *
 * Эта функция итерирует по клеткам побочной диагонали (клетки, где сумма индексов строки и столбца равна
 * размеру поля + 1). Для каждой клетки она проверяет её содержимое и в зависимости от символа ('X' или 'O')
 * увеличивает соответствующий счётчик (`xCount` или `oCount`). После проверки побочной диагонали вызывается
 * функция `diagonalChecker`, которая выполняет дополнительную логику для проверки диагоналей.
 *
 * @returns {void} - Функция не возвращает значений, но изменяет состояние счётчиков.
 */
function sideDiagonalCheck() {
  for (let i = 0; i < rowsAndColumns.value; i++) {
    const cell = document.querySelector(`.grid-index-${i + 1}-${rowsAndColumns.value - i}>span`);
    if (!cell) continue;
    if (cell.textContent === 'O') oCount[i] += 1;
    if (cell.textContent === 'X') xCount[i] += 1;
  }
  diagonalChecker();
}

/**
 * Проверяет победителя на диагоналях, используя счётчики для символов 'X' и 'O'.
 *
 * Эта функция проверяет, если сумма значений в массивах счётчиков `xCount` или `oCount`
 * равна размеру игрового поля (то есть если одна из диагоналей полностью занята одним символом).
 * Если сумма для символа 'X' равна размеру поля, это означает победу для 'X', и значение `wonFlag`
 * устанавливается в 1. Если сумма для символа 'O' равна размеру поля, это означает победу для 'O',
 * и значение `wonFlag` устанавливается в -1.
 *
 * @returns {void} - Функция не возвращает значений, но изменяет состояние флага победителя `wonFlag`.
 */
function diagonalChecker() {
  if (xCount.reduce((a, b) => a + b, 0) === rowsAndColumns.value) {
    wonFlag.value = 1;
  }
  if (oCount.reduce((a, b) => a + b, 0) === rowsAndColumns.value) {
    wonFlag.value = -1;
  }
}

/**
 * Проверяет горизонтальные линии на игровом поле и обновляет счётчики для символов 'X' и 'O'.
 *
 * Эта функция проходит по всем клеткам на поле, проверяя каждую строку (горизонтальную линию).
 * Для каждой клетки она обновляет соответствующий счётчик в массивах `xCount` и `oCount` в зависимости от того,
 * какой символ находится в клетке ('X' или 'O'). После завершения проверки она вызывает функцию `lineChecker`,
 * чтобы проверить, не заполнил ли какой-либо игрок горизонталь.
 *
 * @returns {void} - Функция не возвращает значений, но изменяет состояние счётчиков для символов 'X' и 'O'.
 */
function horizontalCheck() {
  for (let i = 1; i <= rowsAndColumns.value; i++) {
    for (let j = 1; j <= rowsAndColumns.value; ++j) {
      const cell = document.querySelector(`.grid-index-${i}-${j}>span`);
      if (!cell) continue;
      if (cell.textContent === 'O') oCount[j - 1] += 1;
      if (cell.textContent === 'X') xCount[j - 1] += 1;
    }
  }
  lineChecker();
}

/**
 * Проверяет вертикальные линии на игровом поле и обновляет счётчики для символов 'X' и 'O'.
 *
 * Эта функция проходит по всем клеткам на поле, проверяя каждую колонку (вертикальную линию).
 * Для каждой клетки она обновляет соответствующий счётчик в массивах `xCount` и `oCount` в зависимости от того,
 * какой символ находится в клетке ('X' или 'O'). После завершения проверки она вызывает функцию `lineChecker`,
 * чтобы проверить, не заполнил ли какой-либо игрок вертикаль.
 *
 * @returns {void} - Функция не возвращает значений, но изменяет состояние счётчиков для символов 'X' и 'O'.
 */
function verticalCheck() {
  for (let i = 1; i <= rowsAndColumns.value; i++) {
    for (let j = 1; j <= rowsAndColumns.value; ++j) {
      const cell = document.querySelector(`.grid-index-${i}-${j}>span`);
      if (!cell) continue;
      if (cell.textContent === 'O') oCount[i - 1] += 1;
      if (cell.textContent === 'X') xCount[i - 1] += 1;
    }
  }
  lineChecker();
}

/**
 * Проверяет, есть ли победная линия на игровом поле для символов 'X' или 'O'.
 *
 * Функция анализирует массивы счётчиков `xCount` и `oCount`, чтобы определить,
 * заполнил ли кто-либо всю линию (вертикаль, горизонталь или диагональ) символом 'X' или 'O'.
 * Если такой линии нет, функция ничего не делает. Если одна из сторон победила (например, игрок 'X' заполнил все клетки в линии),
 * она устанавливает флаг победы `wonFlag` в значение 1 для игрока 'X' или -1 для игрока 'O'.
 *
 * @returns {void} - Функция не возвращает значений, но изменяет флаг победы `wonFlag`.
 */
function lineChecker() {
  if (Math.max(...xCount) === rowsAndColumns.value) {
    wonFlag.value = 1;
  }
  if (Math.max(...oCount) === rowsAndColumns.value) {
    wonFlag.value = -1;
  }
}
/**
 * Проверяет, закончилась ли игра ничьей.
 *
 * Функция анализирует все клетки игрового поля. Если все клетки заняты (не пустые),
 * и никто не победил (флаг `wonFlag.value` равен 0), функция устанавливает флаг ничьей
 * в переменную `wonFlag.value` со значением -2.
 *
 * @returns {void} - Функция не возвращает значений, но изменяет флаг `wonFlag` на -2, если игра закончилась ничьей.
 */
function checkDraw() {
  gameEnd.value = 0;
  for (let i = 1; i <= rowsAndColumns.value; i++) {
    for (let j = 1; j <= rowsAndColumns.value; ++j) {
      const cell = document.querySelector(`.grid-index-${i}-${j}>span`);
      if (cell && cell.textContent !== '') {
        gameEnd.value += 1;
      }
    }
  }
  if (gameEnd.value === rowsAndColumns.value * rowsAndColumns.value && wonFlag.value === 0) {
    wonFlag.value = -2;
  }
}

/**
 * Сбрасывает состояние игры.
 *
 * Функция выполняет сброс флагов и счетчиков, а также очищает игровое поле.
 * Это полезно для начала новой игры после завершения текущей или для сброса игры по каким-либо причинам.
 *
 * Описание действий:
 * 1. Сбрасывает флаг начала игры в `gameStarted.value` на 0.
 * 2. Сбрасывает флаг победы в `wonFlag.value` на 0.
 * 3. Сбрасывает счетчики для подсчета победных линий с помощью функции `resetCounting()`.
 * 4. Очищает клетки игрового поля с помощью функции `resetGameBoardCells()`.
 *
 * @returns {void} - Функция не возвращает значения.
 */
export function resetGame() {
  gameStarted.value = 0;
  wonFlag.value = 0;
  resetCounting();
  resetGameBoardCells();
}
