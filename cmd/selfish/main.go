package main

import (
	"context"
	"log"
	"os"

	"github.com/podhmo/selfish"
	"github.com/podhmo/selfish/cmd/selfish/internal"
	"github.com/podhmo/selfish/pkg/commithistory"
	"github.com/podhmo/structflag"
)

// Option ...
type Option struct {
	Alias  string   `flag:"alias" help:"alias name of uploaded gists"`
	Delete bool     `flag:"delete" help:"delete uploaded gists"`
	Silent bool     `flag:"silent" help:"don't open gist pages with browser, after uploading"`
	Args   []string `flag:"-"`
}

func main() {
	opt := &Option{}
	b := structflag.NewBuilder()
	fs := b.Build(opt)
	if err := fs.Parse(os.Args[1:]); err != nil {
		fs.Usage()
		os.Exit(2)
	}
	opt.Args = fs.Args()

	if err := run(opt); err != nil {
		log.Fatalf("!!%+v", err)
	}
}

func run(opt *Option) error {
	v4 := commithistory.New("selfish")
	v5, err := selfish.LoadConfig(v4)
	if err != nil {
		return err
	}
	v6 := selfish.NewClient(v5)
	v7 := internal.NewApp(v4, v6, v5, opt.Silent, opt.Delete, opt.Alias)
	return v7.Run(context.Background(), opt.Args)
}
