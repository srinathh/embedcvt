# embedcvt
Command embedcvt creates a Go program file that embeds the contents of the specified 
asset file(s) as simple byte-array variables.

## Background
embedcvt is meant for simple data embedding use cases like inlining test data,
inserting a few assets etc and does not do any compression, file
system emulation, asset discovery etc. For more complex use cases,
see tools like [go-bindata](https://github.com/jteeuwen/go-bindata).

## Usage
```
embedcvt [-p packageName] asset1 [asset2]...

Options:
  -p string
    	package name for the generated Go program file. (default "main")
  asset string
    	paths to the files to embed (required)
```

## Installation
```
go install github.com/srinathh/embedcvt
```

## License
Apache 2.0 License
