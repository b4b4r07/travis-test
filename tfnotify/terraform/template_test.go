package terraform

import (
	"testing"
)

func TestFmtTemplateExecute(t *testing.T) {
	testCases := []struct {
		name   string
		result ParseResult
		resp   string
		err    error
	}{
		{
			name: "",
			result: ParseResult{
				Result:   "There is diff in your .tf file (need to be formatted)",
				ExitCode: 0,
				Error:    nil,
			},
			resp: "\na\n\nb\n\n\n\nc\n",
			err:  nil,
		},
	}
	for _, testCase := range testCases {
		template := NewFmtTemplate(DefaultFmtTemplate)
		template.SetValue(CommonTemplate{
			Title:   "a",
			Message: "b",
			Result:  "c",
			Body:    "d",
		})
		resp, err := template.Execute(testCase.result)
		if err != nil {
			t.Fatal(err)
		}
		if resp != testCase.resp {
			t.Errorf("got %q but want %q", resp, testCase.resp)
		}
	}
}

func TestPlanTemplateExecute(t *testing.T) {
	testCases := []struct {
		result ParseResult
		resp   string
		err    error
	}{
		{
			result: ParseResult{
				Result:   "Plan: 1 to add, 0 to change, 0 to destroy.",
				ExitCode: 0,
				Error:    nil,
			},
			// resp: "\na\n\nb\n\n\n<pre><code>c\n</pre></code>\n\n\n<details><summary>Details (Click me)</summary>\n<pre><code>d\n</pre></code></details>\n",
			err: nil,
		},
	}
	for _, testCase := range testCases {
		template := NewPlanTemplate(DefaultPlanTemplate)
		// template.SetValue(CommonTemplate{
		// 	Title:   "",
		// 	Message: "",
		// 	Result:  "",
		// 	Body: "",
		// })
		resp, err := template.Execute(testCase.result)
		if err != nil {
			t.Fatal(err)
		}
		if resp != testCase.resp {
			t.Errorf("got %q but want %q", resp, testCase.resp)
		}
	}
}
