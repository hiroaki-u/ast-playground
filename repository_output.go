package main

import (
	"bytes"
	"embed"
	"strings"
	"text/template"
)

type RepositoryContentBuilder struct {
	MethodTemplate  *template.Template
	NewFuncTemplate *template.Template
}

func NewRepositoryContentBuilder() *RepositoryContentBuilder {
	return &RepositoryContentBuilder{
		MethodTemplate:  methodTemplate,
		NewFuncTemplate: newFuncTemplate,
	}
}

func (rc *RepositoryContentBuilder) Execute(repo Repository, interfacePkg, pkgName string) (string, error) {
	// 型定義やNew関数の作成
	factory, err := rc.createFactory(repo, interfacePkg, pkgName)
	if err != nil {
		return "", err
	}

	// メソッドの作成
	methods, err := rc.createMethod(repo, interfacePkg)
	if err != nil {
		return "", err
	}

	// 合算して返す
	return factory + strings.Join(methods, "\n"), nil
}

// メソッドの作成
func (rc *RepositoryContentBuilder) createMethod(repo Repository, pkgName string) ([]string, error) {
	res := []string{}

	for _, method := range repo.Methods {
		nilList := []string{}
		for i := 0; i < len(method.Returns); i++ {
			nilList = append(nilList, "nil")
		}
		returnStr := "return " + strings.Join(nilList, ", ")
		body := `log.Default().Println("` + repo.Name + "." + method.Name + `")
		` + returnStr

		content := &methodParameter{
			ReceiverValue: "r",
			ReceiverType:  repo.getLowerName(),
			MethodName:    method.Name,
			Args:          method.Args.GetTemplate(pkgName),
			ReturnArgs:    method.Returns.GetTemplate(pkgName),
			Body:          body,
		}
		var buf bytes.Buffer
		if err := rc.MethodTemplate.Execute(&buf, content); err != nil {
			return nil, err
		}
		res = append(res, buf.String())
	}
	return res, nil
}

type methodParameter struct {
	ReceiverValue string
	ReceiverType  string
	MethodName    string
	Args          string
	ReturnArgs    string
	Body          string
}

func (rc *RepositoryContentBuilder) createFactory(repo Repository, interfacePkg, pkgName string) (string, error) {
	content := &factoryParameter{
		LowerName:        repo.getLowerName(),
		Name:             repo.Name,
		InterfacePackage: interfacePkg,
		Package:          pkgName,
	}
	var buf bytes.Buffer
	if err := rc.NewFuncTemplate.Execute(&buf, content); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type factoryParameter struct {
	LowerName        string
	Name             string
	InterfacePackage string
	Package          string
}

type Repository struct {
	Name    string
	Methods []*Method
}

func (r *Repository) getLowerName() string {
	return strings.ToLower(string(r.Name[0])) + r.Name[1:]
}

type Method struct {
	Name    string
	Args    MethodValues
	Returns MethodValues
}

type MethodValues []MethodValue

func (r *MethodValues) GetTemplate(pkgName string) string {
	var s []string
	for _, v := range *r {
		s = append(s, v.GetTemplate(pkgName))
	}
	return strings.Join(s, ", ")
}

type MethodValue struct {
	Type   *MethodType
	Values []string
}

type MethodType struct {
	isSlice        bool   // slice
	isPointer      bool   // ポインタ
	isVariadic     bool   // 可変長引数
	requirePkgName bool   // package名が必要な場合
	Value          string // 型名
}

func (r *MethodType) GetFormatValue(pkgName string) string {
	res := r.Value
	if r.requirePkgName {
		res = pkgName + "." + res
	}
	if r.isPointer {
		res = "*" + res
	}
	if r.isSlice {
		return "[]" + res
	}
	if r.isVariadic {
		res = "..." + res
	}
	return res
}

func (r *MethodValue) GetTemplate(pkgName string) string {
	if r.HasNoValue() {
		return r.Type.GetFormatValue(pkgName)
	}
	return strings.Join(r.Values, ", ") + " " + r.Type.GetFormatValue(pkgName)
}

func (r *MethodValue) AppendValue(v string) {
	r.Values = append(r.Values, v)
}

func (r *MethodValue) HasNoValue() bool {
	return len(r.Values) == 0
}

var (
	//go:embed templates/*
	templates embed.FS

	// 関数用テンプレート
	methodTemplate = template.Must(template.ParseFS(templates, "templates/method.tpl"))
	// New関数用テンプレート
	newFuncTemplate = template.Must(template.ParseFS(templates, "templates/factory.tpl"))
)
