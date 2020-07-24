# gohog
Posthog mock data generator

This is still under dev, but feel free to adapt it


# Benchmarking

You can benchmark the events pipeline in an instance of PostHog by running:

``` bash
POSTHOG_KEY=<your_api_key>; POSTHOG_ENDPOINT=<your_posthog_url> go test -bench .
```
