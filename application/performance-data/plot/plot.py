import csv
import matplotlib.pyplot as plt
import numpy as np

plotTilte = 'Empty minikube cluster with istio memory use in MB'
plotYLabel = 'Memory use (MB)'
fileName = 'data-wacmfozmsttz'

startPath = '../data/wac-output-v2/' + fileName
#hostNames = ['master-01', 'worker-01', 'worker-02']
hostNames = ['']
fileExtension = '.csv'

for hostName in hostNames:
  id = []
  produce = []
  consume = []
  diff = []
  date = []
  with open(startPath + hostName + fileExtension) as csv_file:
      csv_reader = csv.reader(csv_file, delimiter=',')
      line_count = 0
      maxDiff = 0
      for row in csv_reader:
          if line_count == 0:
              print(f'Column names are {", ".join(row)}')
              line_count += 1
          else:
            if int(row[3]) > 0:
              if int(row[3]) < 100000000:
                id.append(int(row[0]))
                produce.append(int(row[1]))
                consume.append(int(row[2]))
                diffVal = int(row[3]) / 1000
                diff.append(diffVal)
                if (diffVal) > maxDiff:
                  maxDiff = diffVal
                print(row[3])

            line_count += 1

  # plot
  fig, ax = plt.subplots()

  ax.plot(id, diff, linewidth=2.0)

  ax.set(xlim=(0, len(diff)),
  ylim=(0, maxDiff))

plt.xlabel('Histogram of process time')
plt.ylabel(plotYLabel)
plt.title(plotTilte)

plt.savefig(fileName + '.png')