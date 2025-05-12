# ReplayAnalyzer (for BL Replays)

## Info

I made this small project in my free time; it certainly is not perfect.

## Features

- Generate a JD config approximation using a players replays
  - Runs regression and builds a [model](https://github.com/HalloTheEngineer/replayAnalyzer/blob/master/example/2169974796454690.jpg), dividing the data into two curves using kmeans first
- more if requested...

## Command Line Arguments

- `fetch [optional player id]` - fetches player replays from BeatLeader
- `generate`
  - `jd-config [optional player id]` - generates a config approximation
- `help` - displays a help message
