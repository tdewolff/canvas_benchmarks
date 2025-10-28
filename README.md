# Performance comparisons of https://github.com/tdewolff/canvas
## Boolean operation: union(europe, chile)
![union(europe,chile)](https://raw.githubusercontent.com/tdewolff/canvas_benchmarks/refs/heads/master/boolean/tdewolff.png)

Benchmarks are performed with the Natural Earth 10m resolutions of countries from the European Union and Chile, projected both to UTM 33N and 19S respectively. The boolean operation of union(Europe,Chile) with different levels of detail (different numbers of segments) is evaluated 10 times and averaged. The results are in seconds.

| Library | 255 | 447 | 935 | 1935 | 3935 | 7692 | 15763 | 34318 | 63809 | 87721 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| [ajohnson1](http://www.angusj.com/delphi/clipper/documentation/Docs/Overview/_Body.htm) | **0.0002** | 0.0005 | 0.0008 | 0.0018 | 0.0044 | 0.0103 | 0.0322 | 0.1032 | 0.2273 | 0.3355 |
| [ajohnson2](https://github.com/AngusJohnson/Clipper2) | 0.0007 | 0.0013 | 0.0025 | 0.0056 | 0.0111 | 0.0230 | 0.0500 | 0.1140 | 0.2028 | 0.2746 |
| [ioverlay](https://github.com/iShape-Rust/iOverlay) | 0.0002 | **0.0004** | **0.0007** | **0.0012** | **0.0022** | **0.0048** | **0.0115** | **0.0154** | **0.0281** | **0.0370** |
| [tdewolff](https://github.com/tdewolff/canvas) | 0.0006 | 0.0012 | 0.0019 | 0.0037 | 0.0074 | 0.0149 | 0.0341 | 0.0861 | 0.1666 | 0.2491 |

![Boolean results graph](https://raw.githubusercontent.com/tdewolff/canvas_benchmarks/refs/heads/master/boolean/results.png)

Benchmark notes:
- ajohnson1 uses a transliteration of C++ in Go and might not accurately display the speed of the original implementation
- ajohnson2 uses the original implementation as an external library where all the work is done, this includes a single (negligible) calling overhead
- ioverlay also provides [benchmarks](https://ishape-rust.github.io/iShape-js/overlay/performance/performance.html) with ajohnson2 and Boost
