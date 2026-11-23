# ![Boids Logo](Documentation/Logo/boid.png) Boids 
**Concurrent Development - Software Development - Year 4, Semester 1 - 2026**  
**South East Technological University**

## Project Description
Boids is an agent-based flocking simulator written in Go. 

The current 2D simulation is in the `boids` folder, split into two packages:

- `simulation`: Contains the boid behaviour and simulation settings
- `ui`: Draws the UI controls 

## Installation Guide
1. Clone the repository:  
   ``` git clone https://github.com/amelia720/Con_Dev_Contagion_Project.git```
2. Navigate to folder:   
   ```cd Con_Dev_Contagion_Project```
## How to Run
**Required: Go 1.27.1 or newer**

```
cd boids
go run .
```

## Project Files
- `boids/program.go`: starts the simulation and draws each frame
- `boids/simulation`: covers boid behaviour, settings and environment
- `boids/ui`: displays simulation controls and grid 
- `boids/tests`: detection range tests and performance benchmarks
- `Documentation`: project images and documentation, jupyter notebooks that run and analyse the benchmark tests

## Tests
- **Detection range:** checks boids inside, exactly at and outside the range.

From the `boids` folder, run:
```bash
go test ./tests/boids_test.go -v
```

## Contributions
- Adam:
- Dorian:
- Amelia: Benchmark tests and results, detection range tests and contribution to GoDoc
- Mark: Handled project architecture, parallelism, plotting serial test results and contributed to GoDoc

## Benchmarks
Measures the time to run the three flocking rules and calculate each bird's next movement. Benchmarks use 25, 50, 100, 500 and 1000 birds.

Serial processes birds one at a time. Parallel divides the work between workers. Both start with the same flock and settings. Grid optimisation and obstacles are off. Creating, resetting and drawing the boids are not timed.

From the `boids` folder, run:

```bash
go test ./tests -run '^$' -bench '^Benchmark(Serial|Parallel)$' -benchtime=100x -count=5
```

Each case runs five times, with 100 flock updates per run. Results were collected on macOS with an Intel Core i7-8569U. The table shows the **median time in milliseconds** for one update of the whole flock.

| Boids | Serial | Parallel: 1 worker | Parallel: 2 workers | Parallel: 4 workers | Parallel: 8 workers | Parallel: 10 workers |
|---:|---:|---:|---:|---:|---:|---:|
| 25 | **0.0126** | 0.0216 | 0.0184 | 0.0194 | 0.0214 | 0.0189 |
| 50 | 0.0676 | 0.0681 | 0.0419 | **0.0358** | 0.0386 | 0.0422 |
| 100 | 0.2184 | 0.2505 | 0.1415 | **0.1032** | 0.1048 | 0.1088 |
| 500 | 4.9378 | 4.9500 | 2.9558 | **2.3857** | 2.6115 | 2.7108 |
| 1000 | 19.3111 | 19.6376 | 12.0778 | **10.3458** | 13.4403 | 15.0843 |

See [BenchmarkResults.md](Documentation/BenchmarkResults.md) for the full benchmark output from all five runs.

### Speed-up and Efficiency

- **Absolute speed-up** = serial time ÷ parallel time.
- **Relative speed-up** = parallel time with 1 worker ÷ parallel time with the chosen number of workers.
- **Efficiency** = speed-up ÷ number of parallel workers × 100%.

**Note:** We calculated absolute speed-up by dividing the serial median time by the best parallel median time: 2 workers for 25 boids and 4 workers for every larger flock.

| Boids | Serial (ns/op) | Parallel: 1 worker (ns/op) | Fastest parallel median (ns/op) | Workers in fastest run | Absolute speed-up | Relative speed-up | Absolute efficiency | Relative efficiency |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| 25 | 12,577 | 21,599 | 18,415 | 2 | 0.68× | 1.17× | 34.1% | 58.6% |
| 50 | 67,607 | 68,098 | 35,825 | 4 | 1.89× | 1.90× | 47.2% | 47.5% |
| 100 | 218,355 | 250,529 | 103,234 | 4 | 2.12× | 2.43× | 52.9% | 60.7% |
| 500 | 4,937,820 | 4,949,988 | 2,385,662 | 4 | 2.07× | 2.07× | 51.7% | 51.9% |
| 1000 | 19,311,108 | 19,637,630 | 10,345,751 | 4 | 1.87× | 1.90× | 46.7% | 47.5% |

## Technologies Used
- GoLang
- Raylib with Golang bindings
- Jupyter Notebooks, matplotlib for analysing benchmark tests



## Authors

Adam Noonan,               Amelia Hamulewicz,            Dorian Nowacki,        Mark Lambert
