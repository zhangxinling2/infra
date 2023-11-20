package httpClient

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHttpClient_Do(t *testing.T) {
	c := NewHttpClient(&Options{timeout: defaultHttpTimeout})
	Convey("http客户端", t, func() {
		req, err := http.NewRequest(http.MethodGet, "http://www.imooc.com", nil)
		So(err, ShouldBeNil)
		So(req, ShouldNotBeNil)
		res, err := c.Do(req)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		defer res.Body.Close()
		d, err := ioutil.ReadAll(res.Body)
		So(err, ShouldBeNil)
		So(d, ShouldNotBeNil)
		fmt.Println(string(d))
	})
}
