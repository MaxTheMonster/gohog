package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/posthog/posthog-go"
)

const (
	endpointDefault = "http://localhost:8000"
)

var (
	endpoint       string = endpointDefault
	posthogKey     string
	posthogKeyFlag string
	numbEvents     int = 50
	funnelDepth    int = 10
	userCount      int = 1000
	sleepTime      int = 0
)

// runOneEventPipeline Enqueues one Event pipeline. Used for benchmarking
func runOneEventPipeline(n int, config posthog.Config, client posthog.Client) {
	s := rand.Intn(funnelDepth)
	userid := rand.Intn(userCount)
	user := fmt.Sprintf("user-%d", userid)
	plan := "Enterprise"
	friends := rand.Intn(42)
	for i := 0; i <= s; i++ {
		step := fmt.Sprintf("step-%d", i)
		client.Enqueue(posthog.Capture{
			DistinctId: user,
			Event:      step,
			Properties: posthog.NewProperties().
				Set("plan", plan).
				Set("friends", friends),
		})
	}
}

// runEventsPipeline runs a burn in test, sending events to a PostHog instance
func runEventsPipeline() {
	config := posthog.Config{
		Endpoint: endpoint,
	}
	client, err := posthog.NewWithConfig(posthogKey, config)
	if err != nil {
		panic("oh no")
	}
	defer client.Close()

	log.Println("~~~~~~~Beginning burn in test~~~~~~~")
	log.Printf("Generating %d events drawing from a pool of %d users\na funnel with a max of %d steps sleeping %d ms between event", numbEvents, userCount, funnelDepth, sleepTime)

	t := 0
	for t <= numbEvents {
		s := rand.Intn(funnelDepth)
		userid := rand.Intn(userCount)
		user := fmt.Sprintf("user-%d", userid)
		plan := "Enterprise"
		friends := rand.Intn(42)
		for i := 0; i <= s; i++ {
			t++
			step := fmt.Sprintf("step-%d", i)
			log.Printf("Logging total message: %d user: %s step: %s plan: %s friends: %d", t, user, step, plan, friends)
			client.Enqueue(posthog.Capture{
				DistinctId: user,
				Event:      step,
				Properties: posthog.NewProperties().
					Set("plan", plan).
					Set("friends", friends),
			})
			time.Sleep(time.Duration(sleepTime) * time.Millisecond)
		}
	}
}

func main() {
	flag.StringVar(&endpoint, "endpoint", endpointDefault, "Endpoint to generate traffic to")
	posthogKeyFlag := os.Getenv("POSTHOG_KEY")
	flag.StringVar(&posthogKey, "key", "", "PostHog Api Key")
	flag.IntVar(&numbEvents, "events", numbEvents, "Total number of events to spawn")
	flag.IntVar(&funnelDepth, "funnel-depth", funnelDepth, "Max depth of funnel")
	flag.IntVar(&userCount, "users", userCount, "user pool to draw from")
	flag.IntVar(&sleepTime, "sleep", sleepTime, "ms to sleep between events")
	flag.Parse()

	fmt.Println(posthogKeyFlag)
	if posthogKeyFlag != "" {
		posthogKey = posthogKeyFlag
	}
	runEventsPipeline()
}
