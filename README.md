## Graph Layout Algorithms in Go

This module provides algorithms for graph visualization in native Go.
As of 2021-11-20, virtually all graph visualization algorithms are bindings to Graphviz dot code which is in C.
This module attempts to provide implementation of latest and best graph visualization algorithms from scratch in Go.
However, given this is very complex task this is work in progress.

## Features

- [x] gonum Isomap
- [x] gonum Eades
- [x] Kozo Sugiyama layers strategy
- [ ] Brandes-Köpf horizontal layers assignment [80% done]
- [ ] Graphviz dot layers algorithm [80% done]
- [x] Gravity force
- [x] Spring force
- [ ] Kozo Sugiyama Magnetic Force
- [ ] Metro Style edges
- [ ] Ports for edges
- [ ] Spline edges
- [ ] Collision avoidance (dot) edge path algorithm

## Contributions

Yes please. These algorithms are hard. If you can, help to finish implementing any of above!

## References

- [Wiki Layered Graph Drawing](https://en.wikipedia.org/wiki/Layered_graph_drawing)
- ["Handbook of Graph Drawing and Visualization"](https://cs.brown.edu/people/rtamassi/gdhandbook/), Roberto Tamassia, Brown, Ch.13
- ["A Technique for Drawing Directed Graphs"](https://ieeexplore.ieee.org/document/221135), Emden R. Gansner Eleftherios Koutsofios Stephen C. North Kiem-Phong Vo, AT&T Bell Laboratories, 1993
- ["Fast and Simple Horizontal Coordinate Assignment"](https://link.springer.com/content/pdf/10.1007/3-540-45848-4_3.pdf), U. Brandes, Boris Köpf, 2002
- "Methods for visual understanding of hierarchical system structures", Sugiyama, Kozo; Tagawa, Shôjirô; Toda, Mitsuhiko, 1981
- "Graph Drawing by the Magnetic Spring Model", Kozo Sugiyama, 1995
