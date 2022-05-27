package godao

import "testing"

func main(m testing.T) {
	var mapTest Map = Map{
		"field":  "some value",
		"field2": 1234,
		"field3": true,
	}

	raw := mapTest.Bytes()
	if raw == nil {
		m.FailNow()
		return
	}
	//convert from bytes
	result := mapTest.FromBytes(raw)
	if result.Value() == nil {
		m.Fail()
	}

}
