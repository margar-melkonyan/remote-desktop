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
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/repository"
)

// Протоколы по названия которых будут возвращаться подключения
const (
	all = "all"
	ssh = "ssh"
	rdp = "rdp"
)

// SessionService предоставляет сервис для работы c подключениями к удаленным серверам.
type SessionService struct {
	repo repository.UserRepository
}

// NewSessionService создает новый экземпляр SessionService.
//
// Параметры:
//   - repo: репозиторий для работы с подключениями
//
// Возвращает:
//   - *SessionService: указатель на созданный сервис
func NewSessionService(repo repository.UserRepository) *SessionService {
	return &SessionService{
		repo: repo,
	}
}

func fetchConnections(guacToken string) ([]*common.GuacamoleRDConnectionResponse, error) {
	getSesions := fmt.Sprintf(
		"%s/%s",
		config.ServerConfig.GuacamoleAPIURL,
		"session/data/postgresql/connectionGroups/ROOT/tree",
	)
	req, err := http.NewRequest(http.MethodGet, getSesions, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Guacamole-Token", strings.TrimSpace(guacToken))
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	var response struct {
		ChildConnections []*common.GuacamoleRDConnectionResponse `json:"childConnections"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return response.ChildConnections, nil
}

func (service *SessionService) GetSession(protocol string, guacToken string) ([]*common.GuacamoleRDConnectionResponse, error) {
	connections, err := fetchConnections(guacToken)
	if protocol == all {
		return connections, err
	}
	var resultConnection []*common.GuacamoleRDConnectionResponse
	resultConnection = make([]*common.GuacamoleRDConnectionResponse, 0)
	for _, connection := range connections {
		if connection.Protocol == protocol {
			resultConnection = append(resultConnection, connection)
		}
	}

	return resultConnection, err
}

func (service *SessionService) CreateConnection(form *common.GuacamoleConnectionRequest, guacToken string) error {
	ignoreCert := "false"

	if form.Protocol == rdp {
		ignoreCert = "true"
	}

	body := common.GuacamoleRDConnectionRequest{
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

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("JSON marshaling failed: %w", err)
	}

	reader := bytes.NewReader(bodyJson)
	getSesions := fmt.Sprintf(
		"%s/%s",
		config.ServerConfig.GuacamoleAPIURL,
		"session/data/postgresql/connections",
	)
	req, err := http.NewRequest(http.MethodPost, getSesions, reader)
	if err != nil {
		return err
	}
	req.Header.Add("Guacamole-Token", strings.TrimSpace(guacToken))
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func (service *SessionService) DestroyConnection(id string, guacToken string) error {
	getSesions := fmt.Sprintf(
		"%s/%s/%s",
		config.ServerConfig.GuacamoleAPIURL,
		"session/data/postgresql/connections",
		id,
	)
	req, err := http.NewRequest(http.MethodDelete, getSesions, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Guacamole-Token", strings.TrimSpace(guacToken))
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}
