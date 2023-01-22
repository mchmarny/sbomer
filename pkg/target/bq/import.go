package bq

import (
	"context"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/mchmarny/sbomer/pkg/doc"
	"github.com/pkg/errors"
)

func Import(ctx context.Context, bom *doc.Document, target string) error {
	t, err := parseTarget(target)
	if err != nil {
		return errors.Wrap(err, "failed to parse target")
	}

	if err := configure(ctx, t); err != nil {
		return errors.Wrap(err, "errors checking target configuration")
	}

	// doc
	if err := insert(ctx, t, tableNameDoc, makeDocRows(bom)); err != nil {
		return errors.Wrapf(err, "failed to insert rows into %s", tableNameDoc)
	}

	// bom
	if err := insert(ctx, t, tableNameBom, makeBomRows(bom)); err != nil {
		return errors.Wrapf(err, "failed to insert rows into %s", tableNameBom)
	}

	// txt
	if err := insert(ctx, t, tableNameCtx, makeCtxRows(bom)); err != nil {
		return errors.Wrapf(err, "failed to insert rows into %s", tableNameCtx)
	}

	return nil
}

func insert(ctx context.Context, target *bqTarget, table string, items interface{}) error {
	if target == nil {
		return errors.New("target must be specified")
	}

	client, err := bigquery.NewClient(ctx, target.projectID)
	if err != nil {
		return errors.Wrap(err, "failed to create bigquery client")
	}
	defer client.Close()

	inserter := client.Dataset(target.datasetID).Table(table).Inserter()
	inserter.SkipInvalidRows = true
	if err := inserter.Put(ctx, items); err != nil {
		return errors.Wrap(err, "failed to insert rows")
	}

	return nil
}

type bqTarget struct {
	projectID string
	datasetID string
	location  string
}

const (
	bqProtocol        = "bq://"
	bqLocationDefault = "us"

	bqQueryParamLocation = "location"

	targetPartNumMin = 1
	targetPartNumMax = 2
)

// ParseImportRequest parses import request.
// e.g. bq://project.dataset?location=us
func parseTarget(target string) (*bqTarget, error) {
	target = strings.Replace(target, bqProtocol, "", 1)

	var query string
	if strings.Contains(target, "?") {
		queryParts := strings.Split(target, "?")
		target = queryParts[0]
		query = queryParts[1]
	}

	parts := strings.Split(target, ".")
	if len(parts) < targetPartNumMin || len(parts) > targetPartNumMax {
		return nil, errors.Errorf("invalid import target: %s", target)
	}

	t := &bqTarget{
		projectID: parts[0],
	}

	if len(parts) == targetPartNumMax {
		t.datasetID = parts[1]
	} else {
		t.datasetID = DatasetNameDefault
	}

	if query != "" {
		qp := strings.Split(query, "&")
		for _, p := range qp {
			if strings.Contains(p, bqQueryParamLocation) {
				t.location = parseValue(p)
			}
		}
	}

	if t.location == "" {
		t.location = bqLocationDefault
	}

	return t, nil
}

func parseValue(v string) string {
	parts := strings.Split(v, "=")
	if len(parts) == targetPartNumMax {
		return parts[1]
	}
	return ""
}
