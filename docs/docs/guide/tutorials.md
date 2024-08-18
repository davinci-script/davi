# Tutorials
This section contains tutorials that demonstrate how to use the `davi` script to create websites and web applications.

::: warning
DaVinci script is still in development, so some features may not work as expected.
:::

### Syntax
The DaVinci script is a PHP-like language, so it has a similar syntax to PHP. The DaVinci script uses the `<?davi` tag to start the script and the `?>` tag to end the script. 

A DAVI script can be placed anywhere in the document.
A DAVI script starts with ```<?davi``` and ends with ```?>```:
```php
<?davi
// Your code here
?>
```

### Hello World
```php
<?davi

echo "Hello World!";

?>
```

### Variables

Simple variables
```php
<?davi
// Strings
$name = "John Doe";
echo "Hello, $name!";

// String concatenation
$age = 30;
echo("You are ",$age," years old");

// String concatenation
$price = "19.99";
echo("The price is", $price);
?>
```
Closure variable
```php
<?davi

// Closure variable
$timeHandler = function() {
    $time = time();
    return($time);
}
echo($timeHandler());

?>
```

### Functions

#### Custom Function declaration
```php
<?davi

// Function declaration
function sayHello($name) {
    echo("Hello,", $name, "!");
}

sayHello("John Doe");

?>
```

#### Built-in Functions

Times functions
```php
<?davi

// Built-in functions
echo("Current time: ", time());

?>
```

Math calculations
```php
<?davi

// In Echo with string 
echo("Test calculation: 1+2*3=", 1 + 2 * 3);

// Only math calculation
echo(1 + 2 * 3);

?>
```

### If Statement
```php
<?davi
$age = 30;
if ($age > 18) {
    echo("You are an adult");
} else {
    echo("You are a child");
}
?>
```

### Arrays
```php
<?davi

// Array declaration
$names = ["John", "Doe", "Jane", "Doe"];

echo($names[0]);

?>
```

### Sort Functions
```php
<?davi

$list = ["Bozhidar", "Veselinov", "Slaveykov", "Asenov"];

sort($list);
echo($list);

sort($list, lower);
echo($list);

?>
```


### For Loop
```php

<?davi

$list = ["Bozhidar", "Veselinov", "Slaveykov", "Asenov"];

for ($x in $list) {
    echo($x);
}

?>
```
