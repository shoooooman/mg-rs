import sys
import json

def isfloat(s):
    try:
        float(s)
    except:
        return False
    return True

args = sys.argv
if len(args) != 5:
    print('the length of args is too small', file=sys.stderr)
    sys.exit(1)

if not args[1].isdecimal():
    print('please input the number of nodes', file=sys.stderr)
    sys.exit(1)
n = int(args[1])

if not isfloat(args[2]):
    print('please input the ratio of normal nodes', file=sys.stderr)
    sys.exit(1)
normal_ratio = float(args[2])
if not (0 <= normal_ratio <= 1.0):
    print('the ratio of normal nodes must be between 0.0 and 1.0', file=sys.stderr)
    sys.exit(1)

if not isfloat(args[3]):
    print('please input the ratio of failure nodes', file=sys.stderr)
    sys.exit(1)
failure_ratio = float(args[3])
if not (0 <= failure_ratio <= 1.0):
    print('the ratio of failure nodes must be between 0.0 and 1.0', file=sys.stderr)
    sys.exit(1)

if not isfloat(args[4]):
    print('please input the ratio of malicious nodes', file=sys.stderr)
    sys.exit(1)
malicious_ratio = float(args[4])
if normal_ratio+malicious_ratio+failure_ratio != 1:
    print('the sum of three ratio should be 1', file=sys.stderr)
    sys.exit(1)

num_normal = int(normal_ratio*n)
num_malicious = int(malicious_ratio*n)
num_failure = n - num_normal - num_malicious
print(num_normal, num_failure, num_malicious)

# first-half: 0.05, second-half: 0.95
def get_var_probs(num_tx):
    first = {"left": 0, "right": num_tx//2, "value": 0.05}
    second = {"left": num_tx//2, "right": num_tx, "value": 0.95}
    return [first, second]

num_tx = 1000
var_probs = get_var_probs(num_tx)

# convert to dict
d = {}
d["behaviors"] = []
for i in range(n):
    if 0 <= i < num_normal:
        normal = {"id": i, "kind": "fixed", "probability": 0.0}
        d["behaviors"].append(normal)
    elif num_normal <= i < num_normal+num_failure:
        tick = 1.0/num_failure
        prob = tick*(i-num_normal+1)
        failure = {"id": i, "kind": "fixed", "probability": prob}
        d["behaviors"].append(failure)
    else:
        malicious = {"id": i, "kind": "variable", "var_probs": var_probs}
        d["behaviors"].append(malicious)

# convert to json and output to config.json
with open('config.json', 'w') as f:
    json.dump(d, f, indent=4)
