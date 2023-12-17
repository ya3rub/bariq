# About Bariq

Bariq is a Golang interpreter for a simple functional programming language. The language incorporates immutable data structures and includes essential data structures such as arrays and hash tables. I implemented a parser from scratch and this all guided by the awesome book by Thorstem ball [WaIGo](https://interpreterbook.com/)

The interpreter supports advanced features such as closures and first-class functions, which allow for more expressive and flexible programming. Additionally, I intend to add asynchronous support, type checking, and module functionality.



# What is new ?

### Asyncronous Prograaming Support  via `async` and `await`

Ex:

```go
  let s = async fn(x) {sleep(5); 5};
  let task = s(5);
  puts("4");
  await(task)
```

output

```
   4
   5
```

#### how does it work ?

- when parser read `async` keyword, it marks that function as Async.
- async function returns a  variable of type`object.Task` which has `Spawned` of type `*sched.Task`,
- the schedular is the responsible for spawning tasks given by the evaluator.
- `await` accepts an epxression, if the evaluated expression is of type `Task` it calls the schedular spwaned task attach to that task in order to  `await` it 



### Generators

generators in bariq is inspired by how Javascript handles generators,

EX:

```js
let s =  fn gen () { yield 2;yield 0;yield 6;yield 1; };
			let genr = s();
			next(genr);
			next(genr);
			next(genr);
// s: object.Iteration{ Val : 6, Done: False}
```

a fancy Ex:

```js
			let w = fn (){1}
			let s =  fn gen () {
					let q = w()
					yield q;
				};
			let genr = s();
			next(genr);
// s: object.Iteration{ Val : 1, Done: False}
```

### how does it work ?

- when parser read `gen` keyword, it marks that function as generator.
- `Generator` type has a reference to the function in addition to the env and index to indicate the position where it left the function.
- when `next()` is called, it checks if the object is of type `Generator` and starts evaluting the function body starting from the `Index` passing the `Env` that passed first at creating the generator.
- if `yield` keyword is found, it sets the `Index` and the `Value` of the generator (you can call it a frame) and reutrn an `Iteration` Object having the state of being `Done` or not and the current value of the generator (the current frame)



