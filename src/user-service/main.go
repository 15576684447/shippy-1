package main

import (
	"github.com/micro/go-micro"
	pb "learn/shippy/src/user-service/proto/user"
	"log"
)

func main() {
	// 连接到数据库
	db, err := CreateConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("connect error: %v\n", err)
	}

	repo := &UserRepository{db}

	// 自动检查 User 结构是否变化
	db.AutoMigrate(&pb.User{})

	srv := micro.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	srv.Init()

	// 获取 broker 实例
	// pubSub := s.Server().Options().Broker
	publisher := micro.NewPublisher(topic, srv.Client())

	t := TokenService{repo}
	pb.RegisterUserServiceHandler(srv.Server(), &handler{repo, &t, publisher})

	if err := srv.Run(); err != nil {
		log.Fatalf("user service error: %v\n", err)
	}

}
