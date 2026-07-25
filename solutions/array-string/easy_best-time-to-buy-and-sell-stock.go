// Best Time to Buy and Sell Stock [Easy]
// https://leetcode.com/problems/best-time-to-buy-and-sell-stock/
// Solved: 2026-07-25

func maxProfit(prices []int) int {
    var maxprofit int
    var dayprofit int 
    mincost:=prices[0]
    for i := 1 ; i <len(prices) ; i++{
        if mincost  > prices[i]{
            mincost = prices [i]
        }
        dayprofit = prices[i]-mincost
        if dayprofit > maxprofit {
            maxprofit = dayprofit 
        }
    }
    return maxprofit
}
