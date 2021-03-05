package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("app_or_package_name")

func main() {
	ctx := context.Background()

	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",
	})
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	filter := otelrestful.OTelFilter("service-name")
	restful.DefaultContainer.Filter(filter)

	ws := &restful.WebService{}
	ws.Route(ws.GET("/profiles/{username}").To(userProfileHandler))
	restful.Add(ws)
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func userProfileHandler(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()

	username := req.PathParameter("username")
	name, err := selectUser(ctx, username)
	if err != nil {
		resp.WriteError(404, err)
		return
	}

	html := fmt.Sprintf(`<html><h1>Hello %s %s </h1></html>`+"\n", username, name)
	resp.Write([]byte(html))
}

func selectUser(ctx context.Context, username string) (string, error) {
	_, span := tracer.Start(ctx, "selectUser")
	defer span.End()

	span.SetAttributes(attribute.String("username", username))

	if username == "admin" {
		return "Joe", nil
	}

	return "", fmt.Errorf("username=%s not found", username)
}
