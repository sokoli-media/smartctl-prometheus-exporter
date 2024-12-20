package prometheus_exporter

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log/slog"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var metricsLabelsNames = []string{"device", "model_name", "serial_number"}

var userCapacityBlocks = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "smartctl_user_capacity_blocks"}, metricsLabelsNames)
var userCapacityBytes = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "smartctl_user_capacity_bytes"}, metricsLabelsNames)

var interfaceSpeedCurrent = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "smartctl_interface_speed_current"}, metricsLabelsNames)
var interfaceSpeedMax = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "smartctl_interface_speed_max"}, metricsLabelsNames)

var smartStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "smartctl_smart_status"}, metricsLabelsNames)

var smartAttributeLabels = append(metricsLabelsNames, "attribute")
var ataSmartAttribute = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "smartctl_ata_smart_attribute_value"}, smartAttributeLabels)

var powerOnTimeHours = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "smartctl_power_on_time_hours"}, metricsLabelsNames)

var powerCycleTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "smartctl_power_cycle_time"}, metricsLabelsNames)

var temperature = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "smartctl_temperature"}, metricsLabelsNames)

var lastUpdate = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "smartctl_last_update"}, metricsLabelsNames)

type SmartCtlDevice struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type SmartCtlScan struct {
	Devices []SmartCtlDevice `json:"devices"`
}

func scan() (*SmartCtlScan, error) {
	out, err := exec.Command("smartctl", "--scan-open", "--json").Output()
	if err != nil {
		return nil, err
	}

	var smartCtlScan SmartCtlScan
	if err := json.Unmarshal(out, &smartCtlScan); err != nil {
		return nil, err
	}

	return &smartCtlScan, nil
}

func fetchDeviceMetrics(logger *slog.Logger, device SmartCtlDevice) error {
	out, err := exec.Command("smartctl", device.Name, "-d", device.Type, "-a", "--json").Output()
	if err != nil {
		return err
	}

	return loadMetricsFromDeviceScan(logger, device, out)
}

func getLabels(device SmartCtlDevice, deviceScan SmartCtlDeviceScan) prometheus.Labels {
	return prometheus.Labels{
		"device":        device.Name,
		"model_name":    deviceScan.ModelName,
		"serial_number": deviceScan.SerialNumber,
	}
}

func normalizeString(input string) string {
	input = strings.ReplaceAll(input, "\n", " ")

	re := regexp.MustCompile(`\s+`)
	input = re.ReplaceAllString(input, " ")

	return strings.TrimSpace(input)
}

func loadMetricsFromDeviceScan(logger *slog.Logger, device SmartCtlDevice, commandOutput []byte) error {
	var smartCtlDeviceScan SmartCtlDeviceScan
	if err := json.Unmarshal(commandOutput, &smartCtlDeviceScan); err != nil {
		return err
	}

	dumpedScan, _ := json.Marshal(smartCtlDeviceScan)
	logger.Info(
		"loaded smartctl device scan",
		"commandOutput",
		normalizeString(string(commandOutput)),
		"loadedScan",
		string(dumpedScan),
	)

	deviceMetricLabels := getLabels(device, smartCtlDeviceScan)

	userCapacityBlocks.With(deviceMetricLabels).Set(float64(smartCtlDeviceScan.UserCapacity.Blocks))
	userCapacityBytes.With(deviceMetricLabels).Set(float64(smartCtlDeviceScan.UserCapacity.Bytes))

	currentInterfaceSpeed := smartCtlDeviceScan.InterfaceSpeed.Current
	maxInterfaceSpeed := smartCtlDeviceScan.InterfaceSpeed.Max
	interfaceSpeedCurrent.With(deviceMetricLabels).Set(float64(currentInterfaceSpeed.BitsPerUnit * currentInterfaceSpeed.UnitsPerSecond))
	interfaceSpeedMax.With(deviceMetricLabels).Set(float64(maxInterfaceSpeed.BitsPerUnit * maxInterfaceSpeed.UnitsPerSecond))

	smartStatusValue := 0
	if smartCtlDeviceScan.SmartStatus.Passed {
		smartStatusValue = 1
	}
	smartStatus.With(deviceMetricLabels).Set(float64(smartStatusValue))

	for _, attribute := range smartCtlDeviceScan.ATASMARTAttributes.Table {
		attributeLabels := getLabels(device, smartCtlDeviceScan)
		attributeLabels["attribute"] = attribute.Name

		ataSmartAttribute.With(attributeLabels).Set(float64(attribute.Value))
	}

	powerOnTimeHours.With(deviceMetricLabels).Set(float64(smartCtlDeviceScan.PowerOnTime.Hours))

	powerCycleTime.With(deviceMetricLabels).Set(float64(smartCtlDeviceScan.PowerCycleCount))

	temperature.With(deviceMetricLabels).Set(float64(smartCtlDeviceScan.Temperature.Current))

	lastUpdate.With(deviceMetricLabels).SetToCurrentTime()

	return nil
}

func fetchSmartCtlMetrics(logger *slog.Logger) {
	logger.Info("looking for devices")
	scanResult, err := scan()
	if err != nil {
		logger.Error("scanning for devices failed", "error", err)
		return
	}

	for _, device := range scanResult.Devices {
		logger.Info("scanning device", "device", device.Name, "type", device.Type)
		err = fetchDeviceMetrics(logger, device)
		if err != nil {
			logger.Error("fetching metrics for device failed", "error", err, "device", device.Name)
		}
	}
}

func CollectSmartCtlStats(logger *slog.Logger) {
	ticker := time.NewTicker(5 * time.Minute)
	for {
		go fetchSmartCtlMetrics(logger)

		<-ticker.C
	}
}
