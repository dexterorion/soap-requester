package main

import (
	"flag"
	"fmt"
	"github.com/soap-requester/soap"

	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

var (
	requestFile  string
	responseFile string
	ws           string
	action       string

	log *zap.Logger
)

func init() {
	flag.StringVar(&requestFile, "xmlRequestPath", "/home/user/request.xml", "path to file with XML to perform request")
	flag.StringVar(&responseFile, "xmlResponsePath", "/home/user/response.xml", "path to store on file XML response")
	flag.StringVar(&ws, "ws", "test.com", "webservice URL value")
	flag.StringVar(&action, "action", "get_list", "webservice action")

	required := []string{"xmlRequestPath", "xmlResponsePath", "ws", "action"}
	flag.Parse()

	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			fmt.Fprintf(os.Stderr, "missing required [-%s] argument/flag\n", req)
			os.Exit(2)
		}
	}

	log, _ = zap.NewProduction()
}

func main() {
	log.Info("Starting soap requester")
	defer log.Sync()

	data, err := ioutil.ReadFile(requestFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file [%s]: [%s]\n", requestFile, err.Error())
		os.Exit(2)
	}

	result, err := soap.SoapCall(ws, action, data)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error requesting webservice [%s] - action [%s] failied with message: [%s]\n", ws, action, err.Error())
		os.Exit(2)
	}

	out, err := os.Create(responseFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating file [%s] : [%s]\n", responseFile, err.Error())
		os.Exit(2)
	}

	nb, err := out.Write(result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing to file [%s] : [%s]\n", responseFile, err.Error())
		os.Exit(2)
	}

	log.Info(fmt.Sprintf("%d bytes written to file %s", nb, responseFile))
	log.Info("Finishing...")
}
