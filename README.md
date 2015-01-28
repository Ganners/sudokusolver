Sudoku Solver in Go (Golang)
============================

A package which will solve Sudoku puzzles pretty quickly. Uses a depth first
approach.

It is currently using a basic recursive approach, it should be easy to modify it
to use goroutines and channels (goroutine for each search, channel receives once
complete). It should also possibly simplify the implementation - though will end
up spawning a lot of threads which will need to be garbage collected.

There may be ways to improve it, I have made modifications and benchmarked them,
current hard puzzles can take about 0.03 seconds, just simplifying it slightly
can yield solutions in 0.001 seconds. I'll be curious to see if spawning threads
will actually make much of a difference!

- Mark Gannaway
