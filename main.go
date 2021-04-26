package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli"
)

type Entry struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Dummy1   string `json:"dummy1"`
	Dummy2   string `json:"dummy2"`
	Dummy3   string `json:"dummy3"`
}

const (
	version = "1.0.1"
	app_url = "https://go-itquiz.herokuapp.com/data/ja"
)

func action(c *cli.Context) error {
	res, err := http.Get(app_url)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return err
	}

	var entries []Entry
	err = json.Unmarshal(body, &entries)
	if err != nil {
		log.Fatal(err)
		return err
	}

	correct := 0

	rand.Seed(time.Now().UnixNano())
	for i, e := range entries {
		fmt.Printf("\x1b[1m# Question.%d\x1b[0m\n%s\n", i+1, e.Question)
		ans := e.Answer

		// shuffle choices
		strs := []*string{&e.Answer, &e.Dummy1, &e.Dummy2, &e.Dummy3}
		for i := range strs {
			j := rand.Intn(i + 1)
			strs[j], strs[i] = strs[i], strs[j]
		}

		for i, s := range strs {
			fmt.Printf(" %d. %s\n", i+1, *s)
		}

		fmt.Printf("Input Number of Answer: ")
		var n int
		for {
			fmt.Scan(&n)
			if n == 1 || n == 2 || n == 3 || n == 4 {
				break
			} else {
				fmt.Println("Enter number of 1-4.")
			}
		}

		if *strs[n-1] == ans {
			fmt.Println("\x1b[1;32mCorrect!!\x1b[0m")
			correct++
		} else {
			fmt.Printf("\x1b[1;31mWrong...\x1b[0m  ")

			for i, s := range strs {
				if *s == ans {
					fmt.Printf("The Correct Answer is %d.\n", i+1)
					break
				}
			}
		}
	}

	fmt.Printf("\x1b[1;34mTotal Score: %d / %d\x1b[0m\n", correct, len(entries))

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "go-itquiz-cli"
	app.Usage = "IT Quiz in CLI"
	app.Version = version
	app.Action = action

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
