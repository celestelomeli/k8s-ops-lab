"""Self-managed Kubernetes on EC2 (kubeadm), built with Terraform + Ansible.

Render from the repo root:
    python docs/diagrams/self-managed.py
Outputs: cluster-internals/docs/self-managed.png
"""
from diagrams import Diagram, Cluster, Edge
from diagrams.aws.compute import EC2
from diagrams.onprem.iac import Terraform, Ansible

graph_attr = {"fontsize": "16", "bgcolor": "white"}

with Diagram(
    "Self-Managed Kubernetes on EC2 (kubeadm)",
    filename="cluster-internals/docs/self-managed",
    show=False,
    direction="LR",
    graph_attr=graph_attr,
):
    tf = Terraform("Terraform\nprovision infra")
    ans = Ansible("Ansible\nkubeadm setup")

    with Cluster("AWS  -  Default VPC (172.31.0.0/16)"):
        with Cluster("EC2 nodes (t3.medium)"):
            cp = EC2("control-plane\nkubeadm init")
            wk = EC2("worker\nkubeadm join")
            cp >> Edge(label="Calico CNI\napi-server :6443") >> wk

    tf >> Edge(label="create") >> cp
    tf >> Edge(label="create") >> wk
    ans >> Edge(label="configure", style="dashed") >> cp
    ans >> Edge(style="dashed") >> wk
