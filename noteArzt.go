package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type serverConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	OfficeID string `json:"officeID"`
}

// reading in the info from the config.json into a serverConfig object
var confvar = loadConfiguration("./config.json")

var dbMockUp = make(map[string][]PatientInfo)

func loadConfiguration(file string) serverConfig {
	var config serverConfig
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

// NotificationInfo contains the appointment information that can be shared with the centralised app
type NotificationInfo struct {
	Time          string         `json:"time"`
	PatientID     string         `json:"patientId"`
	Notifications []Notification `json:"notifications"`
	Estimate      int            `json:"estimatedInMinutes"`
}

type Notification struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

// PatientInfo contains the internal information that is not shared with the centralised app
type PatientInfo struct {
	Time           string         `json:"time"`
	PatientID      string         `json:"patientId"`
	PatientName    string         `json:"patientName"`
	Notifications  []Notification `json:"notifications"`
	Estimate       int            `json:"estimatedInMinutes"`
	Urgent         bool           `json:"urgent"`
	COVIDSuspected bool           `json:"potentialCOVID-19"`
	QueuingApp     bool           `json:"queuingApp"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/view/{day:[0-9]+}", view).Methods("GET")
	router.HandleFunc("/addPatient/{day:[0-9]+}", addPatient).Methods("POST")
	router.HandleFunc("/deletePatient/{day:[0-9]+}/{patientID}", deletePatient).Methods("POST")
	router.HandleFunc("/addList/{day:[0-9]+}", addList).Methods("POST")
	log.Println("Listening at" + confvar.Port + "...")
	log.Fatal(http.ListenAndServe(confvar.Port, router))
}

func view(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	day := vars["day"]
	patients, found := dbMockUp[day]
	if !found {
		// better error handling for later
		log.Println("couldn't find", day)
	}
	resultJSON, _ := json.Marshal(patients)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintln(w, string(resultJSON))
}

func addPatient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	day := vars["day"]
	// receiving json does not really handle errors yet
	var patient PatientInfo
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	patients := dbMockUp[day]
	patients = append(patients, patient)
	dbMockUp[day] = patients
	// sending it to external service
	var info NotificationInfo
	info.Time = patient.Time
	info.PatientID = patient.PatientID
	info.Notifications = patient.Notifications
	info.Estimate = patient.Estimate
	infoJSON, _ := json.Marshal(info)
	// errors should be handled, POST should be encrypted
	extLink := confvar.Host + "/" + confvar.OfficeID + "/" + day
	http.Post(extLink, "application/json", bytes.NewBuffer(infoJSON))
	// simple redirect to view for now
	newLink := "http://localhost:1919/view/" + day
	http.Redirect(w, r, newLink, 301)
}

func deletePatient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	day := vars["day"]
	patientID := vars["patientID"]
	patients := dbMockUp[day]
	// lets not test whether the patient exist for now
	tmpPatients := []PatientInfo{}
	for _, v := range patients {
		if v.PatientID != patientID {
			tmpPatients = append(tmpPatients, v)
		}
	}
	dbMockUp[day] = tmpPatients
	extLink := confvar.Host + "/" + confvar.OfficeID + "/" + day + "/" + patientID
	// errors should be handled
	req, _ := http.NewRequest("DELETE", extLink, nil)
	http.DefaultClient.Do(req)
	// simple redirect to view for now
	newLink := "http://localhost:1919/view/" + day
	http.Redirect(w, r, newLink, 301)
}

func addList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	day := vars["day"]
	// receiving json does not really handle errors yet
	var patients []PatientInfo
	err := json.NewDecoder(r.Body).Decode(&patients)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dbMockUp[day] = patients
	// sending it to external service
	var bulk []NotificationInfo
	for _, v := range patients {
		var info NotificationInfo
		info.Time = v.Time
		info.PatientID = v.PatientID
		info.Notifications = v.Notifications
		info.Estimate = v.Estimate
		bulk = append(bulk, info)
	}
	bulkJSON, _ := json.Marshal(bulk)
	// errors should be handled, POST should be encrypted
	extLink := confvar.Host + "/" + confvar.OfficeID + "/" + day
	http.Post(extLink, "application/json", bytes.NewBuffer(bulkJSON))
	// simple redirect to view for now
	newLink := "http://localhost:1919/view/" + day
	http.Redirect(w, r, newLink, 301)
}
