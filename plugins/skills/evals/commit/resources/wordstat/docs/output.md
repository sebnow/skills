# Output

wordstat prints one line per file, in the order the paths were given:

```
      1       3      14 notes.txt
```

The columns are right-aligned to a width of seven:

1. Lines: the number of newline characters. A last line without a
   trailing newline is not counted.
2. Words: runs of characters between separators.
3. Characters: Unicode code points, not bytes.
4. The path as given on the command line.
