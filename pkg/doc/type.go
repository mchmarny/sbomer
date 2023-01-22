package doc

const (
	ContextKindUndefined        ContextKind = ""
	ContextKindLicense          ContextKind = "license"
	ContextKindCPE              ContextKind = "cpe"
	ContextKindSWID             ContextKind = "swid"
	ContextKindPURL             ContextKind = "purl"
	ContextKindType             ContextKind = "type"
	ContextKindSupplier         ContextKind = "supplier"
	ContextKindAuthor           ContextKind = "author"
	ContextKindPublisher        ContextKind = "publisher"
	ContextKindGroup            ContextKind = "group"
	ContextKindDescription      ContextKind = "description"
	ContextKindScope            ContextKind = "scope"
	ContextKindDownloadLocation ContextKind = "download-location"
	ContextKindHomepage         ContextKind = "homepage"
	ContextKindSourceInfo       ContextKind = "source-info"
	ContextKindSummary          ContextKind = "summary"
	ContextKindPurpose          ContextKind = "purpose"
	ContextKindReleaseDate      ContextKind = "release-date"
	ContextKindBuildDate        ContextKind = "built-date"
	ContextKindVerificationCode ContextKind = "verification-code"
)

type ContextKind string
