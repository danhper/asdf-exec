# asdf-exec

[![CircleCI](https://circleci.com/gh/danhper/asdf-exec.svg?style=svg)](https://circleci.com/gh/danhper/asdf-exec)

Experimental command to find and run executable used by [asdf][asdf] shims.

NOTE: This is not officially supported by asdf and requires to slightly
patch the asdf code

## Installation

I assume `asdf` is installed in `~/.asdf`. If not, please change the commands accordingly.

First, grab the `asdf-exec` executable from [the releases](https://github.com/danhper/asdf-exec/releases/)
and put it in `~/.asdf/bin/private` as `asdf-exec`

```
# for linux
wget https://github.com/danhper/asdf-exec/releases/download/v0.1.2/asdf-exec-linux-x64 -O ~/.asdf/bin/private/asdf-exec
# for macos
wget https://github.com/danhper/asdf-exec/releases/download/v0.1.2/asdf-exec-darwin-x64 -O ~/.asdf/bin/private/asdf-exec

# for both:
chmod +x ~/.asdf/bin/private/asdf-exec
```

Then, patch asdf reshim command code and regenerate all shims.

```
sed -i.bak -e 's|exec $(asdf_dir)/bin/asdf exec|exec $(asdf_dir)/bin/private/asdf-exec|' ~/.asdf/lib/commands/reshim.sh
rm ~/.asdf/shims/*
asdf reshim
```

## Rationale

As asdf is growing in features and complexity, the logic to run a single
command has become fairly involved and unfortunately quite slow.
For example

```
$ time python --version
Python 3.7.2
0.19user 0.04system 0:00.16elapsed 142%CPU (0avgtext+0avgdata 4344maxresident)k
0inputs+0outputs (0major+25175minor)pagefaults 0swaps
```

While there are surely many things we could do better with bash to improve
performance, tuning bash is quite tedious and error-prone.

I decided to give a native command a try. This would be called from shims
and would take care of locating the correct version.
Although this does results in quite a bit of duplication between the bash
code and this one, it does give a speed improvement which might be worth it.

```
$ time python --version
Python 3.7.2
0.00user 0.00system 0:00.01elapsed 100%CPU (0avgtext+0avgdata 4096maxresident)k
0inputs+0outputs (0major+691minor)pagefaults 0swaps
```

This is still totally experimental and there are no plan to merge or use this into
asdf for now, but please feel free to give it a try and let me know what you think.




[asdf]: https://github.com/asdf-vm/asdf
