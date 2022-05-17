import csv
import matplotlib.pyplot as plt
import numpy as np
from os import path
import math
from pathlib import Path


runNames = ["data-nopromwacmfoz", "data-nopromwacmfozz", "data-nopromwacmfozzz"] #, "data-wacmfozz", "data-wacmfozzz"
runValues = ["10Hz", "100Hz", "1000Hz"]

index = 0
for runName in runNames:

  print("index " + str(index))
  plotTilte = 'Histogram ' + runValues[index] 
  plotYLabel = 'Messages'
  plotXLabel = 'Time (ms)'
  startPath = '../data/nopromwac/'
  fileNames = [runName+'mstt', runName+'msttz', runName+'msttzz',runName+'msttzzz'] # 
  fileExtension = '.csv'

  fig, axs = plt.subplots(4, 1, sharey=True, sharex=True, tight_layout=True, figsize=(20,5))
  fig.suptitle(plotTilte, fontsize=16)

  figSmall, axsSmall = plt.subplots(1, 4, sharey=True, sharex=True, tight_layout=True)
  figSmall.suptitle(plotTilte, fontsize=16)

  # figSmall, axsSmall = plt.subplots(1, 4, sharey=True, sharex=True, tight_layout=True)
  # figSmall.suptitle(plotTilte, fontsize=16)

  # figHigh, axsHigh = plt.subplots(1, 4, sharey=True, sharex=True, tight_layout=True)
  # figHigh.suptitle(plotTilte, fontsize=16)

  count = 0
  for fileName in fileNames:
    filePath = startPath + fileName + fileExtension
    currentIndex = 0
    data = []
    ok = 0
    missing = 0
    incomplete = 0

    if (path.exists(filePath)):
      with open(filePath) as csv_file:
          csv_reader = csv.reader(csv_file, delimiter=',')
          line_count = 0
          # limit = 70
          for row in csv_reader:
            # if limit == 0:
            #   continue
            if line_count == 0:
                print(f'Column names are {", ".join(row)}')
                line_count += 1
            else:
              if int(row[0]) == currentIndex + 1:
                if int(row[3]) == -1 or int(row[3]) > 100000:
                  data.append(-1)
                  incomplete = incomplete + 1
                else:
                  data.append(0)
                  ok = ok + 1

                currentIndex = currentIndex + 1
              elif int(row[0]) > currentIndex + 1:
                val = (int(row[0])) - currentIndex
                for v in range(1, val):
                  data.append(1)
                  missing = missing + 1

                currentIndex = currentIndex + val

            line_count += 1
            # limit = limit - 1
            

    print("Current index::: " + str(currentIndex))
    # plt.figure(figsize=(30,5))
    # plt.step(range(len(data)),data,where='post')
    # plt.savefig('step-all-' +runName + '.png')
    size = int(32 * math.pow(10, count))
    axs[count].step(range(len(data)),data,where='post')
    # axs[count].scatter(range(len(data)),data, c=cm.hot(np.abs(data)), edgecolor='none')
    # axs[count].set_xlabel(plotXLabel)
    # axs[count].set_ylabel(plotYLabel)
    axs[count].set_title(str(size) + " byte")

    lab = ["Ok", "Missing", "Incomplete"]
    val = [ok, missing, incomplete]
    axsSmall[count].bar(lab,val)

    # # x =range(len(data))
    # # y = data
    # N = 10
    # x = np.arange(N)
    # # Here are many sets of y to plot vs. x
    # ys = [x + i for i in x]
    # print(ys)
    # fig, ax = plt.subplots()

    # s = [np.column_stack([x, y]) for y in ys]
    # print(s)

    # line_segments = LineCollection(s,
    #                            linewidths=(0.5, 1, 1.5, 2),
    #                            linestyles='solid')

    # line_segments.set_array(x)
    # ax.add_collection(line_segments)
    # axcb = fig.colorbar(line_segments)    

    # fig.savefig('step-all-' +runName + '.png')                               

    # dydx = np.cos(0.5 * (x[:-1] + x[1:]))  # first derivative

    # Create a set of line segments so that we can color them individually
    # This creates the points as a N x 1 x 2 array so that we can stack points
    # together easily to get the segments. The segments array for line collection
    # needs to be (numlines) x (points per line) x 2 (for x and y)
    # points = np.array([x, y]).T.reshape(-1, 1, 2)
    # segments = np.concatenate([points[:-1], points[1:]], axis=1)

    # fig, axs = plt.subplots(sharex=True, sharey=True)

    # Create a continuous norm to map from data points to colors
    # norm = plt.Normalize(data.min(), data.max())
    # lc = LineCollection(segments, cmap='viridis', norm=norm)
    # # Set the values used for colormapping
    # lc.set_array(data)
    # lc.set_linewidth(2)
    # line = axs[0].add_collection(lc)
    # fig.colorbar(line, ax=axs[0])

    # Use a boundary norm instead
    # cmap = ListedColormap(['r', 'g', 'b'])
    # norm = BoundaryNorm([-1, -0.5, 0.5, 1], cmap.N)
    # lc = LineCollection(segments, cmap=cmap, norm=norm)
    # # lc.set_array(data)
    # lc.set_linewidth(2)
    # line = axs.add_collection(lc)
    # fig.colorbar(line, ax=axs)

    # axs[0].set_xlim(x.min(), x.max())
    # axs[0].set_ylim(-1.1, 1.1)


    count = count + 1
  index = index + 1
    # axs[0].set(ylim=(0, len(diff)))
  # fig(figsize=(30,5))
  Path(runName).mkdir(parents=True, exist_ok=True)
  fig.savefig(runName + '/step-all-' +runName + '.png')
  figSmall.savefig(runName + '/step-bar-all-' +runName + '.png')


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

  # plt.savefig('step-all-' +runName + '-' +fileName + '.png')
  # figSmall.savefig('hist-small-' +runName + '.png')
  # figHigh.savefig('hist-high-' +runName + '.png')