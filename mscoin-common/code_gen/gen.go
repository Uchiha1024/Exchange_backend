package code_gen

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"text/template"
)

func connectMysql() *gorm.DB {
	//配置MySQL连接参数
	username := "root"  //账号
	password := "root"  //密码
	host := "127.0.0.1" //数据库地址，可以是Ip或者域名
	port := 3309        //数据库端口
	Dbname := "mscoin"  //数据库名
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	return db
}

type Result struct {
	Field        string
	MessageField string
	Type         string
	Gorm         string
	Json         string
	Form         string
	JsonForm     string
}
type StructResult struct {
	StructName string
	Result     []*Result
}
type MessageResult struct {
	MessageName string
	Result      []*Result
}

func GenModel(table string, name string) {
	GenStruct(table, name)
	GenProtoMessage(table, name)
}

func GenStruct(table string, structName string) {
	db := connectMysql()
	var results []*Result
	db.Raw(fmt.Sprintf("describe %s", table)).Scan(&results)
	for _, v := range results {
		field := v.Field
		name := Name(field)       // 表字段 aa_bb  字段名 AaBb
		tfName := TFName(v.Field) //驼峰命名  aaBb
		v.Field = name
		v.Type = getType(v.Type)
		v.Json = "`json:\"" + tfName + "\"`"
		v.JsonForm = "`json:\"" + tfName + "\" from:\"" + tfName + "\"`"
		v.Gorm = "`gorm:\"column:" + field + "\"`"
	}
	tmpl, err := template.ParseFiles("./struct.tpl")
	log.Println(err)
	tmpl1, err := template.ParseFiles("./struct_gorm.tpl")
	log.Println(err)
	sr := StructResult{StructName: structName, Result: results}
	_, err = os.Stat("./gen")
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("./gen", 0666)
		}
	}
	file, err := os.Create("./gen/" + strings.ToLower(structName) + ".go")
	log.Println(err)
	tmpl.Execute(file, sr)
	defer file.Close()
	file1, err := os.Create("./gen/" + strings.ToLower(structName) + "_gorm.go")
	defer file1.Close()
	log.Println(err)
	tmpl1.Execute(file1, sr)
}

func GenProtoMessage(table string, messageName string) {
	db := connectMysql()
	var results []*Result
	db.Raw(fmt.Sprintf("describe %s", table)).Scan(&results)
	for _, v := range results {
		v.MessageField = TFName(v.Field)
		v.Type = getMessageType(v.Type)
	}
	var fm template.FuncMap = make(map[string]any)
	fm["Add"] = func(v int, add int) int {
		return v + add
	}
	t := template.New("message.tpl")
	t.Funcs(fm)
	tmpl, err := t.ParseFiles("./message.tpl")
	log.Println(err)
	sr := MessageResult{MessageName: messageName, Result: results}
	_, err = os.Stat("./gen")
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("./gen", 0666)
		}
	}
	file, err := os.Create("./gen/" + strings.ToLower(messageName) + ".proto")
	defer file.Close()
	log.Println(err)
	err = tmpl.Execute(file, sr)
	log.Println(err)
}

func getMessageType(t string) string {
	if strings.Contains(t, "bigint") {
		return "int64"
	}
	if strings.Contains(t, "varchar") {
		return "string"
	}
	if strings.Contains(t, "text") {
		return "string"
	}
	if strings.Contains(t, "tinyint") {
		return "int32"
	}
	if strings.Contains(t, "int") &&
		!strings.Contains(t, "tinyint") &&
		!strings.Contains(t, "bigint") {
		return "int32"
	}
	if strings.Contains(t, "double") {
		return "double"
	}
	if strings.Contains(t, "decimal") {
		return "double"
	}
	return ""
}

func getType(t string) string {
	if strings.Contains(t, "bigint") {
		return "int64"
	}
	if strings.Contains(t, "varchar") {
		return "string"
	}
	if strings.Contains(t, "text") {
		return "string"
	}
	if strings.Contains(t, "tinyint") {
		return "int"
	}
	if strings.Contains(t, "int") &&
		!strings.Contains(t, "tinyint") &&
		!strings.Contains(t, "bigint") {
		return "int"
	}
	if strings.Contains(t, "double") {
		return "float64"
	}
	if strings.Contains(t, "decimal") {
		return "float64"
	}
	return ""
}
func TFName(name string) string {
	var names = name[:]
	isSkip := false
	var sb strings.Builder
	for index, value := range names {
		if index == 0 {
			s := names[:index+1]
			s = strings.ToLower(s)
			sb.WriteString(s)
			continue
		}
		if isSkip {
			isSkip = false
			continue
		}
		//95 下划线  user_name
		if value == 95 {
			s := names[index+1 : index+2]
			s = strings.ToUpper(s)
			sb.WriteString(s)
			isSkip = true
			continue
		} else {
			s := names[index : index+1]
			sb.WriteString(s)
		}
	}
	return sb.String()
}
func Name(name string) string {
	var names = name[:]
	isSkip := false
	var sb strings.Builder
	for index, value := range names {
		if index == 0 {
			s := names[:index+1]
			s = strings.ToUpper(s)
			sb.WriteString(s)
			continue
		}
		if isSkip {
			isSkip = false
			continue
		}
		//95 下划线  user_name
		if value == 95 {
			s := names[index+1 : index+2]
			s = strings.ToUpper(s)
			sb.WriteString(s)
			isSkip = true
			continue
		} else {
			s := names[index : index+1]
			sb.WriteString(s)
		}
	}
	return sb.String()
}
