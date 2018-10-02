package gotten

type (
	Unmarshaler interface {
		Unmarshal(data []byte, v interface{}) error
	}

	UnmarshalFunc func(data []byte, v interface{}) error
)

func (fn UnmarshalFunc) Unmarshal(data []byte, v interface{}) error {
	return fn(data, v)
}
