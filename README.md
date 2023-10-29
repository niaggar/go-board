# Galton Board Simulation in Go

> Implementation for solving an IPT 2023 problem

![Galton Board](./img/galton-board.png)

This is a Go application that simulates a Galton Board, also known as a "bean machine" or "quincunx." The Galton Board
is a device that demonstrates the central limit theorem and the normal distribution of particles as they bounce off pins
and fall into bins.

> # IPT Problem
>
> Dropping a set of beads on a board with evenly distributed pegs results in a binomial distribution. Is it possible to
> generate other kinds of distributions by varying some parameters (pegs size, pegs distribution, bead format, etc.)? Is
> it possible to achieve a distribution that does not obey the central limit theorem in an i.i.d. scenario? What happens
> to the distribution when one makes the board vibrate?
>

## Usage

1. Make sure you have Go installed on your system.

2. Clone this repository:

   ```bash
   git clone https://github.com/niaggar/go-board

3. Create a folder named "data/configs" where you will save the configurations for the experiments you will run.

4. Create your first configuration using the "model.config.json" file as a base, and modify fields such as the number of
   particles (performance could be affected).

5. Execute.
   ```bash
   go run main.go

Another way is to run the console files corresponding to the system you need to compile the application, and then open
the executable.