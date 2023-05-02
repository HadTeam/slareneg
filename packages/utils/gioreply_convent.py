#!/bin/env python3
# -*- coding: utf-8 -*-

import os
import sys
import json
import numpy as np

# usage: python3 convert.py [json file] [target file]
if len(sys.argv) < 3:
    print("usage: python3 convert.py [json file] [target file]")

with open(sys.argv[1], 'r') as f:
    data = json.load(f)


w = data['mapWidth']
h = data['mapHeight']

def getPos(blockNum):
    return [blockNum % w, int(blockNum/w)]

m = np.zeros((w,h), dtype=np.uint)

for city in data['cities']:
    pos=getPos(city)
    m[pos[0]][pos[1]]=3
for general in data['generals']:
    pos=getPos(general)
    m[pos[0]][pos[1]]=2

moves = data['moves']
usernames = data['usernames']

moves_by_user = {}
for i in range(len(usernames)):
    moves_by_user[i] = []

for move in moves:
    user_index = move['index']
    username = usernames[user_index]
    prev_location = getPos(move['start'])
    after_location = getPos(move['end'])
    if prev_location[0] > after_location[0]:
        toward = 'left'
    elif prev_location[0] < after_location[0]:
        toward = 'right'
    elif prev_location[1] > after_location[1]:
        toward = 'up'
    elif prev_location[1] < after_location[1]:
        toward = 'down'
    else:
        raise Exception('Error')
    if int(move['is50']):
        number = 65535
    else:
        number = 0
    move_str = f'Move {prev_location[0]} {prev_location[1]} {toward} {number}'
    moves_by_user[user_index].append(move_str)

with open(sys.argv[2], 'w', newline='\n') as f:
    f.write(f'[\n')
    for i,col in enumerate(m):
        f.write(f'  [')
        for j,b in enumerate(col):
            if j!=len(col)-1:
                e=','
            else:
                e=''
            f.write(f'{b}{e}')
        if i!=len(m)-1:
            e=','
        else:
            e=''
        f.write(f']{e}\n')
    f.write(f']\n')
    for i in range(len(usernames)):
        username = usernames[i]
        user_moves = moves_by_user[i]
        if len(user_moves) > 0:
            f.write(f'|{username}:\n')
            f.write('\n'.join(user_moves))
            f.write('\n')
