import csv
import matplotlib.pyplot as plt
import numpy as np
from os import path
import math

runNames = ["data-wacmfoz", "data-wacmfozz", "data-wacmfozzz"] #, "data-wacmfozz", "data-wacmfozzz"
runValues = ["10Hz", "100Hz", "1000Hz"]

index = 0
for runName in runNames:

  print("index " + str(index))
  plotTilte = 'Histogram ' + runValues[index] 
  plotYLabel = 'Messages'
  plotXLabel = 'Time (ms)'
  startPath = '../data/wac-output-v2/'
  fileNames = [runName+'mstt', runName+'msttz', runName+'msttzz',runName+'msttzzz'] # 
  fileExtension = '.csv'

  fig, axs = plt.subplots(1, 4, sharey=True, sharex=True, tight_layout=True)
  fig.suptitle(plotTilte, fontsize=16)

  figSmall, axsSmall = plt.subplots(1, 4, sharey=True, sharex=True, tight_layout=True)
  figSmall.suptitle(plotTilte, fontsize=16)

  figHigh, axsHigh = plt.subplots(1, 4, sharey=True, sharex=True, tight_layout=True)
  figHigh.suptitle(plotTilte, fontsize=16)

  count = 0
  for fileName in fileNames:
    filePath = startPath + fileName + fileExtension
    id = []
    produce = []
    consume = []
    diff = []
    dataSmall = []
    dataHigh = []
    maxDiff = 0
    maxSmallDiff = 0
    maxHighDiff = 0
    if (path.exists(filePath)):
      with open(filePath) as csv_file:
          csv_reader = csv.reader(csv_file, delimiter=',')
          line_count = 0
          for row in csv_reader:
              if line_count == 0:
                  print(f'Column names are {", ".join(row)}')
                  line_count += 1
              else:
                if int(row[3]) > 0:
                  if int(row[3]) < 10000000: 
                    id.append(int(row[0]))
                    produce.append(int(row[1]))
                    consume.append(int(row[2]))
                    diffVal = int(row[3]) / 1000
                    diff.append(diffVal)
                    if (diffVal) > maxDiff:
                      maxDiff = diffVal
                    print(row[3])

                    if diffVal <= 100:
                      if (diffVal) > maxSmallDiff:
                        maxSmallDiff = diffVal                      
                      dataSmall.append(diffVal)

                    if diffVal > 100:
                      if (diffVal) > maxHighDiff:
                        maxHighDiff = diffVal                      
                      dataHigh.append(diffVal)

                line_count += 1

    size = int(32 * math.pow(10, count))
    axs[count].hist(diff, bins=10)
    axs[count].set_xlabel(plotXLabel)
    axs[count].set_ylabel(plotYLabel)
    axs[count].set_title(str(size) + " byte")

    axsSmall[count].hist(dataSmall, bins=25)
    axsSmall[count].set_xlabel(plotXLabel)
    axsSmall[count].set_ylabel(plotYLabel)
    axsSmall[count].set_title(str(size) + " byte")  

    axsHigh[count].hist(dataHigh, bins=25)
    axsHigh[count].set_xlabel(plotXLabel)
    axsHigh[count].set_ylabel(plotYLabel)
    axsHigh[count].set_title(str(size) + " byte")  

    count = count + 1
  index = index + 1
    # axs[0].set(ylim=(0, len(diff)))



    # axs[1].hist(data, bins=50)
    # axs[1].set_xlabel(plotXLabel)
    # axs[1].set_ylabel(plotYLabel)
    # axs[1].set_xlim(1, 100)
    # axs[1].set_ylim(1, 100)
    # # axs[1].set(ylim=(0, len(data)))

    # axs[2].hist(dataBig, bins=50)
    # axs[2].set_xlabel(plotXLabel)
    # axs[2].set_ylabel(plotYLabel)
    # axs[2].set_xlim(1, 1000)
    # axs[2].set_ylim(1, 1000)  
    # # axs[2].set(ylim=(0, len(dataBig)))

    # print(len(diff))
    # print(len(data))
    # print(len(dataBig))

  fig.savefig('hist-all-' +runName + '.png')
  figSmall.savefig('hist-small-' +runName + '.png')
  figHigh.savefig('hist-high-' +runName + '.png')