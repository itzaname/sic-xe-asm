package machine

import "testing"

func TestInstructionByOpCode(t *testing.T) {
	// Check opcode is good
	ins, err := InstructionByOpCode(0x18)
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
	}

	if ins != &insList[0] {
		t.Error("Got wrong instruction:", ins, "needed:", &insList[0])
	}

	// Check twice for cache
	ins, err = InstructionByOpCode(0x18)
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
	}

	if ins != &insList[0] {
		t.Error("Got wrong instruction:", ins, "needed:", &insList[0])
	}

	// Check error working
	_, err = InstructionByOpCode(255)
	if err == nil {
		t.Error("Didn't get error on invalid opcode")
	}
}

func TestInstructionByName(t *testing.T) {
	// Check opcode is good
	ins, err := InstructionByName("ADD")
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
	}

	if ins != &insList[0] {
		t.Error("Got wrong instruction:", ins, "needed:", &insList[0])
	}

	// Check twice for cache
	ins, err = InstructionByName("ADD")
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
	}

	if ins != &insList[0] {
		t.Error("Got wrong instruction:", ins, "needed:", &insList[0])
	}

	// Check error working
	_, err = InstructionByName("INVALID-INSTRUCTION")
	if err == nil {
		t.Error("Didn't get error on invalid name")
	}
}
