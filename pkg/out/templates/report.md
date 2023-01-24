# SBOM Report 

> Generated: {{ .GeneratedOn }}

{{ range $r := .Items }}
### {{ $r.Name }}

* Compliant SBOM: {{ $r.Compliance.Value }}
* SBOM Creation Info: {{ $r.Creation.Value }}

| Package Data     | Ratio of {{ $r.Packages }} packages |
| ---------------- | ----------------------------------- |
| Id (CPE & PURL)  | {{ $r.Identification.Value }}       |
| With CPE         | {{ $r.CPE.Value }}                  |
| With PURL        | {{ $r.PURL.Value }}                 |
| With Version     | {{ $r.Version.Value }}              |
| With License     | {{ $r.License.Value }}              |

{{ end }}
