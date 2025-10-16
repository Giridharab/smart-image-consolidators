package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/moby/moby/api/types"
    "github.com/moby/moby/api/types/container"
	"github.com/docker/docker/client"
	"github.com/moby/moby/client"
)

type PerfMetric struct {
	CPUUsage      string
	MemoryUsage   string
	Storage       string
	EstimatedCost float64
}

// MeasurePerformanceAndCost runs a temporary container from the image,
// collects CPU/Memory stats, calculates storage size, and estimates cost.
func MeasurePerformanceAndCost(image string) (PerfMetric, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return PerfMetric{}, err
	}
	defer cli.Close()

	// Pull image
	_, err = cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return PerfMetric{}, fmt.Errorf("failed to pull image: %v", err)
	}

	// Create temporary container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{"sleep", "10"}, // short-lived container
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		return PerfMetric{}, fmt.Errorf("failed to create container: %v", err)
	}
	containerID := resp.ID

	defer func() {
		cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
	}()

	// Start container
	if err := cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return PerfMetric{}, fmt.Errorf("failed to start container: %v", err)
	}

	// Give container a few seconds to initialize
	time.Sleep(3 * time.Second)

	// Collect stats (stream=false gives a single snapshot)
	stats, err := cli.ContainerStatsOneShot(ctx, containerID)
	if err != nil {
		return PerfMetric{}, fmt.Errorf("failed to get container stats: %v", err)
	}
	defer stats.Body.Close()

	// Decode stats JSON
	var s types.StatsJSON
	err = json.NewDecoder(stats.Body).Decode(&s)
	if err != nil {
		return PerfMetric{}, fmt.Errorf("failed to decode stats: %v", err)
	}

	// Inspect image to get storage size
	inspect, _, err := cli.ImageInspectWithRaw(ctx, image)
	if err != nil {
		return PerfMetric{}, fmt.Errorf("failed to inspect image: %v", err)
	}

	// Convert values
	cpu := fmt.Sprintf("%.2f%%", calculateCPUPercent(&s))
	mem := fmt.Sprintf("%.2fMiB", float64(s.MemoryStats.Usage)/(1024*1024))
	storage := fmt.Sprintf("%.2fMiB", float64(inspect.Size)/(1024*1024))
	estimatedCost := float64(inspect.Size)/(1024*1024*1024) * 0.1 // $0.1 per GB

	return PerfMetric{
		CPUUsage:      cpu,
		MemoryUsage:   mem,
		Storage:       storage,
		EstimatedCost: estimatedCost,
	}, nil
}

// calculateCPUPercent computes the CPU usage percentage
func calculateCPUPercent(s *types.StatsJSON) float64 {
	cpuDelta := float64(s.CPUStats.CPUUsage.TotalUsage - s.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(s.CPUStats.SystemUsage - s.PreCPUStats.SystemUsage)
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		numCPU := float64(len(s.CPUStats.CPUUsage.PercpuUsage))
		return (cpuDelta / systemDelta) * numCPU * 100.0
	}
	return 0.0
}
