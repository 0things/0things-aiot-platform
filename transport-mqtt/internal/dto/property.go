package dto

// ThingPropertyPostPayload defines standard Alink/ICA uplink property report structure:
//
//	{
//	    "id": "123",
//	    "version": "1.0",
//	    "sys": { "ack": 0 },
//	    "params": {
//	        "Power": { "value": "on", "time": 1524448722000 },
//	        "WF": { "value": 23.6, "time": 1524448722000 }
//	    },
//	    "method": "thing.event.property.post"
//	}
type ThingPropertyPostPayload struct {
	ID      string                   `json:"id" validate:"required"`
	Version string                   `json:"version,omitempty"`
	Sys     *ThingSys                `json:"sys,omitempty"`
	Params  map[string]PropertyValue `json:"params" validate:"required,min=1,dive"`
	Method  string                   `json:"method,omitempty"`
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
