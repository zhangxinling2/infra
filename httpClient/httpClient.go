package httpClient

import (
	"bytes"
	"errors"
	"github.com/zhangxinling2/infra/lb"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultHttpTimeout = 30 * time.Second
)

var parseUrl = url.Parse

type Options struct {
	timeout time.Duration
}
type HttpClient struct {
	client  *http.Client
	Options Options
	apps    *lb.Apps
}

func NewHttpClient(apps *lb.Apps, opt *Options) *HttpClient {
	c := &HttpClient{apps: apps}
	if opt == nil {
		c.Options = Options{timeout: defaultHttpTimeout}
	}
	c.Options = *opt
	c.client = &http.Client{Timeout: c.Options.timeout}
	return c
}
func (h *HttpClient) NewHttpRequest(method, url string, body io.Reader, headers http.Header) (*http.Request, error) {
	if method == "" {
		method = http.MethodGet
	}
	u, err := parseUrl(url)
	if err != nil {
		return nil, err
	}
	name := u.Host
	app := h.apps.Get(name)
	if app == nil {
		return nil, errors.New("没有可用的微服务应用，应用名称：" + name + ",请求：" + url)
	}
	ins := app.Get(url)
	if ins == nil {
		return nil, errors.New("没有可用的应用实例，应用名称：" + name + ",请求：" + url)
	}
	u.Host = ins.Address
	url = u.String()
	r, err := http.NewRequest(method, url, body)
	if len(headers) > 0 {
		for key, value := range headers {
			for _, val := range value {
				r.Header.Add(key, val)
			}
		}
	}
	return r, err
}
func (h *HttpClient) Do(request *http.Request) (*http.Response, error) {
	res, err := h.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return res, err
}
