package services

import "testing"

/*
Rule Based Controller Tests grid
Rows: Indoor Temperature
Columns: Outdoor Temperature

________|Colder |Cold   |Neutral|Hot    |Hotter |
|Colder |       |Open   |       |       |       |
|Cold   |Close  |       |Open   |Open   |       |
|Neutral|       |Close  |NIL    |Close  |       |
|Hot    |       |Open   |Open   |       |Open   |
|Hotter |       |       |       |Open   |       |

Test name convention: TestFunction_IndoorTemp_OutdoorTemp
*/
func TestRuleBasedControllerEval_Colder_Cold(t *testing.T) {
	// TODO: Add test logic here
	t.Log("tmp")
}

func TestRuleBasedControllerEval_Cold_Colder(t *testing.T) {
	// TODO: Add test logic here
	t.Log("tmp")
}

func TestRuleBasedControllerEval_Cold_Neutral(t *testing.T) {
	// TODO: Add test logic here
	t.Log("tmp")
}

func TestRuleBasedControllerEval_Cold_Hot(t *testing.T) {
	// TODO: Add test logic here
	t.Log("tmp")
}

func TestRuleBasedControllerEval_Neutral_Cold(t *testing.T) {
	// TODO: Add test logic here
	t.Log("tmp")
}

func TestRuleBasedControllerEval_Neutral_Neutral(t *testing.T) {
	// TODO: Add test logic here
	t.Log("tmp")
}

func TestRuleBasedControllerEval_Neutral_Hot(t *testing.T) {
	// TODO: Add test logic here
	t.Log("tmp")
}

func TestRuleBasedControllerEval_Hot_Cold(t *testing.T) {
	// TODO: Add test logic here
	t.Log("tmp")
}

func TestRuleBasedControllerEval_Hot_Neutral(t *testing.T) {
	// TODO: Add test logic here
	t.Log("tmp")
}

func TestRuleBasedControllerEval_Hot_Hotter(t *testing.T) {
	// TODO: Add test logic here
	t.Log("tmp")
}

func TestRuleBasedControllerEval_Hotter_Hot(t *testing.T) {
	// TODO: Add test logic here
	t.Log("tmp")
}
