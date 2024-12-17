// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package datadogrumreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/datadogrumreceiver"

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/datadogrumreceiver/internal/translator"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
)

type datadogRUMReceiver struct {
	address string
	config  *Config
	params  receiver.Settings

	nextTracesConsumer consumer.Traces
	nextLogsConsumer   consumer.Logs

	server    *http.Server
	lReceiver *receiverhelper.ObsReport
}

func newDataDogRUMReceiver(config *Config, params receiver.Settings) (component.Component, error) {
	instance, err := receiverhelper.NewObsReport(receiverhelper.ObsReportSettings{LongLivedCtx: false, ReceiverID: params.ID, Transport: "http", ReceiverCreateSettings: params})
	if err != nil {
		return nil, err
	}

	return &datadogRUMReceiver{
		params: params,
		config: config,
		server: &http.Server{
			ReadTimeout: config.ReadTimeout,
		},
		lReceiver: instance,
	}, nil
}

func (ddr *datadogRUMReceiver) Start(ctx context.Context, host component.Host) error {
	ddmux := http.NewServeMux()

	ddmux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	if ddr.nextTracesConsumer != nil || ddr.nextLogsConsumer != nil {
		ddmux.HandleFunc("/api/v2/rum", ddr.handleEvent)
	}

	var err error
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://localhost:*", "http://localhost:*"}, // Specify allowed origins
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},                    // Specify allowed methods
		AllowedHeaders:   []string{"Content-Type", "Authorization"},             // Specify allowed headers
		AllowCredentials: true,                                                  // Allow credentials
	}).Handler(ddmux)

	ddr.server, err = ddr.config.ServerConfig.ToServer(
		ctx,
		host,
		ddr.params.TelemetrySettings,
		corsHandler,
	)
	if err != nil {
		return fmt.Errorf("failed to create server definition: %w", err)
	}
	hln, err := ddr.config.ServerConfig.ToListener(ctx)
	if err != nil {
		return fmt.Errorf("failed to create datadog listener: %w", err)
	}

	ddr.address = hln.Addr().String()

	go func() {
		if err := ddr.server.Serve(hln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			//componentstatus.ReportStatus(host, componentstatus.NewFatalErrorEvent(fmt.Errorf("error starting datadog receiver: %w", err)))
		}
	}()
	return nil
}

var bufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func GetBuffer() *bytes.Buffer {
	buffer := bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	return buffer
}

func PutBuffer(buffer *bytes.Buffer) {
	bufferPool.Put(buffer)
}

type RUMPayload struct {
	Type string
}

