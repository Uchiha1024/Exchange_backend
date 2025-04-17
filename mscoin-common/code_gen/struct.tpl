package gen


type {{.StructName}} struct {
{{- range $index,$value := .Result }}
    {{$value.Field}}  {{$value.Type}}  {{$value.JsonForm}}
{{- end}}
}