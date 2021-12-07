package cmd

/*
	已知问题：由于浮点数精度问题，最终出来的肯定会略微不同，但是不影响
*/

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type File_t struct {
	name     string
	size     uint32
	size1    uint32
	content  []byte
	content1 []byte

	content_offset uint32
}

var file_arr []File_t
var out []byte

func file2bytes(src_path string) []byte {
	fp, err := os.Open(src_path)
	if err != nil {
		logrus.Fatalf("open %s failed: %v", src_path, err)
	}

	fi, _ := fp.Stat()
	//fmt.Println(fi.Size())
	bs := make([]byte, fi.Size())
	n, _ := fp.Read(bs)
	if n != int(fi.Size()) {
		logrus.Fatalf("size not same")
	}
	fp.Close()

	return bs
}

var dnt2pakCmd = &cobra.Command{
	Use:   "dnt2pak",
	Short: "dnt to pak patch",

	Run: func(cmd *cobra.Command, args []string) {
		config := MainConfig.Dnt2Pak

		src_dir := config["src_dir"].(string)
		dst_path := config["dst_path"].(string)
		pad := config["pad"].(bool)

		os.MkdirAll(src_dir, os.FileMode(0766))

		dir1 := strings.Split(src_dir, "/")[0]
		os.Chdir(dir1)

		filepath.Walk("./", func(path string, info os.FileInfo, err error) error {

			if strings.Contains(path, ".dnt") {
				path1 := "\\" + path

				fmt.Printf("compress %s\n", path1)

				file_item := &File_t{}
				file_item.name = path1
				file_item.size = uint32(info.Size())
				file_item.content = file2bytes(path)

				var in bytes.Buffer
				w, _ := zlib.NewWriterLevel(&in, zlib.BestSpeed)
				w.Write(file_item.content)
				w.Close()
				file_item.content1 = in.Bytes()
				file_item.size1 = uint32(in.Len())

				file_arr = append(file_arr, *file_item)
			}

			return nil
		})

		tmp_bs := make([]byte, 8)
		var file_offset_index uint32

		out = append(out, "EyedentityGames Packing File 0.1"...)
		out = append(out, bytes.Repeat([]byte{0}, 256-len(out))...)

		out = append(out, "\x0b\x00\x00\x00"...)

		//文件个数
		binary.LittleEndian.PutUint32(tmp_bs, uint32(len(file_arr)))
		out = append(out, tmp_bs[:4]...)

		//文件信息偏移，这里要到后面才能知晓，暂时填充0
		file_offset_index = uint32(len(out))
		out = append(out, "\x00\x00\x00\x00"...)

		//填充0到1024个字节
		out = append(out, bytes.Repeat([]byte{0}, 1024-len(out))...)

		//遍历写入压缩后的文件内容
		for i, v := range file_arr {
			file_arr[i].content_offset = uint32(len(out))
			out = append(out, v.content1...)
		}

		//填充
		if pad {
			sz := 0
			sz += 1024
			for _, v := range file_arr {
				sz += int(v.size1) + 256 + 16 + 44
			}
			pad_sz := 510*1024*1024 - sz
			out = append(out, bytes.Repeat([]byte{0}, pad_sz)...)
		}

		//填充文件信息偏移
		fo := (*uint32)(unsafe.Pointer(&out[file_offset_index]))
		*fo = uint32(len(out))

		//遍历写入文件信息
		for _, v := range file_arr {
			out = append(out, v.name...)
			out = append(out, bytes.Repeat([]byte{0}, 256-len(v.name))...) //文件名占用256字节

			binary.LittleEndian.PutUint32(tmp_bs, v.size1)
			out = append(out, tmp_bs[:4]...)

			binary.LittleEndian.PutUint32(tmp_bs, v.size)
			out = append(out, tmp_bs[:4]...)

			binary.LittleEndian.PutUint32(tmp_bs, v.size1)
			out = append(out, tmp_bs[:4]...)

			binary.LittleEndian.PutUint32(tmp_bs, v.content_offset)
			out = append(out, tmp_bs[:4]...)

			out = append(out, bytes.Repeat([]byte{0}, 44)...)
		}

		fmt.Printf("\nbegin wirting......")
		fp, err := os.OpenFile(dst_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))
		if err != nil {
			logrus.Fatalf("open %s failed: %v", dst_path, err)
		}
		n, err := fp.Write(out)
		if err != nil {
			logrus.Fatalf("write %s failed: %v n:%d", dst_path, err, n)
		}
		fp.Close()

		fmt.Printf("done\n")

		fmt.Printf("write %d bytes to %s\n", n, dst_path)

	},
}

func init() {
	rootCmd.AddCommand(dnt2pakCmd)
}
