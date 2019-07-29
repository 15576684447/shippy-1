package main

import (
	"context"
	"errors"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/broker/nats"
	"golang.org/x/crypto/bcrypt"
	pb "learn/shippy/src/user-service/proto/user"
	"log"
)

const topic = "user.created"

type handler struct {
	repo         Repository
	tokenService Authable
	Publisher    micro.Publisher
	//PubSub      broker.Broker
}

func (h *handler) Create(ctx context.Context, req *pb.User, resp *pb.Response) error {
	log.Printf("Called by user-cli to Create user info [with password bcrypt]")
	// 哈希处理用户输入的密码
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("bcrypt.GenerateFromPassword err: %s\n", err)
		return err
	}
	req.Password = string(hashedPwd)
	if err := h.repo.Create(req); err != nil {
		log.Printf("Create data into database err: %s\n", err)
		return nil
	}
	resp.User = req

	// 发布带有用户所有信息的消息
	//if err := h.publishEvent(req); err != nil {
	//	return err
	//}

	log.Printf("Called by user-cli to Create user success, now publish event to notify email")
	if err := h.Publisher.Publish(ctx, req); err != nil {
		log.Printf("publish message err\n")
		return err
	}
	return nil
}

// 发送消息通知
//func (h *handler) publishEvent(user *pb.User) error {
//	body, err := json.Marshal(user)
//	if err != nil {
//		return err
//	}
//
//	msg := &broker.Message{
//		Header: map[string]string{
//			"id": user.Id,
//		},
//		Body: body,
//	}
//
//	// 发布 user.created topic 消息
//	if err := h.PubSub.Publish(topic, msg); err != nil {
//		log.Fatalf("[pub] failed: %v\n", err)
//	}
//	return nil
//}

func (h *handler) Get(ctx context.Context, req *pb.User, resp *pb.Response) error {
	log.Printf("Called by user-cli to Get user info")
	u, err := h.repo.Get(req.Id)
	if err != nil {
		log.Printf("Get user info from database err: %s\n", err)
		return err
	}
	resp.User = u
	return nil
}

func (h *handler) GetAll(ctx context.Context, req *pb.Request, resp *pb.Response) error {
	log.Printf("Called by user-cli to Get all user info")
	users, err := h.repo.GetAll()
	if err != nil {
		log.Printf("GetAll user info from database err: %s\n", err)
		return err
	}
	resp.Users = users
	return nil
}

func (h *handler) Auth(ctx context.Context, req *pb.User, resp *pb.Token) error {
	log.Printf("Called by user-cli to Auth user info")
	// 在 part3 中直接传参 &pb.User 去查找用户
	// 会导致 req 的值完全是数据库中的记录值
	// 即 req.Password 与 u.Password 都是加密后的密码
	// 将无法通过验证
	u, err := h.repo.GetByEmail(req.Email)
	if err != nil {
		log.Printf("Auth user info err: %s\n", err)
		return err
	}

	// 进行密码验证
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		log.Printf("Auth CompareHashAndPassword err: %s\n", err)
		return err
	}
	t, err := h.tokenService.Encode(u)
	if err != nil {
		log.Printf("Auth Token Encode err: %s\n", err)
		return err
	}
	resp.Token = t
	return nil
}

func (h *handler) ValidateToken(ctx context.Context, req *pb.Token, resp *pb.Token) error {
	log.Printf("Called by user-cli to ValidateToken")
	// Decode token
	claims, err := h.tokenService.Decode(req.Token)
	if err != nil {
		log.Printf("ValidateToken Token Decode err: %s\n", err)
		return err
	}
	if claims.User.Id == "" {
		return errors.New("invalid user")
	}

	resp.Valid = true
	return nil
}
