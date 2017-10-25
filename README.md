sbinet.github.io
================

My personal ramblings.
On the web!

[Here: sbinet.github.io](https://sbinet.github.io)

## Adding new content

To add new content to the `blog` section:

```
git clone https://github.com/sbinet/blog
cd blog
hugo new "posts/a-new-entry.md"
$EDITOR ./src/posts/a-new-entry.md
```

## Testing your changes

```
cd blog
make serve
open http://localhost:8080
```

## Pushing your changes

```
git checkout -b my-branch origin/master
git add -p
git commit -m "content: bla bla"
git push my-branch
```

