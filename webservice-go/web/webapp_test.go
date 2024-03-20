package web

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	_, err := os.Stat(frontendAppDir)
	if os.IsNotExist(err) {
		t.Skip("The frontend application has not been generated in ./app/dist, skipping test ...")
	}
	testServer := httptest.NewServer(Handler())
	defer testServer.Close()
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/index.html", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	indexBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	t.Logf("%+v", string(indexBody))
}
