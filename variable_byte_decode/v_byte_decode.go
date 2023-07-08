package vByteDecode

import (
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"variableByteCode/format_byte"
)

// []byte{0xB8, 0x9E, 0x03}等くるので一つずつ処理
func vByteDecode(data []byte) (uint64, error) {
	var result uint64
	var shift uint
	for _, b := range data {
		// 1バイトずつ読み込んで、下位7ビットをresultに追加
		// 0x7F: 2進数で0111 1111→下位7ビットを取り出すためのマスク
		// b&0x7F: 論理積(0-0 → 0, 0-1 → 0, 1-0 → 0, 1-1 → 1)なので先頭は必ず0になる、下位7ビットが取り出せる
		// shift: 7ビットずつシフトしていく
		// |=: 論理和、+=と意味同じ
		result |= uint64(b&0x7F) << shift
		// 0x80: 2進数で1000 0000→最上位ビットを取り出すためのマスク
		// b&0x80: 論理積(0-0 → 0, 0-1 → 0, 1-0 → 0, 1-1 → 1)なので先頭は必ず0になる、最上位ビットが取り出せる
		// 0x80 == 0: 最上位ビットが0かどうかを判定
		if b&0x80 == 0 {
			return result, nil
		}
		shift += 7
	}
	return 0, fmt.Errorf("invalid vByte encoding")
}

func Check_decodeValue(b []byte) {
	// バイト列をそのまま出力
	fmt.Printf("data: %s\n", formatByte.FormatBytes(b))
	data, err := vByteDecode(b)
	if err != nil {
		fmt.Println("Error decoding entry:", err)
	}
	fmt.Printf("vByteDecode(n): %d\n", data)
}

func DecodeCsv() {
	// 入力ファイルを開く
	inputFile, err := os.Open("./encode.csv")
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer inputFile.Close()

	// CSVリーダーを作成、デリミタとしてタブを指定
	reader := csv.NewReader(inputFile)

	// 出力ファイルを開く
	outputFile, err := os.Create("./decode.csv")
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
		encodeEntries := strings.Split(record[1], ",") // エントリをカンマで分割
		decodedEntries := make([]string, len(encodeEntries))
		prevEntry := uint64(0)
		for i, encodeEntry := range encodeEntries {
		  // encodeされた16進数エントリをバイトスライスに変換
			encoded, err := hex.DecodeString(encodeEntry)
			if err != nil {
				fmt.Println("Error decoding entry:", err)
				return
			}

			// バイトスライスをデコードして元の値を取得
			gap, err := vByteDecode(encoded)
			if err != nil {
				fmt.Println("Error decoding entry:", err)
				return
			}

			// ギャップ値を元のエントリに変換
			entry := prevEntry + gap
			prevEntry = entry

			// デコードされたエントリを文字列に変換
			decodedEntries[i] = strconv.FormatUint(entry, 10)
		}

		// エンコードされたエントリを書き込む
		err := writer.Write([]string{string(record[0]), strings.Join(decodedEntries, ",")})
		if err != nil {
			fmt.Println("Error writing to output file:", err)
			return
		}
	}
}
