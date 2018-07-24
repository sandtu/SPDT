package policies_derivation

import (
	"github.com/Cloud-Pie/SPDT/types"
	"time"
	"math"
	"strconv"
	"github.com/Cloud-Pie/SPDT/util"
	"gopkg.in/mgo.v2/bson"
	db "github.com/Cloud-Pie/SPDT/storage/policies"
)

type SStepRepackPolicy struct {
	algorithm 		string				 //Algorithm's name
	timeWindow 		TimeWindowDerivation //Algorithm used to process the forecasted time serie
}

func (policy SStepRepackPolicy) CreatePolicies(processedForecast types.ProcessedForecast, mapVMProfiles map[string]types.VmProfile, serviceProfile types.ServiceProfile) [] types.Policy {

	policies := []types.Policy{}
	//Compute results for cluster of each type
	maxPolicies := 100
	for maxPolicies > 0 {
		newPolicy := types.Policy{}
		newPolicy.StartTimeDerivation = time.Now()
		configurations := []types.Configuration{}
		totalOverProvision := float32(0.0)
		totalUnderProvision := float32(0.0)
		nPeaksInConfiguration := float32(0.0)
		avgOver := float32(0.0)
		avgUnder := float32(0.0)

		for _, it := range processedForecast.CriticalIntervals {
			requests := it.Requests
			services := make(map[string]int)

			//Select the performance profile that fits better
			performanceProfile := selectProfile(serviceProfile.PerformanceProfiles)
			//Compute number of replicas needed depending on requests
			nProfileCopies := int(math.Ceil(float64(requests) / float64(performanceProfile.TRN)))
			nServiceReplicas := nProfileCopies * performanceProfile.NumReplicas
			services[serviceProfile.Name] = nServiceReplicas

			//Compute under/over provision
			diff := (nProfileCopies * performanceProfile.TRN) - requests
			var over, under float32
			if diff >= 0 {
				over = float32(diff * 100 / requests)
				under = 0
			} else {
				over = 0
				under = -1 * float32(diff*100/requests)
			}
			nPeaksInConfiguration +=1
			avgOver += over
			avgUnder += under

			//Find suitable Vm(s) depending on resources limit and current state
			vms := FindSuitableVMs(mapVMProfiles, nServiceReplicas, performanceProfile.Limit, "")

			state := types.State{}
			state.Services = services
			state.VMs = vms

			//Adjust termination times for resources configuration
			terminationTime := ComputeVMTerminationTime(mapVMProfiles, vms)
			finishTime := it.TimeEnd.Add(time.Duration(terminationTime) * time.Second)

			nConfigurations := len(configurations)
			if nConfigurations >= 1 && state.Equal(configurations[nConfigurations-1].State) {
				configurations[nConfigurations-1].TimeEnd = finishTime
				configurations[nConfigurations-1].Metrics.OverProvision = avgOver/ nPeaksInConfiguration
			} else {
				//Adjust booting times for resources configuration
				transitionTime := ComputeVMBootingTime(mapVMProfiles, vms)                            //TODO: It should include a booting rate
				totalServicesBootingTime := performanceProfile.BootTimeSec //TODO: It should include a booting rate
				startTime := it.TimeStart.Add(-1 * time.Duration(transitionTime) * time.Second)       //Booting time VM
				startTime = startTime.Add(-1 * time.Duration(totalServicesBootingTime) * time.Second) //Start time containers
				state.LaunchTime = startTime
				state.Name = strconv.Itoa(nConfigurations) + "__" + serviceProfile.Name + "__" + startTime.Format(util.TIME_LAYOUT)
				configurations = append(configurations,
					types.Configuration{
						State:     state,
						TimeStart: it.TimeStart,
						TimeEnd:   finishTime,
						Metrics: types.Metrics{
							OverProvision:  over,
							UnderProvision: under,
						},
					})
				nPeaksInConfiguration = 0
				avgOver = 0
			}
			totalOverProvision += over
			totalUnderProvision += under
		}

		totalConfigurations := len(configurations)
		totalPeaks := len(processedForecast.CriticalIntervals)
		//Add new policy
		newPolicy.Configurations = configurations
		newPolicy.FinishTimeDerivation = time.Now()
		newPolicy.Algorithm = policy.algorithm
		newPolicy.ID = bson.NewObjectId()
		newPolicy.Metrics = types.Metrics{
			OverProvision:        totalOverProvision / float32(totalPeaks),
			UnderProvision:       totalUnderProvision / float32(totalPeaks),
			NumberConfigurations: totalConfigurations,
		}
		//store policy
		db.Store(newPolicy)
		policies = append(policies, newPolicy)

		maxPolicies -=1
		}
		return policies
}


