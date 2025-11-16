package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/RuslanBoo/metrics/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Svc *service.MetricService
}

func New(svc *service.MetricService) *Handler {
	return &Handler{Svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/update/:type/:name/:value", h.Update)
	r.GET("/value/:type/:name", h.GetMetric)
	r.GET("/", h.GetAll)
}

func (h *Handler) Update(c *gin.Context) {
	metricType := c.Param("type")
	name := c.Param("name")
	value := c.Param("value")

	if err := h.Svc.UpdateMetric(metricType, name, value); err != nil {
		switch err {
		case service.ErrUnknownMetricType, service.ErrInvalidValue:
			c.String(http.StatusBadRequest, "invalid metric")
		default:
			c.String(http.StatusInternalServerError, "internal error")
		}
		return
	}

	c.String(http.StatusOK, "OK")
}

func (h *Handler) GetMetric(c *gin.Context) {
	metricType := c.Param("type")
	name := c.Param("name")

	val, ok := h.Svc.GetMetric(metricType, name)
	if !ok {
		c.String(http.StatusNotFound, "metric not found")
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("%v", val))
}

func (h *Handler) GetAll(c *gin.Context) {
	tmplStr := `
	<html>
	<head><title>Metrics</title></head>
	<body>
	<h1>All Metrics</h1>
	<ul>
	{{range $name, $val := .Gauges}}
		<li>Gauge {{$name}}: {{$val}}</li>
	{{end}}
	{{range $name, $val := .Counters}}
		<li>Counter {{$name}}: {{$val}}</li>
	{{end}}
	</ul>
	</body>
	</html>
	`
	tmpl := template.Must(template.New("metrics").Parse(tmplStr))

	gauges, counters := h.Svc.GetAllMetrics()
	data := struct {
		Gauges   map[string]interface{}
		Counters map[string]interface{}
	}{
		Gauges:   make(map[string]interface{}, len(gauges)),
		Counters: make(map[string]interface{}, len(counters)),
	}

	for k, v := range gauges {
		data.Gauges[k] = v
	}
	for k, v := range counters {
		data.Counters[k] = v
	}

	c.Writer.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
