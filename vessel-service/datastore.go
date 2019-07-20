package main

import "github.com/go-mgo/mgo"

func CreateSession(host string) (*mgo.Session, error) {
	s, err := mgo.Dial(host)
	if err != nil {
		return nil, err
	}
	s.SetMode(mgo.Monotonic, true)
	return s, nil
}
