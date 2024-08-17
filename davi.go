package main

import (
	"fmt"
	. "github.com/DavinciScript/Davi/lexer"
	"github.com/hokaccha/go-prettyjson"
)

func main() {

	fmt.Println("Running davi.go")

	// Run the lexer
	lexer := NewLexer([]byte(`

<?davi

class Dog
{
    public $name = "Rex";
    public $breed = "Golden Retriever";

    public function bark()
    {
        echo "Woof!";
    }

}

$dog = new Dog();
echo $dog->name; // Rex
echo $dog->breed; // Golden Retriever
$dog->bark(); // Woof!

?>



`))
	for {
		pos, tok, val := lexer.Next()
		if tok == EOF {
			break
		}
		data := map[string]interface{}{
			"posLine":   pos.Line,
			"posColumn": pos.Column,
			"token":     tok,
			"value":     val,
		}

		formatter := prettyjson.NewFormatter()
		output, _ := formatter.Marshal(data)
		fmt.Println(string(output))

		if tok == ILLEGAL {
			break
		}
	}

}
