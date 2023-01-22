package bq

import (
	"cloud.google.com/go/bigquery"
	"github.com/mchmarny/sbomer/pkg/doc"
)

func makeDocRows(in *doc.Document) []*DocumentRow {
	return []*DocumentRow{{Document: in}}
}

type DocumentRow struct {
	*doc.Document
}

func (r *DocumentRow) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"id":              r.ID,
		"subject":         r.Subject,
		"subject_version": r.SubjectVersion,
		"format":          r.Format,
		"format_version":  r.FormatVersion,
		"provider":        r.Provider,
		"created":         r.Created,
	}, "", nil
}

func makeBomRows(in *doc.Document) []*BomRow {
	list := make([]*BomRow, 0)
	for _, r := range in.Items {
		list = append(list, &BomRow{DicID: in.ID, Item: r})
	}
	return list
}

type BomRow struct {
	DicID string
	*doc.Item
}

func (r *BomRow) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"id":      r.ID,
		"doc_id":  r.DicID,
		"name":    r.Name,
		"version": r.Version,
	}, "", nil
}

func makeCtxRows(in *doc.Document) []*CtxRow {
	list := make([]*CtxRow, 0)
	for _, r := range in.Items {
		for _, c := range r.Contexts {
			list = append(list, &CtxRow{BomID: r.ID, Context: c})
		}
	}
	return list
}

type CtxRow struct {
	BomID string
	*doc.Context
}

func (r *CtxRow) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"bom_id":    r.BomID,
		"ctx_type":  r.Type,
		"ctx_key":   r.Key,
		"ctx_value": r.Value,
	}, "", nil
}
