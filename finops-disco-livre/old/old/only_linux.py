#!/usr/bin/env python3
import subprocess
import json
import sys

def run_command(cmd, description):
    print(f"\nExecuting: {description}")
    print(f"Command: {cmd}")
    try:
        # Increased timeout to 120 seconds for SSH commands
        timeout = 120 if 'ssh' in cmd else 30
        result = subprocess.run(cmd, shell=True, capture_output=True, text=True, timeout=timeout)
        print(f"Return code: {result.returncode}")
        print(f"Output: {result.stdout}")
        if result.stderr:
            print(f"Errors: {result.stderr}")
        return result
    except subprocess.TimeoutExpired:
        print(f"Command timed out after {timeout} seconds")
        raise
    except Exception as e:
        print(f"Error executing command: {str(e)}")
        raise

def get_project_instances():
    print("\nStep 1: Getting project ID")
    project_cmd = "gcloud config get-value project"
    project_result = run_command(project_cmd, "Getting project ID")
    project = project_result.stdout.strip()
    print(f"Working with project: {project}")
    
    print("\nStep 2: Getting instances list")
    cmd = f"gcloud compute instances list --project={project} --format=json --quiet"
    instances_result = run_command(cmd, "Getting instances list")
    
    try:
        instances = json.loads(instances_result.stdout or '[]')
        print(f"Successfully parsed JSON data. Found {len(instances)} total instances")
    except json.JSONDecodeError as e:
        print(f"JSON parsing error: {str(e)}")
        print(f"Raw output: {instances_result.stdout}")
        sys.exit(1)
    
    linux_instances = [
        {
            'name': instance['name'],
            'zone': instance['zone'].split('/')[-1]
        }
        for instance in instances
        if not any('windows' in license.lower() 
                  for disk in instance['disks'] 
                  for license in disk.get('licenses', []))
    ]
    
    print(f"Filtered {len(linux_instances)} Linux instances")
    return linux_instances
def get_linux_disk_usage(instance_name, zone):
    print(f"\nStep 3: Getting disk usage for {instance_name}")
    cmd = f"gcloud compute ssh {instance_name} --zone={zone} --tunnel-through-iap --command='df -B1' --quiet"

    try:
        result = run_command(cmd, f"SSH disk check on {instance_name}")
        return parse_linux_df(result.stdout, instance_name)
    except subprocess.TimeoutExpired:
        print(f"Timeout connecting to {instance_name}. Skipping to next instance.")
        return []
    except Exception as e:
        print(f"Error connecting to {instance_name}: {str(e)}. Skipping to next instance.")
        return []

def parse_linux_df(output, instance_name):
    print(f"\nStep 4: Parsing df output for {instance_name}")
    disks = []
    lines = output.splitlines()
    print(f"Found {len(lines)} lines of output")
    
    for line in lines[1:]:
        if line.startswith('/dev/'):
            parts = line.split()
            try:
                disks.append({
                    'server': instance_name,
                    'mount': parts[5],
                    'total': parts[1],
                    'used': parts[2],
                    'free': parts[3]
                })
            except IndexError as e:
                print(f"Error parsing line: {line}")
                print(f"Error details: {str(e)}")
    
    print(f"Successfully parsed {len(disks)} disk entries")
    return disks

def main():
    print("Starting Linux disk usage collection...")
    instances = get_project_instances()
    all_disks = []
    
    for instance in instances:
        print(f"\nProcessing instance: {instance['name']}")
        disks = get_linux_disk_usage(instance['name'], instance['zone'])
        all_disks.extend(disks)
    
    print("\nGenerating final report...")
    print("\nDisk Usage Report")
    print("-" * 80)
    print(f"{'SERVER':<30}{'MOUNT':<10}{'TOTAL':<22}{'USED':<22}{'FREE':<22}")
    print("-" * 80)
    
    for disk in all_disks:
        print(f"{disk['server']:<30}{disk['mount']:<10}{disk['total']:<22}{disk['used']:<22}{disk['free']:<22}")

    # CSV output
    import csv
    from datetime import datetime
    
    filename = f"disk_usage_report_{datetime.now().strftime('%Y%m%d_%H%M%S')}.csv"
    with open(filename, 'w', newline='') as csvfile:
        fieldnames = ['SERVER', 'MOUNT', 'TOTAL', 'USED', 'FREE']
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        
        writer.writeheader()
        for disk in all_disks:
            writer.writerow({
                'SERVER': disk['server'],
                'MOUNT': disk['mount'],
                'TOTAL': disk['total'],
                'USED': disk['used'],
                'FREE': disk['free']
            })
    
    print(f"\nReport saved to {filename}")

if __name__ == "__main__":
    main()
