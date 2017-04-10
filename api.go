package nightscout

import (
	"time"
)

// The JSON encoding of this information must match what is expected by the
// Nightscout upload API in https://github.com/nightscout/cgm-remote-monitor

const (
	DateStringLayout = "2006-01-02T15:04:05.000-0700"
)

type (
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

	// Struct used to unmarshal just the date field.
	EntryTime struct {
		Date int64 `json:"date"` // Unix time in milliseconds
	}

	DeviceStatus struct {
		Device   string   `json:"device"`
		Openaps  Openaps  `json:"openaps,omitempty"`
		Pump     Pump     `json:"pump,omitempty"`
		Uploader Uploader `json:"uploader,omitempty"`
	}

	Openaps struct {
		Iob Iob `json:"iob"`
	}

	Iob struct {
		Iob Insulin `json:"iob"`
	}

	Pump struct {
		Battery   Battery   `json:"battery"`
		Clock     time.Time `json:"clock"`
		Reservoir Insulin   `json:"reservoir"`
		Status    Status    `json:"status"`
	}

	Battery struct {
		Voltage Voltage `json:"voltage"`
	}

	Status struct {
		Status    string `json:"status"`
		Bolusing  bool   `json:"bolusing"`
		Suspended bool   `json:"suspended"`
	}

	Uploader struct {
		BatteryLevel   int     `json:"battery"`
		BatteryVoltage Voltage `json:"batteryVoltage,omitempty"`
		RawBattery     int     `json:"rawBattery,omitempty"`
	}

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

	// Structure used to unmarshal just the created_at field.
	TreatmentTime struct {
		CreatedAt time.Time `json:"created_at"`
	}

	Profile struct {
		Id             string                 `json:"_id"`
		CreatedAt      time.Time              `json:"created_at"`
		StartDate      time.Time              `json:"startDate"`
		DefaultProfile string                 `json:"defaultProfile"`
		Store          map[string]ProfileData `json:"store"`
	}

	// Structure used to unmarshal just the _id field.
	ProfileId struct {
		Id string `json:"_id"`
	}

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

	Schedule []TimeValue

	TimeValue struct {
		Time  string      `json:"time"`
		Value interface{} `json:"value"`
	}

	// These correspond to types defined in the medtronic package.
	Glucose int
	Insulin float64
	Voltage float64
)
