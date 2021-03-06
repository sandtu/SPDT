package util


const ENDPOINT_FORECAST = "/predict"

const ENDPOINT_VMS_PROFILES = "/api/vms"
const ENDPOINT_STATES = "/api/states"
const ENDPOINT_INVALIDATE_STATES = "/api/invalidate/{timestamp}"
const ENDPOINT_CURRENT_STATE = "/api/current"

const ENDPOINT_SERVICE_PROFILES = "/getRegressionTRNsMongoDBAll/{apptype}/{appname}/{mainservicename}"
const ENDPOINT_VM_TIMES = "/getPerVMTypeOneBootShutDownData"
const ENDPOINT_ALL_VM_TIMES = "/getPerVMTypeAllBootShutDownData"

const ENDPOINT_SERVICE_PROFILE_BY_MSC = "/getPredictedRegressionReplicas/{apptype}/{appname}/{mainservicename}/{msc}/{numcoresutil}/{numcoreslimit}/{nummemlimit}"
const ENDPOINT_SERVICE_PROFILE_BY_REPLICAS = "/getPredictedRegressionTRN/{apptype}/{appname}/{mainservicename}/{replicas}/{numcoresutil}/{numcoreslimit}/{nummemlimit}"

const ENDPOINT_SUBSCRIBE_NOTIFICATIONS = "/subscribe"
const ENDPOINT_RECIVE_NOTIFICATIONS = "/api/forecast"

