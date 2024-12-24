import networkx as nx


with open("input-0.txt") as fp:
    connections = fp.read().split("\n")


m = {}
for c in connections:
    nodes = c.split("-")
    for n in nodes:
        m[n] = True

nodes = list(m.keys())


g = nx.Graph()
g.add_nodes_from(nodes)

for c in connections:
    parts = c.split("-")
    g.add_edge(parts[0], parts[1])


nx.draw(g)
