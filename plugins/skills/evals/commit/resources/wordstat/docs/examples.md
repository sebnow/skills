# Examples

Count one file:

```sh
$ wordstat notes.txt
     12      80     431 notes.txt
```

Count several files at once; each gets its own line:

```sh
$ wordstat notes.txt todo.txt
     12      80     431 notes.txt
      3       9      52 todo.txt
```

Sort a directory's text files by word count with standard tools:

```sh
wordstat *.txt | sort -k2 -n
```
