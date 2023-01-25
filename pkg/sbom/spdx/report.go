package spdx

import (
	"fmt"
	"strings"

	"github.com/mchmarny/sbomer/pkg/sbom"
	"github.com/rs/zerolog/log"
	spdx_common "github.com/spdx/tools-golang/spdx/common"

	"regexp"
)

var isNumeric = regexp.MustCompile(`\d`)

var missingPackages = sbom.ReportValue{
	Ratio: 0,
	Value: "0%",
}

type SpdxReport struct {
	doc      Document
	docError error
	valid    bool

	totalPackages int
	totalFiles    int
	hasLicense    int
	hasPackDigest int
	hasPurl       int
	hasCPE        int
	hasFileDigest int
	hasPackVer    int
}

func (r *SpdxReport) GetName() string {
	return r.doc.GetName()
}

func (r *SpdxReport) IsSpecCompliant() *sbom.ReportValue {
	if r.docError != nil {
		log.Error().Err(r.docError).Msg("Invalid SPDX document")
		return &sbom.ReportValue{Ratio: 0, Value: "No"}
	}
	return &sbom.ReportValue{Ratio: 1, Value: "Yes"}
}

const rationPrecision = 2

func (r *SpdxReport) PackageIdentification() *sbom.ReportValue {
	if r.totalPackages == 0 {
		return &missingPackages
	}
	return sbom.PrettyPercent(r.hasPurl+r.hasCPE, r.totalPackages*rationPrecision)
}

func (r *SpdxReport) PackageCPE() *sbom.ReportValue {
	return sbom.PrettyPercent(r.hasCPE, r.totalPackages)
}

func (r *SpdxReport) PackagePURL() *sbom.ReportValue {
	return sbom.PrettyPercent(r.hasPurl, r.totalPackages)
}

func (r *SpdxReport) PackageVersions() *sbom.ReportValue {
	if r.totalPackages == 0 {
		return &missingPackages
	}
	return sbom.PrettyPercent(r.hasPackVer, r.totalPackages)
}

func (r *SpdxReport) PackageLicenses() *sbom.ReportValue {
	if r.totalPackages == 0 {
		return &missingPackages
	}
	return sbom.PrettyPercent(r.hasLicense, r.totalPackages)
}

func (r *SpdxReport) GetTotalPackages() int {
	return r.totalPackages
}

const reportRationPrecision = .2

func (r *SpdxReport) CreationInfo() *sbom.ReportValue {
	foundTool := false
	hasVersion := false

	if r.doc == nil || r.doc.GetCreationInfo() == nil {
		return &missingPackages
	}

	for _, creator := range r.doc.GetCreationInfo().Creators {
		if creator.CreatorType == "Tool" {
			foundTool = true
			if isNumeric.MatchString(creator.Creator) {
				hasVersion = true
			}
		}
	}

	if !foundTool {
		return &missingPackages
	}

	if !hasVersion {
		return &sbom.ReportValue{
			Ratio: reportRationPrecision,
			Value: "20%",
		}
	}

	return &sbom.ReportValue{
		Ratio: 1,
		Value: "100%",
	}
}

const (
	noLicenseShort = "NONE"
	noLicenseLong  = "NOASSERTION"
)

func GetSpdxReport(b []byte) sbom.SbomReport {
	sr := SpdxReport{}
	doc, err := LoadDocument(b)
	if err != nil {
		fmt.Printf("loading document: %v\n", err)
		return &sr
	}

	// try to load the SPDX file's contents as a json file, version 2.2
	sr.doc = doc
	sr.docError = err
	sr.valid = err == nil
	if sr.doc != nil {
		packages := sr.doc.GetPackages()

		for _, p := range packages {
			sr.totalPackages += 1
			if p.PackageLicenseConcluded != noLicenseShort &&
				p.PackageLicenseConcluded != noLicenseLong &&
				p.PackageLicenseConcluded != "" {
				sr.hasLicense += 1
			} else if p.PackageLicenseDeclared != noLicenseShort &&
				p.PackageLicenseDeclared != noLicenseLong &&
				p.PackageLicenseDeclared != "" {
				sr.hasLicense += 1
			}

			if len(p.PackageChecksums) > 0 {
				sr.hasPackDigest += 1
			}

			for _, ref := range p.PackageExternalReferences {
				if ref.RefType == spdx_common.TypePackageManagerPURL {
					sr.hasPurl += 1
				}
			}

			for _, ref := range p.PackageExternalReferences {
				if strings.HasPrefix(ref.RefType, "cpe") {
					sr.hasCPE += 1
					break
				}
			}

			if p.PackageVersion != "" {
				sr.hasPackVer += 1
			}
		}

		for _, file := range sr.doc.GetFiles() {
			sr.totalFiles += 1
			if len(file.Checksums) > 0 {
				sr.hasFileDigest += 1
			}
		}
	}
	return &sr
}
