The Uni programming language is a simple programming language for educational purposes.

## How to install
In order to build it from the source code, you need `git`, and `go` version 1.19 or later. Run the following commands to build from the source code:
```sh
git clone git@github.com:parsaakbari1209/university.git
cd university/interpreter
go build -o uni
```
---
## How to use
The uni-lang uses `.uni` suffix(e.g. `main.uni`). Use the following command to execute a uni-lang source code:
```sh
// Linux
./uni main.uni
```
---
## Syntax
### Comments
```
# This is a comment!
```
### Boolean
```
true
!false
true or false
true and true
false == false
true != false
```
### Number
```
1
-1.0
1.0 + 2
1.0 - 2
1.0 * 2
1.0 / 2
1.0 < 2
1.0 > 2
1.0 <= 2
1.0 >= 2
1.0 == 2
1.0 != 2

```
### String
```
"Hello World!"
"Hello" + " " + "World" + "!"
```
### Variable
```
var a = 0
a = 0.0
a = "Hello World!"

```
### Array
```
var num = [0, 1, 2]
num[0]

var str = ["Hello", "World", "!"]
str[0]

var mix = [1, "Hello", 1.5, "World"]
mix[0]
```
### Map
```
var data = {"slug": "Hello World!", "version": 1}
data["slug"]
```
### Condition
```
if a == b {
    #...
}

if a == b {
    #...
} else {
    #...
}
```
### Loop
```
while true {
    #...
}

for k, v in "Hello World" {
    #...
}

for k, v in ["Hello", "World", "!"] {
    #...
}

for k, v in {"one": 1, "two": 2} {
    #...
}
```
### Function
```
fn sum(a, b) {
    return a + b
}
sum(1, 2)
```
### Built-in
```
len("Hello World!")
len(["Hello", "World", "!"])
len({1: "Hello", 2: "World", 3: "!"})
print("Hello World!")
println("Hello World!")
```
---
## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.  
Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
