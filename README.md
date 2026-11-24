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
   ``` git clone https://github.com/Team-A-D-A-M/Con_Dev_Contagion_Project.git```
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
| 25 | **0.0123** | 0.0182 | 0.0126 | 0.0152 | 0.0165 | 0.0157 |
| 50 | 0.0704 | 0.0652 | 0.0391 | **0.0309** | 0.0344 | 0.0356 |
| 100 | 0.2078 | 0.2352 | 0.1483 | **0.1098** | 0.1259 | 0.1245 |
| 500 | 5.1581 | 5.8250 | 3.3775 | **2.9472** | 3.4078 | 3.1867 |
| 1000 | 21.2251 | 24.0790 | 14.0662 | **11.1900** | 15.8381 | 13.5273 |

See [BenchmarkResults.md](Documentation/BenchmarkResults.md) for the full benchmark output.

### Speed-up and Efficiency

- **Absolute speed-up** = serial time ÷ parallel time.
- **Relative speed-up** = parallel time with 1 worker ÷ parallel time with the chosen number of workers.
- **Efficiency** = speed-up ÷ number of parallel workers × 100%.

**Note:** We calculated absolute speed-up by dividing the serial median time by the best parallel median time: 2 workers for 25 boids and 4 workers for every larger flock.

| Boids | Serial (ns/op) | Parallel: 1 worker (ns/op) | Lowest parallel median (ns/op) | Parallel workers | Absolute speed-up | Relative speed-up | Absolute efficiency | Relative efficiency |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| 25 | 12,257 | 18,223 | 12,598 | 2 | 0.97× | 1.45× | 48.6% | 72.3% |
| 50 | 70,414 | 65,167 | 30,897 | 4 | 2.28× | 2.11× | 57.0% | 52.7% |
| 100 | 207,789 | 235,158 | 109,791 | 4 | 1.89× | 2.14× | 47.3% | 53.5% |
| 500 | 5,158,066 | 5,824,954 | 2,947,194 | 4 | 1.75× | 1.98× | 43.8% | 49.4% |
| 1000 | 21,225,115 | 24,078,982 | 11,190,008 | 4 | 1.90× | 2.15× | 47.4% | 53.8% |

## Technologies Used
- GoLang
- Raylib with Golang bindings
- Jupyter Notebooks, pandas and matplotlib for analysing benchmark tests



## Authors

Adam Noonan,               Amelia Hamulewicz,            Dorian Nowacki,        Mark Lambert
