package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/matryer/filedb"
)

/*
  backup command

  usage:

    backup -db=./backupdata.db add {path} [{path} {path}...]
    backup -db=./backupdata.db remove {path} [{path} {path}...]
    backup -db=./backupdata.db list

*/

type path struct {
	Path string
	Hash string
}

func (p path) String() string {
	return fmt.Sprintf("%s [%s]", p.Path, p.Hash)
}

func main() {
	var fatalErr error
	defer func() {
		if fatalErr != nil {
			flag.PrintDefaults()
			log.Fatalln(fatalErr)
		}
	}()
	var (
		dbpath = flag.String("db", "./backupdata", "path to database directory")
	)
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fatalErr = errors.New("invalid usage; must specify command")
		return
	}
	db, err := filedb.Dial(*dbpath)
	if err != nil {
		fatalErr = err
		return
	}
	defer db.Close()
	col, err := db.C("paths")
	if err != nil {
		fatalErr = err
		return
	}
	switch strings.ToLower(args[0]) {
	case "list":
		var path path
		col.ForEach(func(i int, data []byte) bool {
			err := json.Unmarshal(data, &path)
			if err != nil {
				fatalErr = err
				return false
			}
			fmt.Printf("= %s\n", path)
			return false
		})
	case "add":
		if len(args[1:]) == 0 {
			fatalErr = errors.New("must specify path to add")
			return
		}
		for _, p := range args[1:] {
			path := &path{Path: p, Hash: "Not yet archived"}
			if err := col.InsertJSON(path); err != nil {
				fatalErr = err
				return
			}
			fmt.Printf("+ %s\n", path)
		}
	case "remove":
		var path path
		col.RemoveEach(func(i int, data []byte) (bool, bool) {
			err := json.Unmarshal(data, &path)
			if err != nil {
				fatalErr = err
				return false, true
			}
			for _, p := range args[1:] {
				if path.Path == p {
					fmt.Printf("- %s\n", path)
					return true, false
				}
			}
			return false, false
		})
	}
}
