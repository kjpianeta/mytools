package version

import (
	"bytes"
	"runtime"
	"strings"
	"text/template"
)

// Build information. Populated at build-time.
var (
	Version     string
	GitCommit   string
	GitBranch   string
	GitState    string
	GitSummary  string
	BuildDate   string
	BuildUserID string
	GoVersion   = runtime.Version()
)

var versionInfoTmpl = `
{{.program}}, version {{.version}} (branch: {{.branch}}, revision: {{.revision}})
  build user:       {{.buildUser}}
  build date:       {{.buildDate}}
  git state:        {{.gitState}}
  git summary:      {{.gitSummary}}
  go version:       {{.goVersion}}
`

// Print returns version information.
func Print(program string) string {
	var m = map[string]string{
		"program":    program,
		"version":    Version,
		"revision":   GitCommit,
		"branch":     GitBranch,
		"buildUser":  BuildUserID,
		"gitState":   GitState,
		"gitSummary": GitSummary,
		"buildDate":  BuildDate,
		"goVersion":  GoVersion,
	}
	t := template.Must(template.New("version").Parse(versionInfoTmpl))

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "version", m); err != nil {
		panic(err)
	}
	return strings.TrimSpace(buf.String())
}
