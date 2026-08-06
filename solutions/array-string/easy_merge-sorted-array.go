// Merge Sorted Array [Easy]
// https://leetcode.com/problems/merge-sorted-array/
// Solved: 2026-08-06

func merge(nums1 []int, m int, nums2 []int, n int)  {
    k := m+n - 1
    i, j := m-1, n-1; 
    for i >= 0 && j >= 0{
        if nums1[i]>nums2[j]{
            nums1[k] = nums1[i]
            i-=1
            k-=1
            
        }else{
            nums1[k]=nums2[j]
            j-=1
            k-=1
        }

    }
    for j >= 0  {
        nums1[k]=nums2[j]
        k-=1
        j-=1
    }
}
