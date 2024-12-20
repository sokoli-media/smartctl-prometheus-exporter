package prometheus_exporter

type SmartCtlDeviceScan struct {
	JsonFormatVersion   []int              `json:"json_format_version"`
	Smartctl            SmartctlInfo       `json:"smartctl"`
	LocalTime           LocalTimeInfo      `json:"local_time"`
	Device              DeviceInfo         `json:"device"`
	ModelName           string             `json:"model_name"`
	SerialNumber        string             `json:"serial_number"`
	WWN                 WWNInfo            `json:"wwn"`
	FirmwareVersion     string             `json:"firmware_version"`
	UserCapacity        CapacityInfo       `json:"user_capacity"`
	LogicalBlockSize    int                `json:"logical_block_size"`
	PhysicalBlockSize   int                `json:"physical_block_size"`
	RotationRate        int                `json:"rotation_rate"`
	FormFactor          FormFactorInfo     `json:"form_factor"`
	Trim                TrimInfo           `json:"trim"`
	InSmartctlDatabase  bool               `json:"in_smartctl_database"`
	ATAVersion          VersionInfo        `json:"ata_version"`
	SATAVersion         VersionInfo        `json:"sata_version"`
	InterfaceSpeed      InterfaceSpeedInfo `json:"interface_speed"`
	SmartSupport        SmartSupportInfo   `json:"smart_support"`
	SmartStatus         SmartStatusInfo    `json:"smart_status"`
	ATASMARTData        ATASMARTData       `json:"ata_smart_data"`
	ATASCTCapabilities  SCTCapabilities    `json:"ata_sct_capabilities"`
	ATASMARTAttributes  SMARTAttributes    `json:"ata_smart_attributes"`
	PowerOnTime         PowerOnTimeInfo    `json:"power_on_time"`
	PowerCycleCount     int                `json:"power_cycle_count"`
	Temperature         TemperatureInfo    `json:"temperature"`
	ATASMARTErrorLog    SMARTErrorLog      `json:"ata_smart_error_log"`
	ATASMARTSelfTestLog SMARTSelfTestLog   `json:"ata_smart_self_test_log"`
}

type SmartctlInfo struct {
	Version              []int         `json:"version"`
	PreRelease           bool          `json:"pre_release"`
	SVNRevision          string        `json:"svn_revision"`
	PlatformInfo         string        `json:"platform_info"`
	BuildInfo            string        `json:"build_info"`
	Argv                 []string      `json:"argv"`
	DriveDatabaseVersion VersionString `json:"drive_database_version"`
	ExitStatus           int           `json:"exit_status"`
}

type VersionString struct {
	String string `json:"string"`
}

type LocalTimeInfo struct {
	TimeT   int64  `json:"time_t"`
	Asctime string `json:"asctime"`
}

type DeviceInfo struct {
	Name     string `json:"name"`
	InfoName string `json:"info_name"`
	Type     string `json:"type"`
	Protocol string `json:"protocol"`
}

type WWNInfo struct {
	NAA int   `json:"naa"`
	OUI int   `json:"oui"`
	ID  int64 `json:"id"`
}

type CapacityInfo struct {
	Blocks int64 `json:"blocks"`
	Bytes  int64 `json:"bytes"`
}

type FormFactorInfo struct {
	ATAValue int    `json:"ata_value"`
	Name     string `json:"name"`
}

type TrimInfo struct {
	Supported bool `json:"supported"`
}

type VersionInfo struct {
	String     string `json:"string"`
	MajorValue int    `json:"major_value,omitempty"`
	MinorValue int    `json:"minor_value,omitempty"`
}

type InterfaceSpeedInfo struct {
	Max     SpeedDetail `json:"max"`
	Current SpeedDetail `json:"current"`
}

type SpeedDetail struct {
	SATAValue      int    `json:"sata_value"`
	String         string `json:"string"`
	UnitsPerSecond int    `json:"units_per_second"`
	BitsPerUnit    int    `json:"bits_per_unit"`
}

type SmartSupportInfo struct {
	Available bool `json:"available"`
	Enabled   bool `json:"enabled"`
}

type SmartStatusInfo struct {
	Passed bool `json:"passed"`
}

type ATASMARTData struct {
	OfflineDataCollection OfflineDataCollection `json:"offline_data_collection"`
	SelfTest              SelfTest              `json:"self_test"`
	Capabilities          Capabilities          `json:"capabilities"`
}

type OfflineDataCollection struct {
	Status            StatusDetail `json:"status"`
	CompletionSeconds int          `json:"completion_seconds"`
}

type StatusDetail struct {
	Value  int    `json:"value"`
	String string `json:"string"`
}

type SelfTest struct {
	Status         SelfTestStatus `json:"status"`
	PollingMinutes PollingTimes   `json:"polling_minutes"`
}

type SelfTestStatus struct {
	Value  int    `json:"value"`
	String string `json:"string"`
	Passed bool   `json:"passed"`
}

type PollingTimes struct {
	Short    int `json:"short"`
	Extended int `json:"extended"`
}

type Capabilities struct {
	Values                     []int `json:"values"`
	ExecOfflineImmediate       bool  `json:"exec_offline_immediate_supported"`
	OfflineIsAbortedUponNewCmd bool  `json:"offline_is_aborted_upon_new_cmd"`
	OfflineSurfaceScan         bool  `json:"offline_surface_scan_supported"`
	SelfTestsSupported         bool  `json:"self_tests_supported"`
	ConveyanceSelfTest         bool  `json:"conveyance_self_test_supported"`
	SelectiveSelfTest          bool  `json:"selective_self_test_supported"`
	AttributeAutosaveEnabled   bool  `json:"attribute_autosave_enabled"`
	ErrorLoggingSupported      bool  `json:"error_logging_supported"`
	GPLoggingSupported         bool  `json:"gp_logging_supported"`
}

type SCTCapabilities struct {
	Value                int  `json:"value"`
	ErrorRecoveryControl bool `json:"error_recovery_control_supported"`
	FeatureControl       bool `json:"feature_control_supported"`
	DataTableSupported   bool `json:"data_table_supported"`
}

type SMARTAttributes struct {
	Revision int              `json:"revision"`
	Table    []SMARTAttribute `json:"table"`
}

type SMARTAttribute struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Value      int        `json:"value"`
	Worst      int        `json:"worst"`
	Thresh     int        `json:"thresh"`
	WhenFailed string     `json:"when_failed"`
	Flags      SMARTFlags `json:"flags"`
	Raw        RawData    `json:"raw"`
}

type SMARTFlags struct {
	Value         int    `json:"value"`
	String        string `json:"string"`
	Prefailure    bool   `json:"prefailure"`
	UpdatedOnline bool   `json:"updated_online"`
	Performance   bool   `json:"performance"`
	ErrorRate     bool   `json:"error_rate"`
	EventCount    bool   `json:"event_count"`
	AutoKeep      bool   `json:"auto_keep"`
}

type RawData struct {
	Value  int    `json:"value"`
	String string `json:"string"`
}

type PowerOnTimeInfo struct {
	Hours int `json:"hours"`
}

type TemperatureInfo struct {
	Current int `json:"current"`
}

type SMARTErrorLog struct {
	Summary ErrorSummary `json:"summary"`
}

type ErrorSummary struct {
	Revision int `json:"revision"`
	Count    int `json:"count"`
}

type SMARTSelfTestLog struct {
	Standard TestSummary `json:"standard"`
}

type TestSummary struct {
	Revision int `json:"revision"`
	Count    int `json:"count"`
}
