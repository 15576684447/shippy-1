package main

import (
	"context"
	"gopkg.in/mgo.v2"
	pb "learn/shippy/src/consignment-service/proto/consignment"
	vesselPb "learn/shippy/src/vessel-service/proto/vessel"
	"log"
)

// 微服务服务端 struct handler 必须实现 protobuf 中定义的 rpc 方法
// 实现方法的传参等可参考生成的 consignment.pb.go
type handler struct {
	session *mgo.Session
	vesselClient vesselPb.VesselServiceClient
}

// 从主会话中 Clone() 出新会话处理查询
func (h *handler)GetRepo()Repository  {
	return &ConsignmentRepository{h.session.Clone()}
}

func (h *handler)CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) error {
	defer h.GetRepo().Close()
	log.Printf("Called by consignment-cli to Create Consignment\n")
	// 检查是否有适合的货轮
	vReq := &vesselPb.Specification{
		Capacity:  int32(len(req.Containers)),
		MaxWeight: req.Weight,
	}
	log.Printf("Call vessel-service find available vessel")
	vResp, err := h.vesselClient.FindAvailable(context.Background(), vReq)
	if err != nil {
		log.Printf("Consignment-service Find Available from Vessel-service err: %s\n", err)
		return err
	}

	// 货物被承运
	log.Printf("found vessel: %s\n", vResp.Vessel.Name)
	req.VesselId = vResp.Vessel.Id
	//consignment, err := h.repo.Create(req)
	err = h.GetRepo().Create(req)
	if err != nil {
		log.Printf("Consignment-service Create Vessel Info err: %s\n", err)
		return err
	}
	resp.Created = true
	resp.Consignment = req
	return nil
}

func (h *handler)GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error {
	defer h.GetRepo().Close()
	log.Printf("Called by consignment-cli to Get Consignments\n")
	consignments, err := h.GetRepo().GetAll()
	if err != nil {
		log.Printf("Consignment-service GetAll Vessel Info err: %s\n", err)
		return err
	}
	resp.Consignments = consignments
	return nil
}