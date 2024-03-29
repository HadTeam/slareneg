#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import sys
import json
import numpy as np

# usage: python3 convert.py [json file] [target file]
if len(sys.argv) < 3:
    print("usage: python3 gioreply_convent.py [json file] [target file]")
    sys.exit()

with open(sys.argv[1], 'r') as json_file:
    data = json.load(json_file)

width = data['mapWidth']
height = data['mapHeight']

def get_position(block_num):
    return [block_num % width, int(block_num/width)]

map_data = np.zeros((height, width, 3), dtype=np.uint)

for i, city in enumerate(data['cities']):
    pos = get_position(city)
    map_data[pos[1]][pos[0]] = [3, 0, data['cityArmies'][i]]

for u, general in enumerate(data['generals']):
    pos = get_position(general)
    map_data[pos[1]][pos[0]] = [2, u + 1, 0]

moves = data['moves']
usernames = data['usernames']

moves_by_user = {}
turns_by_user = {}

for i in range(len(usernames)):
    moves_by_user[i] = []
    turns_by_user[i] = 1

for move in moves:
    user_index = move['index']
    username = usernames[user_index]
    prev_location = get_position(move['start'])
    after_location = get_position(move['end'])
    if prev_location[0] > after_location[0]:
        direction = 'left'
    elif prev_location[0] < after_location[0]:
        direction = 'right'
    elif prev_location[1] > after_location[1]:
        direction = 'up'
    elif prev_location[1] < after_location[1]:
        direction = 'down'
    else:
        raise Exception('Error')
    if int(move['is50']):
        move_number = 65535
    else:
        move_number = 0
    move_str = f'Move {prev_location[0]+1} {prev_location[1]+1} {direction} {move_number}'
    while turns_by_user[user_index] < move['turn'] -1:
        turns_by_user[user_index] += 1
        moves_by_user[user_index].append("")
    turns_by_user[user_index] += 1
    moves_by_user[user_index].append(move_str)

with open(sys.argv[2], 'w', newline='\n') as target_file:
    json.dump(map_data.tolist(), target_file)
    target_file.write(f'\n')
    for i in range(len(usernames)):
        username = usernames[i]
        user_moves = moves_by_user[i]
        if len(user_moves) > 0:
            target_file.write(f'|{username}:\n')
            target_file.write(f'\n'.join(user_moves))
            target_file.write('\n')