package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"google.golang.org/grpc"
)

const (
	socketPath = "/var/lib/kubelet/device-plugins/metal-gpu.sock"
	resourceName = "quest1.io/gpu"
)

type MetalDevicePlugin struct {
	server *grpc.Server
}

func NewMetalDevicePlugin() *MetalDevicePlugin {
	return &MetalDevicePlugin{}
}

// ListAndWatch - Returns available Metal GPUs
func (m *MetalDevicePlugin) ListAndWatch(_ *pluginapi.Empty, stream pluginapi.DevicePlugin_ListAndWatchServer) error {
	devices := []*pluginapi.Device{
		{ID: "metal-gpu-0", Health: pluginapi.Healthy},
	}
	stream.Send(&pluginapi.ListAndWatchResponse{Devices: devices})
	for {
		time.Sleep(10 * time.Second)
		stream.Send(&pluginapi.ListAndWatchResponse{Devices: devices})
	}
}

// Allocate - Assigns Metal GPU to a container
func (m *MetalDevicePlugin) Allocate(_ context.Context, req *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	responses := []*pluginapi.ContainerAllocateResponse{}
	for range req.ContainerRequests {
		responses = append(responses, &pluginapi.ContainerAllocateResponse{})
	}
	return &pluginapi.AllocateResponse{ContainerResponses: responses}, nil
}

// GetDevicePluginOptions - Returns plugin options (default: no special options)
func (m *MetalDevicePlugin) GetDevicePluginOptions(_ context.Context, _ *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
    return &pluginapi.DevicePluginOptions{}, nil
}

func (m *MetalDevicePlugin) GetPreferredAllocation(_ context.Context, req *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
    // Create a response object
    response := &pluginapi.PreferredAllocationResponse{
        ContainerResponses: []*pluginapi.ContainerPreferredAllocationResponse{},
    }

    // Loop through each container request
    for _, containerReq := range req.ContainerRequests {
        containerResponse := &pluginapi.ContainerPreferredAllocationResponse{
            DeviceIDs: containerReq.AvailableDeviceIDs, // Use the available device list
        }
        response.ContainerResponses = append(response.ContainerResponses, containerResponse)
    }

    return response, nil
}


// PreStartContainer - Called before the container starts (Metal GPUs don't need special setup)
func (m *MetalDevicePlugin) PreStartContainer(_ context.Context, _ *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
    return &pluginapi.PreStartContainerResponse{}, nil
}



func (m *MetalDevicePlugin) Register() error {
	conn, err := grpc.Dial("unix:///var/lib/kubelet/device-plugins/kubelet.sock", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pluginapi.NewRegistrationClient(conn)
	_, err = client.Register(context.Background(), &pluginapi.RegisterRequest{
		Version:      pluginapi.Version,
		Endpoint:     filepath.Base(socketPath),
		ResourceName: resourceName,
	})
	return err
}

func (m *MetalDevicePlugin) Start() error {
	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}

	m.server = grpc.NewServer()
	pluginapi.RegisterDevicePluginServer(m.server, m)
	go m.server.Serve(listener)

	return m.Register()
}

func main() {
	plugin := NewMetalDevicePlugin()
	if err := plugin.Start(); err != nil {
		log.Fatalf("Failed to start Metal Device Plugin: %v", err)
	}
	fmt.Println("Metal Device Plugin running...")
	select {}
}


