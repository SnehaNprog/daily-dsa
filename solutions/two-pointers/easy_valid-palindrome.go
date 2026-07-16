// Valid Palindrome [Easy]
// https://leetcode.com/problems/valid-palindrome/
// solved 2026-07-17

import "unicode"

func isPalindrome(s string) bool {

    var results []rune
    for _ , i := range s {
        
        if unicode.IsLetter(i) || unicode.IsDigit(i){
            results = append(results , unicode.ToLower(i))

        }
        
    }
    clean := string(results)


    
    left := 0 
    right := len(clean)-1 
    for left <right {
        if clean[left]!= clean[right]{
            return false 

        }

        left ++
        right --
        
    }
    return true

        

    
    
    
}



