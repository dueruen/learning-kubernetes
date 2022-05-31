import matplotlib.pyplot as plt
import numpy as np
from os import path
import math
from pathlib import Path
from pandas import read_csv

Path("bars/").mkdir(parents=True, exist_ok=True)

names = ["wac", "niwac", "nopromwac", "wfnc"] #
titleText = ["With instrumentation", "Without Cilium and instrumentation", "With instrumentation but no prometheus", "With flannel instead of Cilium"]
subTitles = ['Average Latency', 'Median Latency', 'Latency Std', 'Latency Min', 'Latency Max']

figAvg, axsAvg = plt.subplots(len(names), 1, sharey=True, tight_layout=True, figsize=(14,4))
figAvg.suptitle('Average Latency all', fontsize=16)

figMedian, axsMedian = plt.subplots(len(names), 1, sharey=True, tight_layout=True, figsize=(14,4))
figMedian.suptitle('Median Latency all', fontsize=16)

figStd, axsStd = plt.subplots(len(names), 1, sharey=True, tight_layout=True, figsize=(14,4))
figStd.suptitle('Latency Std all', fontsize=16)

figMin, axsMin = plt.subplots(len(names), 1, sharey=True, tight_layout=True, figsize=(14,4))
figMin.suptitle('Latency Min all', fontsize=16)

figMax, axsMax = plt.subplots(len(names), 1, sharey=True, tight_layout=True, figsize=(14,4))
figMax.suptitle('Latency Max all', fontsize=16)

figAvgWithout, axsAvgWithout = plt.subplots(len(names), 1, sharey=True, tight_layout=True, figsize=(14,4))
figAvgWithout.suptitle('Average Latency without first 5 minutes', fontsize=16)

figMedianWithout, axsMedianWithout = plt.subplots(len(names), 1, sharey=True, tight_layout=True, figsize=(14,4))
figMedianWithout.suptitle('Median Latency without first 5 minutes', fontsize=16)

figStdWithout, axsStdWithout = plt.subplots(len(names), 1, sharey=True, tight_layout=True, figsize=(14,4))
figStdWithout.suptitle('Latency Std without first 5 minutes', fontsize=16)

figMinWithout, axsMinWithout = plt.subplots(len(names), 1, sharey=True, tight_layout=True, figsize=(14,4))
figMinWithout.suptitle('Latency Min without first 5 minutes', fontsize=16)

figMaxWithout, axsMaxWithout = plt.subplots(len(names), 1, sharey=True, tight_layout=True, figsize=(14,4))
figMaxWithout.suptitle('Latency Max without first 5 minutes', fontsize=16)

