package client

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/ebi-fujizuku/sample_go_grpc_graphql/article/common"
	"github.com/ebi-fujizuku/sample_go_grpc_graphql/article/pb"
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



type Client_ServerStream struct{
	Service pb.ArticleServiceClient
}
func (c_ss *Client_ServerStream)List() {
	//streamレスポンスを受ける
	stream,err := c_ss.Service.ListArticle(
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
