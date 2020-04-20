package machine

import "testing"

func TestRegisterByName(t *testing.T) {
	reg, err := RegisterByName("A")
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
	}

	if reg != &regList[0] {
		t.Error("Got wrong register:", reg, "needed:", &regList[0])
	}

	_, err = RegisterByName("INVALID")
	if err == nil {
		t.Error("Didn't get error on invalid name")
	}
}

func TestRegisterByNumber(t *testing.T) {
	reg, err := RegisterByNumber(0)
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
	}

	if reg != &regList[0] {
		t.Error("Got wrong register:", reg, "needed:", &regList[0])
	}

	_, err = RegisterByNumber(255)
	if err == nil {
		t.Error("Didn't get error on invalid number")
	}
}
