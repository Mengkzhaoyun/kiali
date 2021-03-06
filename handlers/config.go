package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/common/model"
	"gopkg.in/yaml.v2"

	"github.com/kiali/kiali/config"
	"github.com/kiali/kiali/log"
	"github.com/kiali/kiali/prometheus"
	"github.com/kiali/kiali/status"
)

const (
	defaultPrometheusGlobalScrapeInterval = 15 // seconds
)

type ThreeScaleConfig struct {
	AdapterName    string `json:"adapterName"`
	AdapterPort    string `json:"adapterPort"`
	AdapterService string `json:"adapterService"`
	Enabled        bool   `json:"enabled"`
	TemplateName   string `json:"templateName"`
}
type Iter8Config struct {
	Enabled   bool   `json:"enabled"`
	Namespace string `json:"namespace"`
}
type Extensions struct {
	ThreeScale ThreeScaleConfig `json:"threescale,omitempty"`
	Iter8      Iter8Config      `json:"iter8,omitempty"`
}

// PrometheusConfig holds actual Prometheus configuration that is useful to Kiali.
// All durations are in seconds.
type PrometheusConfig struct {
	GlobalScrapeInterval int64 `json:"globalScrapeInterval,omitempty"`
	StorageTsdbRetention int64 `json:"storageTsdbRetention,omitempty"`
}

// PublicConfig is a subset of Kiali configuration that can be exposed to clients to
// help them interact with the system.
type PublicConfig struct {
	Extensions               Extensions                      `json:"extensions,omitempty"`
	InstallationTag          string                          `json:"installationTag,omitempty"`
	IstioStatusEnabled       bool                            `json:"istioStatusEnabled,omitempty"`
	IstioIdentityDomain      string                          `json:"istioIdentityDomain,omitempty"`
	IstioNamespace           string                          `json:"istioNamespace,omitempty"`
	IstioComponentNamespaces config.IstioComponentNamespaces `json:"istioComponentNamespaces,omitempty"`
	IstioLabels              config.IstioLabels              `json:"istioLabels,omitempty"`
	IstioTelemetryV2         bool                            `json:"istioTelemetryV2"`
	KialiFeatureFlags        config.KialiFeatureFlags        `json:"kialiFeatureFlags,omitempty"`
	Prometheus               PrometheusConfig                `json:"prometheus,omitempty"`
}

// Config is a REST http.HandlerFunc serving up the Kiali configuration made public to clients.
func Config(w http.ResponseWriter, r *http.Request) {
	defer handlePanic(w)

	// Note that determine the Prometheus config at request time because it is not
	// guaranteed to remain the same during the Kiali lifespan.
	promConfig := getPrometheusConfig()
	config := config.Get()
	publicConfig := PublicConfig{
		Extensions: Extensions{
			ThreeScale: ThreeScaleConfig{
				AdapterName:    config.Extensions.ThreeScale.AdapterName,
				AdapterPort:    config.Extensions.ThreeScale.AdapterPort,
				AdapterService: config.Extensions.ThreeScale.AdapterService,
				Enabled:        config.Extensions.ThreeScale.Enabled,
				TemplateName:   config.Extensions.ThreeScale.TemplateName,
			},
			Iter8: Iter8Config{
				Enabled:   config.Extensions.Iter8.Enabled,
				Namespace: config.Extensions.Iter8.Namespace,
			},
		},
		InstallationTag:          config.InstallationTag,
		IstioStatusEnabled:       config.ExternalServices.Istio.ComponentStatuses.Enabled,
		IstioIdentityDomain:      config.ExternalServices.Istio.IstioIdentityDomain,
		IstioNamespace:           config.IstioNamespace,
		IstioComponentNamespaces: config.IstioComponentNamespaces,
		IstioLabels:              config.IstioLabels,
		KialiFeatureFlags:        config.KialiFeatureFlags,
		Prometheus: PrometheusConfig{
			GlobalScrapeInterval: promConfig.GlobalScrapeInterval,
			StorageTsdbRetention: promConfig.StorageTsdbRetention,
		},
		IstioTelemetryV2: status.IsMixerDisabled(),
	}

	RespondWithJSONIndent(w, http.StatusOK, publicConfig)
}

type PrometheusPartialConfig struct {
	Global struct {
		Scrape_interval string
	}
}

func getPrometheusConfig() PrometheusConfig {
	promConfig := PrometheusConfig{
		GlobalScrapeInterval: defaultPrometheusGlobalScrapeInterval,
	}

	client, err := prometheus.NewClient()
	if !checkErr(err, "") {
		log.Error(err)
		return promConfig
	}

	configResult, err := client.GetConfiguration()
	if checkErr(err, "Failed to fetch Prometheus configuration") {
		var config PrometheusPartialConfig
		if checkErr(yaml.Unmarshal([]byte(configResult.YAML), &config), "Failed to unmarshal Prometheus configuration") {
			scrapeIntervalString := config.Global.Scrape_interval
			scrapeInterval, err := model.ParseDuration(scrapeIntervalString)
			if checkErr(err, fmt.Sprintf("Invalid global scrape interval [%s]", scrapeIntervalString)) {
				promConfig.GlobalScrapeInterval = int64(time.Duration(scrapeInterval).Seconds())
			}
		}
	}

	flags, err := client.GetFlags()
	if checkErr(err, "Failed to fetch Prometheus flags") {
		if retentionString, ok := flags["storage.tsdb.retention"]; ok {
			retention, err := model.ParseDuration(retentionString)
			if checkErr(err, fmt.Sprintf("Invalid storage.tsdb.retention [%s]", retentionString)) {
				promConfig.StorageTsdbRetention = int64(time.Duration(retention).Seconds())
			}
		}
	}

	return promConfig
}

func checkErr(err error, message string) bool {
	if err != nil {
		log.Errorf("%s: %v", message, err)
		return false
	}
	return true
}
