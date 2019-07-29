package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	pb "learn/shippy/src/vessel-service/proto/vessel"
	"log"
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
	var v *pb.Vessel
	err := repo.collection().Find(bson.M{
		"capacity":  bson.M{"$gte": spec.Capacity},
		"maxweight": bson.M{"$gte": spec.MaxWeight},
	}).One(&v)
	if err != nil {
		log.Printf("FindAvailable Vessels from mongodb err: %s\n", err)
		return nil, err
	}
	return v, nil
}

func (repo *VesselRepository) Create(v *pb.Vessel) error {
	log.Printf("mongodb: Create\n")
	return repo.collection().Insert(v)
}

func (repo *VesselRepository) Close() {
	log.Printf("mongodb: Close\n")
	repo.session.Close()
}

func (repo *VesselRepository) collection() *mgo.Collection {
	return repo.session.DB(DB_NAME).C(VESSEL_COLLECTION)
}
