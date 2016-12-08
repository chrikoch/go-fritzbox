package auth

import (
	"testing"
)

func TestResponseCalculation(t *testing.T) {
	expected := "1234567z-9e224a41eeefa284df7bb0f26c2913e2"

	realResult := calculateResponse("1234567z", "äbc")

	if realResult != expected {
		t.Errorf("Result <%v> does not match expected result <%v>\n", realResult, expected)
	}

}

func TestRuneReplacement(t *testing.T) {

	expected := "123...456"
	input := "123日本語456"

	realResult := replaceInvalidChallengeRunes(input)
	if realResult != expected {
		t.Errorf("Result <%v> does not match expected result <%v>\n", realResult, expected)
	}
}
