# traffic

[![GoDoc](https://godoc.org/github.com/spencer-p/traffic?status.svg)](https://godoc.org/github.com/spencer-p/traffic)

## Overview

Package traffic implements an abstract simulation of traffic. A simulation is
comprised of a collection of Edges and Agents. A graph is constructed from the
Edges, which the Agents travel on over time. The traffic library separates
complex path finding and meticulous agent management from the actual edge and
agent implementations. This, and the interfaces for Agent and Edge, allow great
creative freedom and quick development times for any simple traffic simulation.
