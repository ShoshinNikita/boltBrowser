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
		Path         string
		pathAnswer   string
		Records      []record
	}{
		"Hello",
		1515,
		"Hello world Some <script>alert(5);</script>",
		"Hello world Some &lt;script&gt;alert(5);&lt;/script&gt;",
		[]record{
			{
				"bucket",
				"Some key", "Some key",
				"Some value", "Some value",
			},
			{
				"bucket",
				"Some <script>alert(5);</script>", "Some &lt;script&gt;alert(5);&lt;/script&gt;",
				"Some <body></body>", "Some &lt;body&gt;&lt;/body&gt;",
			},
			{
				"bucket",
				"Some key!\"#$%^:)", "Some key!&#34;#$%^:)",
				"value", "value",
			},
		},
	}

	if err := web.EscapeRecords(&tt); err != nil {
		t.Error(err)
		return
	}
	for _, r := range tt.Records {
		if r.Key != r.keyAnswer {
			t.Errorf("Bad key. Want: %s Got: %s", r.keyAnswer, r.Key)
		}
		if r.Value != r.valueAnswer {
			t.Errorf("Bad value. Want: %s Got: %s", r.valueAnswer, r.Value)
		}
	}
	if tt.Path != tt.pathAnswer {
		t.Errorf("Bad path. Want: %s Got: %s", tt.pathAnswer, tt.Path)
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
			{
				"bucket",
				"Some <script>alert(Hello, World!);</script>", "Some &lt;script&gt;alert(Hello, World!);&lt;/script&gt;",
				"Some <body></body>", "Some &lt;body&gt;&lt;/body&gt;",
			},
			{
				"record",
				"Some ()&^%@ke`~y!\"#$%^:)", "Some ()&amp;^%@ke`~y!&#34;#$%^:)",
				"VaLuE", "VaLuE",
			},
		},
	}

	if err := web.EscapeRecords(&tt); err != nil {
		t.Error(err)
		return
	}
	for _, r := range tt.Records {
		if r.Key != r.keyAnswer {
			t.Errorf("Bad key. Want: %s Got: %s", r.keyAnswer, r.Key)
		}
		if r.Value != r.valueAnswer {
			t.Errorf("Bad value. Want: %s Got: %s", r.valueAnswer, r.Value)
		}
	}
}
