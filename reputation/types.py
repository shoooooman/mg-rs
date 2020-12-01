import sys
import random
import json

args = sys.argv
if len(args) != 3:
    print('the length of args must be 3', file=sys.stderr)
    sys.exit(1)

if not args[1].isdecimal():
    print('please input the number of nodes', file=sys.stderr)
    sys.exit(1)
n = int(args[1])

if not args[2].isdecimal():
    print('please input the number of dishonest nodes', file=sys.stderr)
    sys.exit(1)
d = int(args[2])

if d > n:
    print('the number of dishonest nodes is too large')
    sys.exit(1)

l = range(n)
# exclude ID:0 from dishonest candinates
dl = random.sample(l[1:], d)
ds = set(dl)
print(ds)

# convert to dict
dic = {}
dic["types"] = []
for i in range(n):
    if i in ds:
        reverse = {"id": i, "type": "reverse"}
        dic["types"].append(reverse)
    else:
        honest = {"id": i, "type": "honest"}
        dic["types"].append(honest)

# convert to json and output to config.json
with open('config.json', 'w') as f:
    json.dump(dic, f, indent=4)
