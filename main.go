package main

import (
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:        "createRepository",
			Description: "未読ステータスを設定する",
			Action:      createRepository,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "input_file",
				},
				cli.StringFlag{
					Name: "output_file",
				},
			},
		},
	}
	app.Run(os.Args)
}

func createRepository(cCtx *cli.Context) {
	// inputファイルを取得する
	in := cCtx.String("input_file")
	if in == "" {
		log.Fatal("input_file is required")
		return
	}
	// outputファイルを取得する
	out := cCtx.String("output_file")
	if out == "" {
		log.Fatal("output_file is required")
		return
	}

	// inputファイルを読み込み、interfaceとそのメソッドの情報を取得する
	repo, err := parseRepositoryStructure(in)
	if err != nil {
		log.Fatal(err)
		return
	}

	// inputファイルのパッケージ名を取得
	pkgName, err := getPackageName(in)
	if err != nil {
		log.Fatal(err)
		return
	}

	// 取得したrepositoryの構造体から、repositoryの実装ファイルに記載する内容を作成する
	builder := NewRepositoryContentBuilder()
	ss, err := builder.Execute(repo, pkgName, getDirectoryName(out))
	if err != nil {
		log.Fatal(err)
		return
	}

	// outputファイルに書き込む
	p, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer p.Close()
	if _, err := p.Write([]byte(ss)); err != nil {
		log.Fatal(err)
		return
	}
}

// outputファイルのディレクトリ名を取得する
// 例：example/domain/user.goとなっている場合は、domainを返す
func getDirectoryName(outputFileName string) string {
	s := strings.Split(outputFileName, "/")
	return s[len(s)-2]
}
