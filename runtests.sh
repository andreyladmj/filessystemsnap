go test -bench . -benchmem main_test.go
go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 main_test.go
go test -bench=.