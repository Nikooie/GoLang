package main

import (
	"encoding/json" /** it is used for sending/receiving JSON objects */
	"fmt"           /** it is used for formatting fuctions such as printline */
	"net/http"      /** because of REST API we are using http protocol */
	"strconv"       /** it is used because of converting string to float */
)

/** @brief This is the group of all constant variable */
const (
	HOST_NAME = "localhost"
	HOST_PORT = 7777
	SECTOR_ID = 2
)

/** @brief DNSRequest is a JSON object which is sent from drone to DNS */
type DNSRequest struct {
	X   string `json:"x"`   //x is coord of drone
	Y   string `json:"y"`   //y is coord of drone
	Z   string `json:"z"`   //z is coord of drone
	Vel string `json:"vel"` //vel is velocity of drone
}

/** @brief DNSResponse is a JSON object which is sent as response from DNS to drone */
type DNSResponse struct {
	Loc float64 `json:"loc"`
}

var requestNumer int64
/**
 *@brief This function is used for check the return value as error
 *@param
 *	err: error parameter as input
 *@return
 *	is a bool type means there is an error at the function or not
 */
func checkErr(err error) bool {
	if err != nil {
		fmt.Println("error in receiving message: %s", err)
		return false
	} else {
		return true
	}
}

/**
 *@brief This function is used for convert string format to float64
 *@param
 *	str: string format is an input param
 *	number: float64 format is an output param
 *@return
 *	err: is a bool type means there is an error at the function or not
 */
func convertStrToFloat(str string, number *float64) bool {
	var err error

	*number, err = strconv.ParseFloat(str, 64) //convert input string value to float64
	return checkErr(err)
}

/**
 *@brief GelLoc is used for calculating Loc parameter as response value
		The math formula is loc = x * sectorID + y * sectorID + z * sectorID + vel
 *@param
 *	drone: DNSRequest is an input param
 *@return
 *	LOC: is a float64
*/
func GetLoc(drone DNSRequest) float64 {
	var x, y, z, vel float64
	var LOC float64

	if !convertStrToFloat(drone.X, &x) {
		return -1
	}
	if !convertStrToFloat(drone.Y, &y) {
		return -1
	}
	if !convertStrToFloat(drone.Z, &z) {
		return -1
	}
	if !convertStrToFloat(drone.Vel, &vel) {
		return -1
	}

	LOC = SECTOR_ID*x + SECTOR_ID*y + SECTOR_ID*z + vel

	fmt.Printf("%.2f", LOC) /** LOC is with 2 digits precision*/
	return LOC
}

/**
 *@brief calcLoc is a method at url: localhost:7777 for responding to the drones
 *@param
 *	w: type http.ResponseWriter as an output param for writing the response
 * 	req: type http.Request as an input parameter for getting the client request
 */
func calcLoc(w http.ResponseWriter, req *http.Request) {
	var drone DNSRequest
	var location DNSResponse

	/** there is no JSON object in request */
	if req.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	/** Decode the request body to drone JSON object */
	err := json.NewDecoder(req.Body).Decode(&drone)
	if err != nil {
		fmt.Println("error in receiving message")
		return
	}
	fmt.Printf("requestNumer: %d, DNSRequest: x: %s, Y: %s, Z: %s, Vel: %s\n", requestNumer, drone.X, drone.Y, drone.Z, drone.Vel)
	requestNumer++

	/** create the output JSON object */
	w.Header().Set("Content-Type", "application/json")
	location.Loc = GetLoc(drone)
	response, err := json.Marshal(location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

/** @brief handle incomming request to address localhost:7777/
 *	The first one is calculate Loc params and send it back as JSON object,
 *	But if there is other methods such as new formula for other drone types
 * 	We have to defone the specific method here
 */
func handleRequest() {
	http.HandleFunc("/calcLoc", calcLoc)
}

/** @brief The start point of DNS */
func main() {
	/** Define the server address */
	fmt.Printf("Please Connect to %s:%d\n", HOST_NAME, HOST_PORT)
	server := &http.Server{
		Addr: "localhost:7777",
	}

	/** Go routines means running threads for meeting concurrency at the server side*/
	go handleRequest()

	/**  server is listening on defined port */
	server.ListenAndServe()
}
