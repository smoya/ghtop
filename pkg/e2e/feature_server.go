package testp

import (
	"net/http"
	"net/http/httptest"

	"encoding/json"

	"fmt"

	"reflect"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/pkg/errors"
	"github.com/smoya/ghtop/pkg/contributor"
	"github.com/smoya/ghtop/pkg/server"
)

type serverFeature struct {
	server          *server.Server
	config          server.Config
	contributorRepo *contributorInMemoryRepository
	response        *httptest.ResponseRecorder
}

func (f *serverFeature) prepareScenario(s *godog.Suite) {
	s.Step(`^a server$`, f.aServer)
	s.Step(`^response is successful$`, f.responseIsSuccessful)
	s.Step(`^response is error (\d+)$`, f.responseIsError)
	s.Step(`^response is JSON$`, f.responseIsJSON)
	s.Step(`^I "([^"]*)" request to "([^"]*)"$`, f.iRequestTo)
	s.Step(`^response data should match JSON list:$`, f.responseDataShouldMatchJSONList)

	s.BeforeSuite(func() {
		f.config = server.NewConfig(8080, "e2e", "", "")
	})

	s.BeforeScenario(func(interface{}) {
		f.contributorRepo = &contributorInMemoryRepository{make([]contributor.Contributor, 0)}
	})
}

func (f *serverFeature) iRequestTo(method, url string) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	if f.server == nil {
		return errors.New("No server")
	}

	f.response = server.Request(f.server, req)

	return nil
}

func (f *serverFeature) responseIsSuccessful() error {
	if f.response == nil {
		return errors.New("No response")
	}

	if http.StatusOK != f.response.Code {
		return fmt.Errorf("response status is not successful, code: %d, msg: %s", f.response.Code, f.response.Body.String())
	}

	return nil
}

func (f *serverFeature) responseIsError(expectedErrorCode int) error {
	if f.response == nil {
		return errors.New("No response")
	}

	if f.response.Code != expectedErrorCode {
		return fmt.Errorf("response status should be error, expected: %v, received: %v", expectedErrorCode, f.response.Code)
	}

	return nil
}

func (f *serverFeature) responseIsJSON() error {
	if "application/json; charset=utf-8" != f.response.Header().Get("Content-type") {
		return errors.New("Content-type is not the expected")
	}

	return nil
}

func (f *serverFeature) responseDataShouldMatchJSONList(body *gherkin.DocString) (err error) {
	var currentJSON, expectedJSON []map[string]interface{}

	if err = json.Unmarshal(f.response.Body.Bytes(), &currentJSON); err != nil {
		return errors.Wrap(err, "Invalid JSON response")
	}

	if err = json.Unmarshal([]byte(body.Content), &expectedJSON); err != nil {
		return errors.Wrap(err, "Invalid JSON specification")
	}

	if len(currentJSON) != len(expectedJSON) {
		return fmt.Errorf("different elements amount on the JSON list. Expected: %v - Received: %v", len(expectedJSON), len(currentJSON))
	}

	for k, expectedElement := range expectedJSON {
		err = assertEqualJSON(expectedElement, currentJSON[k], "root")
		if err != nil {
			return
		}
	}

	return nil
}

func assertEqualJSON(expectedKV, actualKV map[string]interface{}, rootKey string) (err error) {
	if _, err = json.MarshalIndent(expectedKV, "", "  "); err != nil {
		return errors.Wrap(err, "Invalid JSON spec")
	}
	var actualDataJ []byte
	if actualDataJ, err = json.MarshalIndent(actualKV, "", "  "); err != nil {
		return errors.Wrap(err, "Invalid JSON response")
	}

	if len(expectedKV) != len(actualKV) {
		return errors.Errorf("JSON length mismatch at %s, expected: %d, got: %d\n%s",
			rootKey,
			len(expectedKV),
			len(actualKV),
			string(actualDataJ),
		)
	}

	for k, expectedVal := range expectedKV {

		actualVal, ok := actualKV[k]
		if !ok {
			return errors.Errorf("JSON lacks key: \"%s\" at %s \n%s", k, rootKey, string(actualDataJ))
		}

		expectedValKV, ok := expectedVal.(map[string]interface{})
		if ok {
			actualValKV, ok := actualVal.(map[string]interface{})
			if !ok {
				return errors.Errorf("JSON key: \"%s\" at %s has not KV, actual data:\n%s", k, rootKey, string(actualDataJ))
			}

			// recurse
			err = assertEqualJSON(expectedValKV, actualValKV, fmt.Sprintf("%s.%s", rootKey, k))
			if err != nil {
				return err
			}

		} else if !reflect.DeepEqual(actualVal, expectedVal) {
			return errors.Errorf("JSON val for \"%s\" at %s:\n%s\n\tDoes not match expected:\n%s", k, rootKey, actualVal, expectedVal)
		}
	}

	return nil
}
