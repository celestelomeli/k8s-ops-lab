"""Managed Kubernetes on AWS (EKS) with CI/CD.

Render from the repo root:
    python docs/diagrams/eks.py
Outputs: docs/eks.png
"""
from diagrams import Diagram, Cluster, Edge
from diagrams.aws.compute import EKS, EC2, ECR
from diagrams.aws.network import ELB, NATGateway, InternetGateway
from diagrams.aws.security import SecretsManager
from diagrams.aws.storage import EBS
from diagrams.aws.management import Cloudwatch
from diagrams.k8s.compute import Deployment, Pod
from diagrams.k8s.clusterconfig import HPA
from diagrams.onprem.gitops import ArgoCD
from diagrams.onprem.ci import GithubActions
from diagrams.onprem.vcs import Github
from diagrams.onprem.client import Users

graph_attr = {"fontsize": "16", "bgcolor": "white"}

with Diagram(
    "Managed Kubernetes on AWS (EKS) with CI/CD",
    filename="docs/eks",
    show=False,
    direction="LR",
    graph_attr=graph_attr,
):
    dev = Users("developer")

    # ---- CI/CD pipeline (outside the cluster) ----
    with Cluster("CI / CD"):
        repo = Github("GitHub repo\nmanifests")
        ci = GithubActions("GitHub Actions\nbuild + push")
        ecr = ECR("ECR\nimages")

    # ---- AWS managed services ----
    sm = SecretsManager("Secrets Manager\nDB password")
    cw = Cloudwatch("CloudWatch\nContainer Insights")
    eks_cp = EKS("EKS control plane\n(managed, IP-locked)")

    # ---- VPC ----
    with Cluster("VPC 10.0.0.0/16  (2 AZs)"):
        igw = InternetGateway("IGW")
        nat = NATGateway("NAT")

        with Cluster("Public subnets"):
            elb = ELB("LoadBalancer\ninventory-api")

        with Cluster("Private subnets - EKS node group (t3.medium x2)"):
            argocd = ArgoCD("ArgoCD\n(GitOps CD)")
            eso = Pod("External Secrets\nOperator (IRSA)")

            with Cluster("namespace: inventory"):
                api = Deployment("inventory-api\nx2 replicas")
                hpa = HPA("HPA")
                notifier = Deployment("notifier")
                pg = Deployment("postgres")
                ebs = EBS("EBS volume\n(PVC)")

    # ---- flows ----
    dev >> Edge(label="git push") >> repo >> Edge(label="trigger") >> ci
    ci >> Edge(label="OIDC, no static keys") >> ecr
    argocd >> Edge(label="watch + sync (pull)", style="bold") >> repo
    ecr >> Edge(label="image pull") >> api

    igw >> elb >> Edge(label=":80") >> api
    api >> Edge(label="5432") >> pg >> Edge(label="persist") >> ebs
    hpa >> Edge(label="scale pods 2-5", style="dashed") >> api
    eso >> Edge(label="GetSecretValue") >> sm
    argocd >> Edge(style="dashed") >> api

    eks_cp >> Edge(style="dotted") >> argocd
    api >> Edge(label="logs + metrics", style="dotted") >> cw
