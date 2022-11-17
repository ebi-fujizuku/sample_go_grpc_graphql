package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/ebi-fujizuku/sample_go_grpc_graphql/article/common"
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
	s := grpc.NewServer(
		grpc.UnaryInterceptor(myUnaryServerInterceptor1),
		grpc.StreamInterceptor(myStreamServerInterceptor1),
	)

	// gRPCサーバーにGreetingServiceを登録
	pb.RegisterArticleServiceServer(s, service)

	// サーバーリフレクションの設定
	reflection.Register(s)

	// 作成したgRPCサーバーを、8080番ポートで稼働させる
	go func() {
		log.Printf("start gRPC server port: %v", port)
		common.PrintDelimiter(1)
		s.Serve(listener)
	}()

	// Ctrl+Cが入力されたらGraceful shutdownされるようにする
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	common.PrintDelimiter(1)
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}

func myUnaryServerInterceptor1(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler)(resp interface{}, err error){
	log.Println("[pre]Unary",info.FullMethod,common.DELIMITER_2) // ハンドラの前に割り込ませる前処理
	res,err := handler(ctx,req) // 本来の処理
	log.Println("[post]Unary",info.FullMethod,common.DELIMITER_2) // ハンドラの後に割り込ませる前処理
	return res,err
}

func myStreamServerInterceptor1(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error{
	// ストリームopen時の前処理
	log.Println("[pre]Stream",info.FullMethod,common.DELIMITER_2)
	err := handler(srv,&myServerStreamWrapper1{ss})
	// ストリームがcloseされる時に行われる後処理
	log.Println("[post]Stream",info.FullMethod,common.DELIMITER_2)
	return err
}

type myServerStreamWrapper1 struct{
	grpc.ServerStream
}

func (s *myServerStreamWrapper1)RecMsg(m interface{})error{
	// ストリームからリクエストを受信
	err := s.ServerStream.RecvMsg(m)
	// 受信したリクエストを、ハンドラで処理する前に差し込む前処理
	if !errors.Is(err,io.EOF){
		log.Println("[pre msg]RecvMsg",m,common.DELIMITER_3)
	}
	return err
}

func (s *myServerStreamWrapper1)SendMsg(m interface{})error{
	// ハンドラで作成したレスポンスを、ストリームから返信する直前に差し込む後処理
	log.Println("[post msg]SendMsg",m,common.DELIMITER_3)
	return s.ServerStream.SendMsg(m)
}
