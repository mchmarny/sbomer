CREATE TABLE IF NOT EXISTS doc (
    id TEXT NOT NULL,
    subject TEXT NOT NULL,
    subject_version TEXT NOT NULL,
    format TEXT NOT NULL,
    format_version TEXT NOT NULL,
    provider TEXT NOT NULL,
    created INTEGER NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS bom (
    id TEXT NOT NULL,
    doc_id TEXT NOT NULL,
    name TEXT NOT NULL,
    version TEXT NOT NULL,
    PRIMARY KEY (id, doc_id)
);

CREATE TABLE IF NOT EXISTS ctx (
    bom_id TEXT NOT NULL,
    ctx_type TEXT NOT NULL,
    ctx_key TEXT NOT NULL,
    ctx_value TEXT NOT NULL,
    PRIMARY KEY (bom_id, ctx_type, ctx_key, ctx_value)
);