import sys
import networkx as nx
import json
import pprint

def isfloat(s):
    try:
        float(s)
    except:
        return False
    return True

args = sys.argv
if len(args) != 3:
    print('the length of args is too small', file=sys.stderr)
    sys.exit(1)

if not args[1].isdecimal():
    print('please input the number of nodes', file=sys.stderr)
    sys.exit(1)
n = int(args[1])

if not isfloat(args[2]):
    print('please input the probability of generating edges', file=sys.stderr)
    sys.exit(1)
p = float(args[2])

# generate a random graph
G = nx.fast_gnp_random_graph(n, p)

# network setting
host = "127.0.0.1"
base_port = 10000

# convert to dict
d = {}
d["nodes"] = []
for i in range(n):
    addr = host + ":" + str(base_port + i)
    peers = list(G.adj[i])
    node = {"id": i, "address": addr, "peers": peers}
    d["nodes"].append(node)

# convert to json and output to config.json
with open('config.json', 'w') as f:
    json.dump(d, f, indent=4)
