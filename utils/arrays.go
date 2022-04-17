package utils

import "fmt"

// "=" is used for shallow copy (overload is not supported in Go) and will not
// be marked using shallow copy tag as Athena++ did. So it's your duty to
// avoid shallow copy. And it's also your duty to avoid two goroutine write
// at the same time (concurrency).

type Array[T any] struct {
	pdata_  []T
	dim_num []int
}

func AthenaArray[T any](nx ...int) Array[T] {
	var this Array[T]
	data_num := 1
	for _, i := range nx {
		this.dim_num = append(this.dim_num, i)
		data_num *= i
	}
	this.pdata_ = make([]T, data_num)
	return this
}

func DeepCopyArray[T any](other Array[T]) Array[T] {
	var this Array[T]
	this.pdata_ = make([]T, len(other.pdata_))
	copy(this.pdata_, other.pdata_)
	this.dim_num = make([]int, len(other.dim_num))
	copy(this.dim_num, other.dim_num)
	return this
}

func (this *Array[T]) GetDim(dim int) int {
	if dim > len(this.dim_num) {
		return 0
	}
	return this.dim_num[dim-1]
}

func (this *Array[T]) IsAllocated() bool {
	return this.pdata_ != nil
}

func (this *Array[T]) GetSize() int {
	if len(this.pdata_) == 0 {
		return 0
	}
	size := 1
	for _, i := range this.dim_num {
		size *= i
	}
	return size
}

// If T is a reference type, note that it will be returned with shallow copying. Remember
// that the first dimention is access using the last parameter as default (in Athena++)
func (this *Array[T]) Get(ij ...int) (T, error) {
	sum, err := this.access(ij)
	var temp T // Make sure return the zero value of type T.
	if err != nil {
		return temp, err
	}
	return this.pdata_[sum], nil
}

func (this *Array[T]) Set(value T, ij ...int) error {
	sum, err := this.access(ij)
	if err != nil {
		return err
	}
	this.pdata_[sum] = value
	return nil
}

// It's a private function.
func (this *Array[T]) access(ij []int) (int, error) {
	if len(ij) != len(this.dim_num) {
		panic("Access Array Error: Number of parameters doesn't equal to array dimention.")
	}
	var sum int
	weight := 1
	// Change order to make sure it's the same as normal (Athena++) version.
	for i, j := 0, len(ij)-1; i < j; i, j = i+1, j-1 {
		ij[i], ij[j] = ij[j], ij[i]
	}
	for i, j := range ij {
		if j > this.dim_num[i] {
			return 0, fmt.Errorf("Access Array Error: Parameter %d exceed the limit.", i)
		}
		sum += j * weight
		weight *= this.dim_num[i]
	}
	return sum, nil
}
