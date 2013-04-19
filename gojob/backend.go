package main

import (
	"github.com/retzkek/gojob"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type Server struct {
	Address   string
	Status    string
	Load      gojob.Load
	Processes []gojob.Process
	Timestamp time.Time
}

type Backend interface {
	AddNewServer(server string) error
	AddServer(server Server) error
	GetServers() ([]Server, error)
}

type Mongo struct {
	Address  string
	Database string
}

func InitMongo(address string, database string) (*Mongo, error) {
	m := new(Mongo)
	m.Address = address
	m.Database = database
	return m, nil
}

func (d *Mongo) AddNewServer(server string) error {
	session, err := mgo.Dial(d.Address)
	if err != nil {
		return err
	}
	defer session.Close()
	c := session.DB(d.Database).C("servers")
	q := c.Find(bson.M{"address": server})
	cnt, err := q.Count()
	if err != nil {
		return err
	}
	if cnt == 0 {
		err = c.Insert(&Server{server, "unknown", gojob.Load{0, 0, 0}, nil, time.Now()})
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Mongo) AddServer(server Server) error {
	session, err := mgo.Dial(d.Address)
	if err != nil {
		return err
	}
	defer session.Close()
	c := session.DB(d.Database).C("servers")
	q := c.Find(bson.M{"address": server})
	cnt, err := q.Count()
	if err != nil {
		return err
	}
	if cnt == 0 {
		err = c.Insert(server)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Mongo) GetServers() ([]Server, error) {
	session, err := mgo.Dial(d.Address)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	c := session.DB(d.Database).C("servers")
	var result []Server
	err = c.Find(nil).Limit(1000).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
