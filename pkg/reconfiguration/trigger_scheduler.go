package reconfiguration

import (
	"fmt"
	"github.com/Cloud-Pie/SPDT/types"
	"github.com/Cloud-Pie/SPDT/rest_clients/scheduler"
	"time"
)

func TriggerScheduler(policy types.Policy, endpoint string){
	for _, conf := range policy.Configurations {
		err := scheduler.CreateState(conf.State, endpoint)
		if err != nil {
			fmt.Printf("The scheduler request failed with error %s\n", err)
			panic(err)
		}
	}
}

func InvalidateStates(timestamp time.Time, endpoint string){
	endpoint = endpoint + timestamp.String()
		err := scheduler.InvalidateStates(endpoint)
		if err != nil {
			fmt.Printf("The scheduler failed with error %s\n", err)
			panic(err)
		}
}