package sbom

import "fmt"

type ReportValue struct {
	Ratio float32 `json:"ratio"`
	Value string  `json:"value"`
}

type SbomReport interface {
	GetName() string
	GetTotalPackages() int
	IsSpecCompliant() *ReportValue
	PackageIdentification() *ReportValue
	PackagePURL() *ReportValue
	PackageCPE() *ReportValue
	PackageVersions() *ReportValue
	PackageLicenses() *ReportValue
	CreationInfo() *ReportValue
}

type ReportResult struct {
	Name       string       `json:"name"`
	Packages   int          `json:"packages"`
	Compliance *ReportValue `json:"compliance"`
	// Identification (both Purl & CPEs)
	Identification *ReportValue `json:"identification"`
	PURL           *ReportValue `json:"purl"`
	CPE            *ReportValue `json:"cpe"`
	Version        *ReportValue `json:"version"`
	License        *ReportValue `json:"license"`
	Creation       *ReportValue `json:"creation"`
}

func GetReport(sr SbomReport) ReportResult {
	rr := ReportResult{
		Name:           sr.GetName(),
		Packages:       sr.GetTotalPackages(),
		Compliance:     sr.IsSpecCompliant(),
		Identification: sr.PackageIdentification(),
		PURL:           sr.PackagePURL(),
		CPE:            sr.PackageCPE(),
		Version:        sr.PackageVersions(),
		License:        sr.PackageLicenses(),
		Creation:       sr.CreationInfo(),
	}

	return rr
}

const percent100 = 100.0

func PrettyPercent(num, denom int) *ReportValue {
	r := &ReportValue{}
	if denom == 0 {
		return r
	}
	r.Ratio = float32(100 * (1.0 * num) / denom)

	f := ((float64(num) / float64(denom)) * percent100)
	r.Value = fmt.Sprintf("%.0f%%", f)

	return r
}
