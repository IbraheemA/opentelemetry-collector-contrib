// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package datadogrumexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/datadogrumexporter"

import (
	"context"
	"encoding/json"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/featuregate"
	"go.opentelemetry.io/collector/pdata/ptrace"
	//"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/datadogrumexporter/internal/clientutil"
	//"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/datadogrumexporter/internal/scrub"
)

var traceCustomHTTPFeatureGate = featuregate.GlobalRegistry().MustRegister(
	"exporter.datadogrumexporter.TraceExportUseCustomHTTPClient",
	featuregate.StageAlpha,
	featuregate.WithRegisterDescription("When enabled, trace export uses the HTTP client from the exporter HTTP configs"),
	featuregate.WithRegisterFromVersion("v0.105.0"),
)

type traceExporter struct {
	params exporter.Settings
	cfg    *Config
	ctx    context.Context // ctx triggers shutdown upon cancellation
}

func newTracesExporter(
	ctx context.Context,
	params exporter.Settings,
	cfg *Config,
) (*traceExporter, error) {
	exp := &traceExporter{
		params: params,
		cfg:    cfg,
		ctx:    ctx,
	}
	return exp, nil
}

var _ consumer.ConsumeTracesFunc = (*traceExporter)(nil).consumeTraces

func (exp *traceExporter) consumeTraces(
	ctx context.Context,
	td ptrace.Traces,
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
	rspans := td.ResourceSpans()
	exp.params.Logger.Debug("&&&&&&&&&& RECEIVED SPAN: ")
	for i := range rspans.Len() {
		rspan := rspans.At(i)
		//s, _ := json.MarshalIndent(rspan.Resource().Attributes().AsRaw(), "", "\t")
		s := rspan.ScopeSpans().At(0).Spans().At(0).Status().Code().String()
		exp.params.Logger.Debug(string(s))
		payloadDump, _ := rspan.Resource().Attributes().Get("payload_dump")
		//prettyPayload, _ := rspan.Resource().Attributes().Get("pretty_payload")

		result, err := json.MarshalIndent(payloadDump.AsString(), "", "\t")
		if err != nil {
			return err
		}

		exp.params.Logger.Debug("&&&&&&&&&& IN EXPORTER, marshalled json: ")
		exp.params.Logger.Debug(string(result))
	}
	//client := &http.Client{
	//	Timeout: 10 * time.Second,
	//}
	//
	//var rawRumData pcommon.Value
	//
	//rawRumData, _ = td.ResourceSpans().At(0).Resource().Attributes().Get("payload_dump")
	//
	//req, err := http.NewRequest("POST", "http://rum.browser-intake-datadoghq.com/api/v2/rum", bytes.NewBuffer(rawRumData.Bytes().AsRaw()))
	//
	//req.Header.Set("Content-Type", "application/json")
	//
	//// XXXXXXX WORKING ON THIS NEXT
	//// Send the request
	//resp, err := client.Do(req)
	//if err != nil {
	//	return fmt.Errorf("failed to send request: %v", err)
	//}
	//defer func(Body io.ReadCloser) {
	//	err := Body.Close()
	//	if err != nil {
	//
	//	}
	//}(resp.Body)
	//
	//// Read the response body
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return fmt.Errorf("failed to read response: %v", err)
	//}
	//
	//// Check the status code of the response
	//if resp.StatusCode != http.StatusOK {
	//	return fmt.Errorf("received non-OK response: %s", body)
	//}

	//fmt.Println("Response:", string(body))
	return nil
}
