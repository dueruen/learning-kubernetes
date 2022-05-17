# import csv
# import matplotlib.pyplot as plt
# import numpy as np



# from pandas import read_csv
# from matplotlib import pyplot
# cols = list(read_csv('../data/metrics/all-wfnc-Memory.csv', nrows =1))
# series = read_csv('../data/metrics/all-wfnc-Memory.csv', header=0, index_col=0, parse_dates=True, squeeze=True, usecols =[i for i in cols if i != 'max capacity'])
# series.plot.area()
# pyplot.legend(bbox_to_anchor=(1.0, 1.0))
# pyplot.savefig('test.png')

from pandas import read_csv
from pandas import DataFrame
from matplotlib import pyplot
cols = list(read_csv('../data/nopromwac/data-nopromwacmfozmstt.csv', nrows =1))
print(cols)
series = read_csv('../data/nopromwac/data-nopromwacmfozmstt.csv', header=0, index_col=0, parse_dates=True, squeeze=True, usecols =[i for i in cols if i == 'diff'])
# series = series[series['diff'] != -1]
series.drop(series.index[series['diff'] != -1], inplace=True)
print(series)
s = DataFrame(
  {
    "diff": series.index
  },
  columns=["diff"]
)
pyplot.figure()
s.plot.hist(stacked=True, bins=20)
# pyplot.legend(bbox_to_anchor=(1.0, 1.0))
pyplot.savefig('test-hist.png')

# plotTilte = 'Empty minikube cluster with istio memory use in MB'
# plotYLabel = 'Memory use (MB)'
# fileName = 'all-nowac-CPU'

# startPath = '../data/metrics/' + fileName
# #hostNames = ['master-01', 'worker-01', 'worker-02']
# hostNames = ['']
# fileExtension = '.csv'

# data = []
# headers = []

# for hostName in hostNames:
#   with open(startPath + hostName + fileExtension) as csv_file:
#       csv_reader = csv.reader(csv_file, delimiter=',')
#       line_count = 0
#       maxDiff = 0
#       for row in csv_reader:
#           if line_count == 0:
#               print(f'Column names are {", ".join(row)}')
#               line_count += 1
#               headers = row
#               for h in headers:
#                 data.append([])
#           else:
#             for index, d in enumerate(row):
#               data[index].append(d)

#             line_count += 1

#   # plot
#   # fig, ax = plt.subplots()

# for index, da in enumerate(data):
#   if index == 0 or index == 1:
#     continue

#   plt.plot(data[0], data[index], linewidth=2.0)

#   # ax.set(xlim=(0, len(diff)),
#   # ylim=(0, maxDiff))

# plt.xlabel('Histogram of process time')
# plt.ylabel(plotYLabel)
# plt.title(plotTilte)

# plt.savefig(fileName + '.png')