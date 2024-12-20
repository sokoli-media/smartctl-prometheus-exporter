package prometheus_exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"
)

var sLogForTesting = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func getLabelsForTest(device SmartCtlDevice) prometheus.Labels {
	return prometheus.Labels{
		"device":        device.Name,
		"model_name":    "WDC WD80EFZZ-68BTXN0",
		"serial_number": "WD-CAZXCVBN",
	}
}

func getGaugeVecValue(t *testing.T, metric *prometheus.GaugeVec, labels prometheus.Labels) float64 {
	var m = &dto.Metric{}
	if err := metric.With(labels).Write(m); err != nil {
		t.Fatalf("couldnt get metric with metricsLabelsNames: %s", err)
	}
	return m.Gauge.GetValue()
}

func Test_LoadUserCapacity(t *testing.T) {
	var output = `{
		"model_name": "WDC WD80EFZZ-68BTXN0",
		"serial_number": "WD-CAZXCVBN",
		"user_capacity": {
			"blocks": 15628053168,
			"bytes": 8001563222016
		}
	}`

	device := SmartCtlDevice{
		Name: "/dev/sdb",
		Type: "sat",
	}

	err := loadMetricsFromDeviceScan(sLogForTesting, device, []byte(output))
	require.NoError(t, err)

	testLabels := getLabelsForTest(device)
	require.Equal(t, float64(15628053168), getGaugeVecValue(t, userCapacityBlocks, testLabels))
	require.Equal(t, float64(8001563222016), getGaugeVecValue(t, userCapacityBytes, testLabels))
}

func Test_LoadInterfaceSpeed(t *testing.T) {
	var output = `{
	  "model_name": "WDC WD80EFZZ-68BTXN0",
	  "serial_number": "WD-CAZXCVBN",
	  "interface_speed": {
		"max": {
		  "sata_value": 14,
		  "string": "6.0 Gb/s",
		  "units_per_second": 60,
		  "bits_per_unit": 100000000
		},
		"current": {
		  "sata_value": 3,
		  "string": "3.0 Gb/s",
		  "units_per_second": 30,
		  "bits_per_unit": 100000000
		}
	  }
	}`

	device := SmartCtlDevice{
		Name: "/dev/sdb",
		Type: "sat",
	}

	err := loadMetricsFromDeviceScan(sLogForTesting, device, []byte(output))
	require.NoError(t, err)

	testLabels := getLabelsForTest(device)
	require.Equal(t, float64(30*100000000), getGaugeVecValue(t, interfaceSpeedCurrent, testLabels))
	require.Equal(t, float64(60*100000000), getGaugeVecValue(t, interfaceSpeedMax, testLabels))
}

func Test_LoadSmartStatus_Passed(t *testing.T) {
	var output = `{
	  "model_name": "WDC WD80EFZZ-68BTXN0",
	  "serial_number": "WD-CAZXCVBN",
	  "smart_status": {
		"passed": true
	  }
	}`

	device := SmartCtlDevice{
		Name: "/dev/sdb",
		Type: "sat",
	}

	err := loadMetricsFromDeviceScan(sLogForTesting, device, []byte(output))
	require.NoError(t, err)

	testLabels := getLabelsForTest(device)
	require.Equal(t, float64(1), getGaugeVecValue(t, smartStatus, testLabels))
}

func Test_LoadSmartStatus_Failed(t *testing.T) {
	var output = `{
	  "model_name": "WDC WD80EFZZ-68BTXN0",
	  "serial_number": "WD-CAZXCVBN",
	  "smart_status": {
		"passed": false
	  }
	}`

	device := SmartCtlDevice{
		Name: "/dev/sdb",
		Type: "sat",
	}

	err := loadMetricsFromDeviceScan(sLogForTesting, device, []byte(output))
	require.NoError(t, err)

	testLabels := getLabelsForTest(device)
	require.Equal(t, float64(0), getGaugeVecValue(t, smartStatus, testLabels))
}

func Test_LoadAtaAttributes(t *testing.T) {
	var output = `{
	  "model_name": "WDC WD80EFZZ-68BTXN0",
	  "serial_number": "WD-CAZXCVBN",
	  "ata_smart_attributes": {
		"table": [
		  {
			"name": "Raw_Read_Error_Rate",
			"value": 200,
			"worst": 201
		  },
		  {
			"name": "Spin_Up_Time",
			"value": 253,
			"worst": 188
		  }
		]
	  }
	}`

	device := SmartCtlDevice{
		Name: "/dev/sdb",
		Type: "sat",
	}

	err := loadMetricsFromDeviceScan(sLogForTesting, device, []byte(output))
	require.NoError(t, err)

	errorRateLabels := getLabelsForTest(device)
	errorRateLabels["attribute"] = "Raw_Read_Error_Rate"
	require.Equal(t, float64(200), getGaugeVecValue(t, ataSmartAttribute, errorRateLabels))

	spinUpLabels := getLabelsForTest(device)
	spinUpLabels["attribute"] = "Spin_Up_Time"
	require.Equal(t, float64(253), getGaugeVecValue(t, ataSmartAttribute, spinUpLabels))
}

func Test_LoadPowerOnTime(t *testing.T) {
	var output = `{
	  "model_name": "WDC WD80EFZZ-68BTXN0",
	  "serial_number": "WD-CAZXCVBN",
	  "power_on_time": {
		"hours": 10973
	  }
	}`

	device := SmartCtlDevice{
		Name: "/dev/sdb",
		Type: "sat",
	}

	err := loadMetricsFromDeviceScan(sLogForTesting, device, []byte(output))
	require.NoError(t, err)

	testLabels := getLabelsForTest(device)
	require.Equal(t, float64(10973), getGaugeVecValue(t, powerOnTimeHours, testLabels))
}

func Test_LoadPowerCycleTime(t *testing.T) {
	var output = `{
	  "model_name": "WDC WD80EFZZ-68BTXN0",
	  "serial_number": "WD-CAZXCVBN",
	  "power_cycle_count": 38
	}`

	device := SmartCtlDevice{
		Name: "/dev/sdb",
		Type: "sat",
	}

	err := loadMetricsFromDeviceScan(sLogForTesting, device, []byte(output))
	require.NoError(t, err)

	testLabels := getLabelsForTest(device)
	require.Equal(t, float64(38), getGaugeVecValue(t, powerCycleTime, testLabels))
}

func Test_LoadTemperature(t *testing.T) {
	var output = `{
	  "model_name": "WDC WD80EFZZ-68BTXN0",
	  "serial_number": "WD-CAZXCVBN",
	  "temperature": {
		"current": 35
	  }
	}`

	device := SmartCtlDevice{
		Name: "/dev/sdb",
		Type: "sat",
	}

	err := loadMetricsFromDeviceScan(sLogForTesting, device, []byte(output))
	require.NoError(t, err)

	testLabels := getLabelsForTest(device)
	require.Equal(t, float64(35), getGaugeVecValue(t, temperature, testLabels))
}
