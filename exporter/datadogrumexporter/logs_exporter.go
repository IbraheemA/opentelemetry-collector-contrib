// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package datadogrumexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/datadogrumexporter"

import (
	"bytes"
	"context"
	"fmt"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	//"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/datadogrumexporter/internal/clientutil"
	//"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/datadogrumexporter/internal/scrub"
)

type logsExporter struct {
	params exporter.Settings
	cfg    *Config
	ctx    context.Context // ctx triggers shutdown upon cancellation
}

func newLogsExporter(
	ctx context.Context,
	params exporter.Settings,
	cfg *Config,
) (*logsExporter, error) {
	exp := &logsExporter{
		params: params,
		cfg:    cfg,
		ctx:    ctx,
	}
	return exp, nil
}

var _ consumer.ConsumeLogsFunc = (*logsExporter)(nil).consumeLogs

func (exp *logsExporter) consumeLogs(
	ctx context.Context,
	td plog.Logs,
) (err error) {
	//defer func() { err = exp.scrubber.Scrub(err) }()
	//header := make(http.Header)
	//for i := 0; i < rspans.Len(); i++ {
	//	rspan := rspans.At(i)
	//	payloadDump, _ := rspan.Resource().Attributes().Get("payload_dump")
	//	//prettyPayload, _ := rspan.Resource().Attributes().Get("pretty_payload")
	//
	//	result, err := json.MarshalIndent(payloadDump.AsString(), "", "\t")
	//	if err != nil {
	//		return err
	//	}
	//
	//	exp.params.Logger.Debug("&&&&&&&&&& IN EXPORTER, marshalled json: ")
	//	exp.params.Logger.Debug(string(result))
	//}
	//return nil
	rlogs := td.ResourceLogs()
	fmt.Printf("&&&&&&&&&& RECEIVED LOGS: ")
	for i := range rlogs.Len() {
		rlog := rlogs.At(i)
		//s, _ := json.MarshalIndent(rspan.Resource().Attributes().AsRaw(), "", "\t")
		//s := rspan.ScopeSpans().At(0).Spans().At(0).Status().Code().String()
		//exp.params.Logger.Debug(string(s))
		//payloadDump, _ := rspan.Resource().Attributes().Get("request_body_dump")
		//prettyPayload, _ := rspan.Resource().Attributes().Get("pretty_payload")

		//result, err := json.MarshalIndent(payloadDump.AsString(), "", "\t")
		//if err != nil {
		//	return err
		//}
		//
		//exp.params.Logger.Debug("&&&&&&&&&& IN EXPORTER, marshalled json: ")
		//exp.params.Logger.Debug(string(result))
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		var rawRumData pcommon.Value

		rattr := rlog.Resource().Attributes()

		rawRumData, _ = rattr.Get("request_body_dump")

		//requestQuery, _ := rattr.Get("request_query")
		//
		//outUrl := &url.URL{
		//	Scheme: "https",
		//	Host:   "browser-intake-datadoghq.com",
		//	Path:   "/api/v2/rum",
		//}

		//outQuery := outUrl.Query()
		//
		//requestQuery.Map().Range(func(key string, v pcommon.Value) bool {
		//	exp.params.Logger.Debug(">setting query param:")
		//	exp.params.Logger.Debug(key)
		//	for i := range v.Slice().Len() {
		//		exp.params.Logger.Debug(v.Slice().At(i).AsString())
		//		outQuery.Add(key, v.Slice().At(i).AsString())
		//	}
		//	return true
		//})
		//
		//outUrl.RawQuery = outQuery.Encode()
		//
		//exp.params.Logger.Debug("&&&&&&&&&& SENDING REQUEST TO: ")
		//exp.params.Logger.Debug(outUrl.String())

		ddforward, _ := rattr.Get("request_ddforward")
		outUrlString := "https://browser-intake-datadoghq.com" +
			ddforward.AsString()

		req, err := http.NewRequest("POST", outUrlString, bytes.NewBuffer(rawRumData.Bytes().AsRaw()))

		headersMap, _ := rattr.Get("request_headers")
		headersMap.Map().Range(func(key string, v pcommon.Value) bool {
			exp.params.Logger.Debug("setting header:")
			exp.params.Logger.Debug(key)
			for i := range v.Slice().Len() {
				exp.params.Logger.Debug(v.Slice().At(i).AsString())
				req.Header.Set(key, v.Slice().At(i).AsString())
			}
			return true
		})

		req.Header.Set("Content-Type", "application/json")

		// XXXXXXX WORKING ON THIS NEXT
		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send request: %v", err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %v", err)
		}

		// Check the status code of the response
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			return fmt.Errorf("received non-OK response: status: %s, body: %s", resp.Status, body)
		}

		fmt.Println("Response:", string(body))
	}
	return nil
}
