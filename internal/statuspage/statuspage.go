package statuspage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/TicketsBot-cloud/status-updates/internal/config"
	"github.com/TicketsBot-cloud/status-updates/internal/model"
	"go.uber.org/zap"
)

type StatusPageClient struct {
	logger *zap.Logger
	config config.Config
}

func NewClient(logger *zap.Logger) StatusPageClient {
	return StatusPageClient{
		logger: logger,
		config: config.Conf,
	}
}

func (s *StatusPageClient) GetIncidents() ([]model.Incident, error) {
	url := fmt.Sprintf("https://api.statuspage.io/v1/pages/%s/incidents", s.config.StatusPage.PageId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		s.logger.Error("Error creating request", zap.Error(err))
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", s.config.StatusPage.ApiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("Error making request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("Unexpected status code", zap.Int("status_code", resp.StatusCode))
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Error reading response body", zap.Error(err))
		return nil, err
	}

	var incidents []model.Incident
	err = json.Unmarshal(body, &incidents)
	if err != nil {
		s.logger.Error("Error decoding JSON", zap.Error(err))
		return nil, err
	}

	return incidents, nil
}
