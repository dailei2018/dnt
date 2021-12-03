package cmd

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/dailei2018/dnt/lib"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
)

type Item_t struct {
	tp  int
	val interface{}
}

func dump_res_arr(res_arr [][]Item_t) {
	for _, v := range res_arr {
		for _, v1 := range v {
			fmt.Printf("%v ", v1.val)
		}
		fmt.Printf("\n")
	}
}

func Items2dnt(res_arr [][]Item_t, dst_path string) {
	i := 0
	tmp_bs := make([]byte, 8)
	bs := make([]byte, 0, 4)
	bs = append(bs, "\x00\x00\x00\x00"...)
	i += 4

	col_n := uint16(len(res_arr[0]))
	row_n := uint32(len(res_arr))

	binary.LittleEndian.PutUint16(tmp_bs, col_n-1) //不算pkid列
	bs = append(bs, tmp_bs[0:2]...)
	binary.LittleEndian.PutUint32(tmp_bs, row_n-1) //不算标题行
	bs = append(bs, tmp_bs[0:4]...)

	//写入字段名，跳过pkid
	for i := 1; i < len(res_arr[0]); i++ {
		binary.LittleEndian.PutUint16(tmp_bs, uint16(len(res_arr[0][i].val.(string))))
		bs = append(bs, tmp_bs[0:2]...)                //标题长度
		bs = append(bs, res_arr[0][i].val.(string)...) //标题值
		bs = append(bs, byte(res_arr[0][i].tp))        //标题类型

		//fmt.Println(res_arr[0][i].val.(string))
	}

	for i := 1; i < len(res_arr); i++ {
		for j := 0; j < len(res_arr[i]); j++ {
			item := res_arr[i][j]

			switch item.tp {
			case 1:
				binary.LittleEndian.PutUint16(tmp_bs, uint16(len(item.val.(string))))
				bs = append(bs, tmp_bs[0:2]...)
				bs = append(bs, item.val.(string)...)
			case 2:
				binary.LittleEndian.PutUint32(tmp_bs, uint32(item.val.(int32)))
				bs = append(bs, tmp_bs[0:4]...)
			case 3:
				binary.LittleEndian.PutUint32(tmp_bs, uint32(item.val.(int32)))
				bs = append(bs, tmp_bs[0:4]...)
			case 4:
				binary.LittleEndian.PutUint32(tmp_bs, math.Float32bits(item.val.(float32)))
				bs = append(bs, tmp_bs[0:4]...)
			case 5:
				binary.LittleEndian.PutUint32(tmp_bs, math.Float32bits(item.val.(float32)))
				bs = append(bs, tmp_bs[0:4]...)
			default:
				logrus.Fatalln(item.tp)
			}
		}
	}

	bs = append(bs, "\x05THEND"...)

	fp, err := os.OpenFile(dst_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		logrus.Fatalf("open %s failed: %v", dst_path, err)
	}
	n, err := fp.Write(bs)
	if err != nil {
		logrus.Fatalf("write %s failed: %v n:%d", dst_path, err, n)
	}
	fp.Close()

	fmt.Printf("write %d bytes to %s\n", n, dst_path)
}

var excel2dntCmd = &cobra.Command{
	Use:   "excel2dnt",
	Short: "excel to dnt",

	Run: func(cmd *cobra.Command, args []string) {
		var res_arr [][]Item_t

		config := MainConfig.Excel2dnt
		names := MainConfig.Names

		src_dir := config["src_dir"].(string)
		dst_dir := config["dst_dir"].(string)

		lib.Mk_dir_if_not(src_dir)
		lib.Mk_dir_if_not(dst_dir)

		for _, v := range names {
			src_path := src_dir + "/" + v + ".xlsx"
			dst_path := dst_dir + "/" + v + ".dnt"

			f, err := excelize.OpenFile(src_path)
			if err != nil {
				logrus.Fatalln(err)
			}

			rows, err := f.GetRows("Sheet1")
			if err != nil {
				logrus.Fatalln(err)
				return
			}

			res_arr = make([][]Item_t, len(rows))

			row1 := rows[0]
			rows = rows[1:]

			res_arr[0] = make([]Item_t, len(row1))
			res_arr[0][0].tp = 3
			res_arr[0][0].val = "pkid"

			for i := 1; i < len(row1); i++ {
				arr := strings.Split(row1[i], "((")
				res_arr[0][i].tp, _ = strconv.Atoi(arr[1])
				res_arr[0][i].val = arr[0]
			}

			for i, row := range rows {
				res_arr[i+1] = make([]Item_t, len(row1))
				for j, col_v := range row {
					res_arr[i+1][j].tp = res_arr[0][j].tp
					switch res_arr[i+1][j].tp {
					case 1:
						res_arr[i+1][j].val = col_v
					case 2:
						res_arr[i+1][j].val, _ = strconv.Atoi(col_v)
					case 3:
						res_arr[i+1][j].val, _ = strconv.Atoi(col_v)
					case 4:
						v64, _ := strconv.ParseFloat(col_v, 32)
						res_arr[i+1][j].val = float32(v64)
					case 5:
						v64, _ := strconv.ParseFloat(col_v, 32)
						res_arr[i+1][j].val = float32(v64)
					default:
						logrus.Fatalln(res_arr[i+1][j].tp)
					}
				}
			}

			Items2dnt(res_arr, dst_path)
		}

	},
}

func init() {
	rootCmd.AddCommand(excel2dntCmd)
}
