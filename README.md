# sample_go_grpc_graphql
## 概要
+ protocのバージョン3.21.8で[作ってわかる！ はじめてのgRPC](https://zenn.dev/hsaki/books/golang-grpc-starting)やった感じです。
## 参考サイト
+ [GoでgRPCサーバーを立ててみる](https://zenn.dev/k88t76/books/f3892660871ab2)
	+ このあと、GraphQLも学べるとのことで、これを作る。
	+ ただし、gRPCが古いのか、このサイト通りに作れなかったため、改造。
+ [作ってわかる！ はじめてのgRPC](https://zenn.dev/hsaki/books/golang-grpc-starting)
	+ Stepごとに作れるのでよくわかりやすかった。
## 現在状況
	grpcまでしかやってません。GraphQLはそのうち。
## 宣伝
	下記の配信で作っていたものです。
+ [gRPCを学ぶんだ！](https://youtube.com/playlist?list=PL7eUTvd9iQoWKhVHa7-pAwJlcXnbUeK9u)

## 利用する場合
### 必須
+ /sample_go_grpc_graphql/article/configに.envを作成し、下記のようにコーディング。
	~~~bash:.env
	SQLITE3_PATH=あなたのローカルのclone先/sample_go_grpc_graphql/article/mydb.sqlite3
	~~~
### 必要ならば
+ もし、これを利用して、自身のGitHubにアップロードする場合、go moduleのパスは変更しましょう。
	```bash
	# go moduleのパス変更
	go mod edit -module github.com/あなたのGitHubアカウント名/sample_go_grpc_graphql/article
	# go moduleを更新
	go mod tidy
	```
	参考: [Goメモ-260 (go.modのモジュール名を変更)(go mod edit -module)](https://devlights.hatenablog.com/entry/2022/10/24/073000)
