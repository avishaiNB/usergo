package rabbitmq2

import (
	"encoding/json"
)

// Message is the basic unit used to send and received data thought rabbitmq
// or more information about the URN http://masstransit-project.com/MassTransit/architecture/interoperability.html
type Message struct {
	URN            string                 `json:"-"`
	Data           json.RawMessage        `json:"data,omitempty"`
	CorrelationID  string                 `json:"correlationId"`
	AdditionalData map[string]interface{} `json:"additionalData"`
}

/*
MarshalJSON will override the default marshalJSON when doing json.Marshal().
We need to override this in order to flat the data field of the messages.

From:

	"message": {
		"data": {
			"mapInsuranceCompanyDrawInformationData": {
				"drawRef": 152677,
				"insuranceCompanyRef": 2
			},
		}
		"correlationId": "b4890000-e89f-e454-ad90-08d81f435381",
		"additionalData": {}
  	},

To:

	"message": {
		"mapInsuranceCompanyDrawInformationData": {
			"drawRef": 152677,
			"insuranceCompanyRef": 2
		},
		"correlationId": "b4890000-e89f-e454-ad90-08d81f435381",
		"additionalData": {}
  	},

*/
func (e *Message) MarshalJSON() ([]byte, error) {
	type duplicate Message
	ev := duplicate(*e)

	dataJSON, err := json.Marshal(ev.Data)
	if err != nil {
		return nil, err
	}

	ev.Data = nil
	eventJSON, err := json.Marshal(ev)
	if err != nil {
		return nil, err
	}

	// Flatten struct
	s1 := string(eventJSON[:len(eventJSON)-1])
	s2 := string(dataJSON[1:])
	flattenJSON := s1 + ", " + s2
	return []byte(flattenJSON), nil
}

// UnmarshalJSON will override the default json.Marshal in order to do the reverse process of MarshalJSON()
func (e *Message) UnmarshalJSON(b []byte) error {
	type duplicate Message
	ev := duplicate(*e)
	data := e.Data

	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	ev.Data = nil
	err = json.Unmarshal(b, &ev)
	if err != nil {
		return err
	}

	ev.Data = data
	*e = Message(ev)
	return nil
}
