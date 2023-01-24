package cdx

import (
	"bytes"

	cdx "github.com/CycloneDX/cyclonedx-go"
	"github.com/mchmarny/sbomer/pkg/sbom"
	"github.com/rs/zerolog/log"
)

type CycloneDXReport struct {
	valid    bool
	docError error
	name     string

	creationToolName    int
	creationToolVersion int

	totalPackages  int
	hasLicense     int
	hasPackVersion int
	hasPackDigest  int
	hasPurl        int
	hasCPE         int
}

func (r *CycloneDXReport) GetName() string {
	return r.name
}

func (r *CycloneDXReport) IsSpecCompliant() *sbom.ReportValue {
	if r.docError != nil {
		log.Error().Err(r.docError).Msg("Invalid CycloneDX document")
		return &sbom.ReportValue{
			Ratio: 0,
			Value: "No",
		}
	}
	return &sbom.ReportValue{
		Ratio: 1,
		Value: "Yes",
	}
}

const rationPrecision = 2

func (r *CycloneDXReport) PackageIdentification() *sbom.ReportValue {
	return sbom.PrettyPercent(r.hasPurl+r.hasCPE, r.totalPackages*rationPrecision)
}

func (r *CycloneDXReport) PackageCPE() *sbom.ReportValue {
	return sbom.PrettyPercent(r.hasCPE, r.totalPackages)
}

func (r *CycloneDXReport) PackagePURL() *sbom.ReportValue {
	return sbom.PrettyPercent(r.hasPurl, r.totalPackages)
}

func (r *CycloneDXReport) PackageVersions() *sbom.ReportValue {
	return sbom.PrettyPercent(r.hasPackVersion, r.totalPackages)
}

func (r *CycloneDXReport) PackageDigests() *sbom.ReportValue {
	return sbom.PrettyPercent(r.hasPackDigest, r.totalPackages)
}

func (r *CycloneDXReport) PackageLicenses() *sbom.ReportValue {
	return sbom.PrettyPercent(r.hasLicense, r.totalPackages)
}

func (r *CycloneDXReport) CreationInfo() *sbom.ReportValue {
	return &sbom.ReportValue{
		Ratio: 1,
		Value: "100%",
	}
}

func (r *CycloneDXReport) GetTotalPackages() int {
	return r.totalPackages
}

func GetCycloneDXReport(b []byte) sbom.SbomReport {
	r := CycloneDXReport{}
	formats := []cdx.BOMFileFormat{cdx.BOMFileFormatJSON, cdx.BOMFileFormatXML}

	bom := new(cdx.BOM)
	for _, format := range formats {
		decoder := cdx.NewBOMDecoder(bytes.NewReader(b), format)
		if err := decoder.Decode(bom); err != nil {
			r.valid = false
			r.docError = err
		} else {
			r.valid = true
			r.docError = nil
			break
		}
	}

	if !r.valid {
		return &r
	}

	if bom.Metadata.Component != nil {
		r.name = bom.Metadata.Component.Name
	}

	if bom.Metadata.Tools != nil {
		for _, t := range *bom.Metadata.Tools {
			if t.Name != "" {
				r.creationToolName += 1
			}
			if t.Version != "" {
				r.creationToolVersion += 1
			}
		}
	}

	if bom.Components != nil {
		for _, p := range *bom.Components {
			r.totalPackages += 1
			if p.Licenses != nil && len(*p.Licenses) > 0 {
				r.hasLicense += 1
			}
			if p.Hashes != nil && len(*p.Hashes) > 0 {
				r.hasPackDigest += 1
			}
			if p.Version != "" {
				r.hasPackVersion += 1
			}
			if p.PackageURL != "" {
				r.hasPurl += 1
			}
			if p.CPE != "" {
				r.hasCPE += 1
			}
		}
	}

	return &r
}
