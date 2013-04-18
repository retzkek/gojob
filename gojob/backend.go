package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type Server struct {
	Address   string
	Status    string
	Timestamp time.Time
}

type Backend interface {
	AddServer(server string) error
	GetServers() ([]Server, error)
}

type Mongo struct {
	Address  string
	Database string
	Session  *mgo.Session
}

func InitMongo(address string, database string) (*Mongo, error) {
	m := new(Mongo)
	m.Address = address
	m.Database = database
	return m, nil
}

func (d *Mongo) Close() {
	d.Session.Close()
}

func (m *Mongo) AddServer(server string) error {
	session, err := mgo.Dial(m.Address)
	if err != nil {
		return err
	}
	defer session.Close()
	c := session.DB(m.Database).C("servers")
	q := c.Find(bson.M{"address": server})
	cnt, err := q.Count()
	if err != nil {
		return err
	}
	if cnt == 0 {
		err = c.Insert(&Server{server, "unknown", time.Now()})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Mongo) GetServers() ([]Server, error) {
	session, err := mgo.Dial(m.Address)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	c := session.DB(m.Database).C("servers")
	var result []Server
	err = c.Find(nil).Limit(1000).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
