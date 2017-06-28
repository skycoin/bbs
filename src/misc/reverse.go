package misc

func ReverseStringSlice(lst []string) chan string {
	ret := make(chan string)
	go func() {
		for i := range lst {
			ret <- lst[len(lst)-1-i]
		}
		close(ret)
	}()
	return ret
}

func ReverseIntSlice(lst []int) chan int {
	ret := make(chan int)
	go func() {
		for i := range lst {
			ret <- lst[len(lst)-1-i]
		}
		close(ret)
	}()
	return ret
}
