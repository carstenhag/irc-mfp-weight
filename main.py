#!/usr/bin/env python3

import myfitnesspal

client = myfitnesspal.Client('moter8')

weightList = list(client.get_measurements('Weight').items())

weight = weightList[0][1]
print(weight)
