package consul

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

type ConsulClient interface {
	RegisterService(registration *api.AgentServiceRegistration) error
	DeregisterService(serviceID string) error
	DiscoverServices(serviceName string, queryOptions *api.QueryOptions) ([]*api.CatalogService, error)
	DiscoverServiceByName(serviceName string) ([]*api.ServiceEntry, error)
}

type consulClient struct {
	client *api.Client
}

func NewConsulClient(address string) (*consulClient, error) {
	config := api.DefaultConfig()
	config.Address = address
	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %w", err)
	}
	return &consulClient{client: client}, nil
}

func (c *consulClient) RegisterService(registration *api.AgentServiceRegistration) error {
	if err := c.client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("failed to register service with Consul: %w", err)
	}

	log.Printf("Service %s successfully registered with Consul", registration.Name)
	return nil
}

func (c *consulClient) DeregisterService(serviceID string) error {
	if serviceID == "" {
		return fmt.Errorf("serviceID cannot be empty")
	}

	if err := c.client.Agent().ServiceDeregister(serviceID); err != nil {
		return fmt.Errorf("failed to deregister service from Consul: %w", err)
	}

	log.Printf("Service %s successfully deregistered from Consul", serviceID)
	return nil
}

func (c *consulClient) DiscoverServices(serviceName string, queryOptions *api.QueryOptions) ([]*api.CatalogService, error) {
	if serviceName == "" {
		return nil, fmt.Errorf("serviceName cannot be empty")
	}

	services, _, err := c.client.Catalog().Service(serviceName, "", queryOptions)
	if err != nil {
		return nil, fmt.Errorf("error discovering services: %w", err)
	}
	return services, nil
}

func (c *consulClient) DiscoverServiceByName(serviceName string) ([]*api.ServiceEntry, error) {
	if serviceName == "" {
		return nil, fmt.Errorf("serviceName cannot be empty")
	}

	specificService, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		log.Fatal(err)
	}

	return specificService, nil
}
