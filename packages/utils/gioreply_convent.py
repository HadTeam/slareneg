#!/bin/env python3
#-*- coding: utf-8 -*-

import os
import sys
import json

# 用法: python3 convert.py [json文件] [目标文件]
if len(sys.argv) < 3:
    print("用法: python3 convert.py [json文件] [目标文件]")

# 读取 json 文件
with open(sys.argv[1], 'r') as f:
    data = json.load(f)

moves = data['moves']

# 打开目标文件
with open(sys.argv[2], 'w') as f:
    # 逐行写入
    for move in moves:
        w = data['mapWidth']
        h = data['mapHeight']
        prev_location = [move['start'] % w, int(move['start'] / w)]
        after_location = [move['end'] % w, int(move['end'] / w)]
        if prev_location[0] > after_location[0]:
            toward = 'left'
            distance = prev_location[0] - after_location[0]
        elif prev_location[0] < after_location[0]:
            toward = 'right'
            distance = after_location[0] - prev_location[0]
        elif prev_location[1] > after_location[1]:
            toward = 'up'
            distance = prev_location[1] - after_location[1]
        elif prev_location[1] < after_location[1]:
            toward = 'down'
            distance = after_location[1] - prev_location[1]
        else:
            raise Exception('Error')
        f.write(f'Move {prev_location[0]} {prev_location[1]} {toward} {distance}\n')