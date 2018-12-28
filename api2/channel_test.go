package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	. "github.com/onsi/gomega"
)

var ttAddChannel = []struct {
	name string
	body *bytes.Buffer
	code int
}{
	{
		name: "fails with empty body",
		body: bytes.NewBufferString("{}"),
		code: http.StatusBadRequest,
	},
	{
		name: "passes with valid fields",
		body: bytes.NewBufferString(
			fmt.Sprintf(`{"id": "%s", "owner": "%s", "name": "%s"}`,
				UUIDRecal("id-1"),
				UUIDRecal("owner-1"),
				"ch-name",
			)),
		code: http.StatusOK,
	},
}

func TestAddChannel(t *testing.T) {
	for _, tt := range ttAddChannel {
		t.Run(tt.name, func(t *testing.T) {
			var (
				g       = NewGomegaWithT(t)
				channel = &Channel{}
				router  = mux.NewRouter()
				handler = channel.Register(router)
			)

			handler.Use(JSONMiddleWare)
			server := httptest.NewServer(handler)
			defer server.Close()

			url := fmt.Sprintf("%s/%s", server.URL, "channel/")

			req, err := http.NewRequest(http.MethodPut, url, tt.body)
			g.Expect(err).ShouldNot(HaveOccurred())

			rsp, err := http.DefaultClient.Do(req)
			g.Expect(err).ShouldNot(HaveOccurred())
			g.Expect(rsp.StatusCode).Should(Equal(tt.code))

		})
	}

}

var _uuids = map[string]string{}

func UUIDRecal(keys ...string) string {

	if len(keys) == 0 {
		return uuid.New().String()
	}

	key := keys[0]

	if value, ok := _uuids[key]; ok {
		return value
	}

	_uuids[key] = uuid.New().String()

	return _uuids[key]
}
