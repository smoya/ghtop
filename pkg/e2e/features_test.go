package testp

import (
	"fmt"

	"os"
	"testing"

	"github.com/DATA-DOG/godog"
)

func TestMain(m *testing.M) {
	format := "progress"
	if isVerboseMode() {
		format = "pretty"
	}

	options := godog.Options{
		Format: format,
		Paths:  []string{featuresDir},
	}

	// Scenario should not share any state.
	// Running this way will ensure there is no state corruption or race conditions between scenarios.
	if format == "progress" {
		options.Concurrency = 4
	}

	var status int
	for _, suite := range features {
		featureFile := fmt.Sprintf("%s/%s.feature", featuresDir, suite.name)

		options.Paths = []string{featureFile}

		status = godog.RunWithOptions(suite.name, suite.contextInitializer, options)
		if status > 0 {
			break
		}
	}

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
func isVerboseMode() bool {
	for _, arg := range os.Args[1:] {
		// go test transforms -v option to this
		if arg == "-test.v=true" {
			return true
		}
	}

	return false
}
