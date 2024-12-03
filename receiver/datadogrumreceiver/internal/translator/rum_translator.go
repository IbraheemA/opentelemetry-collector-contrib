// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package translator // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/datadogreceiver/internal/translator"

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	semconv "go.opentelemetry.io/collector/semconv/v1.16.0"
)

func ToLogs(payload map[string]any, req *http.Request, reqBytes []byte) plog.Logs {
	results := plog.NewLogs()
	rl := results.ResourceLogs().AppendEmpty()
	rl.SetSchemaUrl(semconv.SchemaURL)
	ParseRUMRequestIntoResource(rl.Resource(), payload, req, reqBytes)

	in := rl.ScopeLogs().AppendEmpty()
	in.Scope().SetName("Datadog")

	newLogRecord := in.LogRecords().AppendEmpty()
	newLogRecord.Attributes().PutBool("should_tail_sample", rand.Intn(2) == 1)

	fmt.Println("%%%%% successful log parse!")
	return results
}

func ToTraces(payload map[string]any, req *http.Request, reqBytes []byte) ptrace.Traces {
	results := ptrace.NewTraces()
	rs := results.ResourceSpans().AppendEmpty()
	rs.SetSchemaUrl(semconv.SchemaURL)
	ParseRUMRequestIntoResource(rs.Resource(), payload, req, reqBytes)

	in := rs.ScopeSpans().AppendEmpty()
	in.Scope().SetName("Datadog")
	//var metadata map[string]any
	//var ok bool
	//if metadata, ok = payload["_dd"].(map[string]any); !ok {
	//	fmt.Println("failed to parse metadata")
	//	return results
	//}
	//traceID := uint64(metadata["trace_id"].(float64))
	//traceIDString := req.Header.Get("X-Datadog-Trace-Id")
	//traceIDString := payload["_dd"].(map[string]any)["trace_id"].(string)
	traceIDString, ok := payload["_dd"].(map[string]any)["trace_id"].(string)
	if !ok {
		fmt.Println("failed to retrieve traceID from payload")
		return results
	}
	traceID, err := strconv.Atoi(traceIDString)
	if err != nil {
		fmt.Println("failed to parse traceID")
		return results
	}
	//spanID := uint64(metadata["span_id"].(float64))
	//spanID := uint64(1)
	//spanIDString := req.Header.Get("X-Datadog-Parent-Id")
	spanIDString := payload["_dd"].(map[string]any)["span_id"].(string)
	spanID, err := strconv.Atoi(spanIDString)
	if err != nil {
		fmt.Println("failed to parse parent ID")
		return results
	}

	//traceIDString := metadata["trace_id"].(string)
	//spanIDString := metadata["span_id"].(string)
	//var spanID, traceID uint64
	//if t, err := strconv.Atoi(traceIDString); err != nil {
	//	fmt.Println("found trace_id:")
	//	fmt.Println(t)
	//	traceID = 0
	//	return results
	//} else {
	//	traceID = uint64(t)
	//}
	//
	//if t, err := strconv.Atoi(spanIDString); err != nil {
	//	fmt.Println("found span_id:")
	//	fmt.Println(t)
	//	spanID = 0
	//	return results
	//} else {
	//	spanID = uint64(t)
	//}
	fmt.Println("b1")
	newSpan := in.Spans().AppendEmpty()
	if rand.Intn(2) == 0 {
		newSpan.Status().SetCode(ptrace.StatusCodeError)
	} else {
		newSpan.Status().SetCode(ptrace.StatusCodeOk)
	}
	//newSpan.SetTraceID(uInt64ToTraceID(uint64(traceID), uint64(traceID)))
	newSpan.SetTraceID(uInt64ToTraceID(0, uint64(traceID)))
	newSpan.SetSpanID(uInt64ToSpanID(uint64(spanID)))
	newSpan.Attributes().PutBool("should_tail_sample", rand.Intn(2) == 1)

	fmt.Println("%%%%% successful trace parse!")
	return results
}

func ParseRUMRequestIntoResource(res pcommon.Resource, payload map[string]any, req *http.Request, reqBytes []byte) {
	formattedPayload, _ := json.MarshalIndent(payload, "", "\t")
	res.Attributes().PutStr("pretty_payload", string(formattedPayload))
	res.Attributes().PutStr(semconv.AttributeServiceName, "browser-rum-sdk")
	rand.Seed(time.Now().UnixNano())
	//fmt.Println("pretty_payload:")
	//fmt.Println(string(formattedPayload))

	emptyDumpBytes := res.Attributes().PutEmptyBytes("request_body_dump")
	emptyDumpBytes.FromRaw(reqBytes)

	requestHeadersMap := res.Attributes().PutEmptyMap("request_headers")
	for key, values := range req.Header {
		valueSlice := requestHeadersMap.PutEmptySlice(key)
		for _, value := range values {
			valueSlice.AppendEmpty().SetStr(value)
		}
	}

	requestQueryMap := res.Attributes().PutEmptyMap("request_query")
	for key, values := range req.URL.Query() {
		valueSlice := requestQueryMap.PutEmptySlice(key)
		for _, value := range values {
			valueSlice.AppendEmpty().SetStr(value)
		}
	}

	res.Attributes().PutStr("request_ddforward", req.URL.Query().Get("ddforward"))
}

func uInt64ToTraceID(high, low uint64) pcommon.TraceID {
	traceID := [16]byte{}
	binary.BigEndian.PutUint64(traceID[:8], high)
	binary.BigEndian.PutUint64(traceID[8:], low)
	return traceID
}

func uInt64ToSpanID(id uint64) pcommon.SpanID {
	spanID := [8]byte{}
	binary.BigEndian.PutUint64(spanID[:], id)
	return spanID
}
