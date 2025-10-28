# Performance comparisons of https://github.com/tdewolff/canvas
## Boolean operation: union(europe, chile)
![union(europe,chile)](https://raw.githubusercontent.com/tdewolff/canvas_benchmarks/refs/heads/master/boolean/tdewolff.png)

Benchmarks are performed with the Natural Earth 10m resolutions of countries from the European Union and Chile, projected both to UTM 33N and 19S respectively. The boolean operation of union(Europe,Chile) with different levels of detail (different numbers of segments) is evaluated 10 times and averaged. The results are in seconds.

| Library | 255 | 447 | 935 | 1935 | 3935 | 7692 | 15763 | 34318 | 63809 | 87721 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| [ajohnson1](http://www.angusj.com/delphi/clipper/documentation/Docs/Overview/_Body.htm) | **0.000170** | 0.000515 | 0.000806 | 0.001813 | 0.004367 | 0.010307 | 0.032215 | 0.103214 | 0.227305 | 0.335477 |
| [ajohnson2](https://github.com/AngusJohnson/Clipper2) | 0.000650 | 0.001318 | 0.002538 | 0.005605 | 0.011127 | 0.023019 | 0.050037 | 0.113970 | 0.202780 | 0.274606 |
| [ioverlay](https://github.com/iShape-Rust/iOverlay) | 0.000249 | **0.000355** | **0.000723** | **0.001238** | **0.002222** | **0.004795** | **0.011537** | **0.015399** | **0.028063** | **0.036972** |
| [tdewolff](https://github.com/tdewolff/canvas) | 0.000572 | 0.001218 | 0.001896 | 0.003683 | 0.007419 | 0.014916 | 0.034088 | 0.086138 | 0.166645 | 0.249121 |

![Boolean results graph](https://raw.githubusercontent.com/tdewolff/canvas_benchmarks/refs/heads/master/boolean/results.png)

Benchmark notes:
- ajohnson1 uses a transliteration of C++ in Go and might not accurately display the speed of the original implementation
- ajohnson2 uses the original implementation as an external library where all the work is done, this includes a single (negligible) calling overhead
- ioverlay also provides [benchmarks](https://ishape-rust.github.io/iShape-js/overlay/performance/performance.html) with ajohnson2 and Boost
