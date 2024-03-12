package web

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	testServer := httptest.NewServer(Handler())
	defer testServer.Close()
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/index.html", nil)
	assert.Nil(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	indexBody, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	t.Logf("%+v", string(indexBody))
}
