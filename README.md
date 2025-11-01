# Performance comparisons of https://github.com/tdewolff/canvas
## Boolean operation: union(europe, chile)
![union(europe,chile)](https://raw.githubusercontent.com/tdewolff/canvas_benchmarks/refs/heads/master/boolean/tdewolff.png)

Benchmarks are performed with the Natural Earth 10m resolutions of countries from the European Union and Chile, projected both to UTM 33N and 19S respectively. The boolean operation of union(Europe,Chile) with different levels of detail (different numbers of segments) is evaluated 10 times and averaged. The results are in seconds.

| Segments | 255 | 447 | 935 | 1935 | 3935 | 7692 | 15763 | 34318 | 63809 | 87721 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| [ajohnson1](http://www.angusj.com/delphi/clipper/documentation/Docs/Overview/_Body.htm) | **0.0002** | 0.0006 | 0.0008 | 0.0021 | 0.0047 | 0.0103 | 0.0320 | 0.1022 | 0.2185 | 0.3362 |
| [ajohnson2](https://github.com/AngusJohnson/Clipper2) | 0.0006 | 0.0012 | 0.0026 | 0.0052 | 0.0112 | 0.0224 | 0.0512 | 0.1117 | 0.2041 | 0.2693 |
| [ioverlay](https://github.com/iShape-Rust/iOverlay) | 0.0003 | **0.0004** | **0.0007** | **0.0011** | **0.0026** | **0.0050** | **0.0118** | **0.0157** | **0.0285** | **0.0368** |
| [tdewolff](https://github.com/tdewolff/canvas) | 0.0006 | 0.0010 | 0.0017 | 0.0035 | 0.0071 | 0.0144 | 0.0337 | 0.0902 | 0.1681 | 0.2563 |

![Boolean results graph](https://raw.githubusercontent.com/tdewolff/canvas_benchmarks/refs/heads/master/boolean/results.png)

Benchmark notes:
- ajohnson1 uses a transliteration of C++ in Go and might not accurately display the speed of the original implementation
- ajohnson2 uses the original implementation as an external library where all the work is done, this includes a single (negligible) calling overhead
- ioverlay also provides [benchmarks](https://ishape-rust.github.io/iShape-js/overlay/performance/performance.html) with ajohnson2 and Boost
