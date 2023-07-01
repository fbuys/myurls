package acceptance

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/fbuys/myurls/internal/myurls"
	"github.com/fbuys/myurls/test/acceptance"
)

type test struct {
	path string
	url  string
}

func TestServer(t *testing.T) {
	cleanup, err := acceptance.LaunchTestProgram("../../../cmd/server/main.go")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cleanup)

	tests := []test{{path: "/unrecognized-path", url: "fbuys.dev/"}}
	for _, url := range myurls.GetAllUrls() {
		tests = append(tests, test{path: url.Id, url: trimHttp(url.Address)})
	}

	for _, tt := range tests {
		resp, err := http.Get(acceptance.Url + tt.path)
		if err != nil {
			t.Fatalf("The request failed and returned an error (%v).", err)
		}
		defer resp.Body.Close()
		gotUrl := resp.Request.URL.Host + resp.Request.URL.Path
		if tt.url != gotUrl {
			t.Errorf("Redirected to wrong URL, wanted: %v, got: %v", tt.url, gotUrl)
		}
	}
}

func trimHttp(url string) string {
	re := regexp.MustCompile(`https?://`)
	return re.ReplaceAllString(url, "")
}
