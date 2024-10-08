package main

import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "os"
)

var emojiAlphabet = []rune{
    '😀', '😃', '😄', '😁', '😆', '😅', '😂', '🤣',
    '😊', '😇', '🙂', '🙃', '😉', '😌', '😍', '🥰',
    '😘', '😗', '😙', '😚', '😋', '😛', '😝', '😜',
    '🤪', '🤨', '🧐', '🤓', '😎', '🥳', '😏', '😒',
    '😞', '😔', '😟', '😕', '🙁', '☹',  '😣', '😖',
    '😫', '😩', '🥺', '😢', '😭', '😤', '😠', '😡',
    '🤬', '🤯', '😳', '🥵', '🥶', '😱', '😨', '😰',
    '😥', '😓', '🤗', '🤔', '🤭', '🤫', '🤥', '😶',
}

var emojiToValue = make(map[rune]int)

func init() {
    for i, emoji := range emojiAlphabet {
        emojiToValue[emoji] = i
    }
}

// EncodeToEmojiStream reads from reader and writes encoded emojis to writer
func EncodeToEmojiStream(reader io.Reader, writer io.Writer) error {
    bufReader := bufio.NewReader(reader)
    bufWriter := bufio.NewWriter(writer)
    defer bufWriter.Flush()

    var buffer uint64 = 0
    var bitsInBuffer int = 0

    const chunkSize = 4096
    data := make([]byte, chunkSize)

    for {
        n, err := bufReader.Read(data)
        if n > 0 {
            for i := 0; i < n; i++ {
                buffer = (buffer << 8) | uint64(data[i])
                bitsInBuffer += 8

                for bitsInBuffer >= 6 {
                    bitsInBuffer -= 6
                    index := (buffer >> bitsInBuffer) & 0x3F
                    emoji := emojiAlphabet[index]
                    _, err := bufWriter.WriteRune(emoji)
                    if err != nil {
                        return err
                    }
                }
                buffer &= (1 << bitsInBuffer) - 1 // Keep only the remaining bits
            }
        }
        if err != nil {
            if err == io.EOF {
                break
            }
            return err
        }
    }

    // Handle remaining bits
    if bitsInBuffer > 0 {
        buffer <<= (6 - bitsInBuffer)
        index := buffer & 0x3F
        emoji := emojiAlphabet[index]
        _, err := bufWriter.WriteRune(emoji)
        if err != nil {
            return err
        }
    }

    return nil
}

// DecodeFromEmojiStream reads encoded emojis from reader and writes decoded data to writer
func DecodeFromEmojiStream(reader io.Reader, writer io.Writer) error {
    bufReader := bufio.NewReader(reader)
    bufWriter := bufio.NewWriter(writer)
    defer bufWriter.Flush()

    var buffer uint64 = 0
    var bitsInBuffer int = 0

    for {
        char, _, err := bufReader.ReadRune()
        if err != nil {
            if err == io.EOF {
                break
            }
            return err
        }

        value, exists := emojiToValue[char]
        if !exists {
            return fmt.Errorf("invalid character: %c", char)
        }

        buffer = (buffer << 6) | uint64(value)
        bitsInBuffer += 6

        for bitsInBuffer >= 8 {
            bitsInBuffer -= 8
            byteValue := byte((buffer >> bitsInBuffer) & 0xFF)
            err := bufWriter.WriteByte(byteValue)
            if err != nil {
                return err
            }
            buffer &= (1 << bitsInBuffer) - 1 // Keep only the remaining bits
        }
    }

    return nil
}

func main() {
    // Define flags for encoding and decoding
    encodeFlag := flag.Bool("e", false, "Encode input from stdin to emojis")
    decodeFlag := flag.Bool("d", false, "Decode emoji input from stdin")
    flag.Parse()

    // Ensure that exactly one of -e or -d is provided
    if (*encodeFlag && *decodeFlag) || (!*encodeFlag && !*decodeFlag) {
        fmt.Fprintln(os.Stderr, "Usage: program -e | -d")
        fmt.Fprintln(os.Stderr, "  -e    Encode input from stdin to emojis")
        fmt.Fprintln(os.Stderr, "  -d    Decode emoji input from stdin")
        os.Exit(1)
    }

    if *encodeFlag {
        // Encode from stdin to stdout
        err := EncodeToEmojiStream(os.Stdin, os.Stdout)
        if err != nil {
            fmt.Fprintln(os.Stderr, "Encoding Error:", err)
            os.Exit(1)
        }
    } else if *decodeFlag {
        // Decode from stdin to stdout
        err := DecodeFromEmojiStream(os.Stdin, os.Stdout)
        if err != nil {
            fmt.Fprintln(os.Stderr, "Decoding Error:", err)
            os.Exit(1)
        }
    }
}
