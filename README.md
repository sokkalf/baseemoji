# baseemoji - encode to/from emojis

## Examples

### Encoding

```
~ $ echo "Hello, World." | ./baseemoji -e
😲😻🐱🐱🐰😭😞🤠🐰🐼🐱😹😠🙂
```

### Decoding

```
~ $ echo -ne "😲😻🐱🐱🐰😭😞🤠🐰🐼🐱😹😠🙂" | ./baseemoji -d
Hello, World.
```
