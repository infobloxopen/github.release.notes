package githubclient

const (
	defaultTitle = `[{{.Tag}}] {{.TagComment}} ({{.TagDate}})`

	defaultBody = `{{if .ChangeLogLink -}}
[Full Changelog]({{.ChangeLogLink}})

{{end -}}
**New commits and merged pull requests:**
{{range .Commits}}
- {{.CommitMessage}}
{{- if .CommitPR}} {{.CommitPR}} {{end}}
{{- if and (not .CommitPR) .CommitURL}} [[link]({{.CommitURL}})] {{end -}}
([{{.CommitAuthor}}]({{.CommitAuthorURL}})) {{.CommitDate}}
{{- end}}`
)
