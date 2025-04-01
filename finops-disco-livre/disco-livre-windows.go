package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Instance struct {
	Name string
	Zone string
}

type Disk struct {
	Server string
	Mount  string
	Total  string
	Used   string
	Free   string
}

type DiskResult struct {
	Instance Instance
	Disks    []Disk
	Error    error
}

func runCommand(cmd string, description string) (string, error) {
	fmt.Printf("\nExecuting: %s\n", description)
	fmt.Printf("Command: %s\n", cmd)

	// Create a context with timeout
	// Use 120 seconds for SSH commands, 30 seconds for other commands
	timeout := 30 * time.Second
	if strings.Contains(cmd, "ssh") {
		timeout = 120 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Use shell to execute the command to handle quotes and special characters
	command := exec.CommandContext(ctx, "sh", "-c", cmd)
	command.Env = os.Environ()

	output, err := command.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("command timed out after %v seconds", timeout.Seconds())
		}
		fmt.Printf("Error executing command: %v\n", err)
		return "", err
	}

	return string(output), nil
}

func getProjectInstances() ([]Instance, error) {
	fmt.Println("\nStep 1: Getting project ID")
	projectOutput, err := runCommand("gcloud config get-value project", "Getting project ID")
	if err != nil {
		return nil, err
	}
	// Debug output
	fmt.Printf("Raw projectOutput: %q\n", projectOutput)

	// Split by newline and remove empty strings
	lines := strings.Split(strings.TrimSpace(projectOutput), "\n")
	fmt.Printf("Lines after split: %q\n", lines)

	// Get the last non-empty line
	var project string
	for i := len(lines) - 1; i >= 0; i-- {
		if trimmed := strings.TrimSpace(lines[i]); trimmed != "" {
			project = trimmed
			break
		}
	}

	if project == "" {
		return nil, fmt.Errorf("could not parse project ID from output")
	}

	fmt.Printf("Working with project: %s\n", project)

	fmt.Println("\nStep 2: Getting instances list")
	cmd := fmt.Sprintf("gcloud compute instances list --project=%s --format=json --quiet", project)
	instancesOutput, err := runCommand(cmd, "Getting instances list")
	if err != nil {
		return nil, err
	}

	var instances []map[string]interface{}
	if err := json.Unmarshal([]byte(instancesOutput), &instances); err != nil {
		return nil, err
	}

	var linuxInstances []Instance
	for _, instance := range instances {
		name := instance["name"].(string)
		zone := strings.Split(instance["zone"].(string), "/")
		zoneStr := zone[len(zone)-1]

		// Skip Windows instances
		isWindows := false
		if disks, ok := instance["disks"].([]interface{}); ok {
			for _, disk := range disks {
				if diskMap, ok := disk.(map[string]interface{}); ok {
					if licenses, ok := diskMap["licenses"].([]interface{}); ok {
						for _, license := range licenses {
							if strings.Contains(strings.ToLower(license.(string)), "windows") {
								isWindows = true
								break
							}
						}
					}
				}
			}
		}

		if !isWindows {
			linuxInstances = append(linuxInstances, Instance{Name: name, Zone: zoneStr})
		}
	}

	fmt.Printf("Filtered %d Linux instances\n", len(linuxInstances))
	return linuxInstances, nil
}

func getLinuxDiskUsage(instance Instance) []Disk {
	fmt.Printf("\nStep 3: Getting disk usage for %s\n", instance.Name)
	cmd := fmt.Sprintf("gcloud compute ssh %s --zone=%s --tunnel-through-iap --command='df -B1' --quiet",
		instance.Name, instance.Zone)

	output, err := runCommand(cmd, fmt.Sprintf("SSH disk check on %s", instance.Name))
	if err != nil {
		fmt.Printf("Error connecting to %s: %v. Skipping to next instance.\n", instance.Name, err)
		return nil
	}

	return parseLinuxDf(output, instance.Name)
}

func parseLinuxDf(output string, instanceName string) []Disk {
	fmt.Printf("\nStep 4: Parsing df output for %s\n", instanceName)
	var disks []Disk
	lines := strings.Split(output, "\n")
	fmt.Printf("Found %d lines of output\n", len(lines))

	for _, line := range lines[1:] {
		if strings.HasPrefix(line, "/dev/") {
			parts := strings.Fields(line)
			if len(parts) >= 6 {
				disks = append(disks, Disk{
					Server: instanceName,
					Mount:  parts[5],
					Total:  parts[1],
					Used:   parts[2],
					Free:   parts[3],
				})
			}
		}
	}

	fmt.Printf("Successfully parsed %d disk entries\n", len(disks))
	return disks
}

// Modify main() to use goroutines
func main() {
	fmt.Println("Starting Linux disk usage collection...")
	instances, err := getProjectInstances()
	if err != nil {
		fmt.Printf("Error getting instances: %v\n", err)
		os.Exit(1)
	}

	// Create a channel to receive results
	results := make(chan DiskResult, len(instances))

	// Launch goroutine for each instance
	for _, instance := range instances {
		go func(inst Instance) {
			fmt.Printf("\nProcessing instance: %s\n", inst.Name)
			disks := getLinuxDiskUsage(inst)
			results <- DiskResult{
				Instance: inst,
				Disks:    disks,
			}
		}(instance)
	}

	// Collect results
	var allDisks []Disk
	for i := 0; i < len(instances); i++ {
		result := <-results
		if result.Disks != nil {
			allDisks = append(allDisks, result.Disks...)
		}
	}

	// Console output
	fmt.Println("\nGenerating final report...")
	fmt.Println("\nDisk Usage Report")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-30s%-10s%-22s%-22s%-22s\n", "SERVER", "MOUNT", "TOTAL", "USED", "FREE")
	fmt.Println(strings.Repeat("-", 80))

	for _, disk := range allDisks {
		fmt.Printf("%-30s%-22s%-22s%-22s%-22s\n",
			disk.Server, disk.Mount, disk.Total, disk.Used, disk.Free)
	}

	// CSV output
	filename := fmt.Sprintf("disk_usage_report_%s.csv",
		time.Now().Format("20060102_150405"))
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating CSV file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"SERVER", "MOUNT", "TOTAL", "USED", "FREE"})

	// Write data
	for _, disk := range allDisks {
		writer.Write([]string{
			disk.Server,
			disk.Mount,
			disk.Total,
			disk.Used,
			disk.Free,
		})
	}

	fmt.Printf("\nReport saved to %s\n", filename)
}
