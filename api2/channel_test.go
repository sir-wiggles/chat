package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/onsi/gomega"
)

func setupChannelTest() *httptest.Server {

	return httptest.NewServer(handler)
}

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
				UUIDRecal(0),
				UUIDRecal(1),
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
				ctrl    = gomock.NewController(t)
				mdb     = NewMockDatabaseController(ctrl)
			)
			channel.database = mdb

			handler.Use(JSONMiddleWare)
			server := httptest.NewServer(handler)
			defer server.Close()

			url := fmt.Sprintf("%s/%s", server.URL, "channel/")

			req, err := http.NewRequest(http.MethodPut, url, tt.body)
			g.Expect(err).ShouldNot(HaveOccurred())

			rsp, err := http.DefaultClient.Do(req)
			g.Expect(err).ShouldNot(HaveOccurred())
			g.Expect(rsp.StatusCode).Should(Equal(tt.code))

			ci := &ChannelInfo{}
			mdb.EXPECT().CreateChannel(ci)

		})
	}

}
