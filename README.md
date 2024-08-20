# Da Vinci Script (DAVI) -  Superfast web language based on GO 

Davi is a superfast web language based on GO. It is a simple and easy to use language that can be used to create web applications. Davi is a compiled language, which means that it is converted into machine code before it is run. This makes it very fast and efficient.

```php
<?davi

// Simple Hello World program

echo("Hello, World!");

?>
```

## Features
- Simple and easy to use
- Fast and efficient
- Supports multiple platforms
- Similar syntax to PHP
- Support for websockets
- Support for databases
- Support for RESTful APIs
- Support for JSON
- Easy to learn

## Installation
To install Davi, you need to download last release here.

## Getting Started
To get started with Davi, you need to create a new file with the .davi extension. You can then write your code in this file and run it using the Davi interpreter.

Here is an example of a simple Davi program:

```php
<?davi

// This is a comment
$test = "Hello, World!";
echo($test);
echo("Time is:", time());
echo("Test calculation: 1+2*3=", 1 + 2 * 3);
echo(1 + 2 * 3);


// Variable declaration
$timeHandler = function() {
    $time = time();
    return($time);
}
echo($timeHandler());


// Variable assignment
$calculationHandler = function() {
    return(5 + 5);
}
echo $calculationHandler();
echo time();


// Array declaration
$names = ["John", "Doe", "Jane", "Doe"];
echo($names[0]);


// If statement
$age = 30;
if ($age > 18) {
    echo("You are an adult");
} else {
    echo("You are a child");
}


// For loop
$list = ["Bozhidar", "Veselinov", "Slaveykov", "Asenov"];
sort($list);
echo($list);
sort($list, lower);
for ($x in $list) {
    echo($x);
}


// Define a class
class Person {
    public function greet() {
        echo("Hello, my name is Bozhidar!");
    }
}

// Create an instance of the class
$person = new Person();
$person->greet();

// output: "Hello, my name is Bozhidar!"


?>
```

To run this program, save it to a file called hello.davi and run the following command:

```bash
davi hello.davi
```

## Documentation
For more information on how to use Davi, you can check out the official documentation here.

## Contributing
If you would like to contribute to Davi, you can do so by forking the repository and submitting a pull request. You can also report any issues or bugs that you find by opening an issue on the repository.
