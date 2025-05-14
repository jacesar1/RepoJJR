package main
import ( 
	"strings"
	"fmt"
)

var msg = `**************************

        *    BUY NOW, SAVE 10%   *

        **************************`

func CleanupMessage(oldMsg string) string {
	return strings.Trim(strings.ReplaceAll(oldMsg, "*", ""), " ")
}


func main() {
	fmt.Println(CleanupMessage(msg))
}