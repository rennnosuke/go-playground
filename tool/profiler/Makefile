.PHONY: build-prof-flip-filter ls-prof

build-prof-flip-filter:
	go test -bench FlipFilter -benchmem -o pprof/test.bin -cpuprofile pprof/cpu.out -memprofile pprof/mem.out

ls-prof-flip-filter:
	go tool pprof -text pprof/test.bin pprof/cpu.out
	go tool pprof -text pprof/test.bin pprof/mem.out

print-prof-flip-filter:
	go tool pprof -pdf pprof/test.bin pprof/cpu.out > pprof/cpu.pdf
	go tool pprof -pdf pprof/test.bin pprof/mem.out > pprof/mem.pdf