package kubernetes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	Scheme     string
	Host       string
	HTTPClient *http.Client
}

func (c *Client) request(method string, path string, params url.Values, headers http.Header, body interface{}, v interface{}) error {
	if c.Scheme == "" {
		c.Scheme = "http"
	}

	if c.Host == "" {
		c.Host = "127.0.0.1:8001"
	}

	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}

	u := url.URL{
		Scheme:   c.Scheme,
		Host:     c.Host,
		Path:     path,
		RawQuery: params.Encode(),
	}

	var bodyReader io.ReadWriter = nil
	if body != nil {
		bodyReader = new(bytes.Buffer)
		encoder := json.NewEncoder(bodyReader)
		err := encoder.Encode(body)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return err
	}

	if headers != nil {
		for key, header := range headers {
			for _, value := range header {
				fmt.Printf("%s:%s\n", key, value)
				req.Header.Add(key, value)
			}
		}
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	return nil
}