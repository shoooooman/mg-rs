import sys
import json
from os import path

args = sys.argv
if len(args) != 2:
    print('please input experiment number')
    sys.exit(1)

if not args[1].isdecimal():
    print('please input the number', file=sys.stderr)
    sys.exit(1)
n = int(args[1])

dirname = path.dirname(__file__)
configs = path.join(dirname, 'configs.json')
with open(configs) as f:
    d = json.load(f)

if len(d) <= n:
    print('the experiment number is wrong')
    sys.exit(1)

config = path.join(dirname, 'config.json')
with open(config, 'w') as f:
    json.dump(d[n], f, indent=4)
