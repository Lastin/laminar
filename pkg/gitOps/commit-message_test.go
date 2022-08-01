package gitOps

import (
	"github.com/digtux/laminar/pkg/shared"
	"testing"
)

func TestNicerMessage(t *testing.T) {
	regexTests := []struct {
		input  shared.ChangeRequest
		output string
	}{
		{shared.ChangeRequest{
			Image: "1122334455.dkr.ecr.eu-west-2.amazonaws.com/acmecorp/myimage",
			File:  "/some/filename.yaml",
			New:   "develop-123123",
		}, "filename: myimage:develop-123123"},
		{shared.ChangeRequest{
			Image: "1122334455.dkr.ecr.eu-west-2.amazonaws.com/acmecorp/my-image-name",
			File:  "/some/path/staging.yaml",
			New:   "feature-FOO-123123-added-feature-and-made-a-silly-long-branch-name-v1-v6.5.4-3-g0c8df55",
		}, "staging: my-image-name:feature-FOO-123123-added-fea...g0c8df55"},
	}
	for _, test := range regexTests {
		s := NicerMessage(test.input)
		if s != test.output {
			t.Errorf("TestRegex(%s), got: '%s' but expected: '%s'", test.input, s, test.output)
		}
	}
}
func TestTruncateForwardSlash(t *testing.T) {
	regexTests := []struct {
		input  string
		output string
	}{
		{"/home/path/to/some/filename.yaml", "filename"},
		{"/some/route/to/a/file.txt", "file"},
		{"/some/uri/path.html", "path"},
	}
	for _, test := range regexTests {
		s := truncateForwardSlash(test.input)
		if s != test.output {
			t.Errorf("TestRegex(%s), got: '%s' but expected: '%s'", test.input, s, test.output)
		}
	}
}
