package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"time"
)

func readfile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("ファイル読み込みでエラーが発生しました。")
		panic(err)
	}
	return string(content)
}

func convertData(s string) (int, [][]int) { //s=読み込んだファイルの文字列
	spliteds := strings.Split(s, "\n")
	n, err := strconv.Atoi(spliteds[0])
	if err != nil {
		fmt.Println("ファイルの形式が正しくありません。")
		panic(err)
	}
	dt := spliteds[1:]
	dests := make([][]int, 0, n)
	for i := 0; i < n; i++ {
		t := strings.Split(dt[i], " ")
		dest := make([]int, 0, 2)
		for j := 0; j < 2; j++ {
			t2, _ := strconv.Atoi(t[j])
			dest = append(dest, t2)
		}
		dests = append(dests, dest)
	}
	return n, dests
}

func saleman(n int, dests [][]int, plist [][]int) float64 {
	min := math.Inf(1)
	//fmt.Println(plist)
	for i := 0; i < len(plist); i++ {
		t := 0.0
		t2 := plist[i] //対象のPlist ex[0,1,2,3,4]
		for j := 0; j < n-1; j++ {
			idx := t2[j]
			idx_d := t2[j+1]
			t += calc_dist(idx, idx_d, dests, plist)
		}
		t += calc_dist(t2[n-1], t2[0], dests, plist)

		//fmt.Println(t)
		if min > t {
			min = t
		}
	}
	return min
}

func calc_dist(idx int, idx_d int, dests [][]int, plist [][]int) float64 {
	x := float64(dests[idx][0])
	y := float64(dests[idx][1])
	dx := float64(dests[idx_d][0])
	dy := float64(dests[idx_d][1])
	return math.Sqrt(math.Pow(math.Abs(dx-x), 2) + math.Pow(math.Abs(dy-y), 2))
}

func parallel_saleman_controller(n int, dests [][]int, plist [][]int, parallel_num int) float64 {
	ans := math.Inf(1)
	ch1 := make(chan float64)
	ch2 := make(chan float64)
	go p_saleman(n, dests, plist, 0, len(plist)/2, ch1)
	go p_saleman(n, dests, plist, len(plist)/2+1, len(plist)-1, ch2)
	a1 := <-ch1
	a2 := <-ch2
	if a1 > a2 {
		ans = a2
	} else {
		ans = a1
	}
	return ans
}

func p_saleman(n int, dests [][]int, plist [][]int, head int, end int, ch chan float64) {
	min := math.Inf(1)
	//fmt.Println(plist)
	for i := head; i < end; i++ {
		t := 0.0
		t2 := plist[i] //対象のPlist ex[0,1,2,3,4]
		for j := 0; j < n-1; j++ {
			idx := t2[j]
			idx_d := t2[j+1]
			t += calc_dist(idx, idx_d, dests, plist)
		}
		t += calc_dist(t2[n-1], t2[0], dests, plist)

		//fmt.Println(t)
		if min > t {
			min = t
		}
	}
	ch <- min
}

//Permute関連の実装はhttps://raahii.github.io/posts/permutations-in-go/を参照
func Permute(nums []int) [][]int {
	n := factorial(len(nums))
	ret := make([][]int, 0, n)
	permute(nums, &ret)
	return ret
}

func permute(nums []int, ret *[][]int) {
	*ret = append(*ret, makeCopy(nums))

	n := len(nums)
	p := make([]int, n+1)
	for i := 0; i < n+1; i++ {
		p[i] = i
	}
	for i := 1; i < n; {
		p[i]--
		j := 0
		if i%2 == 1 {
			j = p[i]
		}

		nums[i], nums[j] = nums[j], nums[i]
		*ret = append(*ret, makeCopy(nums))
		for i = 1; p[i] == 0; i++ {
			p[i] = i
		}
	}
}

func factorial(n int) int {
	ret := 1
	for i := 2; i <= n; i++ {
		ret *= i
	}
	return ret
}

func makeCopy(nums []int) []int {
	return append([]int{}, nums...)
}

func main() {
	//n=巡回する都市の数
	//座標の２重配列[[x0,y0],[x1,y1],[x2,y2],[x3,y3]]
	fmt.Println(time.Now())
	s := readfile("dataset.txt")
	n, dests := convertData(s)
	prange := make([]int, 0, n)
	for i := 0; i < n; i++ {
		prange = append(prange, i)
	}
	t := Permute(prange)
	fmt.Println("finish make plist array")
	fmt.Println(time.Now())
	ans := saleman(n, dests, t)
	fmt.Println(ans)
	fmt.Println(time.Now())
	parallel_saleman_controller(n, dests, t, 2)
	fmt.Println(time.Now())
}
