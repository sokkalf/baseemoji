package main

import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "os"
)

var bitsPerSymbol = 8 // Change this to 7 for 128 emojis

// Initialize the emoji alphabet with 256 unique emojis
var emojiAlphabet = []string{
    "😀", "😃", "😄", "😁", "😆", "😅", "😂", "🤣",
    "😊", "😇", "🙂", "🙃", "😉", "😌", "😍", "🥰",
    "😘", "😗", "😙", "😚", "😋", "😛", "😝", "😜",
    "🤪", "🤨", "🧐", "🤓", "😎", "🥳", "😏", "😒",
    "😞", "😔", "😟", "😕", "🙁", "☹️", "😣", "😖",
    "😫", "😩", "🥺", "😢", "😭", "😤", "😠", "😡",
    "🤬", "🤯", "😳", "🥵", "🥶", "😱", "😨", "😰",
    "😥", "😓", "🤗", "🤔", "🤭", "🤫", "🤥", "😶",
    "😐", "😑", "😬", "🙄", "😯", "😦", "😧", "😮",
    "😲", "🥱", "😴", "🤤", "😪", "😵", "🤐", "🥴",
    "🤢", "🤮", "🤧", "😷", "🤒", "🤕", "🤑", "🤠",
    "😈", "👿", "🤡", "💩", "👹", "👺", "👻", "👽",
    "👾", "🤖", "😺", "😸", "😹", "😻", "😼", "😽",
    "🙀", "😿", "😾", "🐶", "🐱", "🐭", "🐹", "🐰",
    "🦊", "🐻", "🐼", "🐨", "🐯", "🦁", "🐮", "🐷",
    "🐽", "🐸", "🐵", "🙈", "🙉", "🙊", "🐒", "🐔",
    "🐧", "🐦", "🐤", "🐣", "🐥", "🦆", "🦅", "🦉",
    "🦇", "🐺", "🐗", "🐴", "🦄", "🐝", "🐛", "🦋",
    "🐌", "🐞", "🐜", "🦟", "🦗", "🐢", "🐍", "🦎",
    "🐙", "🦑", "🦐", "🦀", "🐡", "🐠", "🐟", "🐬",
    "🐳", "🐋", "🦈", "🐊", "🐅", "🐆", "🦓", "🦍",
    "🐘", "🦏", "🦛", "🐪", "🐫", "🦒", "🐃", "🐂",
    "🐄", "🐎", "🐖", "🐏", "🐑", "🐐", "🦌", "🐕",
    "🐩", "🐈", "🐓", "🦃", "🕊️", "🐇", "🐁", "🐀",
    "🐿️", "🦔", "🐾", "🐉", "🐲", "🌵", "🎄", "🌲",
    "🌳", "🌴", "🌱", "🌿", "☘️", "🍀", "🎍", "🎋",
    "🍃", "🍂", "🍁", "🍄", "🌾", "💐", "🌷", "🌹",
    "🥀", "🌺", "🌸", "🌼", "🌻", "🌞", "🌝", "🌛",
    "🌚", "🌕", "🌖", "🌗", "🌘", "🌑", "🌒", "🌓",
    "🌔", "🌙", "🌎", "🌍", "🌏", "💫", "⭐", "🌟",
    "✨", "⚡", "🔥", "💥", "🌈", "☄️", "🌪️", "🌨️",
    "🌧️", "🌦️", "☀️", "🌤️", "⛅", "🌥️", "☁️", "🌩️",
    "🌊", "💧", "💦", "☔", "❄️", "🌬️", "🖕", "🌁",
}

var emojiToValue = make(map[string]int)

func init() {
    for i, emoji := range emojiAlphabet {
        emojiToValue[emoji] = i
    }
}

// EncodeToEmojiStream encodes data from reader to emojis and writes to writer
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

                for bitsInBuffer >= bitsPerSymbol {
                    bitsInBuffer -= bitsPerSymbol
                    index := (buffer >> bitsInBuffer) & ((1 << bitsPerSymbol) - 1)
                    emoji := emojiAlphabet[index]
                    _, err := bufWriter.WriteString(emoji)
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
        buffer <<= (bitsPerSymbol - bitsInBuffer)
        index := buffer & ((1 << bitsPerSymbol) - 1)
        emoji := emojiAlphabet[index]
        _, err := bufWriter.WriteString(emoji)
        if err != nil {
            return err
        }
    }

    return nil
}

// DecodeFromEmojiStream decodes emojis from reader and writes original data to writer
func DecodeFromEmojiStream(reader io.Reader, writer io.Writer) error {
    bufReader := bufio.NewReader(reader)
    bufWriter := bufio.NewWriter(writer)
    defer bufWriter.Flush()

    var buffer uint64 = 0
    var bitsInBuffer int = 0

    for {
        // Read the next emoji
        emoji, err := bufReader.ReadString('\n') // Use a delimiter that won't appear in emojis
        if err != nil && err != io.EOF {
            return err
        }

        if len(emoji) == 0 && err == io.EOF {
            break
        }

        // Trim any whitespace or newlines
        emoji = emoji[:len(emoji)-1]

        value, exists := emojiToValue[emoji]
        if !exists {
            // Try to read more if the emoji wasn't complete
            for !exists {
                nextRune, _, err := bufReader.ReadRune()
                if err != nil {
                    return fmt.Errorf("invalid character or incomplete emoji")
                }
                emoji += string(nextRune)
                value, exists = emojiToValue[emoji]
            }
        }

        buffer = (buffer << bitsPerSymbol) | uint64(value)
        bitsInBuffer += bitsPerSymbol

        for bitsInBuffer >= 8 {
            bitsInBuffer -= 8
            byteValue := byte((buffer >> bitsInBuffer) & 0xFF)
            err := bufWriter.WriteByte(byteValue)
            if err != nil {
                return err
            }
            buffer &= (1 << bitsInBuffer) - 1 // Keep only the remaining bits
        }

        if err == io.EOF {
            break
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
