package skylib

import (
	"log"
	"json"
	"flag"
	"os"
	"github.com/ha/doozer"
	"fmt"
	"rand"
	"rpc"
	"expvar"
)


var DC *doozer.Conn
var NS *NetworkServers
var RpcServices []*RpcService


var Port *int = flag.Int("port", 9999, "tcp port to listen")
var Name *string = flag.String("name", "changeme", "name of this server")
var BindIP *string = flag.String("bindaddress", "127.0.0.1", "address to bind")
var LogFileName *string = flag.String("logFileName", "myservice.log", "name of logfile")
var DoozerServer *string = flag.String("doozerServer", "127.0.0.1:8046", "addr:port of doozer server")
var Requests *expvar.Int
var Errors *expvar.Int
var Goroutines *expvar.Int


// This is simple today - it returns the first listed service that matches the request
// Load balancing needs to be applied here somewhere.
func GetRandomClientByProvides(provides string) (*rpc.Client, os.Error) {
	var providesList = make([]*Service, 0)

	var newClient *rpc.Client
	var err os.Error

	for _, v := range NS.Services {
		if v != nil {
			if v.Provides == provides {
				providesList = append(providesList, v)
			}

		}
	}

	if len(providesList) > 0 {
		random := rand.Int() % len(providesList)
		s := providesList[random]

		portString := fmt.Sprintf("%s:%d", s.IPAddress, s.Port)
		newClient, err = rpc.DialHTTP("tcp", portString)
		if err != nil {
			log.Printf("Found %d Clients to service %s request.", len(providesList), provides)
			return nil, NewError(NO_CLIENT_PROVIDES_SERVICE, provides)
		}

	} else {
		return nil, NewError(NO_CLIENT_PROVIDES_SERVICE, provides)
	}
	return newClient, nil
}


func DoozerConnect() {
	var err os.Error
	DC, err = doozer.Dial(*DoozerServer)
	if err != nil {
		log.Panic(err.String())
	}
}

// on startup load the configuration file. 
// After the config file is loaded, we set the global config file variable to the
// unmarshaled data, making it useable for all other processes in this app.
func LoadConfig() {
	data, _, err := DC.Get("/servers/config/networkservers.conf", nil)
	if err != nil {
		log.Panic(err.String())
	}
	if len(data) > 0 {
		log.Println("setting data ", string(data))
		setConfig(data)
		return
	}
	log.Println("Error, loading default config")
	NS = &NetworkServers{}
}

func RemoveServiceAt(i int) {

	newServices := make([]*Service, 0)

	for k, v := range NS.Services {
		if k != i {
			if v != nil {
				newServices = append(newServices, v)
			}
		}
	}
	NS.Services = newServices
	b, err := json.Marshal(NS)
	if err != nil {
		log.Panic(err.String())
	}
	rev, err := DC.Rev()
	if err != nil {
		log.Panic(err.String())
	}
	_, err = DC.Set("/servers/config/networkservers.conf", rev, b)
	if err != nil {
		log.Panic(err.String())
	}

}

func (r *Service) RemoveFromConfig() {

	newServices := make([]*Service, 0)

	for _, v := range NS.Services {
		if v != nil {
			if !v.Equal(r) {
				newServices = append(newServices, v)
			}

		}
	}
	NS.Services = newServices
	b, err := json.Marshal(NS)
	if err != nil {
		log.Panic(err.String())
	}
	rev, err := DC.Rev()
	if err != nil {
		log.Panic(err.String())
	}
	_, err = DC.Set("/servers/config/networkservers.conf", rev, b)
	if err != nil {
		log.Panic(err.String())
	}
}

func (r *Service) AddToConfig() {
	for _, v := range NS.Services {
		if v != nil {
			if v.Equal(r) {
				log.Printf("Skipping adding %s because it alreday exists.", v.Name)
				return // it's there so we don't need an update
			}
		}
	}
	NS.Services = append(NS.Services, r)
	b, err := json.Marshal(NS)
	if err != nil {
		log.Panic(err.String())
	}
	rev, err := DC.Rev()
	if err != nil {
		log.Panic(err.String())
	}
	_, err = DC.Set("/servers/config/networkservers.conf", rev, b)
	if err != nil {
		log.Panic(err.String())
	}
}

// unmarshal data from remote store into global config variable
func setConfig(data []byte) {
	err := json.Unmarshal(data, &NS)
	if err != nil {
		log.Panic(err.String())
	}
}

// Watch for remote changes to the config file.  When new changes occur
// reload our copy of the config file.
// Meant to be run as a goroutine continuously.
func WatchConfig() {
	rev, err := DC.Rev()
	if err != nil {
		log.Panic(err.String())
	}
	for {

		// blocking wait call returns on a change
		ev, err := DC.Wait("/servers/config/networkservers.conf", rev)
		if err != nil {
			log.Panic(err.String())
		}
		log.Println("Received new configuration.  Setting local config.")
		setConfig(ev.Body)

		rev = ev.Rev + 1
	}

}


func initDefaultExpVars(name string) {
	Requests = expvar.NewInt(name + "-processed")
	Errors = expvar.NewInt(name + "-errors")
	Goroutines = expvar.NewInt(name + "-goroutines")
}

func Setup(name string) {
	DoozerConnect()
	LoadConfig()
	if x := recover(); x != nil {
		log.Println("No Configuration File loaded.  Creating One.")
	}

	initDefaultExpVars(name)

	svc := NewService(name)

	svc.AddToConfig()

	go WatchConfig()

	RegisterHeartbeat()

}