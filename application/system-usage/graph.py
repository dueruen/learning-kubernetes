import csv
import matplotlib.pyplot as plt

plotTilte = 'Empty minikube cluster with istio memory use in MB'
plotYLabel = 'Memory use (MB)'
fileName = 'empty-istio-minikube'

startPath = './data/minikube/empty-istio-cluster'
#hostNames = ['master-01', 'worker-01', 'worker-02']
hostNames = ['']
fileExtension = '.txt'

for hostName in hostNames:
  memuse = []
  memtotal = []
  mempro = []
  cpuload = []
  date = []
  with open(startPath + hostName + fileExtension) as csv_file:
      csv_reader = csv.reader(csv_file, delimiter=' ')
      line_count = 0
      for row in csv_reader:
          if line_count == 0:
              print(f'Column names are {", ".join(row)}')
              line_count += 1
          else:
            memuse.append(int(row[0]))
            memtotal.append(int(row[1]))
            mempro.append(float(row[2]))
            cpuload.append(float(row[3]))

            s1 = row[4].split(":")
            s2 = s1[1].split("-")

            t1 = (int(s2[0]) * 3600000000000)
            t2 = int(s2[1]) * 60000000000
            t3 = int(s2[2]) * 1000000000 
            time = t1 + t2 + t3 + int(s2[3])
            date.append(time)
          #     print(f'\t{row[0]} works in the {row[1]} department, and was born in {row[2]}.')
            line_count += 1

  plt.plot(date, memuse, label = hostName)
  #plt.ylim(ymin=0, ymax = 20)
  #plt.plot(date, memtotal)

plt.xlabel('Current time in ns')
plt.ylabel(plotYLabel)
plt.title(plotTilte)

#plt.legend()
plt.ylim(ymin=0, ymax=8000)
plt.savefig(fileName + '.png')