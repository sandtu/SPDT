package storage

import (
	"gopkg.in/mgo.v2"
	"github.com/Cloud-Pie/SPDT/types"
	"gopkg.in/mgo.v2/bson"
	"github.com/op/go-logging"
	"time"
	"os"
)

type ForecastDAO struct {
	Server	string
	Database	string
	Collection  string
	db *mgo.Database
	session *mgo.Session
}
var (
	ForecastDB *ForecastDAO
 	log = logging.MustGetLogger("spdt")
 	DEFAULT_DB_SERVER_FORECAST = os.Getenv("FORECASTDB_HOST")
 	forecastDBHost = []string{ DEFAULT_DB_SERVER_FORECAST, }

)

const(
 	DEFAULT_DB_FORECAST = "Forecast"
	DEFAULT_DB_COLLECTION_FORECAST = "Forecast"
)

//Connect to the database
func (p *ForecastDAO) Connect() (*mgo.Database, error) {
	var err error

	if p.session == nil {
		p.session,  err = mgo.DialWithInfo(&mgo.DialInfo{
			Addrs: forecastDBHost,
			Username: os.Getenv("FORECASTDB_USER"),
			Password: os.Getenv("FORECASTDB_PASS"),
			Timeout:  60 * time.Second,
		})
		if err != nil {
			return nil, err
		}
	}
	p.session = p.session.Clone()
	p.db = p.session.DB(p.Database)
	return p.db,err
}

//Retrieve all the stored elements
func (p *ForecastDAO) FindAll() ([]types.Forecast, error) {
	var forecast []types.Forecast
	err := p.db.C(p.Collection).Find(bson.M{}).All(&forecast)
	return forecast, err
}

//Retrieve the item with the specified ID
func (p *ForecastDAO) FindByID(id string) (types.Forecast, error) {
	var forecast types.Forecast
	err := p.db.C(p.Collection).FindId(bson.ObjectIdHex(id)).One(&forecast)
	return forecast,err
}

//Insert a new forecast
func (p *ForecastDAO) Insert(forecast types.Forecast) error {
	err := p.db.C(p.Collection).Insert(&forecast)
	return err
}

//Delete the specified item
func (p *ForecastDAO) Delete(forecast types.Forecast) error {
	err := p.db.C(p.Collection).Remove(&forecast)
	return err
}

//Delete the specified item
func (p *ForecastDAO) Update(id bson.ObjectId, forecast types.Forecast) error {
	err := p.db.C(p.Collection).Update(bson.M{"_id":id}, forecast)
	return err
}

//Delete all forecast older than a timestamp
func (p *ForecastDAO) DeleteAllBeforeDate(timestamp time.Time) error {
	_,err := p.db.C(p.Collection).RemoveAll(bson.M{"end_time": bson.M{"$lte":timestamp}})
	return err
}

//Retrieve all policies for start time greater than or equal to time t
func (p *ForecastDAO) FindOneByTimeWindow(startTime time.Time, endTime time.Time) (types.Forecast, error) {
	var forecast types.Forecast
	//Search for that retrieves exact time window
	err := p.db.C(p.Collection).
		Find(bson.M{"start_time": bson.M{"$eq":startTime},
					"end_time": bson.M{"$eq":endTime}}).One(&forecast)

	//If user specified search parameters which are not precise, then search the closest time window
	/*if err != nil {
		err = p.db.C(DEFAULT_DB_COLLECTION_FORECAST).
			Find(bson.M{"start_time": bson.M{"$gte":startTime, "$lte":endTime},
			"end_time": bson.M{"$gte":endTime}}).One(&forecast)

		if err != nil {
			err = p.db.C(DEFAULT_DB_COLLECTION_FORECAST).
				Find(bson.M{"start_time": bson.M{"$lte":startTime},
				"end_time": bson.M{"$lte":endTime, "$gte":startTime}}).One(&forecast)

			if err != nil {
				err = p.db.C(DEFAULT_DB_COLLECTION_FORECAST).
					Find(bson.M{"start_time": bson.M{"$lte":startTime},
					"end_time": bson.M{"$gte":endTime}}).One(&forecast)
			}
		}
	}*/
	return forecast,err
}

func GetForecastDAO(serviceName string) *ForecastDAO{
	if ForecastDB == nil {
		ForecastDB = &ForecastDAO {
			Database:DEFAULT_DB_FORECAST,
			Collection:DEFAULT_DB_COLLECTION_FORECAST + "_" + serviceName,
		}
		_,err := ForecastDB.Connect()
		if err != nil {
			log.Fatalf(err.Error())
		}
	} else if ForecastDB.Collection != DEFAULT_DB_COLLECTION_FORECAST + "_" + serviceName {
		ForecastDB = &ForecastDAO {
			Database:DEFAULT_DB_FORECAST,
			Collection:DEFAULT_DB_COLLECTION_FORECAST + "_" + serviceName,
		}
		_,err := ForecastDB.Connect()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
	return ForecastDB
}
