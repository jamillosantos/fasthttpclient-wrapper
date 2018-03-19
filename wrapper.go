package fasthttpclient_wrapper

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
	"net/url"
)

var (
	headerUserAgent     = []byte("User-Agent")
	jsonContentTypeUtf8 = []byte("application/json; charset=utf-8")
)

type Client struct {
	baseURL        string
	baseURLBytes   []byte
	baseURLParsed  *url.URL
	userAgent      string
	userAgentBytes []byte
	Headers        map[string]string
	client         fasthttp.Client
}

func NewClient() *Client {
	return &Client{}
}

func (client *Client) BaseURL() string {
	return client.baseURL
}

func (client *Client) SetBaseURL(baseURL string) *Client {
	client.baseURL = baseURL
	client.baseURLBytes = []byte(baseURL)
	client.baseURLParsed, _ = url.Parse(baseURL)
	return client
}

func (client *Client) UserAgent() string {
	return client.userAgent
}

func (client *Client) SetUserAgent(userAgent string) *Client {
	client.userAgent = userAgent
	client.userAgentBytes = []byte(userAgent)
	return client
}

func (client *Client) Header(name string) (string, bool) {
	v, ok := client.Headers[name]
	return v, ok
}

func (client *Client) AddHeader(name, value string) *Client {
	client.Headers[name] = value
	return client
}

func (client *Client) DelHeader(name string) *Client {
	delete(client.Headers, name)
	return client
}

func (client *Client) setupRequest(req *fasthttp.Request) {
	req.Header.SetUserAgentBytes(client.userAgentBytes)
}

func (client *Client) setupRequestJSON(req *fasthttp.Request) {
	client.setupRequest(req)
	req.Header.SetContentTypeBytes(jsonContentTypeUtf8)
}

func (client *Client) RequestJson(method, path string, queryParams *url.Values, body interface{}, out interface{}, headers map[string]string) (int, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client.setupRequestJSON(req)

	req.Header.SetHost(client.baseURLParsed.Host)
	req.Header.SetMethod(method)

	reqURI := client.baseURL
	if len(reqURI) > 0 && reqURI[len(reqURI)-1] != '/' && path != "" && path[0] != '/' {
		reqURI += "/"
	}
	reqURI += path

	if queryParams != nil {
		if len(reqURI) > 0 && reqURI[len(reqURI)-1] != '?' {
			reqURI += "?" + queryParams.Encode()
		}
	}
	req.SetRequestURI(reqURI)

	if body != nil {
		bodyJson, err := json.Marshal(body)
		if err != nil {
			return -1, err
		}
		req.AppendBody(bodyJson)
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	err := client.client.Do(req, resp)
	if err != nil {
		return -1, err
	}

	err = json.Unmarshal(resp.Body(), out)
	if err != nil {
		return resp.StatusCode(), err
	}

	return resp.StatusCode(), nil
}

func (client *Client) GetJSON(path string, queryParams *url.Values, out interface{}, headers map[string]string) (int, error) {
	return client.RequestJson("GET", path, queryParams, nil, out, headers)
}

func (client *Client) PostJSON(path string, data interface{}, out interface{}, headers map[string]string) (int, error) {
	return client.RequestJson("POST", path, nil, data, out, headers)
}

func (client *Client) PutJSON(path string, data interface{}, out interface{}, headers map[string]string) (int, error) {
	return client.RequestJson("PUT", path, nil, data, out, headers)
}

func (client *Client) DeleteJSON(path string, data interface{}, out interface{}, headers map[string]string) (int, error) {
	return client.RequestJson("DELETE", path, nil, data, out, headers)
}
