package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/neonlabsorg/neon-service-framework/pkg/logger"
	"github.com/neonlabsorg/neon-service-framework/pkg/service/configuration"
)

type PrometheusAdapter struct {
	cfg        *configuration.PrometheusAdapterConfig
	clientName string
	sendURL    string
	httpClient *http.Client
	context    *Context
	log        logger.Logger
}

func NewPrometheusAdapter(
	cfg *configuration.PrometheusAdapterConfig,
	context *Context,
	log logger.Logger,
) *PrometheusAdapter {
	return &PrometheusAdapter{
		cfg:        cfg,
		sendURL:    cfg.URL + "/api/v2/alerts",
		httpClient: &http.Client{},
		clientName: "AlertManagerClient",
		context:    context,
		log:        log,
	}
}

func (s *PrometheusAdapter) GetName() string {
	return "prometheus"
}

func (s *PrometheusAdapter) Send(alert Alert) error {
	if s.cfg.IsDemo {
		s.log.Info().Msgf(
			"%s: allert was succesful sent on demo mode",
			s.GetName(),
			alert.GetName(),
		)
		return nil
	}

	for i := 1; i <= s.cfg.Attempts; i++ {
		err := s.send(alert)
		if err == nil {
			return nil
		}

		s.log.Error().Err(err).Msgf(
			"%s: attempt #%d to send alert %s failed: %s",
			s.GetName(),
			i,
			alert.GetName(),
		)

		if i != s.cfg.Attempts {
			time.Sleep(s.cfg.Interval)
		}
	}

	return fmt.Errorf(
		"%s: attempts to send alert %s exceeded",
		s.GetName(),
		alert.GetName(),
	)
}

func (s *PrometheusAdapter) send(alert Alert) error {
	reqBodyJSON, err := s.buildRequestBodyJSON(alert)
	if err != nil {
		return err
	}

	resp, err := s.httpClient.Post(s.sendURL, "application/json", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response status code - %d; response body: %s", resp.StatusCode, respBody)
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *PrometheusAdapter) buildRequestBodyJSON(alert Alert) ([]byte, error) {
	alertID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	annotations := map[string]string{
		"summary":     alert.GetSummary(),
		"description": alert.GetDescription(),
	}
	for k, v := range alert.GetAdditionalAnnotations() {
		annotations[k] = v
	}

	labels := map[string]string{
		"alertname": alert.GetName().String(),
		"code":      alert.GetCode().String(),
		"severity":  alert.GetSeverity().String(),
		"project":   s.context.GetProject(),
		"service":   s.context.GetService(),
		"instance":  s.context.GetInstance(),
		"sender":    s.clientName,
		"alertID":   alertID.String(),
	}
	for k, v := range alert.GetAdditionalLabels() {
		labels[k] = v
	}

	reqBody := []map[string]interface{}{
		{
			"startsAt":    alert.GetDate(),
			"annotations": annotations,
			"labels":      labels,
		},
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	return reqBodyJSON, nil
}
