// Package service реализует бизнес-логику приложения.
package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/common"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/config"
)

// Протоколы по названия которых будут возвращаться подключения
const (
	all = "all"
	ssh = "ssh"
	rdp = "rdp"
)

const (
	indexURL       = "session/data/postgresql/connectionGroups/ROOT/tree"
	connectionsURL = "session/data/postgresql/connections"
)

// SessionService предоставляет сервис для работы c подключениями к удаленным серверам.
type SessionService struct {
	client http.Client
}

// NewSessionService создает новый экземпляр SessionService.
//
// Параметры:
//   - repo: репозиторий для работы с подключениями
//
// Возвращает:
//   - *SessionService: указатель на созданный сервис
func NewSessionService() *SessionService {
	return &SessionService{
		client: http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (service *SessionService) fetchConnections(guacToken string) ([]*common.GuacamoleRDConnectionResponse, error) {
	var response struct {
		ChildConnections []*common.GuacamoleRDConnectionResponse `json:"childConnections"`
	}

	if err := service.makeGuacamoleRequest(
		http.MethodGet,
		indexURL,
		guacToken,
		nil,
		&response,
	); err != nil {
		return nil, fmt.Errorf("failed to fetch connections: %w", err)
	}

	return response.ChildConnections, nil
}

func (service *SessionService) GetSession(protocol string, guacToken string) ([]*common.GuacamoleRDConnectionResponse, error) {
	connections, err := service.fetchConnections(guacToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if protocol == all {
		return connections, nil
	}

	result := make([]*common.GuacamoleRDConnectionResponse, 0, len(connections))
	for _, conn := range connections {
		if conn.Protocol == string(protocol) {
			result = append(result, conn)
		}
	}

	return result, nil
}

func (service *SessionService) EditConnection(id string, guacToken string) (*common.GuacamoleConnectionRequest, error) {
	// Получаем параметры соединения
	var params common.Parameters
	if err := service.makeGuacamoleRequest(
		http.MethodGet,
		fmt.Sprintf("%s/%s/parameters", connectionsURL, id),
		guacToken,
		nil,
		&params,
	); err != nil {
		return nil, fmt.Errorf("failed to get connection parameters: %w", err)
	}

	// Получаем основную информацию о соединении
	var connectionInfo common.GuacamoleRDConnectionRequest
	if err := service.makeGuacamoleRequest(
		http.MethodGet,
		fmt.Sprintf("%s/%s", connectionsURL, id),
		guacToken,
		nil,
		&connectionInfo,
	); err != nil {
		return nil, fmt.Errorf("failed to get connection info: %w", err)
	}

	// Объединяем результаты
	connectionInfo.Parameters = params
	return &common.GuacamoleConnectionRequest{
		Id:       connectionInfo.Id,
		Name:     connectionInfo.Name,
		HostName: connectionInfo.Parameters.HostName,
		Username: connectionInfo.Parameters.Username,
		Password: connectionInfo.Parameters.Password,
		Port:     connectionInfo.Port,
		Protocol: connectionInfo.Protocol,
	}, nil
}

func (service *SessionService) CreateConnection(form *common.GuacamoleConnectionRequest, guacToken string) error {
	ignoreCert := "false"
	if form.Protocol == rdp {
		ignoreCert = "true"
	}

	requestBody := common.GuacamoleRDConnectionRequest{
		Name:             form.Name,
		Protocol:         form.Protocol,
		ParentIdentifier: "ROOT",
		Parameters: common.Parameters{
			HostName:   form.HostName,
			Username:   form.Username,
			Password:   form.Password,
			IgnoreCert: ignoreCert,
			Port:       form.Port,
		},
	}

	if err := service.makeGuacamoleRequest(
		http.MethodPost,
		connectionsURL,
		guacToken,
		requestBody,
		nil,
	); err != nil {
		return fmt.Errorf("failed to create connection: %w", err)
	}

	return nil
}

func (service *SessionService) UpdateConnection(
	id string,
	form *common.GuacamoleConnectionRequest,
	guacToken string,
) error {
	ignoreCert := "false"
	if form.Protocol == rdp {
		ignoreCert = "true"
	}

	path := fmt.Sprintf("%s/%s", connectionsURL, id)

	if err := service.makeGuacamoleRequest(
		http.MethodPut,
		path,
		guacToken,
		common.GuacamoleRDConnectionRequest{
			Id:               id,
			Name:             form.Name,
			Protocol:         form.Protocol,
			ParentIdentifier: "ROOT",
			Parameters: common.Parameters{
				HostName:   form.HostName,
				Username:   form.Username,
				Password:   form.Password,
				IgnoreCert: ignoreCert,
				Port:       form.Port,
			},
		},
		nil,
	); err != nil {
		return fmt.Errorf("failed to create connection: %w", err)
	}

	return nil
}

func (service *SessionService) DestroyConnection(id string, guacToken string) error {
	path := fmt.Sprintf("%s/%s", connectionsURL, id)

	if err := service.makeGuacamoleRequest(
		http.MethodDelete,
		path,
		guacToken,
		nil,
		nil,
	); err != nil {
		return fmt.Errorf("failed to destroy connection: %w", err)
	}

	return nil
}

// makeGuacamoleRequest универсальная функция для выполнения запросов к Guacamole API
func (service *SessionService) makeGuacamoleRequest(
	method string,
	path string,
	guacToken string,
	requestBody interface{},
	responseTarget interface{},
) error {
	url := fmt.Sprintf("%s/%s", config.ServerConfig.GuacamoleAPIURL, path)

	var body io.Reader
	if requestBody != nil {
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Guacamole-Token", strings.TrimSpace(guacToken))
	if body != nil {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}

	resp, err := service.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		errorBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(errorBody))
	}

	if responseTarget != nil {
		if err := json.NewDecoder(resp.Body).Decode(responseTarget); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
