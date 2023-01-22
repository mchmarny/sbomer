package data

import (
	"database/sql"

	"github.com/mchmarny/sbomer/pkg/doc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	sqlSelectDoc = `SELECT
		id, 
		subject, 
		subject_version,
		format,
		format_version,
		provider,
		created
	FROM doc WHERE id = ?`

	sqlSelectBom = `SELECT
		id, 
		name,
		version
	FROM bom WHERE doc_id = ?`

	sqlSelectCtx = `SELECT
		ctx_type, 
		ctx_key,
		ctx_value
	FROM ctx WHERE bom_id = ?`
)

func Get(id string) (*doc.Document, error) {
	if db == nil {
		return nil, errors.New("database not initialized")
	}
	if id == "" {
		return nil, errors.New("id is empty")
	}

	stmtDoc, err := db.Prepare(sqlSelectDoc)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to prepare select doc statement")
	}

	stmtBom, err := db.Prepare(sqlSelectBom)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to prepare select bom statement")
	}

	stmtCtx, err := db.Prepare(sqlSelectCtx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to prepare select ctx statement")
	}

	var d doc.Document
	row := stmtDoc.QueryRow(id)
	if err = row.Scan(&d.ID, &d.Subject, &d.SubjectVersion, &d.Format,
		&d.FormatVersion, &d.Provider, &d.Created); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug().Err(err).Msgf("no doc with id: %s", id)
			return nil, nil
		}
		return nil, errors.Wrapf(err, "failed to scan row")
	}
	d.Items = make([]*doc.Item, 0)

	itemRows, err := stmtBom.Query(d.ID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrapf(err, "failed to execute select bom statement")
	}
	defer itemRows.Close() // nolint: errcheck

	for itemRows.Next() {
		var item doc.Item
		if err := itemRows.Scan(&item.ID, &item.Name, &item.Version); err != nil {
			return nil, errors.Wrapf(err, "failed to scan item row")
		}

		ctxRows, err := stmtCtx.Query(item.ID)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrapf(err, "failed to execute select ctx statement")
		}
		defer ctxRows.Close() // nolint: errcheck

		item.Contexts = make([]*doc.Context, 0)
		for ctxRows.Next() {
			var ctx doc.Context
			if err := ctxRows.Scan(&ctx.Type, &ctx.Key, &ctx.Value); err != nil {
				return nil, errors.Wrapf(err, "failed to scan ctx row")
			}
			item.Contexts = append(item.Contexts, &ctx)
		}

		d.Items = append(d.Items, &item)
	}

	return &d, nil
}
