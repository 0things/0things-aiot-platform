package dto

// ThingPropertyPostPayload defines standard Alink/ICA uplink property report structure:
//
//	{
//	    "version": "1.0",
//	    "sys": { "ack": 0 },
//	    "params": {
//	        "Power": { "value": "on", "time": 1524448722000 },
//	        "WF": { "value": 23.6, "time": 1524448722000 }
//	    }
//	}
type ThingPropertyPostPayload struct {
	Version string                   `json:"version,omitempty"`
	Sys     *ThingSys                `json:"sys,omitempty"`
	Params  map[string]PropertyValue `json:"params" validate:"required,min=1,dive"`
}

// ThingSys represents optional system control parameters.
type ThingSys struct {
	Ack int `json:"ack,omitempty"`
}

// PropertyValue represents a single property metric value and its sampling timestamp (in milliseconds).
type PropertyValue struct {
	Value interface{} `json:"value" validate:"required"`
	Time  int64       `json:"time" validate:"required,gt=0"` // Millisecond Unix timestamp
}
