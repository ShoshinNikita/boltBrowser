package web_test

import (
	"testing"

	"github.com/ShoshinNikita/boltBrowser/internal/web"
)

type record struct {
	T           string
	Key         string
	keyAnswer   string
	Value       string
	valueAnswer string
}

func TestEscape(t *testing.T) {
	test1(t)
	test2(t)
}

func test1(t *testing.T) {
	tt := struct {
		SomeField    string
		AnotherField int
		Records      []record
	}{
		"Hello",
		1515,
		[]record{
			record{"bucket", "Some key", "Some key", "Some value", "Some value"},
			record{"bucket",
				"Some <script>alert(5);</script>", "Some &lt;script&gt;alert(5);&lt;/script&gt;",
				"Some <body></body>", "Some &lt;body&gt;&lt;/body&gt;"},
			record{"bucket", "Some key!\"#$%^:)", "Some key!&#34;#$%^:)",
				"value", "value"},
		},
	}

	err := web.EscapeRecords(tt)
	if err != nil {
		t.Error(err)
		return
	}
	for _, r := range tt.Records {
		if r.Key != r.keyAnswer {
			t.Errorf("Bad key. Want: %s Got: %s", r.keyAnswer, r.Key)
		}
		if r.Value != r.valueAnswer {
			t.Errorf("Bad key. Want: %s Got: %s", r.valueAnswer, r.Value)
		}
	}
}

func test2(t *testing.T) {
	tt := struct {
		One     string
		two     uint
		three   []int64
		Records []record
	}{
		"",
		48,
		[]int64{56, 28},
		[]record{
			record{"bucket",
				"Some <script>alert(Hello, World!);</script>", "Some &lt;script&gt;alert(Hello, World!);&lt;/script&gt;",
				"Some <body></body>", "Some &lt;body&gt;&lt;/body&gt;"},
			record{"record",
				"Some ()&^%@ke`~y!\"#$%^:)", "Some ()&amp;^%@ke`~y!&#34;#$%^:)",
				"VaLuE", "VaLuE"},
		},
	}

	err := web.EscapeRecords(tt)
	if err != nil {
		t.Error(err)
	}
	for _, r := range tt.Records {
		if r.Key != r.keyAnswer {
			t.Errorf("Bad key. Want: %s Got: %s", r.keyAnswer, r.Key)
		}
		if r.Value != r.valueAnswer {
			t.Errorf("Bad key. Want: %s Got: %s", r.valueAnswer, r.Value)
		}
	}
}
