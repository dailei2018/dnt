package cmd

/*
	已知问题：由于浮点数精度问题，最终出来的肯定会略微不同，但是不影响
*/

import (
	"fmt"
	"os"
	"strings"
	"unsafe"

	"encoding/binary"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func dnt2item(src_path string) [][]Item_t {
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

	bs = bs[4:]

	col_n := binary.LittleEndian.Uint16(bs)
	bs = bs[2:]
	row_n := binary.LittleEndian.Uint32(bs)
	bs = bs[4:]

	//fmt.Println(col_n, row_n)

	patch_arr := make([][]Item_t, row_n+1)
	patch_arr[0] = make([]Item_t, col_n+1)

	patch_arr[0][0].tp = 3
	patch_arr[0][0].val = "pkid"

	//遍历标题
	for i := 1; i < int(col_n+1); i++ {
		len := binary.LittleEndian.Uint16(bs)
		bs = bs[2:]

		patch_arr[0][i].val = string(bs[:len])
		bs = bs[len:]

		patch_arr[0][i].tp = int(bs[0])
		bs = bs[1:]
	}

	//遍历行
	for i := 1; i < int(row_n+1); i++ {
		patch_arr[i] = make([]Item_t, col_n+1)

		patch_arr[i][0].tp = 3
		patch_arr[i][0].val = *(*int32)(unsafe.Pointer(&bs[0]))
		bs = bs[4:]

		for j := 1; j < int(col_n+1); j++ {

			patch_arr[i][j].tp = patch_arr[0][j].tp

			switch patch_arr[0][j].tp {
			case 1:
				//边长字符串
				flen := binary.LittleEndian.Uint16(bs)
				bs = bs[2:]

				patch_arr[i][j].val = string(bs[:flen])
				bs = bs[flen:]

			case 2:
				//bool
				patch_arr[i][j].val = *(*int32)(unsafe.Pointer(&bs[0]))
				bs = bs[4:]
			case 3:
				//int
				patch_arr[i][j].val = *(*int32)(unsafe.Pointer(&bs[0]))
				bs = bs[4:]
			case 4:
				//float 百分比
				patch_arr[i][j].val = *(*float32)(unsafe.Pointer(&bs[0]))
				bs = bs[4:]
			case 5:
				//float
				patch_arr[i][j].val = *(*float32)(unsafe.Pointer(&bs[0]))
				bs = bs[4:]
			default:
				logrus.Fatalf("unkwon type %d", patch_arr[0][j].tp)
			}
		}
	}

	return patch_arr
}

var dnt2dntCmd = &cobra.Command{
	Use:   "dnt2dnt",
	Short: "dnt to dnt patch",

	Run: func(cmd *cobra.Command, args []string) {
		config := MainConfig.Dnt2Dnt
		names := MainConfig.Names

		src_dir := config["src_dir"].(string)
		dst_dir := config["dst_dir"].(string)

		for _, v := range names {
			src_path := src_dir + "/" + v + ".dnt"
			dst_path := dst_dir + "/" + v + ".dnt"

			items := dnt2item(src_path)

			//打补丁
			switch v {
			case "playerweighttable":
				items = playerweighttable(items)
			case "rebootplayerweighttable":
				items = rebootplayerweighttable(items)
			case "enchanttable":
				items = enchanttable(items)
			case "playerleveltable":
				items = playerleveltable(items)
			case "stageentertable":
				items = stageentertable(items)
			case "maptable":
				items = maptable(items)
			case "stagerewardtable":
				items = stagerewardtable(items)
			case "monstertable":
				items = monstertable(items)
			default:
				if strings.Contains(v, "skillleveltable_character") {
					items = skillleveltable_character(items)
				} else if strings.Contains(v, "skilltable_character") {
					items = skilltable_character(items)
				} else {
					fmt.Println("skip----", v)
					continue
				}
			}
			Items2dnt(items, dst_path)

			//fmt.Println("patch write to ", dst_path)

			//dump_res_arr(items)
			//os.Exit(0)
			//fmt.Println(col_n, row_n, patch_arr)
		}

	},
}

func init() {
	rootCmd.AddCommand(dnt2dntCmd)
}
