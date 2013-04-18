package main

import (
	"fmt"
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
	session, err := mgo.Dial(address)
	if err != nil {
		return nil, err
	}
	m.Session = session
	return m, nil
}

func (d *Mongo) Close() {
	d.Session.Close()
}

func (d *Mongo) AddServer(server string) error {
	if d.Session == nil {
		return fmt.Errorf("database session not established")
	}
	c := d.Session.DB(d.Database).C("servers")
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

func (d *Mongo) GetServers() ([]Server, error) {
	if d.Session == nil {
		return nil, fmt.Errorf("database session not established")
	}
	c := d.Session.DB(d.Database).C("servers")
	var result []Server
	err := c.Find(nil).Limit(1000).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
