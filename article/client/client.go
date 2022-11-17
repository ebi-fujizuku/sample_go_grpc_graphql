package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/ebi-fujizuku/sample_go_grpc_graphql/article/common"
	"github.com/ebi-fujizuku/sample_go_grpc_graphql/article/pb"
	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

type Client struct{
	Service pb.ArticleServiceClient
}

func (c *Client)Create() {
	common.PrintStart("")

	articleInfo := &pb.ArticleInput{
		Author:  "fujito",
		Title:   "my hero academia",
		Content: "so nice Mange and Animetion",
	}

	req := &pb.CreateArticleRequest{
		ArticleInput: articleInfo,
	}
	res, err := c.Service.CreateArticle(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to CreateArticle: %v\n",err)
	}
	fmt.Printf("CreateArticle Response: %v\n",res.GetArticle())
	common.PrintEnd("")
}

func (c *Client)Read() {
	var id int64 = 2
	res,err := c.Service.ReadArticle(
		context.Background(),
		&pb.ReadArticleRequest{Id: id},
	)
	if err != nil{
		log.Fatalf("Failde to ReadArticle: %v\n",err)
	}
	fmt.Printf("ReadArticle Response: %v\n",res)
}

func (c *Client)Update() {
	var id int64= 2
	input := &pb.ArticleInput{
		Author:  "Izuku",
		Title:   "smile on ranway",
		Content: "chiyuki is so cool & cute",
	}
	res,err := c.Service.UpdateArticle(
		context.Background(),
		&pb.UpdateArticleRequest{Id: id,ArticleInput: input},
	)
	if err != nil{
		log.Fatalf("Failed to UpdateArticle: %v\n",err)
	}
	fmt.Printf("UpdateArticle Response: %v\n",res)
}

func (c *Client)Delete() {
	var id int64= 1
	res,err := c.Service.DeleteArticle(
		context.Background(),
		&pb.DeleteArticleRequest{Id: id},
	)
	if err != nil{
		log.Fatalf("Failed to DeleteArticle: %v\n",err)
	}
	fmt.Printf("DeleteArticle Response: %v\n",res)
}

func (c *Client)List() {
	//streamレスポンスを受ける
	stream,err := c.Service.ListArticle(
		context.Background(),
		&pb.ListArticleRequest{},
	)
	if err != nil{
		log.Fatalf("Failed to ListArticle: %v\n",err)
	}

	// Server Streamで渡されたレスポンスを1つ1つ出力
	for{
		res,err := stream.Recv()
		if err == io.EOF{
			fmt.Println("all the responses have already received.")
			break
		}
		if err != nil{
			log.Fatalf("Failed to Server Streaming: %v\n",err)
		}
		fmt.Println(res)
	}
}

func (c *Client)Creates() {
	stream,err := c.Service.CreateArticles(
		context.Background(),
	)
	if err != nil{
		log.Fatalf("Failed to CreatesArticle: %v\n",err)
		return
	}

	articles := make([]*pb.ArticleInput,0,3)
	articles = append(articles,&pb.ArticleInput{
		Author: "Kohei Horikoshi",
		Title: "my hero academia",
		Content: "Izuku Midoriya is Nice!",
	})
	articles = append(articles,&pb.ArticleInput{
		Author: "Kotoba Inotani",
		Title: "smile on runway",
		Content: "Chiyuki Fujito is Cute & Cool!",
	})
	articles = append(articles,&pb.ArticleInput{
		Author: "ONE",
		Title: "mob psycho 100",
		Content: "Mob is String!",
	})

	for _,article := range articles{
		if err := stream.Send(&pb.CreateArticlesRequest{
			ArticleInput: article,
		}); err != nil{
			fmt.Println(err)
			return
		}
	}

	res,err := stream.CloseAndRecv()
	if err != nil{
		fmt.Println(err)
	}else{
		fmt.Println(res.GetId())
	}

	fmt.Printf("CreatesArticle Response: %v\n",res)
}

func (c *Client)FreeReadArticles() {
	stream,err := c.Service.FreeReadArticles(context.Background())
	if err != nil{
		log.Fatalf("Failde to FreeReadArticles: %v\n",err)
		return
	}
	sendNum := []int{35,36,37}

	var sendEnd,recvEnd bool
	sendCount := 0
	const MAX_COUNT = 3
	for !(sendEnd && recvEnd){
		// 送信処理
		if !sendEnd{
			// 送信
			if err := stream.Send(&pb.FreeReadArticlesRequest{
				Id: int64(sendNum[sendCount]),
			});err != nil{
				fmt.Println(err)
				sendEnd = true
			}
			// sendEnd判定
			sendCount++
			if sendCount == MAX_COUNT{
				sendEnd = true
				if err := stream.CloseSend();err != nil{
					fmt.Println(err)
				}
			}
		}
		// 受信処理
		if !recvEnd{
			res,err := stream.Recv()
			if err != nil{
				if !errors.Is(err,io.EOF){
					fmt.Println(err)
				}
				recvEnd = true
			}else{
				fmt.Println(res.GetArticle())
			}
		}
	}
}

func (c *Client)ErrorArticle(){
	id := int64(1)
	res,err := c.Service.ErrorArticle(
		context.Background(),
		&pb.ReadArticleRequest{Id: id},
	)
	if err != nil{
		if stat,ok := status.FromError(err); ok{
			fmt.Printf("code: %s\n",stat.Code())
			fmt.Printf("message: %s\n",stat.Message())
		}else{
			fmt.Println()
		}
	}else{
		fmt.Printf("Success: %v\n",res.GetArticle())
	}
}

func (c *Client)RichErrorArticle(){
	id := int64(1)
	res,err := c.Service.ErrorArticle(
		context.Background(),
		&pb.ReadArticleRequest{Id: id},
	)
	if err != nil{
		if stat,ok := status.FromError(err); ok{
			fmt.Printf("code: %s\n",stat.Code())
			fmt.Printf("message: %s\n",stat.Message())
			fmt.Printf("details: %s\n",stat.Details())
		}else{
			fmt.Println()
		}
	}else{
		fmt.Printf("Success: %v\n",res.GetArticle())
	}
}
