package main

import (
	"context"
	"gopkg.in/mgo.v2"
	pb "learn/shippy/src/vessel-service/proto/vessel"
	"log"
)

// 实现微服务的服务端
type handler struct {
	session *mgo.Session
}

func (h *handler) GetRepo() Repository {
	return &VesselRepository{h.session.Clone()}
}

func (h *handler) FindAvailable(ctx context.Context, req *pb.Specification, resp *pb.Response) error {
	defer h.GetRepo().Close()
	log.Printf("Called by consignment-service to Find Available")
	log.Printf("Call mongodb to find available data")
	v, err := h.GetRepo().FindAvailable(req)
	if err != nil {
		return err
	}
	resp.Vessel = v
	return nil
}

func (h *handler) Create(ctx context.Context, req *pb.Vessel, resp *pb.Response) error {
	defer h.GetRepo().Close()
	log.Printf("Called by consignment-service to Create")
	log.Printf("Call mongodb to Create")
	if err := h.GetRepo().Create(req); err != nil {
		return err
	}
	resp.Vessel = req
	resp.Created = true
	return nil
}
