package sbom

import (
	"fmt"
	"strings"
)

type ReportValue struct {
	Ratio     float32 `json:"ratio"`
	Reasoning string  `json:"reasoning"`
}

const validPoints = 25
const generationPoints = 15
const packageSectionWeight = 20

type SbomReport interface {
	GetName() string
	IsSpecCompliant() ReportValue
	PackageIdentification() ReportValue
	PackageVersions() ReportValue
	PackageLicenses() ReportValue
	CreationInfo() ReportValue
	Report() string
}

type ScoreValue struct {
	ReportValue
	MaxPoints float32 `json:"maxPoints"`
}

func (sv *ScoreValue) Score() float32 {
	return sv.Ratio * sv.MaxPoints
}

type ReportResult struct {
	Name                  string     `json:"name"`
	Compliance            ScoreValue `json:"compliance"`
	PackageIdentification ScoreValue `json:"identification"`
	PackageVersions       ScoreValue `json:"version"`
	PackageLicenses       ScoreValue `json:"license"`
	CreationInfo          ScoreValue `json:"creation"`
	Total                 ScoreValue `json:"total"`
}

func reportValueToScore(rv ReportValue, maxPoints float32) ScoreValue {
	sv := ScoreValue{
		MaxPoints: maxPoints,
	}
	sv.Ratio = rv.Ratio
	sv.Reasoning = rv.Reasoning
	return sv
}

func GetReport(sr SbomReport) ReportResult {
	rr := ReportResult{
		Name:                  sr.GetName(),
		Compliance:            reportValueToScore(sr.IsSpecCompliant(), validPoints),
		PackageIdentification: reportValueToScore(sr.PackageIdentification(), packageSectionWeight),
		PackageVersions:       reportValueToScore(sr.PackageVersions(), packageSectionWeight),
		PackageLicenses:       reportValueToScore(sr.PackageLicenses(), packageSectionWeight),
		CreationInfo:          reportValueToScore(sr.CreationInfo(), generationPoints),
	}
	var totalPoints float32
	var maxPoints float32

	maxPoints += rr.Compliance.MaxPoints
	totalPoints += rr.Compliance.Score()

	maxPoints += rr.PackageIdentification.MaxPoints
	totalPoints += rr.PackageIdentification.Score()

	maxPoints += rr.PackageVersions.MaxPoints
	totalPoints += rr.PackageVersions.Score()

	maxPoints += rr.PackageLicenses.MaxPoints
	totalPoints += rr.PackageLicenses.Score()

	maxPoints += rr.CreationInfo.MaxPoints
	totalPoints += rr.CreationInfo.Score()

	rr.Total = ScoreValue{
		MaxPoints: maxPoints,
	}
	rr.Total.Ratio = totalPoints / maxPoints

	return rr
}

func getReportValueInfo(title string, sv ScoreValue) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s: %d/%d", title, int(sv.Score()), int(sv.MaxPoints)))
	if sv.Reasoning != "" {
		sb.WriteString(fmt.Sprintf(" (%s)", sv.Reasoning))
	}
	sb.WriteString("\n")
	return sb.String()
}

func Grade(sr SbomReport) string {
	sv := GetReport(sr)
	var sb strings.Builder

	sb.WriteString(getReportValueInfo("Spec Compliance", sv.Compliance))
	sb.WriteString(getReportValueInfo("Package ID", sv.PackageIdentification))
	sb.WriteString(getReportValueInfo("Package Versions", sv.PackageVersions))
	sb.WriteString(getReportValueInfo("Package Licenses", sv.PackageLicenses))
	sb.WriteString(getReportValueInfo("Creation Info", sv.CreationInfo))

	sb.WriteString(fmt.Sprintf("Total points: %d/%d or %d%%\n", int(sv.Total.Score()), int(sv.Total.MaxPoints), PrettyPercent(int(sv.Total.Score()), int(sv.Total.MaxPoints))))

	return sb.String()
}

func PrettyPercent(num, denom int) int {
	if denom == 0 {
		return 0
	}
	return 100 * (1.0 * num) / denom
}
