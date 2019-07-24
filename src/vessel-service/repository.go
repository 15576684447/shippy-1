package main

import (
	"github.com/micro/go-micro/util/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	pb "learn/shippy/src/vessel-service/proto/vessel"
)

const (
	DB_NAME           = "vessels"
	VESSEL_COLLECTION = "vessels"
)

type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
	Create(*pb.Vessel) error
	Close()
}

type VesselRepository struct {
	session *mgo.Session
}

// 接口实现
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	// 选择最近一条容量、载重都符合的货轮
	log.Infof("Called by Client to Find available vessel, capacity: %d, weight: %d", spec.Capacity, spec.MaxWeight)
	var v *pb.Vessel
	err := repo.collection().Find(bson.M{
		"capacity":  bson.M{"$gte": spec.Capacity},
		"maxweight": bson.M{"$gte": spec.MaxWeight},
	}).One(&v)
	if err != nil {
		log.Info("FindAvaliable error")
		return nil, err
	}
	return v, nil
}

func (repo *VesselRepository) Create(v *pb.Vessel) error {
	return repo.collection().Insert(v)
}

func (repo *VesselRepository) Close() {
	repo.session.Close()
}

func (repo *VesselRepository) collection() *mgo.Collection {
	return repo.session.DB(DB_NAME).C(VESSEL_COLLECTION)
}