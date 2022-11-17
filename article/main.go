package main

import (
	"fmt"
	"log"

	"github.com/ebi-fujizuku/sample_go_grpc_graphql/article/client"
	"github.com/ebi-fujizuku/sample_go_grpc_graphql/article/common"
	"github.com/ebi-fujizuku/sample_go_grpc_graphql/article/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func main() {
	common.PrintStart("gRPC Client")

	// // クライアント作成
	// err := client.NewClient(port)
	// if err != nil{
		// 	panic(err)
	// }

	// gRPCサーバーとのコネクションを確立
	port := "8080"
	address := fmt.Sprintf("localhost:%s",port)
	fmt.Println("address",address)
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("Connection failed.")
		fmt.Println("NewClient: end")
		// return nil,err
		panic(err)
	}
	defer conn.Close()

	// UnaryのgRPCクライアントを生成
	c :=&client.Client{
		Service: pb.NewArticleServiceClient(conn),
	}

	fmt.Println("--------------------------------")

	// 単方向Stream
	c.Create()
	c.Read()
	c.Update()
	c.Delete()
	fmt.Println("--------------------------------")

	// ServerStream
	c.List()
	fmt.Println("--------------------------------")

	// ClientStream
	c.Creates()
	fmt.Println("--------------------------------")

	// 双方向Stream
	c.FreeReadArticles()
	fmt.Println("--------------------------------")

	// エラー
	c.ErrorArticle()
	fmt.Println("--------------------------------")

	common.PrintEnd("gRPC Client")
}
