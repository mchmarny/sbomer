package spdx

import (
	"fmt"
	"strings"

	"github.com/mchmarny/sbomer/pkg/sbom"
	spdx_common "github.com/spdx/tools-golang/spdx/common"

	"regexp"
)

var isNumeric = regexp.MustCompile(`\d`)

var missingPackages = sbom.ReportValue{
	Ratio:     0,
	Reasoning: "No packages",
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

func (r *SpdxReport) Report() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d total packages\n", r.totalPackages))
	sb.WriteString(fmt.Sprintf("%d total files\n", r.totalFiles))
	sb.WriteString(fmt.Sprintf("%d%% have licenses.\n", sbom.PrettyPercent(r.hasLicense, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%% have package digest.\n", sbom.PrettyPercent(r.hasPackDigest, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%% have package versions.\n", sbom.PrettyPercent(r.hasPackVer, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%% have purls.\n", sbom.PrettyPercent(r.hasPurl, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%% have CPEs.\n", sbom.PrettyPercent(r.hasCPE, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%% have file digest.\n", sbom.PrettyPercent(r.hasFileDigest, r.totalFiles)))
	sb.WriteString(fmt.Sprintf("Spec valid? %v\n", r.valid))
	sb.WriteString(fmt.Sprintf("Has creation info? %v\n", r.CreationInfo().Ratio == 1))

	return sb.String()
}

func (r *SpdxReport) IsSpecCompliant() sbom.ReportValue {
	if r.docError != nil {
		return sbom.ReportValue{
			Ratio:     0,
			Reasoning: r.docError.Error(),
		}
	}
	return sbom.ReportValue{Ratio: 1}
}

const rationPrecision = 2

func (r *SpdxReport) PackageIdentification() sbom.ReportValue {
	if r.totalPackages == 0 {
		return missingPackages
	}
	purlPercent := sbom.PrettyPercent(r.hasPurl, r.totalPackages)
	cpePercent := sbom.PrettyPercent(r.hasCPE, r.totalPackages)
	return sbom.ReportValue{
		// What percentage has both Purl & CPEs?
		Ratio:     float32(r.hasPurl+r.hasCPE) / float32(r.totalPackages*rationPrecision),
		Reasoning: fmt.Sprintf("%d%% have purls and %d%% have CPEs", purlPercent, cpePercent),
	}
}

func (r *SpdxReport) PackageVersions() sbom.ReportValue {
	if r.totalPackages == 0 {
		return sbom.ReportValue{
			Ratio:     0,
			Reasoning: "No packages",
		}
	}
	return sbom.ReportValue{
		Ratio: float32(r.hasPackVer) / float32(r.totalPackages),
	}
}

func (r *SpdxReport) PackageLicenses() sbom.ReportValue {
	if r.totalPackages == 0 {
		return sbom.ReportValue{
			Ratio:     0,
			Reasoning: "No packages",
		}
	}
	return sbom.ReportValue{
		Ratio: float32(r.hasLicense) / float32(r.totalPackages),
	}
}

const reportRationPrecision = .2

func (r *SpdxReport) CreationInfo() sbom.ReportValue {
	foundTool := false
	hasVersion := false

	if r.doc == nil || r.doc.GetCreationInfo() == nil {
		return sbom.ReportValue{
			Ratio:     0,
			Reasoning: "No creation info found",
		}
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
		return sbom.ReportValue{
			Ratio:     0,
			Reasoning: "No tool was used to create the sbom",
		}
	}

	if !hasVersion {
		return sbom.ReportValue{
			Ratio:     reportRationPrecision,
			Reasoning: "The tool used to create the sbom does not have a version",
		}
	}

	return sbom.ReportValue{
		Ratio: 1,
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
