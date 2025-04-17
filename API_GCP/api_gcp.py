from google.cloud import compute_v1, resourcemanager_v3
from google.oauth2 import service_account
import os

# Configuração
SERVICE_ACCOUNT_JSON = "caminho/para/sua-chave.json"
ORGANIZATION_ID = "123456789012"  # Substitua pelo ID da sua organização

# Autenticação
credentials = service_account.Credentials.from_service_account_file(
    SERVICE_ACCOUNT_JSON,
    scopes=["https://www.googleapis.com/auth/cloud-platform"]
)

# 1. Listar todos os projetos da organização
resourcemanager_client = resourcemanager_v3.ProjectsClient(credentials=credentials)
request = resourcemanager_v3.ListProjectsRequest(
    parent=f"organizations/{ORGANIZATION_ID}"
)
projects = resourcemanager_client.list_projects(request=request)

# 2. Para cada projeto, listar VMs
compute_client = compute_v1.InstancesClient(credentials=credentials)

for project in projects:
    project_id = project.project_id  # Ex: "meu-projeto-123"
    
    try:
        # Listar todas as zonas (necessário para listar VMs)
        zones = compute_client.list_zones(project=project_id)
        
        for zone in zones:
            zone_name = zone.name  # Ex: "us-central1-a"
            instances = compute_client.list(project=project_id, zone=zone_name)
            
            print(f"\nProjeto: {project_id} | Zona: {zone_name}")
            for instance in instances:
                print(f"- {instance.name} ({instance.status})")
    
    except Exception as e:
        print(f"Erro no projeto {project_id}: {e}")