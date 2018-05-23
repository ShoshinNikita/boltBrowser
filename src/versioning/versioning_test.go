package versioning_test

import (
	"testing"

	"versioning"
)

func TestGetChanges(t *testing.T) {
	tests := []struct {
		text   string
		answer []string
	}{
		{"+ Hello, dear WORLD!!! \r\n * HOW ARe YOu?\r\n + It's fine",
			[]string{"Hello, dear WORLD!!! ", "HOW ARe YOu?", "It's fine"}},
		{"Changes:\r\n\r\n+ Added pages (issue #8)\r\n+ Code was refactored\r\n+ Updated README.md",
			[]string{"Added pages (issue #8)", "Code was refactored", "Updated README.md"}},
		{"+ Another test.", []string{"Another test."}},
		{"+ ?!\"#$%&'()*+,-./0123456789:;< = >?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_abcdefghijklmnopqrstuvwxyz{}|",
		[]string{"?!\"#$%&'()*+,-./0123456789:;< = >?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_abcdefghijklmnopqrstuvwxyz{}|"}},
	}

	for i, test := range tests {
		res := versioning.GetChanges(test.text)

		if len(res) != len(test.answer) {
			t.Errorf("Test #%d Wrong answer Want: %v Get: %v", i, test.answer, res)
			continue
		}

		for i := range res {
			if res[i] != test.answer[i] {
				t.Errorf("Test #%d Wrong answer Want: %v Get: %v", i, test.answer, res)
			}
		}
	}

}