namesCount = 0
for name in names:

  runNames = ["data-"+ name + "mfoz", "data-"+ name + "mfozz", "data-"+ name + "mfozzz"] #, "data-wacmfozz", "data-wacmfozzz"
  runValues = ["10Hz", "100Hz", "1000Hz"]
  runSize = ["32B", "320B", "3.2KB", "32KB"]

  barLabels = []
  avgValAll = []
  stdValAll = []
  medianValAll = []
  maxValAll = []
  minValAll = []

  avgValWithoutFirst = []
  stdValWithoutFirst = []
  medianValWithoutFirst = []
  maxValWithoutFirst = []
  minValWithoutFirst = []

  index = 0
  for runName in runNames:

    print("index " + str(index))
    plotTilte = runValues[index] 
    plotYLabel = 'Count / []'
    plotXLabel = 'Time / [ms]'
    startPath = '../data/fixed/' + name + '/output/'
    fileNames = [runName+'mstt', runName+'msttz', runName+'msttzz',runName+'msttzzz'] # 
    fileExtension = '.csv'

    fig, axs = plt.subplots(len(fileNames), 1, sharey=True, sharex=True, tight_layout=True)
    fig.suptitle('Successes ' + plotTilte, fontsize=16)

    figAll, axsAll = plt.subplots(1, len(fileNames), sharey=True, sharex=True, tight_layout=True)
    figAll.suptitle('Latency histogram ' + plotTilte, fontsize=16)

    figSmall, axsSmall = plt.subplots(1, len(fileNames), sharey=True, sharex=True, tight_layout=True)
    figSmall.suptitle('Latency histogram ' + plotTilte + " min to 500ms", fontsize=16)

    figHigh, axsHigh = plt.subplots(1, len(fileNames), sharey=True, sharex=True, tight_layout=True)
    figHigh.suptitle('Latency histogram ' + plotTilte + " 500ms to max", fontsize=16)

    figLine, axsLine = plt.subplots(len(fileNames), 1, sharey=True, sharex=True, tight_layout=True, figsize=(6,8))
    figLine.suptitle('Latency ' + plotTilte, fontsize=16)

    figWithoutFirst, axsWithoutFirst = plt.subplots(len(fileNames), 1, sharey=True, sharex=True, tight_layout=True)
    figWithoutFirst.suptitle('Latency excluding first 5 min ' + plotTilte, fontsize=16)  

    count = 0
    for fileName in fileNames:
      filePathProducer = startPath + fileName + "-producer" + fileExtension
      filePathConsumer = startPath + fileName + "-consumer"+ fileExtension

      currentIndex = 0
      data = []
      ok = 0
      missing = 0
      incomplete = 0

      print(filePathProducer)
      if (path.exists(filePathProducer)) and (path.exists(filePathConsumer)):
        producer = read_csv(filePathProducer, header=0, index_col=0, parse_dates=True, squeeze=True)
        consumer = read_csv(filePathConsumer, header=0, index_col=0, parse_dates=True, squeeze=True)

        n = producer.to_frame().merge(consumer.to_frame(), on='id', how='left')
        diff = (n["time_y"] - n["time_x"]) / 1000
        n["diff"] = diff

        ss = n["diff"].apply(lambda x: 0 if (x > 0) else -1)
        ok = ss[ss == 0]
        incomplete = ss[ss < 0]

        print(n.index.max())
        ii = np.array(n.index)
        rr = np.arange(n.index.max())
        print(len(rr) - len(np.intersect1d(ii, rr))) 


        # producer.sort_values(by=['time'])
        # consumer.sort_values(by=['time'])
        # print(producer.head())
        # print(consumer.head())
        # print(n.head())
        # print(n.info)
        Path(name + "/" + runName).mkdir(parents=True, exist_ok=True)
        f = open(name + "/" + runName + '/data.txt', "a")
        
        size = int(32 * math.pow(10, count))

        high = n[n["diff"] >= 500]
        axsHigh[count].hist(high["diff"], bins=25) #, label="avg: " + str(high["diff"].mean())
        axsHigh[count].set_xlabel(plotXLabel)
        axsHigh[count].set_ylabel(plotYLabel)
        axsHigh[count].set_title(str(size) + " byte")      
        # axsHigh[count].legend(loc="upper right")
        figHigh.savefig(name + "/" + runName + '/hist-high-' +runName + '.pdf')
        f.write("highAvg,highMedian,highStd,highMin,highMax\n")
        f.write(str(high["diff"].mean()) + "," + str(high["diff"].median()) + "," + str(high["diff"].std()) + "," + str(high["diff"].min()) + "," + str(high["diff"].max()) + "\n")

        small = n[n["diff"] < 500]
        axsSmall[count].hist(small["diff"], bins=25)
        axsSmall[count].set_xlabel(plotXLabel)
        axsSmall[count].set_ylabel(plotYLabel)
        axsSmall[count].set_title(str(size) + " byte")
        figSmall.savefig(name + "/" + runName + '/hist-small-' +runName + '.pdf')
        f.write("smallAvg,smallMedian,smallStd,smallMin,smallMax\n")
        f.write(str(small["diff"].mean()) + "," + str(small["diff"].median()) + "," + str(small["diff"].std()) + "," + str(small["diff"].min()) + "," + str(small["diff"].max()) + "\n")

        # plt.savefig(runName + '/hist-small-' +runName + '.png')

        # n["diff"].plot.hist()
        axsAll[count].hist(n["diff"], bins=25) #.tolist()
        axsAll[count].set_xlabel(plotXLabel)
        axsAll[count].set_ylabel(plotYLabel)
        axsAll[count].set_title(str(size) + " byte")          
        figAll.savefig(name + "/" + runName + '/hist-all-' +runName + '.pdf')
        f.write("Avg,Median,Std,Min,Max\n")
        f.write(str(n["diff"].mean()) + "," + str(n["diff"].median()) + "," + str(n["diff"].std()) + "," + str(n["diff"].min()) + "," + str(n["diff"].max()) + "\n")

        labelName = runValues[index] + "-" + runSize[count]
        # if count != 3:
        barLabels.append(labelName)
        avgValAll.append(n["diff"].mean())
        medianValAll.append(n["diff"].median())
        stdValAll.append(n["diff"].std())
        maxValAll.append(n["diff"].max())
        minValAll.append(n["diff"].min())   

        #plt.step(range(len(ss)),ss,where='post')
        axs[count].step(range(len(ss)),ss,where='post')
        axs[count].set_xlabel(plotYLabel)
        axs[count].set_ylabel("Received / []")
        axs[count].set_title(str(size) + " byte")      
        fig.savefig(name + "/" + runName + '/step-' +runName + '.pdf')

        # lab = ["Ok", "Incomplete"]
        # val = [len(ok), len(incomplete)]
        # plt.bar(lab,val)

        axsLine[count].plot(range(len(n["diff"])),n["diff"])
        axsLine[count].set_xlabel(plotYLabel)
        axsLine[count].set_ylabel(plotXLabel)
        axsLine[count].set_title(str(size) + " byte")      
        figLine.savefig(name + "/" + runName + '/diff-line-' +runName + '.pdf')

        # n.to_csv("./data-niwac.csv")


        withoutFirst = n[n["time_x"] > n["time_x"].min() + (1000000 * 60 * 5) ]
        axsWithoutFirst[count].plot(range(len(withoutFirst["diff"])),withoutFirst["diff"])
        axsWithoutFirst[count].set_xlabel(plotYLabel)
        axsWithoutFirst[count].set_ylabel(plotXLabel)
        axsWithoutFirst[count].set_title(str(size) + " byte")      
        figWithoutFirst.savefig(name + "/" + runName + '/diff-withoutfirstfive-' +runName + '.pdf')

        # if count != 3:
        avgValWithoutFirst.append(withoutFirst["diff"].mean())
        medianValWithoutFirst.append(withoutFirst["diff"].median())
        stdValWithoutFirst.append(withoutFirst["diff"].std())
        maxValWithoutFirst.append(withoutFirst["diff"].max())
        minValWithoutFirst.append(withoutFirst["diff"].min())             

        f.close()

        count = count + 1

    axsAvg[namesCount].bar(barLabels,avgValAll)
    #axsAvg[namesCount].set_xlabel(plotYLabel)
    axsAvg[namesCount].set_ylabel(plotXLabel)
    axsAvg[namesCount].set_title(titleText[namesCount])      
    figAvg.savefig('bars/avg-bar.pdf')    

    axsMedian[namesCount].bar(barLabels,medianValAll)
    #axsMedian[namesCount].set_xlabel(plotYLabel)
    axsMedian[namesCount].set_ylabel(plotXLabel)
    axsMedian[namesCount].set_title(titleText[namesCount])      
    figMedian.savefig('bars/median-bar.pdf')  

    axsStd[namesCount].bar(barLabels,stdValAll)
    #axsStd[namesCount].set_xlabel(plotYLabel)
    axsStd[namesCount].set_ylabel(plotXLabel)
    axsStd[namesCount].set_title(titleText[namesCount])      
    figStd.savefig('bars/std-bar.pdf')      

    axsMin[namesCount].bar(barLabels,minValAll)
    #axsMin[namesCount].set_xlabel(plotYLabel)
    axsMin[namesCount].set_ylabel(plotXLabel)
    axsMin[namesCount].set_title(titleText[namesCount])      
    figMin.savefig('bars/min-bar.pdf')    

    axsMax[namesCount].bar(barLabels,maxValAll)
    #axsMax[namesCount].set_xlabel(plotYLabel)
    axsMax[namesCount].set_ylabel(plotXLabel)
    axsMax[namesCount].set_title(titleText[namesCount])      
    figMax.savefig('bars/max-bar.pdf')    

    axsAvgWithout[namesCount].bar(barLabels,avgValWithoutFirst)
    #axsAvgWithout[namesCount].set_xlabel(plotYLabel)
    axsAvgWithout[namesCount].set_ylabel(plotXLabel)
    axsAvgWithout[namesCount].set_title(titleText[namesCount])      
    figAvgWithout.savefig('bars/avg-bar-without.pdf')    

    axsMedianWithout[namesCount].bar(barLabels,medianValWithoutFirst)
    #axsMedianWithout[namesCount].set_xlabel(plotYLabel)
    axsMedianWithout[namesCount].set_ylabel(plotXLabel)
    axsMedianWithout[namesCount].set_title(titleText[namesCount])      
    figMedianWithout.savefig('bars/median-bar-without.pdf')  

    axsStdWithout[namesCount].bar(barLabels,stdValWithoutFirst)
    #axsStdWithout[namesCount].set_xlabel(plotYLabel)
    axsStdWithout[namesCount].set_ylabel(plotXLabel)
    axsStdWithout[namesCount].set_title(titleText[namesCount])      
    figStdWithout.savefig('bars/std-bar-without.pdf')      

    axsMinWithout[namesCount].bar(barLabels,minValWithoutFirst)
    #axsMinWithout[namesCount].set_xlabel(plotYLabel)
    axsMinWithout[namesCount].set_ylabel(plotXLabel)
    axsMinWithout[namesCount].set_title(titleText[namesCount])      
    figMinWithout.savefig('bars/min-bar-without.pdf')    

    axsMaxWithout[namesCount].bar(barLabels,maxValWithoutFirst)
    #axsMaxWithout[namesCount].set_xlabel(plotYLabel)
    axsMaxWithout[namesCount].set_ylabel(plotXLabel)
    axsMaxWithout[namesCount].set_title(titleText[namesCount])      
    figMaxWithout.savefig('bars/max-bar-without.pdf')                    

    index = index + 1  
  namesCount = namesCount + 1  