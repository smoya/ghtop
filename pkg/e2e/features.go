package testp

import (
	"github.com/DATA-DOG/godog"
)

const featuresDir = "features"

type feature struct {
	name               string
	contextInitializer func(suite *godog.Suite)
}

var features = []feature{
	{
		"top",
		func(s *godog.Suite) { topFeatureContext(s) },
	},
}
