
func ({{ .ReceiverValue }} *{{ .ReceiverType }}) {{ .MethodName }}({{ .Args }}) ({{ .ReturnArgs }}) {
  {{ .Body.Log }}

  {{ .Body.ReturnValue }}
}
