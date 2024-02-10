 # Compiler in go
 
## Following Thorsten's Ball : https://compilerbook.com/
## This is the sequel to : https://interpreterbook.com/

### Prequel repository here : https://github.com/Neal-C/interpreter-in-go

 
###  Learned and done :
- Defined OpCodes
- Built bytecode
- Made a Stack machine with a symbol table
- Written a compiler
- Written a Virtual Machine (VM)

All done in TDD fashion


Fun facts 😀 :
- the LLVM & Clang projects currently consist of around 3 million lines of code. 
- The GNU Compiler Collection, GCC, is even bigger. 15 million lines of code 😱


Amazing read. 

Huge thanks to [Thorsten Ball](https://github.com/mrnugget)

To try it out:

Requirements : Go >= 1.21.3 or Docker

```shell
git clone git@github.com:Neal-C/compiler-in-go.git
cd compiler-in-go
go run . 
# build executable command : go build -o ./bin/compiler-in-go
# run executable : ./bin/compiler-in-go
```

#### or trying via Docker by running my image

```shell
git clone git@github.com:Neal-C/compiler-in-go.git
cd compiler-in-go
docker build -t nealc:compiler-in-go .
# builds the image
docker run -it --name nealc-compiler nealc:compiler-in-go
# runs the image
```

- features include : common data types, recursive functions, and closures
- Check the test cases in ./**/*_test.go files to see what other behaviors and features are supported

```shell
puts("Hello!")
# Hello!
# null
puts(1234)
# 1234
# null
let people = [{"name": "Alice", "age": 24},{"name": "Neal-C", "age": 999}, {"name": "Anna", "age": 22}];
people[0]["name"];
# Alice
len(people)
# 3
first(people)
# {"name": "Alice", "age": 24}
last(people)
# {"name": "Anna", "age": 22} 
if (true) { 42 } else { "never" };
# 42
if (false) { 42 } else {  "false" }
# "false"
let a = 20
let b = 22;
a == b
# false
a + b;
# 42
let sum = fn(x,y) { return x + y };
# CLOSURE[0xc000140060] 
sum(a,b)
# 42
flex
# Whoops! compilation failed:
# undefined variable : flex
```

To benchmark speed difference between an interpreter and a byte code Virtual Machine:

(requires a go local installation)

```shell
git clone git@github.com:Neal-C/compiler-in-go.git
cd compiler-in-go
go run ./benchmark --engine=eval
# engine=eval, result=9227465, duration=14.88122763s
go run ./benchmark --engine=vm
# engine=vm, result=9227465, duration=3.337770149s
```


