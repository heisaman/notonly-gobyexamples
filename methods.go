package main

import "fmt"

type rectOld struct {
	width, height float64
}

func (r rectOld) area() float64 {
	return r.width * r.height
}
func (r rectOld) perim() float64 {
	return 2*r.width + 2*r.height
}

func main() {
	r := rectOld{width: 10, height: 5}

	// 调用方法时，Go 会自动处理值和指针之间的转换。 想要避免在调用方法时产生一个拷贝，或者想让方法可以修改接受结构体的值， 你都可以使用指针来调用方法
	fmt.Println("area: ", r.area())
	fmt.Println("perim:", r.perim())

	rp := &r
	fmt.Println("area: ", rp.area())
	fmt.Println("perim:", rp.perim())
}