import json

obj = None
with open("cntsPyCovering.json", "r") as fd:
    obj = json.loads(fd.read())
a = obj
with open("cntsTsCovering.json", "r") as fd:
    obj = json.loads(fd.read())
b = obj
d = {}
for [k, v] in a:
    d[k] = [v, 0, 0]
for [k, v] in b:
    d[k][1] = v

for k in d:
    d[k][2] = d[k][1] - d[k][0]
items = list(d.items())
items.sort(key=lambda e: e[1][2])
for [k, e] in items:
    if e[2] != 0:
        print(k, e)
