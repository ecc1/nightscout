package nightscout

import (
	"time"

	"github.com/ecc1/medtronic"
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
		Sgv        uint16  `json:"sgv,omitempty"`
		Direction  string  `json:"direction,omitempty"`
		Filtered   uint32  `json:"filtered,omitempty"`
		Unfiltered uint32  `json:"unfiltered,omitempty"`
		Rssi       uint16  `json:"rssi,omitempty"`
		Slope      float64 `json:"slope,omitempty"`
		Intercept  float64 `json:"intercept,omitempty"`
		Scale      float64 `json:"scale,omitempty"`
		Mbg        uint16  `json:"mbg,omitempty"`
	}

	// Struct used to unmarshal just the date field.
	EntryTime struct {
		Date int64 `json:"date"` // Unix time in milliseconds
	}

	DeviceStatus struct {
		Device  string  `json:"device"`
		Openaps Openaps `json:"openaps"`
		Pump    Pump    `json:"pump"`
	}

	Openaps struct {
		Iob Iob `json:"iob"`
	}

	Iob struct {
		Iob medtronic.Insulin `json:"iob"`
	}

	Pump struct {
		Battery   Battery           `json:"battery"`
		Clock     time.Time         `json:"clock"`
		Reservoir medtronic.Insulin `json:"reservoir"`
		Status    Status            `json:"status"`
	}

	Battery struct {
		Voltage medtronic.Voltage `json:"voltage"`
	}

	Status struct {
		Status    string `json:"status"`
		Bolusing  bool   `json:"bolusing"`
		Suspended bool   `json:"suspended"`
	}

	Treatment struct {
		EventTime time.Time          `json:"eventTime"`
		EventType string             `json:"eventType"`
		EnteredBy string             `json:"enteredBy,omitempty"`
		Glucose   *medtronic.Glucose `json:"glucose,omitempty"`
		Absolute  *medtronic.Insulin `json:"absolute,omitempty"`
		Duration  *int               `json:"duration,omitempty"` // minutes
		Insulin   *medtronic.Insulin `json:"insulin,omitempty"`
	}

	// Structure used to unmarshal just the created_at field.
	TreatmentTime struct {
		CreatedAt time.Time `json:"created_at"`
	}

	Profile struct {
		Id             string         `json:"_id"`
		CreatedAt      time.Time      `json:"created_at"`
		StartDate      time.Time      `json:"startDate"`
		DefaultProfile string         `json:"defaultProfile"`
		Store          DefaultProfile `json:"store"`
	}

	DefaultProfile struct {
		Default ProfileData
	}

	ProfileData struct {
		InsulinAction int         `json:"dia"` // hours
		Basal         []TimeValue `json:"basal"`
		CarbRatio     []TimeValue `json:"carbratio"`
		Sens          []TimeValue `json:"sens"`
		TargetLow     []TimeValue `json:"target_low"`
		TargetHigh    []TimeValue `json:"target_high"`
		TimeZone      string      `json:"timezone"`
		Units         string      `json:"units"`
	}

	TimeValue struct {
		Time  string      `json:"time"`
		Value interface{} `json:"value"`
	}

	// Structure used to unmarshal just the _id field.
	ProfileId struct {
		Id string `json:"_id"`
	}
)
