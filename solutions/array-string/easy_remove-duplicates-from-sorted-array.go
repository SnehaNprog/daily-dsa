// Remove Duplicates from Sorted Array [Easy]
// https://leetcode.com/problems/remove-duplicates-from-sorted-array/
// Solved: 2026-07-26

func removeDuplicates(nums []int) int {
    k:=1
    for i :=1 ; i<len(nums);i++{
        if nums[i-1]!=nums[i] {
            nums[k]=nums[i]
            k+=1
        }else{
            continue
        }
    }
    return k
}
