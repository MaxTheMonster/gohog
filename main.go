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

func main() {
  var endpoint string
  flag.StringVar(&endpoint, "endpoint", endpointDefault, "Endpoint to generate traffic to")
  posthogKey := os.Getenv("POSTHOG_KEY") 
  numbEvents := flag.Int("events", 1000, "Total number of events to spawn")
  funnelDepth := flag.Int("funnel-depth", 10, "Max depth of funnel")
  userCount := flag.Int("users", 10000, "user pool to draw from")
  sleepTime := flag.Int("sleep", 50, "ms to sleep between events")
  flag.Parse()
 
  log.Println("~~~~~~~Beginning burn in test~~~~~~~") 
  log.Printf("Generating %d events drawing from a pool of %d users\na funnel with a max of %d steps sleeping %d ms between event", *numbEvents, *userCount, *funnelDepth, *sleepTime)

  config := posthog.Config{
    Endpoint: endpoint,
  }

  client, err := posthog.NewWithConfig(posthogKey, config)
  if err != nil {
    panic("oh no")
  } 
  defer client.Close()

  t := 0
  for i := 0; i <= *numbEvents; i++ {
    s := rand.Intn(*funnelDepth)
    userid := rand.Intn(*userCount)
    user := fmt.Sprintf("user-%d", userid)  
    plan := "Enterprise"
    friends := rand.Intn(42) 
    for i := 0; i <= s ; i++ {
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
      time.Sleep(time.Duration(*sleepTime) * time.Millisecond)
    } 
  }
}