package nightscout

import (
	"time"
)

// The JSON encoding of this information must match what is expected by the
// Nightscout upload API in https://github.com/nightscout/cgm-remote-monitor

const (
	// DateStringLayout is a time format accepted by Nightscout.
	DateStringLayout = time.RFC3339
)

type (
	// Entry represents data for the Nightscout entries API.
	Entry struct {
		Type       string  `json:"type"`
		Date       int64   `json:"date"` // Unix time in milliseconds
		DateString string  `json:"dateString"`
		Device     string  `json:"device,omitempty"`
		SGV        int     `json:"sgv,omitempty"`
		Direction  string  `json:"direction,omitempty"`
		Filtered   int     `json:"filtered,omitempty"`
		Unfiltered int     `json:"unfiltered,omitempty"`
		RSSI       int     `json:"rssi,omitempty"`
		Noise      int     `json:"noise,omitempty"`
		Slope      float64 `json:"slope,omitempty"`
		Intercept  float64 `json:"intercept,omitempty"`
		Scale      float64 `json:"scale,omitempty"`
		MBG        int     `json:"mbg,omitempty"`
	}

	// Entries represents a sequence of Entry values.
	Entries []Entry

	// DeviceStatus represents data for the Nightscout devicestatus API.
	DeviceStatus struct {
		CreatedAt time.Time `json:"created_at"`
		Device    string    `json:"device"`
		Openaps   Openaps   `json:"openaps,omitempty"`
		Pump      Pump      `json:"pump,omitempty"`
		Uploader  Uploader  `json:"uploader,omitempty"`
	}

	// Openaps represents the openaps data in a DeviceStatus record.
	Openaps struct {
		IOB IOB `json:"iob"`
	}

	// IOB represents the insulin-on-board data in an Openaps record.
	IOB struct {
		IOB Insulin `json:"iob"`
	}

	// Pump represents the pump data in a DeviceStatus record.
	Pump struct {
		Battery   Battery   `json:"battery"`
		Clock     time.Time `json:"clock"`
		Reservoir Insulin   `json:"reservoir"`
		Status    Status    `json:"status"`
	}

	// Battery represents the battery data in a Pump record.
	Battery struct {
		Voltage Voltage `json:"voltage"`
	}

	// Status represents the status data in a Pump record.
	Status struct {
		Status    string `json:"status"`
		Bolusing  bool   `json:"bolusing"`
		Suspended bool   `json:"suspended"`
	}

	// Uploader represents the uploader data in a DeviceStatus record.
	Uploader struct {
		BatteryLevel   int     `json:"battery"`
		BatteryVoltage Voltage `json:"batteryVoltage,omitempty"`
		RawBattery     int     `json:"rawBattery,omitempty"`
	}

	// Treatment represents data for the Nightscout treatments API.
	Treatment struct {
		CreatedAt time.Time `json:"created_at"`
		EventType string    `json:"eventType"`
		EnteredBy string    `json:"enteredBy,omitempty"`
		Glucose   *Glucose  `json:"glucose,omitempty"`
		Absolute  *Insulin  `json:"absolute,omitempty"`
		Duration  *int      `json:"duration,omitempty"` // minutes
		Insulin   *Insulin  `json:"insulin,omitempty"`
		Units     string    `json:"units,omitempty"`
	}

	// TreatmentTime is used to unmarshal just the CreatedAt field of a Treatment.
	TreatmentTime struct {
		CreatedAt time.Time `json:"created_at"`
	}

	// Profile represents data for the Nightscout profile API.
	Profile struct {
		ID             string                 `json:"_id"`
		CreatedAt      time.Time              `json:"created_at"`
		StartDate      time.Time              `json:"startDate"`
		DefaultProfile string                 `json:"defaultProfile"`
		Store          map[string]ProfileData `json:"store"`
	}

	// ProfileID is used to unmarshal just the ID field of a Profile.
	ProfileID struct {
		ID string `json:"_id"`
	}

	// ProfileData represents the information in a Profile record.
	ProfileData struct {
		DIA        int      `json:"dia"` // hours
		Basal      Schedule `json:"basal"`
		CarbRatio  Schedule `json:"carbratio"`
		Sens       Schedule `json:"sens"`
		TargetLow  Schedule `json:"target_low"`
		TargetHigh Schedule `json:"target_high"`
		TimeZone   string   `json:"timezone"`
		Units      string   `json:"units"`
	}

	// Schedule represents a sequence of times and values.
	Schedule []TimeValue

	// TimeValue represents a value with an associated time.
	TimeValue struct {
		Time  string      `json:"time"`
		Value interface{} `json:"value"`
	}

	// The following types are defined here to avoid
	// a circular dependency with the medtronic package.

	// Glucose corresponds to the medtronic.Glucose type.
	Glucose int
	// Insulin corresponds to the medtronic.Insulin type.
	Insulin float64
	// Voltage corresponds to the medtronic.Voltage type.
	Voltage float64
)

// Values for the Entry Type field.
const (
	SGVType = "sgv"
	MBGType = "mbg"
	CalType = "cal"
)
