package data

import (
	"github.com/mchmarny/sbomer/pkg/doc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	sqlInsertDoc = `INSERT INTO doc (
			id, 
			subject, 
			subject_version,
			format,
			format_version,
			provider,
			created
		) VALUES (?, ?, ?, ?, ?, ?, ?)`

	sqlInsertBom = `INSERT INTO bom (
			id, 
			doc_id, 
			name,
			version
		) VALUES (?, ?, ?, ?)`

	sqlInsertCtx = `INSERT INTO ctx (
			bom_id, 
			ctx_type, 
			ctx_key,
			ctx_value
		) VALUES (?, ?, ?, ?)`
)

func Save(doc *doc.Document) error {
	if db == nil {
		return errors.New("database not initialized")
	}
	if doc == nil {
		return errors.New("document is nil")
	}

	stmtDoc, err := db.Prepare(sqlInsertDoc)
	if err != nil {
		return errors.Wrapf(err, "failed to prepare batch doc statement")
	}

	stmtBom, err := db.Prepare(sqlInsertBom)
	if err != nil {
		return errors.Wrapf(err, "failed to prepare batch bom statement")
	}

	stmtTxt, err := db.Prepare(sqlInsertCtx)
	if err != nil {
		return errors.Wrapf(err, "failed to prepare batch ctx statement")
	}

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrapf(err, "failed to begin transaction")
	}

	if _, err = tx.Stmt(stmtDoc).
		Exec(doc.ID,
			doc.Subject,
			doc.SubjectVersion,
			doc.Format,
			doc.FormatVersion,
			doc.Provider,
			doc.Created); err != nil {
		log.Error().Msgf("failed to execute batch doc statement: %v", err)
		if err = tx.Rollback(); err != nil {
			return errors.Wrapf(err, "failed to rollback transaction")
		}
		return errors.Wrapf(err, "failed to execute batch doc statement")
	}

	for _, item := range doc.Items {
		if _, err = tx.Stmt(stmtBom).
			Exec(item.ID, doc.ID, item.Name, item.Version); err != nil {
			log.Error().Msgf("failed to execute batch item statement: %v", err)
			if err = tx.Rollback(); err != nil {
				return errors.Wrapf(err, "failed to rollback transaction")
			}
			return errors.Wrapf(err, "failed to execute batch statement")
		}

		for _, ctx := range item.Contexts {
			if _, err = tx.Stmt(stmtTxt).
				Exec(item.ID, ctx.Type, ctx.Key, ctx.Value); err != nil {
				log.Error().Msgf("failed to execute batch context statement: %v", err)
				if err = tx.Rollback(); err != nil {
					return errors.Wrapf(err, "failed to rollback transaction")
				}
				return errors.Wrapf(err, "failed to execute batch statement")
			}
		}
	}

	if err = tx.Commit(); err != nil {
		log.Error().Msgf("failed to commit transaction: %v", err)
		return errors.Wrapf(err, "failed to commit transaction")
	}

	log.Debug().Msgf("saved doc %s", doc.ID)

	return nil
}
