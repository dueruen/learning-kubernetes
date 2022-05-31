from pandas import read_csv
from pathlib import Path
import matplotlib.pyplot as plt
from os import path
import numpy as np

Path("resources/").mkdir(parents=True, exist_ok=True)
names = ["wac", "niwac", "wfnc"]
unit = ["cpu", "memory"]
unitTitle = ["CPU / [cpu units]", "Memory / [GB]"]
titleText = ["With instrumentation", "Without instrumentation", "With flannel instead of Cilium"]
# frequencies = ["mfoz", "mfozz", "mfozzz"]

unitCount = 0
for unitName in unit:
  count = 0
  for name in names:
    print(name + " - " + unitName + " - " + titleText[count])
    filePath = '../data/fixed/' + name + '/metric/' + name +'-' + unitName + '.csv'

    if (path.exists(filePath) == False):
      print("Has no metrics matching path")
      continue

    cols = list(read_csv(filePath, nrows =1))
    series = read_csv(filePath, header=0, index_col=0, parse_dates=True, squeeze=True, usecols =[i for i in cols if i != 'max capacity']) #and ('consumer-wacmfozmstt-' in i or i == "Time")
    if unitName == 'memory':
      series = series.apply(lambda x: x / (10**9))
    
    plt.rcParams["figure.figsize"] = (20,15)

    series.plot.area(cmap=plt.get_cmap('tab20c'))

    plt.legend(bbox_to_anchor=(1.05, 1), loc=2, borderaxespad=0.)
    plt.xlabel("Time / [clock]")
    plt.ylabel(unitTitle[unitCount])
    plt.title(titleText[count] + " " + unitName + " usage")
    plt.savefig('resources/' + name + '-' + unitName + '-recources.pdf', bbox_inches='tight')

    count = count + 1
  unitCount = unitCount + 1