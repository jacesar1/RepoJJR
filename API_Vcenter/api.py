import requests

vcenter_url = "https://elb1180.int.eletrobras.gov.br/rest/vcenter"
username = "jacesar1@int.eletrobras.gov.br"
password = "##7Ujm,kio"

# Autenticar e obter sessão
session = requests.post(
    f"{vcenter_url}/com/vmware/cis/session",
    auth=(username, password),
    verify=False  # Desativar verificação SSL (não recomendado em produção)
)
session_id = session.json()["value"]

# Listar VMs
vms = requests.get(
    f"{vcenter_url}/vm",
    headers={"vmware-api-session-id": session_id},
    verify=False
)
print(vms.json())