package jwts

import (
	"fmt"
	"testing"
)

func TestGenToken(t *testing.T) {
	token, _ := GenToken(JwtPayload{UserId: 1, Username: "test", Role: "1"}, "lhy", 1)
	fmt.Println(token)
}
func TestParseToken(t *testing.T) {
	payload, _ := ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJyb2xlIjoiMSIsImV4cCI6MTcxNTQxMjg3Mn0.4mADPZoLxkKHvysev-t5c6Pld2foXeeJxeVp0T-gg6c", "lhy")
	fmt.Println(payload)
}
