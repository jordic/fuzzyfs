
# fu

Is a commandline search tool, that trys to traverse filesystem names, quick.
It started as a fuzzy search pattern, working with a index in RAM.. but later, 
as I get good results, I start developing a commandline tool.

But something happened, when I try to search in my home... there are around 200K files, and its almost impossible, to get good results on a reasonable time... less than 500ms

Making some tests and prof concepts, /benchs, I discouvered that the neck bottle is the hard drive io. And later I decided to keep a "cache" of index files, ( that must be updated manually ) for getting as performance as i want.

## Install 

go get github.com/jordic/fuzzyfs
go build github.com/jordic/fuzzyfs/cmd/fu

## Config 

I use the tool in two ways.. as a dir finder, with peco... ( I only want dir results), and as a locate file.. by name.. And i integrate in my bash enviroment, with this scripts:

```bash
cdx() {
    cd $(fu -p="/Users/jordi" -q="$@" -index=".dirs.gob" | peco)
}

fuu() {
    fu -p="/Users/jordi/" -q="$@"
}
```

But first of all, you need to generate the indexes:

For the dir list:
```bash
fu -p="/Users/jordi/" -reindex -depth=8 -method=2 -index=".dirs.gob"
```

For the file index:
```bash
fu -p="/Users/jordi/" -reindex -depth=5 
```

Then, you can use, cdx with peco, oferring, your best directory matches... 

cdx fu .. 




