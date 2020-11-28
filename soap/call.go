package soap

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

func soapCallHandleResponse(ws string, action string, payloadInterface interface{}, result interface{}) error {
	body, err := soapCall(ws, action, payloadInterface)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	return nil
}

func soapCall(ws string, action string, payloadInterface interface{}) ([]byte, error) {
	v := soapRQ{
		XMLNsSoap: "http://schemas.xmlsoap.org/soap/envelope/",
		XMLNsXSD:  "http://www.w3.org/2001/XMLSchema",
		XMLNsXSI:  "http://www.w3.org/2001/XMLSchema-instance",
		Body: soapBody{
			Payload: payloadInterface,
		},
	}
	payload, err := xml.MarshalIndent(v, "", "  ")

	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", ws, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/xml, multipart/related")
	req.Header.Set("SOAPAction", action)
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	_, err = httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return bodyBytes, nil
}
