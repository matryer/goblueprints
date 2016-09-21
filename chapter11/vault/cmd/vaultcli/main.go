package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/matryer/goblueprints/chapter10/vault"
	grpcclient "github.com/matryer/goblueprints/chapter10/vault/client/grpc"
	"google.golang.org/grpc"
)

/*
	Usage

		vaultcli hash password

		vaultcli validate password hash

*/

func main() {
	var (
		grpcAddr = flag.String("addr", ":8081", "gRPC address")
	)
	flag.Parse()
	ctx := context.Background()
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		log.Fatalln("gRPC dial:", err)
	}
	defer conn.Close()
	vaultService := grpcclient.New(conn)
	args := flag.Args()
	var cmd string
	cmd, args = pop(args)
	switch cmd {
	case "hash":
		var password string
		password, args = pop(args)
		hash(ctx, vaultService, password)
	case "validate":
		var password, hash string
		password, args = pop(args)
		hash, args = pop(args)
		validate(ctx, vaultService, password, hash)
	default:
		log.Fatalln("unknown command", cmd)
	}
}

func hash(ctx context.Context, service vault.Service, password string) {
	h, err := service.Hash(ctx, password)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(h)
}

func validate(ctx context.Context, service vault.Service, password, hash string) {
	valid, err := service.Validate(ctx, password, hash)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if !valid {
		fmt.Println("invalid")
		os.Exit(1)
	}
	fmt.Println("valid")
}

func pop(s []string) (string, []string) {
	if len(s) == 0 {
		return "", s
	}
	return s[0], s[1:]
}
