package observer

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Observer struct {
	trace.Tracer
	metric.Meter
	*slog.Logger
	sf shutdownFunc
}

func (obs *Observer) Http(handle func(pattern string, handler http.Handler)) func(pattern string, handler http.HandlerFunc) {

	return func(pattern string, handler http.HandlerFunc) {
		o := otelhttp.WithRouteTag(pattern, handler)
		handle(pattern, o)
	}

	// var root http.Handler
	// for _, r := range mux.Routes() {
	// 	root = otelhttp.WithRouteTag(r.Pattern, r.Handlers[r.Pattern])
	// 	mux.Handle(r.Pattern, root)
	// }

}

func (obs *Observer) RegisterHttp(mux *chi.Mux) func(method, pattern string, handler http.HandlerFunc) {
	return func(method, pattern string, handler http.HandlerFunc) {
		h := otelhttp.WithRouteTag(pattern, handler)
		mux.Method(method, pattern, h)
	}

}
func (obs *Observer) Shutdown(ctx context.Context) error {
	return obs.sf(ctx)
}

type shutdownFunc func(ctx context.Context) error

func New(ctx context.Context) (*Observer, error) {
	var sf []shutdownFunc
	var err error
	shutdownFuncs := func(inner context.Context) error {
		var err error
		for _, fn := range sf {
			err = errors.Join(err, fn(inner))
		}

		return err
	}

	//call shutdwon and make sure all errors are returned
	handleErr := func(in error) {
		err = errors.Join(in, shutdownFuncs(ctx))
	}

	// propagator
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// set up tracer provider
	tp, err := newTraceProvider()
	if err != nil {
		handleErr(err)
		return nil, err
	}
	sf = append(sf, tp.Shutdown)
	otel.SetTracerProvider(tp)

	// set ip meter provider
	mp, err := newMeterProvider()
	if err != nil {
		handleErr(err)
		return nil, err
	}
	sf = append(sf, mp.Shutdown)
	otel.SetMeterProvider(mp)

	// set up logger provider
	lp, err := newLogProvider()
	if err != nil {
		handleErr(err)
		return nil, err
	}

	sf = append(sf, lp.Shutdown)
	global.SetLoggerProvider(lp)

	const name = "sone"

	return &Observer{
		otel.Tracer(name),
		otel.Meter(name),
		otelslog.NewLogger(name),
		shutdownFuncs,
	}, nil

}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider() (*sdkTrace.TracerProvider, error) {
	te, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, err
	}

	tp := sdkTrace.NewTracerProvider(
		sdkTrace.WithBatcher(
			te,
			sdkTrace.WithBatchTimeout(time.Second),
		),
	)
	return tp, err
}

func newMeterProvider() (*sdkMetric.MeterProvider, error) {
	me, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	mp := sdkMetric.NewMeterProvider(
		sdkMetric.WithReader(
			sdkMetric.NewPeriodicReader(
				me,
				sdkMetric.WithInterval(3*time.Second),
			),
		),
	)
	return mp, nil
}

func newLogProvider() (*log.LoggerProvider, error) {
	le, err := stdoutlog.New()
	if err != nil {
		return nil, err
	}

	lp := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(le)),
	)
	return lp, nil
}
