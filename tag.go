package main

import (
	"flag"
	"fmt"
	"os"
	"tag/config"
	"tag/print"
	"tag/service"
)

func usage() {
	fmt.Println("usage: tag -op [rc,rd] -m 'comment' -vt [1,2,3] -ext rc")
	fmt.Println("\t *require -op: operate,rc(release version, -m require) rd(delete tag, -tag require), default:rc ")
	fmt.Println("\t -m: comment ")
	fmt.Println("\t -vt: version type, 1.major 2.minor 3.amendment default:3 ")
	fmt.Println("\t -ext: version number extra string, ignorable")
	fmt.Println("\t -tag: specified value")
	fmt.Println("e.g:")
	fmt.Println("\t tag -op rc -m 'submit gpstream-888' ")
	fmt.Println("\t tag -op rc -m 'submit gpstream-888' -tag 'v1.1.1'")
	fmt.Println("\t tag -op rc -m 'submit gpstream-888' -ext 'rc'")
	fmt.Println("\t tag -op rd -tag 'v1.1.1'")
}

func initPrint() {
	print.SetLevel("debug")
}

func initFlag() {
	var (
		operate      string
		comment      string
		version_type string
		version_ext  string // 版本号额外标识
		tag          string // 指定tag
	)
	flag.StringVar(&operate, "op", "rc", "operate")
	flag.StringVar(&comment, "m", "", "comment")
	flag.StringVar(&version_type, "vt", "", "version type")
	flag.StringVar(&version_ext, "ext", "", "version numer extra string")
	flag.StringVar(&tag, "tag", "", "tag, specified value")
	flag.Parse()

	memconf := config.NewMemConfig()
	memconf.Set("op", operate)
	memconf.Set("comment", comment)
	memconf.Set("version_type", version_type)
	memconf.Set("version_ext", version_ext)
	memconf.Set("tag", tag)
}

func main() {
	if len(os.Args) <= 1 {
		usage()
		return
	}
	initFlag()
	initPrint()

	op := config.NewMemConfig().Get("op")
	switch op {
	case "rc":
		service.NewGitSvc().Release()
	case "rd":
		service.NewGitSvc().DeleteTag()
	default:
		usage()
	}
}
