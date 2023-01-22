package bq

import (
	"context"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
)

const (
	DatasetNameDefault = "sbomer"

	tableNameDoc = "doc"
	tableNameBom = "bom"
	tableNameCtx = "ctx"
)

var (
	tables = map[string]bigquery.Schema{
		tableNameDoc: docSchema,
		tableNameBom: bomSchema,
		tableNameCtx: ctxSchema,
	}
)

var (
	docSchema = bigquery.Schema{
		{Name: "id", Type: bigquery.StringFieldType, Required: true},
		{Name: "subject", Type: bigquery.StringFieldType, Required: true},
		{Name: "subject_version", Type: bigquery.StringFieldType, Required: true},
		{Name: "format", Type: bigquery.StringFieldType, Required: true},
		{Name: "format_version", Type: bigquery.StringFieldType, Required: true},
		{Name: "provider", Type: bigquery.StringFieldType, Required: true},
		{Name: "created", Type: bigquery.TimestampFieldType, Required: true},
	}

	bomSchema = bigquery.Schema{
		{Name: "id", Type: bigquery.StringFieldType, Required: true},
		{Name: "doc_id", Type: bigquery.StringFieldType, Required: true},
		{Name: "name", Type: bigquery.StringFieldType, Required: true},
		{Name: "version", Type: bigquery.StringFieldType, Required: true},
	}

	ctxSchema = bigquery.Schema{
		{Name: "bom_id", Type: bigquery.StringFieldType, Required: true},
		{Name: "ctx_type", Type: bigquery.StringFieldType, Required: true},
		{Name: "ctx_key", Type: bigquery.StringFieldType, Required: true},
		{Name: "ctx_value", Type: bigquery.StringFieldType, Required: true},
	}
)

func configure(ctx context.Context, target *bqTarget) error {
	if target == nil {
		return errors.New("target required")
	}

	log.Debug().
		Str("project", target.projectID).
		Str("dataset", target.datasetID).
		Msg("configuring target")

	exists, err := datasetExists(ctx, target)
	if err != nil {
		return errors.Wrap(err, "failed to check if dataset exists")
	}

	if !exists {
		if err := createDataset(ctx, target); err != nil {
			return errors.Wrap(err, "failed to create dataset")
		}
	}

	for t, s := range tables {
		exists, err = tableExists(ctx, target, t)
		if err != nil {
			return errors.Wrap(err, "failed to check if table exists")
		}

		if !exists {
			if err := createTable(ctx, target, t, s); err != nil {
				return errors.Wrapf(err, "failed to create table %s in: %s.%s", t, target.projectID, target.datasetID)
			}
		}
	}

	return nil
}

func createTable(ctx context.Context, target *bqTarget, table string, schema bigquery.Schema) error {
	client, err := bigquery.NewClient(ctx, target.projectID)
	if err != nil {
		return errors.Wrapf(err, "failed to create bigquery client for project %s", target.projectID)
	}
	defer client.Close()

	metaData := &bigquery.TableMetadata{
		Schema: schema,
	}

	tableRef := client.Dataset(target.datasetID).Table(table)
	if err := tableRef.Create(ctx, metaData); err != nil {
		return errors.Wrapf(err, "failed to create table %s", table)
	}
	return nil
}

func createDataset(ctx context.Context, target *bqTarget) error {
	client, err := bigquery.NewClient(ctx, target.projectID)
	if err != nil {
		return errors.Wrapf(err, "failed to create bigquery client for project %s", target.projectID)
	}
	defer client.Close()

	meta := &bigquery.DatasetMetadata{Location: target.location}
	if err := client.Dataset(target.datasetID).Create(ctx, meta); err != nil {
		return errors.Wrapf(err, "failed to create dataset %s", target.datasetID)
	}
	return nil
}

func datasetExists(ctx context.Context, target *bqTarget) (bool, error) {
	client, err := bigquery.NewClient(ctx, target.projectID)
	if err != nil {
		return false, errors.Wrapf(err, "failed to create bigquery client for project %s", target.projectID)
	}
	defer client.Close()

	it := client.Datasets(ctx)
	for {
		dataset, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return false, errors.Wrapf(err, "failed to list datasets for project %s", target.projectID)
		}

		if strings.EqualFold(dataset.DatasetID, target.datasetID) {
			return true, nil
		}
	}
	return false, nil
}

func tableExists(ctx context.Context, target *bqTarget, table string) (bool, error) {
	client, err := bigquery.NewClient(ctx, target.projectID)
	if err != nil {
		return false, errors.Wrapf(err, "failed to create bigquery client for project %s", target.projectID)
	}
	defer client.Close()

	it := client.Dataset(target.datasetID).Tables(ctx)
	for {
		t, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return false, errors.Wrapf(err, "failed to list datasets for project %s", target.projectID)
		}

		if strings.EqualFold(t.TableID, table) {
			return true, nil
		}
	}
	return false, nil
}
