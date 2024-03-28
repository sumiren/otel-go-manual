package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func dbquery(c *gin.Context, db *sql.DB) []map[string]interface{} {
	//_, span := tracer.Start(c.Request.Context(), "dbquery", trace.WithAttributes(attribute.String("id", "hogehoge")))
	//defer span.End()

	tracer := otel.GetTracerProvider()
	_, span := tracer.Tracer("application").Start(c.Request.Context(), "custom dbquery method", trace.WithAttributes(attribute.String("id", "hogehoge")))
	defer span.End()

	rows, _ := db.QueryContext(c.Request.Context(), "SELECT id, name FROM book_sample")

	defer rows.Close()

	var (
		id   string
		name string
	)

	results := []map[string]interface{}{}
	for rows.Next() {
		rows.Scan(&id, &name)
		results = append(results, map[string]interface{}{"id": id, "name": name})
	}

	return results
}
