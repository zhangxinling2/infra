package lb

import (
	"fmt"
	"github.com/tietang/go-eureka-client/eureka"
	"strings"
)

type Apps struct {
	client *eureka.Client
}

func (a *Apps) Get(appName string) *App {
	var app eureka.Application
	for _, a := range a.client.Applications.Applications {
		if a.Name == strings.ToUpper(appName) {
			app = a
			break
		}
	}
	na := &App{
		Name:      app.Name,
		Instances: make([]*ServerInstance, 0),
		lb:        nil,
	}
	for _, instance := range app.Instances {
		var port int
		if instance.SecurePort.Enabled {
			port = instance.SecurePort.Port
		} else {
			port = instance.Port.Port
		}
		si := &ServerInstance{
			InstanceId: instance.InstanceId,
			AppName:    appName,
			Address:    fmt.Sprintf("%s:%d", instance.IpAddr, port),
			Status:     Status(instance.Status),
			Metadata:   make(map[string]string),
		}
		si.Metadata["rpcAddr"] = fmt.Sprintf("%s:%s", instance.IpAddr, instance.Metadata.Map["rpcPort"])
		na.Instances = append(na.Instances, si)
	}
	return na
}

type App struct {
	Name      string
	Instances []*ServerInstance
	lb        Balancer
}

func (a *App) Get(key string) *ServerInstance {
	ins := a.lb.Next(key, a.Instances)
	return ins
}

type Status string

const (
	StatusEnabled  Status = "enabled"
	StatusDisabled Status = "disabled"
)

type ServerInstance struct {
	InstanceId string
	AppName    string
	Address    string
	Status     Status
	Metadata   map[string]string
}
