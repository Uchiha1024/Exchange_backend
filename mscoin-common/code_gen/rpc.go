package code_gen

import (
	"html/template"
	"log"
	"os"
	"strings"
)

type RpcCommon struct {
	PackageName string
	GrpcPackage string
	ModuleName  string
	ServiceName string
}
type Rpc struct {
	FunName string
	Req     string
	Resp    string
}
type RpcResult struct {
	RpcCommon RpcCommon
	Rpc       []Rpc
	ParamList []string
}

func GenZeroRpc(result RpcResult) {
	t := template.New("client.tpl")
	tmpl, err := t.ParseFiles("./client.tpl")
	log.Println(err)
	_, err = os.Stat("./gen")
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("./gen", 0666)
		}
	}
	var pl []string
	for _, v := range result.Rpc {
		if !isContain(pl, v.Req) {
			pl = append(pl, v.Req)
		}
		if !isContain(pl, v.Resp) {
			pl = append(pl, v.Resp)
		}
	}
	result.ParamList = pl
	file, err := os.Create("./gen/" + strings.ToLower(result.RpcCommon.ServiceName) + ".go")
	defer file.Close()
	log.Println(err)
	err = tmpl.Execute(file, result)
	log.Println(err)
}

func isContain(pl []string, str string) bool {
	for _, p := range pl {
		if p == str {
			return true
		}
	}
	return false
}
