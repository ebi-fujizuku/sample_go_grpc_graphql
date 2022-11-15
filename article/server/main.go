package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	myConf "github.com/ebi-fujizuku/sample_go_grpc_graphql/article/config"
	"github.com/ebi-fujizuku/sample_go_grpc_graphql/article/pb"
	"github.com/ebi-fujizuku/sample_go_grpc_graphql/article/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	myConf.NewConfig("../config/.env")
	// 8080番portのLisnterを作成
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// Service作成
	service,err := service.NewService()
	if err != nil{
		panic(err)
	}

	// gRPCサーバーを作成
	s := grpc.NewServer()

	// gRPCサーバーにGreetingServiceを登録
	pb.RegisterArticleServiceServer(s, service)

	// サーバーリフレクションの設定
	reflection.Register(s)

	// 作成したgRPCサーバーを、8080番ポートで稼働させる
	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	// Ctrl+Cが入力されたらGraceful shutdownされるようにする
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}
