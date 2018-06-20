package types

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

/*Service keeps the name and scale of the scaled service*/
type Service struct {
	Name  string	`json:Name`
	Scale int	`json:Scale`
}

/*State is the metadata of the state expected to scale to*/
type State struct {
	ISODate  time.Time  `json:ISODate`
	Services [] Service `json:Services`
	Name     string     `json:Name`
	VMs      [] VmScale `json:VMs`
	Time     string     `json:Time`

}
/*
func  (state State) MapTypesScale() map[string] int{
	mapTypesScale := make(map[string]int)
	for _,vm := range state.VMs {
		mapTypesScale [vm.Type] = vm.Scale
	}
	return  mapTypesScale
}*/

/*VmScale is the factor for which a type of VM is scales*/
type VmScale struct {
	Type string `json:"Type"`
	Scale int `json:"Scale"`
}

/*Resource configuration*/
type Configuration struct {
	TransitionCost float64
	State	State
	TimeStart time.Time
	TimeEnd	time.Time
}

//Policy states the scaling transitions
type Policy struct {
	ID     bson.ObjectId 	  `bson:"_id" json:"id"`
	Algorithm 	string		  `json:"algorithm" bson:"algorithm"`
	TotalCost float64		`json:"total_cost" bson:"total_cost"`
	Configurations    [] Configuration	`json:"configuration" bson:"configuration"`
	StartTimeDerivation	time.Time	`json:"start_derivation_time" bson:"start_derivation_time"`
	FinishTimeDerivation time.Time	`json:"finish_derivation_time" bson:"finish_derivation_time"`
}