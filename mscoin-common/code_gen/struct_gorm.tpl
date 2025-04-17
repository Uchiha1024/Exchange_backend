package gen


type {{.StructName}}_Gorm struct {
{{- range $index,$value := .Result }}
    {{$value.Field}}  {{$value.Type}}  {{$value.Gorm}}
{{- end}}
}