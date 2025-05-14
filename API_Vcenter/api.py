import requests
import urllib3
from pprint import pprint

# Suppress SSL warnings (only if absolutely necessary in test environment)
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

# Base URL should not include /rest/vcenter
vcenter_url = "https://elb1180.int.eletrobras.gov.br"
username = "jacesar1@int.eletrobras.gov.br"
password = "##8Ik,.lop"

try:
    # Authenticate and get session
    session = requests.post(
        f"{vcenter_url}/rest/com/vmware/cis/session",  # Updated endpoint path
        auth=(username, password),
        verify=False
    )
    session.raise_for_status()
    session_id = session.json().get("value")

    if not session_id:
        raise ValueError("Failed to get valid session ID")

    # List VMs
    vms = requests.get(
        f"{vcenter_url}/rest/vcenter/vm",
        headers={"vmware-api-session-id": str(session_id)},
        verify=False
    )
    vms.raise_for_status()
    
    # Get VM data and format it nicely
    vm_list = vms.json().get("value", [])
    print(f"\nTotal VMs found: {len(vm_list)}")
    
    for vm in vm_list:
        print("\nVM Details:")
        print(f"Name: {vm.get('name')}")
        print(f"Power State: {vm.get('power_state')}")
        print(f"CPU Count: {vm.get('cpu_count')}")
        print(f"Memory Size (MB): {vm.get('memory_size_MiB')}")

except requests.exceptions.RequestException as e:
    print(f"API request failed: {e}")
except ValueError as e:
    print(f"Error: {e}")
except Exception as e:
    print(f"Unexpected error: {e}")