package vByteEncode

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"variableByteCode/format_byte"
)

//Variable Byte Encoding（VByteエンコーディング）は、非負整数をバイト配列として効率的に格納するための方法。
//数値を7ビットのグループに分割し、最後のバイト以外のすべてのバイトに「続くバイトがある」ことを示すために最上位ビットを設定。
func vByteEncode(n uint64) []byte {
	var b bytes.Buffer
	for n >= 0x80 {  // nが128以上である限りループを続ける(0x80: 16進数で128)
		// nの下位7ビットをバッファに書き込む
		b.WriteByte(byte(n) | 0x80) // byte()は下位8ビットを取り出す, 255までしか表現できない、かつ最上位ビットを1にする
		n >>= 7  // nを7ビット右シフトして、次の7ビットを処理する準備をする(左側に0が詰められる)
	}
	b.WriteByte(byte(n))  // 最後の7ビット（もしくはそれ以下）をバッファに書き込む
	return b.Bytes()
}

func Check_encodeValue(n uint64) {
	fmt.Printf("n: %d\n", n)
	// バイト列をそのまま出力
	fmt.Printf("data: %s\n", formatByte.FormatBytes(vByteEncode(n)))
	fmt.Printf("vByteEncode(n): %x\n", vByteEncode(n))
}

// "tag" "整数,整数,整数,整数" という形式のCSVファイルを読み込み、vByteで圧縮したものを出力する
func EncodeCsv() {
	// 入力ファイルを開く
	inputFile, err := os.Open("./test.csv")
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer inputFile.Close()

	// CSVリーダーを作成、デリミタとしてタブを指定
	reader := csv.NewReader(inputFile)
	reader.Comma = '\t' // デフォルトはカンマ区切りなので、タブで区切る

	// 出力ファイルを開く
	outputFile, err := os.Create("./encode.csv")
	if err != nil {
		fmt.Println("Error opening output file:", err)
		return
	}
	defer outputFile.Close()

	// CSVライターを作成
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// 入力ファイルを一行ずつ読み込む
	// ex)["Windows", "1409171,14711245,18265928,21590872"] タグと整数列のスライス
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	for _, record := range records {
		entries := strings.Split(record[1], ",") // エントリをカンマで分割
		entryNums := make([]uint64, len(entries))
		// エントリをuint64に変換し、リストをソート
		for i, entry := range entries {
			entryNums[i], err = strconv.ParseUint(entry, 10, 64)
			if err != nil {
				fmt.Println("Error converting entry to uint64:", err)
				return
			}
		}
		// エントリをソート
		sort.Slice(entryNums, func(i, j int) bool { return entryNums[i] < entryNums[j] })

		// ギャップエンコーディングを適用
		var encodedEntries []string
		var prevEntry uint64 = 0
		for _, entryNum := range entryNums {
			// エントリをVByteエンコード
			gap := entryNum - prevEntry
			prevEntry = entryNum
			encoded := vByteEncode(gap)
			// エンコードされたエントリを16進数文字列に変換(10進数だと1バイトにつき最大3文字必要)
			encodedEntries = append(encodedEntries, fmt.Sprintf("%x", encoded))
		}

		// エンコードされたエントリを書き込む
		err := writer.Write([]string{ record[0], strings.Join(encodedEntries, ",")})
		if err != nil {
			fmt.Println("Error writing to output file:", err)
			return
		}
	}
}
