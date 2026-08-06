// Majority Element [Easy]
// https://leetcode.com/problems/majority-element/
// Solved: 2026-08-06

func majorityElement(nums []int) int {
    freq := make(map[int]int)
    for _,num := range nums{
        freq[num]++
        if freq[num]>len(nums)/2{
            return num

        }
    }
    return 0
}
