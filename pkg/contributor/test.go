package contributor

import (
	"os"

	"encoding/json"

	"fmt"

	"github.com/pkg/errors"
)

var fixtures []Contributor

// Fixtures loads the fixtures located in testdata.
func Fixtures() []Contributor {
	if len(fixtures) > 0 {
		return fixtures
	}

	f, err := os.Open(
		fmt.Sprintf("%s/src/github.com/smoya/ghtop/pkg/contributor/testdata/fixtures.json", os.Getenv("GOPATH")),
	)
	if err != nil {
		panic(errors.Wrap(err, "Error loading test fixtures"))
	}

	defer f.Close()

	decoder := json.NewDecoder(f)

	var con []Contributor
	err = decoder.Decode(&con)
	if err != nil {
		panic(errors.Wrap(err, "Error decoding json test fixtures"))
	}

	fixtures = con

	return con
}
