package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	flag "github.com/spf13/pflag"
	"golang.org/x/net/context"
)

var (
	flags  = flag.NewFlagSet("", flag.ContinueOnError)
	images = flags.String("images", "", "command seperated list of tagged images to pull.")
)

func main() {
	flags.Parse(os.Args)
	if *images == "" {
		log.Fatalf("Specify --images=foo:1.0,bar:2.0 etc")
	}
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.21", nil, defaultHeaders)
	if err != nil {
		log.Fatal(err)
	}

	for _, image := range strings.Split(*images, ",") {
		if len(strings.Split(image, ":")) != 2 {
			log.Printf("Skipping untagged image %v", image)
			continue
		}
		log.Printf("Pulling %v", image)
		resp, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
		if err != nil {
			log.Fatalf("Failed to pull %v: %v", image, err)
		}
		var lines interface{}
		for lineReader := bufio.NewReader(resp); err != io.EOF; {
			err = json.NewDecoder(lineReader).Decode(&lines)
			if m, ok := lines.(map[string]interface{}); ok {
				log.Printf("%+v", m)
			}
		}
	}
}
