package v1

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/1Asi1/metric-track.git/internal/server/config"
	"github.com/1Asi1/metric-track.git/internal/server/repository/memory"
	"github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/1Asi1/metric-track.git/internal/server/transport/rest"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
)

func Example_updateMetric() {
	l := newLogger()
	st := memory.New(l, config.Config{})
	se := service.New(st, l)

	router := chi.NewRouter()
	h := rest.Handler{
		Mux:     router,
		Service: se}
	New(h, "")

	s := httptest.NewServer(router)
	defer s.Close()

	type args struct {
		metricType  string
		metricName  string
		metricValue string
	}

	reqModel := args{metricType: "gauge",
		metricName:  "Test",
		metricValue: "3.14",
	}

	url := fmt.Sprintf("%s/update/%s/%s/%s", s.URL, reqModel.metricType, reqModel.metricName, reqModel.metricValue)

	req := resty.New().R()
	req.Method = http.MethodPost
	req.URL = url

	res, err := req.Send()
	if err != nil {
		l.Err(err).Msg("req.Send()")
	}

	fmt.Println(res)

	// Output:
	// 3.14
}
