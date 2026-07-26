// Remove Element [Easy]
// https://leetcode.com/problems/remove-element/
// Solved: 2026-07-26

func removeElement(nums []int, val int) int {
  var k int
  for i :=0 ; i < len(nums) ; i++{
    if nums[i]!=val {
        nums[k]=nums[i]
        k+=1
    }else{
        continue
    }
  }
  return k
}
