package {{ .Package }}

type {{ .LowerName }} struct {
}

func New{{ .Name }}() {{ .InterfacePackage }}.{{ .Name }} {
	return & {{ .LowerName }}{}
}