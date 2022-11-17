package service

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/ebi-fujizuku/sample_go_grpc_graphql/article/common"
	"github.com/ebi-fujizuku/sample_go_grpc_graphql/article/pb"
	"github.com/ebi-fujizuku/sample_go_grpc_graphql/article/repository"
)

// 自作サービスのインターフェース
type Service struct{
	pb.UnimplementedArticleServiceServer
	repository repository.Repository
}

// 自作サービスの構造体のコンストラクタ
func NewService()(*Service,error){
	repo,err := repository.NewSqliteRepo()
	if err != nil{
		return nil,err
	}
	return &Service{
		repository: repo,
	},nil
}

func (s *Service)CreateArticle(ctx context.Context, req *pb.CreateArticleRequest)(*pb.CreateArticleResponse,error){
	common.PrintStart("")
	// INSERTする記事のInputを取得
	input := req.ArticleInput

	// 記事をDBにINSERTし、INSERTした記事のIDを返す
	id,err := s.repository.InsertArticle(input)
	if err != nil{
		return nil,err
	}
	common.PrintEnd("")
	return &pb.CreateArticleResponse{
		Article: &pb.Article{
			Id:      id,
			Author:  input.Author,
			Title:   input.Title,
			Content: input.Content,
		},
	},nil
}

func (s *Service)ReadArticle(ctx context.Context, req *pb.ReadArticleRequest)(*pb.ReadArticleResponse,error){
	// INSERTする記事のInputを取得
	id := req.GetId()

	// DBから該当IDの記事を取得
	a,err := s.repository.SelectArticleByID(id)
	if err != nil{
		return nil,err
	}

	// 取得した記事をレスポンスとして返す
	return &pb.ReadArticleResponse{
		Article: &pb.Article{
			Id:      id,
			Author:  a.Author,
			Title:   a.Title,
			Content: a.Content,
		},
	},nil
}

func (s *Service)UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest)(*pb.UpdateArticleResponse,error){
	id    := req.GetId()
	input := req.GetArticleInput()

	// 該当記事をUPDATE
	if err := s.repository.UpdateArticle(id,input); err != nil{
		return nil,err
	}

	// 取得した記事をレスポンスとして返す
	return &pb.UpdateArticleResponse{
		Article: &pb.Article{
			Id:      id,
			Title:   input.Title,
			Author:  input.Author,
			Content: input.Content,
		},
	},nil
}

func (s *Service)DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest)(*pb.DeleteArticleResponse,error){
	id    := req.GetId()

	// 該当記事をUPDATE
	if err := s.repository.DeleteArticle(id); err != nil{
		return nil,err
	}

	// 取得した記事をレスポンスとして返す
	return &pb.DeleteArticleResponse{Id: id},nil
}

func (s *Service)ListArticle(
	req *pb.ListArticleRequest,
	stream pb.ArticleService_ListArticleServer,
)error{
	// 全記事を取得
	rows,err := s.repository.SelectAllArticles()
	if err != nil{
		return err
	}

	// 取得した記事を1つ1つレスポンスとしてServer Streamingで返す
	for rows.Next(){
		var a pb.Article
		err := rows.Scan(&a.Id,&a.Author,&a.Title,&a.Content)
		if err != nil{
			return err
		}
		stream.Send(&pb.ListArticleResponse{Article: &a})
	}

	return nil
}

func (s *Service)CreateArticles(stream pb.ArticleService_CreateArticlesServer)error{
	common.PrintStart("")
	articleList := make([]*pb.ArticleInput,0,3)
	for{
		req,err := stream.Recv();
		if errors.Is(err,io.EOF){
			fmt.Println("Create inputs:")
			ids := make([]int64,0,3)
			for _,article := range articleList{
				fmt.Println(article)
				id,err := s.repository.InsertArticle(article)
				if err != nil{
					return err
				}
				ids = append(ids, id)
			}
			return stream.SendAndClose(&pb.CreateArticlesResponse{
				Id: ids,
			})
		}
		if err != nil{
			return err
		}
		articleList = append(articleList, req.GetArticleInput())
	}
}

func (s *Service)FreeReadArticles(stream pb.ArticleService_FreeReadArticlesServer)error{
	for{
		req,err := stream.Recv()
		if errors.Is(err,io.EOF){
			return nil
		}
		if err != nil{
			return err
		}
		fmt.Println(req.GetId()," is Nice!")
		a,err := s.repository.SelectArticleByID(req.GetId())
		if err != nil{
			return err
		}

		if err = stream.Send(&pb.FreeReadArticlesResponse{
			Article:&pb.Article{
				Id: a.GetId(),
				Author: a.GetAuthor(),
				Title: a.GetTitle(),
				Content: a.GetContent(),
			},
		});err != nil{
			return err
		}
	}

}
