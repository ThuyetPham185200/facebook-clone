package fanoutmanager

import (
	"feedservice/internal/core/fanoutmanager/workerpocessor"
	"feedservice/internal/infra/redisclient"
	"feedservice/internal/model"
)

type FanoutManager struct {
	fanoutworkers []workerpocessor.FanoutWorker
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DBNumber int
}

func NewFanoutManager(numWorkers int, redisclient_ *redisclient.RedisClient, newposteventqueue_ *chan model.NewPostEvent) *FanoutManager {
	m := FanoutManager{}
	m.fanoutworkers = make([]workerpocessor.FanoutWorker, numWorkers)
	return &m
}

func (m *FanoutManager) Start() error {
	for _, e := range m.fanoutworkers {
		err := e.Start()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *FanoutManager) Stop() error {
	for _, e := range m.fanoutworkers {
		err := e.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}
