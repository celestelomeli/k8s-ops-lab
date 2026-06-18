"""Local development architecture (Docker Compose).

Render from the repo root:
    python docs/diagrams/local.py
Outputs: inventory-api/docs/local.png
"""
from diagrams import Diagram, Cluster, Edge
from diagrams.programming.language import Go
from diagrams.onprem.database import Postgresql
from diagrams.generic.storage import Storage

graph_attr = {"fontsize": "16", "bgcolor": "white"}

with Diagram(
    "Local Development (Docker Compose)",
    filename="inventory-api/docs/local",
    show=False,
    direction="LR",
    graph_attr=graph_attr,
):
    with Cluster("Laptop  -  Docker Compose"):
        api = Go("inventory-api\n:8080")
        db = Postgresql("postgres:16\n:5432")
        vol = Storage("postgres-data\nvolume")

        api >> Edge(label="DB_HOST=postgres") >> db
        db >> Edge(label="persists", style="dashed") >> vol
