// Package terraform holds static policy checks over terraform/main.
//
// They run in CI with no AWS credentials and no terraform binary: the point is to catch
// an over-broad network rule in review rather than after it reaches an account.
package terraform

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	openCIDR = "0.0.0.0/0"
	tfDir    = "../../terraform/main"
)

// reviewedOpenCIDRLines are the only lines under terraform/main allowed to name the whole
// internet, each with the reason it is safe. Anything else fails until a human adds it
// here. Commented-out lines count: this template is meant to be uncommented as-is, so an
// allow-all rule parked in a comment becomes a real one the moment somebody uses it.
var reviewedOpenCIDRLines = map[string]string{
	`#     cidr_blocks = ["0.0.0.0/0"]`: "Egress from the optional VPC security group example, " +
		"limited to TCP 443. Guarded by TestLambdaEgressExampleIsHTTPSOnly.",
}

func TestOpenCIDRsAppearOnlyAtReviewedSites(t *testing.T) {
	for _, file := range terraformFiles(t) {
		for i, line := range readLines(t, file) {
			if !strings.Contains(line, openCIDR) {
				continue
			}
			trimmed := strings.TrimRight(line, " \t")
			if _, ok := reviewedOpenCIDRLines[trimmed]; !ok {
				t.Errorf("%s:%d unreviewed use of %s\n\t%s\n"+
					"Restrict it, or add the line to reviewedOpenCIDRLines with a reason.",
					filepath.Base(file), i+1, openCIDR, trimmed)
			}
		}
	}
}

// TestLambdaEgressExampleIsHTTPSOnly pins the one reviewed exception: the security group
// example may reach the internet, but only on 443.
func TestLambdaEgressExampleIsHTTPSOnly(t *testing.T) {
	lines := readLines(t, filepath.Join(tfDir, "lambda.tf"))

	found := false
	for i, line := range lines {
		if !strings.Contains(line, openCIDR) {
			continue
		}
		found = true
		block := enclosingEgressBlock(lines, i)
		for _, want := range []string{`protocol    = "tcp"`, "from_port   = 443", "to_port     = 443"} {
			if !blockHas(block, want) {
				t.Errorf("lambda.tf:%d egress to %s must be HTTPS only; missing %q in:\n%s",
					i+1, openCIDR, want, strings.Join(block, "\n"))
			}
		}
	}
	if !found {
		t.Fatalf("expected lambda.tf to contain the reviewed HTTPS egress example")
	}
}

// TestNoSecurityGroupAllowsWorldIngress catches an inbound rule open to the internet,
// commented or not. Nothing in this template should ever accept unsolicited traffic.
func TestNoSecurityGroupAllowsWorldIngress(t *testing.T) {
	for _, file := range terraformFiles(t) {
		lines := readLines(t, file)
		inIngress := false
		for i, line := range lines {
			switch {
			case strings.Contains(line, "ingress"):
				inIngress = true
			case strings.Contains(uncomment(line), "}"):
				inIngress = false
			}
			if inIngress && strings.Contains(line, openCIDR) {
				t.Errorf("%s:%d ingress open to %s: %s", filepath.Base(file), i+1, openCIDR, strings.TrimSpace(line))
			}
		}
	}
}

func terraformFiles(t *testing.T) []string {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(tfDir, "*.tf"))
	if err != nil {
		t.Fatalf("globbing %s: %v", tfDir, err)
	}
	if len(files) == 0 {
		t.Fatalf("no .tf files found under %s", tfDir)
	}
	return files
}

func readLines(t *testing.T, path string) []string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return strings.Split(string(data), "\n")
}

// enclosingEgressBlock returns the lines of the egress block containing index, commented
// or not, by walking out to the nearest "egress {" and its closing brace.
func enclosingEgressBlock(lines []string, index int) []string {
	start := index
	for start > 0 && !strings.Contains(lines[start], "egress") {
		start--
	}
	end := index
	for end < len(lines)-1 && !strings.Contains(uncomment(lines[end]), "}") {
		end++
	}
	return lines[start : end+1]
}

func blockHas(block []string, want string) bool {
	for _, line := range block {
		if strings.Contains(line, want) {
			return true
		}
	}
	return false
}

func uncomment(line string) string {
	return strings.TrimPrefix(strings.TrimSpace(line), "#")
}
