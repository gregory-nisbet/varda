# varda
varda looks for stuff, like grep.

```
let &grepprg = "varda"
```

```
varda pattern file-or-directory file-or-directory file-or-directory
```

The eventual goal of this project is to implement good ideas from ack, ripgrep,
the silver searcher, git grep, and other programmer tools.

I'm also a big fan of the emacs plugin dumb-jump and plan to eventually support
a flag or something for limiting results to things that look like definitions.

I use grep constantly, mostly from within vim. I've had a few ideas over the
years for things that could potentially improve it as a program to be invoked
from an editor, like finding the repo root, automatically including library
locations for the language that you're using, listing results in breadth-first
order, limiting the number of files scanned to preserve responsiveness, and
other things.

For now, the tool does extremely little. It just searches and defaults to
searching the current working directory. It's basically just a non-performant
grep that is always recursive. As it turns out, that alone is enough to make it
usable as grepprg in vim with decent ergonomics.
