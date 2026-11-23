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
- `Documentation`: project images and documentation

## Tests
- **Detection range:** checks boids inside, exactly at and outside the range.

From the `boids` folder, run:
```bash
go test ./tests -v
```

## Benchmarks
Measures the time to run the three flocking rules and calculate each bird's next movement. Tests use 100, 500 and 1000 birds.

Serial processes birds one at a time. Parallel divides the work between workers. Both use the same starting data and settings. Graphics are not timed.

Currently, only the serial benchmark is available. From the `boids` folder, run:
```bash
go test ./tests -run '^$' -bench '^BenchmarkSerial$' -count=5
```

Each benchmark runs five times. The median results will be used to compare serial and parallel times and calculate speed-up and efficiency.

Parallel benchmarks: to be added 

Results and findings: to be put in documentation


## Technologies Used
- GoLang
- Raylib with Golang bindings



## Authors

Adam Noonan,               Amelia Hamulewicz,            Dorian Nowacki,        Mark Lambert
