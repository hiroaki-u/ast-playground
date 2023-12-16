package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// ファイルを走査し、ファイル内で見つかったrepositoryを返す
func parseRepositoryStructure(file string) (Repository, error) {
	// ファイルをパースして、astを取得する
	repo := Repository{}
	f, err := parser.ParseFile(token.NewFileSet(), file, nil, parser.Mode(0))
	if err != nil {
		return repo, err
	}

	// astを走査して、必要な情報を抽出
	ast.Inspect(f, func(n ast.Node) bool {
		methods := []*Method{}
		switch x := n.(type) {
		case *ast.TypeSpec:
			// 対象がinterfaceであるため、interfaceの情報を取得する
			it, ok := x.Type.(*ast.InterfaceType)
			if !ok {
				return true
			}
			repo.Name = x.Name.Name

			// 関数の情報を取得
			for _, field := range it.Methods.List {
				funcType, ok := field.Type.(*ast.FuncType)
				if !ok {
					continue
				}
				method := Method{}
				method.Name = field.Names[0].Name

				// 引数と返り値を取得する
				method.Args = ExtractMethodValues(funcType.Params.List)
				method.Returns = ExtractMethodValues(funcType.Results.List)
				methods = append(methods, &method)
			}
			repo.Methods = methods
		}
		return true
	})
	return repo, nil
}

// 変数と型を統合する（引数と返り値に利用）
// 例：第一返り値：[]string{ctx context.Context, id int} , 第二返り値 []string{ctx, id}
func ExtractMethodValues(list []*ast.Field) []MethodValue {
	mvs := []MethodValue{}
	for _, param := range list {
		mv := MethodValue{}
		mt := &MethodType{}
		IdentifyNodeType(param.Type, mt)
		mv.Type = mt
		if param.Names != nil {
			for _, p := range param.Names {
				mv.AppendValue(p.Name)
			}
		}
		mvs = append(mvs, mv)
	}
	return mvs
}

// 各引数や返り値の型を特定して、MethodTypeに格納する
func IdentifyNodeType(t ast.Expr, mt *MethodType) {
	switch t.(type) {
	// sliceの場合
	case *ast.ArrayType:
		se := t.(*ast.ArrayType).Elt
		mt.isSlice = true
		IdentifyNodeType(se, mt)
	// pointer型の場合
	case *ast.StarExpr:
		se, _ := t.(*ast.StarExpr).X.(*ast.Ident)
		if se != nil {
			if !isPrimitive(se.Name) {
				mt.isPointer = true
				mt.requirePkgName = true
				mt.Value = se.Name
			} else {
				mt.isPrimitive = true
			}
		} else {
			se, _ := t.(*ast.StarExpr).X.(*ast.SelectorExpr)
			x := se.X.(*ast.Ident)
			sel := se.Sel
			mt.isPointer = true
			mt.Value = x.Name + "." + sel.Name
		}
	// シンプルな型の場合（primitive型やstruct）
	case *ast.Ident:
		se := t.(*ast.Ident)
		mt.Value = se.Name
		if !isPrimitive(t.(*ast.Ident).Name) {
			mt.requirePkgName = true
		} else {
			mt.isPrimitive = true
		}
	// package + structの場合
	case *ast.SelectorExpr:
		x := t.(*ast.SelectorExpr).X.(*ast.Ident)
		sel := t.(*ast.SelectorExpr).Sel
		mt.Value = x.Name + "." + sel.Name
	// 可変引数
	case *ast.Ellipsis:
		se := t.(*ast.Ellipsis).Elt
		mt.Value = "..." + se.(*ast.Ident).Name
	default:
	}
}

func getPackageName(fileName string) (string, error) {
	pkg, err := parser.ParseFile(token.NewFileSet(), fileName, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", err
	}
	return pkg.Name.Name, nil
}