func (ddr *datadogRUMReceiver) handleEvent(w http.ResponseWriter, req *http.Request) {
	obsCtx := ddr.lReceiver.StartTracesOp(req.Context())
	var err error
	var eventCount int
	defer func(eventCount *int) {
		ddr.lReceiver.EndTracesOp(obsCtx, "datadog", *eventCount, err)
	}(&eventCount)

	// XXX MOVE THIS TRANSLATION LOGIC TO ITS OWN "Handle" FUNCTION
	defer func() {
		_, errs := io.Copy(io.Discard, req.Body)
		err = errors.Join(err, errs, req.Body.Close())
	}()

	buf := GetBuffer()
	defer PutBuffer(buf)
	_, err = io.Copy(buf, req.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshal reqs", http.StatusBadRequest)
		ddr.params.Logger.Error("Unable to unmarshal reqs", zap.Error(err))
		return
	}
	reqBytes := buf.Bytes()

	//printBuf := GetBuffer()
	//defer PutBuffer(buf)
	//io.Copy(printBuf, req.MultipartReader())
	//printBytes := printBuf.Bytes()

	//traceID := req.Header.Get("X-Datadog-Trace-Id")
	//spanID := req.Header.Get("X-Datadog-Span-Id")

	// check errors
	fmt.Printf("&&&&&&&&&& RECEIVED REQUEST BODY: " + fmt.Sprintf("%v", buf.String()))
	//ddr.params.Logger.Debug("&&&&&&&&&& RECEIVED TraceID: " + fmt.Sprintf("%v", traceID))
	//ddr.params.Logger.Debug("&&&&&&&&&& RECEIVED SpanID: " + fmt.Sprintf("%v", spanID))

	//postfixBytes := []byte(`{"kind": "receiver", "name": "datadogrum", "data_type": "logs"}`) // Convert the postfix string to a byte slice
	//if bytes.HasSuffix(reqBytes, postfixBytes) {
	//	reqBytes = reqBytes[:len(reqBytes)-len(postfixBytes)]
	//} else {
	//	http.Error(w, "Unexpected RUM event body format - Unable to unmarshal reqs", http.StatusBadRequest)
	//	ddr.params.Logger.Error("Unable to unmarshal reqs", zap.Error(err))
	//}

	var jsonEvents []map[string]any
	decoder := json.NewDecoder(buf)
	for {
		var event map[string]any
		if err := decoder.Decode(&event); err != nil {
			if err.Error() == "EOF" {
				break
			}
			http.Error(w, "Unable to unmarshal reqs", http.StatusBadRequest)
			ddr.params.Logger.Error("Unable to unmarshal reqs", zap.Error(err))
			return
		}
		jsonEvents = append(jsonEvents, event)
	}

	//ddr.params.Logger.Debug("&&&&&&&&&& PARSED TO: ")
	//for _, event := range jsonEvents {
	//	s, _ := json.MarshalIndent(event, "", "\t")
	//	ddr.params.Logger.Debug(string(s))
	//}

	for _, event := range jsonEvents {
		_, ok := event["_dd"].(map[string]any)["trace_id"].(string)
		if !ok {
			fmt.Println("failed to retrieve traceID from RUM event payload; treating as log instead")
			otelLogs := translator.ToLogs(event, req, reqBytes)
			if ddr.nextLogsConsumer != nil {
				err = ddr.nextLogsConsumer.ConsumeLogs(obsCtx, otelLogs)
			}
		} else {
			otelTraces := translator.ToTraces(event, req, reqBytes)
			if ddr.nextTracesConsumer != nil {
				err = ddr.nextTracesConsumer.ConsumeTraces(obsCtx, otelTraces)
			}
		}
		if err != nil {
			http.Error(w, "Log consumer errored out", http.StatusInternalServerError)
			ddr.params.Logger.Error("Log consumer errored out", zap.Error(err))
			return
		}
	}

	//obsCtx := ddr.lReceiver.StartTracesOp(req.Context())
	//var err error
	//var spanCount int
	//defer func(spanCount *int) {
	//	ddr.lReceiver.EndTracesOp(obsCtx, "datadog", *spanCount, err)
	//}(&spanCount)
	//
	//var ddTraces []*pb.TracerPayload
	//ddTraces, err = translator.HandleTracesPayload(req)
	//if err != nil {
	//	http.Error(w, "Unable to unmarshal reqs", http.StatusBadRequest)
	//	ddr.params.Logger.Error("Unable to unmarshal reqs", zap.Error(err))
	//	return
	//}
	//for _, ddTrace := range ddTraces {
	//	otelTraces := translator.ToTraces(ddTrace, req)
	//	spanCount = otelTraces.SpanCount()
	//	err = ddr.nextTracesConsumer.ConsumeTraces(obsCtx, otelTraces)
	//	if err != nil {
	//		http.Error(w, "Trace consumer errored out", http.StatusInternalServerError)
	//		ddr.params.Logger.Error("Trace consumer errored out", zap.Error(err))
	//		return
	//	}
	//}
	//
	_, _ = w.Write([]byte("OK"))
}

func (ddr *datadogRUMReceiver) Shutdown(ctx context.Context) (err error) {
	return ddr.server.Shutdown(ctx)
}
