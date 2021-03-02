package icon

// Icon struct
type Icon struct {
	Data []byte
	Name string
}

// Base bytes array icon representation
var Base = &Icon{
	Data: baseIcon,
	Name: "base",
}
