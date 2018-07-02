package policies_derivation

import (
	"github.com/Cloud-Pie/SPDT/types"
	"time"
	"math"
)

type MediumStepOverprovision struct {
	PoIList []types.PoI
	NIntervals	int
}

func (derivationStrategy MediumStepOverprovision) WindowDerivation(values []int, times [] time.Time) (types.ProcessedForecast,[][]int,[][]time.Time){

	intervals := []types.CriticalInterval{}
	valuesIntervals := [][]int{}
	timeIntervals := [][]time.Time{}

	for _,item := range derivationStrategy.PoIList {
		interval := types.CriticalInterval{}
		interval.Requests = values[item.Index]
		interval.TimePeak = times[item.Index]

		interSize := len(intervals)
		var startIndex int
		if (interSize >= 1){
			//start time is equal to the end time from previous interval
			startIndex = interSize-1
			interval.TimeStart = intervals[startIndex].TimeEnd
		}else {
			startIndex = item.Start.Index
			interval.TimeStart = times[int(startIndex)]
		}
		//Calculate End Time using the ips_left
		endIndex := item.End.Index
		timeValleyIpsLeft := adjustTime(times[int(endIndex)], item.End.Left_ips - math.Floor(item.End.Left_ips))
		if timeValleyIpsLeft.After(interval.TimePeak) {
			interval.TimeEnd = timeValleyIpsLeft
		} else {
			interval.TimeEnd = times[int(endIndex)]
		}
		interval.AboveThreshold = item.Peak
		intervals = append(intervals, interval)

		valuesIntervals = append(valuesIntervals,values[startIndex:endIndex-1])
		timeIntervals = append(timeIntervals, times[startIndex:endIndex-1])
	}
	derivationStrategy.NIntervals = len(intervals)
	processedForecast := types.ProcessedForecast{}
	processedForecast.CriticalIntervals = intervals

	return processedForecast, valuesIntervals, timeIntervals
}

func (derivationStrategy MediumStepOverprovision) NumberIntervals() int{
	return derivationStrategy.NIntervals
}
