# baseemoji - encode to/from emojis

## Building

```
$ go build baseemoji.go
```

## Examples

### Encoding

```
$ echo "Hello, World." | ./baseemoji -e
😲😻🐱🐱🐰😭😞🤠🐰🐼🐱😹😠🙂
```

### Decoding

```
$ echo -ne "😲😻🐱🐱🐰😭😞🤠🐰🐼🐱😹😠🙂" | ./baseemoji -d
Hello, World.
```
