# traffic

[![GoDoc](https://godoc.org/github.com/spencer-p/traffic?status.svg)](https://godoc.org/github.com/spencer-p/traffic)

## Overview

From the GoDoc:

> Package traffic implements an abstract simulation of traffic. A simulation is
> comprised of a collection of Edges and Agents. A graph is constructed from
> the Edges, which the Agents travel on over time. The traffic library
> separates complex path finding and meticulous agent management from the
> actual edge and agent implementations. This, and the interfaces for Agent and
> Edge, allow great creative freedom and quick development times for any simple
> traffic simulation.

## Background

### Why

I wrote this project as part of my final project in a class about
[futorology](https://en.wikipedia.org/wiki/Futures_studies). The original goal
was to simulate a representation of Los Angeles traffic. Unfortunately, this
library never got to work with a dataset that big because we lacked the
expertise in GIS to make it happen. However, this project was integral in
gathering a large amount of data that we used to write a paper about the train
system in LA.

### How

This library treats traffic simulation quite literally. Each agent is actually
computing a shortest path to their destination on the graph, and the idea is
that this mimics how real people might consult their phones for a fastest path
(although this library does not require that edge weights correspond to literal
time).  Over time, there is a lot of emergent behaviour as the edge weights
dynamically change with agents travelling along them. This then allows the
designers to ask high level questions about what kinds of trade-offs lessen
congestion the most.

Modelling traffic like this is computationally expensive, and how we'd be able
to pull this off with a graph that represents all of LA was a concern from the
start. I chose to write this in Go because of these concerns.  Go's speed
relative to its ease of development and garbage collection was a strong factor
in choosing it.  Go also has fantastic builtin concurrency support which I
wanted to take advantage of.

In the end, the real meat of this library is in running Dijkstra's algorithm
repeatedly across multiple cores as fast as possible. We never reached a point
where additional optimisations or heuristics were necessary. With the inputs we
were working with, generating large amounts of data and graphics in a few
seconds was a breeze.

### More

This project is the core of several other auxiliary programs for the same
project.

 * [bussim](https://github.com/spencer-p/bussim) - the main client program
 * [webtraffic](https://github.com/spencer-p/webtraffic) - a webserver that acts
	 as a frontend for bussim
 * An as-of-yet unpublished program to generate heatmaps from text data
