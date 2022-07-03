package json

import "io"

type Decodable interface {
	FromJSON()
}

type Encodable interface {
	WriteJSON(w io.Writer)
}

type Codable interface {
	Decodable
	Encodable
}
