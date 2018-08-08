package policy_evaluation

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/Cloud-Pie/SPDT/types"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("spdt")

func ParsePricesFile(configFile string) (types.PriceModel,error) {
	prices := types.PriceModel{}
	source, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(source), &prices)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return prices,err
}
