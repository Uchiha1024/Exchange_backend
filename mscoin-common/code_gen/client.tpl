package {{.RpcCommon.PackageName}}

import (
    "context"
    "github.com/zeromicro/go-zero/zrpc"
    "google.golang.org/grpc"
    "grpc-common/{{.RpcCommon.GrpcPackage}}/types/{{.RpcCommon.GrpcPackage}}"
)

type (
{{- range $index,$value := .ParamList}}
    {{$value}}      = {{$.RpcCommon.GrpcPackage}}.{{$value}}
{{- end}}

    {{.RpcCommon.ServiceName}} interface {
{{- range $index,$value := .Rpc}}
        {{$value.FunName}}(ctx context.Context, in *{{$value.Req}}, opts ...grpc.CallOption) (*{{$value.Resp}}, error)
{{- end}}
    }

    default{{.RpcCommon.ServiceName}} struct {
        cli zrpc.Client
    }
)

func New{{.RpcCommon.ServiceName}}(cli zrpc.Client) *default{{.RpcCommon.ServiceName}} {
    return &default{{.RpcCommon.ServiceName}}{
        cli: cli,
    }
}
{{- range $index,$value := .Rpc}}
func (m *default{{$.RpcCommon.ServiceName}}) {{$value.FunName}}(ctx context.Context, in *{{$value.Req}}, opts ...grpc.CallOption) (*{{$value.Resp}}, error) {
    client := market.New{{$.RpcCommon.ServiceName}}Client(m.cli.Conn())
    return client.{{$value.FunName}}(ctx, in, opts...)
}
{{- end}}
