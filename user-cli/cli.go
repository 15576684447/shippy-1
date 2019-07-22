package main

import (
	"github.com/micro/go-micro/config/cmd"
	"log"
	"os"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	microclient "github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	pb "learn/shippy/user-service/proto/user"
)

func main() {

	cmd.Init()

	// 创建 user-service 微服务的客户端
	client := pb.NewUserServiceClient("go.micro.srv.user", microclient.DefaultClient)

	// 设置命令行参数
	service := micro.NewService(
		micro.Flags(
			cli.StringFlag{
				Name:  "name",
				Usage: "dsl",
			},
			cli.StringFlag{
				Name:  "email",
				Usage: "m15576684447@163.com",
			},
			cli.StringFlag{
				Name:  "password",
				Usage: "xxx",
			},
			cli.StringFlag{
				Name:  "company",
				Usage: "bytedance",
			},
		),
	)

	service.Init(
		micro.Action(func(c *cli.Context) {
			name := c.String("name")
			email := c.String("email")
			password := c.String("password")
			company := c.String("company")

			r, err := client.Create(context.TODO(), &pb.User{
				Name:     name,
				Email:    email,
				Password: password,
				Company:  company,
			})
			if err != nil {
				log.Fatalf("Could not create: %v", err)
			}
			log.Printf("Created: %v", r.User.Id)

			getAll, err := client.GetAll(context.Background(), &pb.Request{})
			if err != nil {
				log.Fatalf("Could not list users: %v", err)
			}
			for _, v := range getAll.Users {
				log.Println(v)
			}

			os.Exit(0)
		}),
	)

	// 启动客户端
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
