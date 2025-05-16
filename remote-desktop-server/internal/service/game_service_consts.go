// Package service реализует бизнес-логику приложения.
package service

// Параметры по умолчанию для игровой доски.
const (
	// DEFAULT_BORDER_SIZE определяет размер доски по умолчанию.
	DEFAULT_BORDER_SIZE = 3
	// DEFAULT_PLAYER задаёт символ игрока по умолчанию.
	DEFAULT_PLAYER = "X"
)

// Названия действий для взаимодействия по WebSocket.
const (
	stepAction                = "step"
	syncSymbolAction          = "sync symbol"
	chooseSymbolAction        = "choose symbol"
	getPositionsAction        = "get positions"
	selectSymbolAction        = "select symbol"
	selectedSymbolAction      = "selected symbol"
	resizeAction              = "resize"
	resetGameAction           = "reset game"
	gameEndAction             = "game end"
	restartGameAction         = "restart game"
	closeRoomAction           = "close room"
	exitRoomAction            = "exit room"
	newConnectionToRoomAction = "new connection to room"
)

// game statuses
const (
	chooseSymbolStatus = "choose symbol"
	inProcessStatus    = "in process"
	gameEndStatus      = "game end"
)
