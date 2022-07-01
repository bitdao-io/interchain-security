import json

o = None
with open("trace_BAD.json", "r") as fd:
    o = json.loads(fd.read())

transitions = o["transitions"]
actions = [t["action"] for t in transitions]
cons = [t["consequence"] for t in transitions]
blocks = o["blocks"]
pblocks = blocks["provider"]
cblocks = blocks["consumer"]
for a in actions:
    print(a)
pblocks = {int(p[0]): p[1] for p in pblocks}
cblocks = {int(p[0]): p[1] for p in cblocks}
print(pblocks[14])
print(cblocks[0])
