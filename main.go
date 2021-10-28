package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/valyala/fasthttp"
)

var (
	addr     = flag.String("addr", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", false, "response compression")
	file     = flag.String("file", "mocky.json", "Location mock json file")

	resp interface{}

	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

func main() {
	flag.Parse()
	err := ParseJson(*file)
	if err != nil {
		log.Fatal(err)
	}

	h := Handler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	fmt.Println("Starting server on", *addr)

	err = fasthttp.ListenAndServe(*addr, h)
	if err != nil {
		panic(err)
	}
}

func Handler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)

	statusCode := http.StatusOK
	ctx.Response.SetStatusCode(statusCode)

	if err := json.NewEncoder(ctx).Encode(resp); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

func ParseJson(file string) error {
	jsonFile, err := os.Open(file)
	if err != nil {
		return err
	}

	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(jsonFile)

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &resp)
	if err != nil {
		return err
	}

	return nil
}
