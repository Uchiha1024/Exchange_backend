syntax = "proto3";

message {{.MessageName}}{
{{- range $index,$value := .Result}}
    {{$value.Type}} {{$value.MessageField}}  = {{Add $index 1}};
{{- end}}
}