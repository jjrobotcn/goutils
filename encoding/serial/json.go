package serial

import "encoding/json"

type jsonImpl struct {
	v interface{}
}

func (j *jsonImpl) UnmarshalSerial(b []byte) error {
	return json.Unmarshal(b, j.v)
}

func (j *jsonImpl) MarshalSerial() ([]byte, error) {
	b, err := json.Marshal(j.v)
	if err != nil {
		return nil, err
	}
	return pack(ProtocolTypeJSON, b)
}
