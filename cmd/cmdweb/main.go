package main

import (
	"flag"
	"fmt"
	cyoa "gophercisesChoseYourOwnAdventure"
	"log"
	"net/http"
	"os"
)

// write a code to add 2 numbers

func getAndSetFlags() (*int, *string) {
	port := flag.Int("port", 3000, "the port to start the CYOA web app on")
	fileName := flag.String("file", "story.json", "The json file with CYOA story")
	flag.Parse()
	return port, fileName
}
func main() {
	port, fileName := getAndSetFlags()
	_ = port
	fmt.Println(*fileName)

	file, err := os.Open(*fileName)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JsonStoryParser(file)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("%v", story)

	// using function options
	//tpl2 := template.Must(template.New("").Parse("Hello!!"))
	//h := cyoa.NewHandler(story, cyoa.WithTemplate(tpl2))
	h := cyoa.NewHandler(story)
	_ = h
	//
	//fmt.Println(*port)
	//
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))

}
