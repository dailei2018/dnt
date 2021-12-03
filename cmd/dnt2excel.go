package cmd

/*
	已知问题：由于浮点数精度问题，最终出来的肯定会略微不同，但是不影响
*/

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"

	"encoding/binary"

	"github.com/dailei2018/dnt/lib"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
)

type field_t struct {
	flen uint16
	name string
	tp   uint8
}

var wh_arr = [...]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ", "BA", "BB", "BC", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BK", "BL", "BM", "BN", "BO", "BP", "BQ", "BR", "BS", "BT", "BU", "BV", "BW", "BX", "BY", "BZ", "CA", "CB", "CC", "CD", "CE", "CF", "CG", "CH", "CI", "CJ", "CK", "CL", "CM", "CN", "CO", "CP", "CQ", "CR", "CS", "CT", "CU", "CV", "CW", "CX", "CY", "CZ", "DA", "DB", "DC", "DD", "DE", "DF", "DG", "DH", "DI", "DJ", "DK", "DL", "DM", "DN", "DO", "DP", "DQ", "DR", "DS", "DT", "DU", "DV", "DW", "DX", "DY", "DZ", "EA", "EB", "EC", "ED", "EE", "EF", "EG", "EH", "EI", "EJ", "EK", "EL", "EM", "EN", "EO", "EP", "EQ", "ER", "ES", "ET", "EU", "EV", "EW", "EX", "EY", "EZ", "FA", "FB", "FC", "FD", "FE", "FF", "FG", "FH", "FI", "FJ", "FK", "FL", "FM", "FN", "FO", "FP", "FQ", "FR", "FS", "FT", "FU", "FV", "FW", "FX", "FY", "FZ", "GA", "GB", "GC", "GD", "GE", "GF", "GG", "GH", "GI", "GJ", "GK", "GL", "GM", "GN", "GO", "GP", "GQ", "GR", "GS", "GT", "GU", "GV", "GW", "GX", "GY", "GZ", "HA", "HB", "HC", "HD", "HE", "HF", "HG", "HH", "HI", "HJ", "HK", "HL", "HM", "HN", "HO", "HP", "HQ", "HR", "HS", "HT", "HU", "HV", "HW", "HX", "HY", "HZ", "IA", "IB", "IC", "ID", "IE", "IF", "IG", "IH", "II", "IJ", "IK", "IL", "IM", "IN", "IO", "IP", "IQ", "IR", "IS", "IT", "IU", "IV", "IW", "IX", "IY", "IZ", "JA", "JB", "JC", "JD", "JE", "JF", "JG", "JH", "JI", "JJ", "JK", "JL", "JM", "JN", "JO", "JP", "JQ", "JR", "JS", "JT", "JU", "JV", "JW", "JX", "JY", "JZ", "KA", "KB", "KC", "KD", "KE", "KF", "KG", "KH", "KI", "KJ", "KK", "KL", "KM", "KN", "KO", "KP", "KQ", "KR", "KS", "KT", "KU", "KV", "KW", "KX", "KY", "KZ", "LA", "LB", "LC", "LD", "LE", "LF", "LG", "LH", "LI", "LJ", "LK", "LL", "LM", "LN", "LO", "LP", "LQ", "LR", "LS", "LT", "LU", "LV", "LW", "LX", "LY", "LZ", "MA", "MB", "MC", "MD", "ME", "MF", "MG", "MH", "MI", "MJ", "MK", "ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS", "MT", "MU", "MV", "MW", "MX", "MY", "MZ", "NA", "NB", "NC", "ND", "NE", "NF", "NG", "NH", "NI", "NJ", "NK", "NL", "NM", "NN", "NO", "NP", "NQ", "NR", "NS", "NT", "NU", "NV", "NW", "NX", "NY", "NZ", "OA", "OB", "OC", "OD", "OE", "OF", "OG", "OH", "OI", "OJ", "OK", "OL", "OM", "ON", "OO", "OP", "OQ", "OR", "OS", "OT", "OU", "OV", "OW", "OX", "OY", "OZ", "PA", "PB", "PC", "PD", "PE", "PF", "PG", "PH", "PI", "PJ", "PK", "PL", "PM", "PN", "PO", "PP", "PQ", "PR", "PS", "PT", "PU", "PV", "PW", "PX", "PY", "PZ", "QA", "QB", "QC", "QD", "QE", "QF", "QG", "QH", "QI", "QJ", "QK", "QL", "QM", "QN", "QO", "QP", "QQ", "QR", "QS", "QT", "QU", "QV", "QW", "QX", "QY", "QZ", "RA", "RB", "RC", "RD", "RE", "RF", "RG", "RH", "RI", "RJ", "RK", "RL", "RM", "RN", "RO", "RP", "RQ", "RR", "RS", "RT", "RU", "RV", "RW", "RX", "RY", "RZ"}

var dnt2excelCmd = &cobra.Command{
	Use:   "dnt2excel",
	Short: "dnt to excel",

	Run: func(cmd *cobra.Command, args []string) {
		config := MainConfig.Dnt2excel
		names := MainConfig.Names

		src_dir := config["src_dir"].(string)
		dst_dir := config["dst_dir"].(string)

		lib.Mk_dir_if_not(src_dir)
		lib.Mk_dir_if_not(dst_dir)

		//fmt.Println(len(wh_arr))

		for _, v := range names {
			src_path := src_dir + "/" + v + ".dnt"

			var dst_path string
			if strings.Contains(dst_dir, "tmp") {
				dst_path = dst_dir + "/" + v + "1.xlsx"
			} else {
				dst_path = dst_dir + "/" + v + ".xlsx"
			}

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

			field_arr := make([]field_t, col_n)

			f := excelize.NewFile()
			f.NewSheet("Sheet1")

			f.SetCellValue("Sheet1", wh_arr[0]+strconv.Itoa(1), "pkid")
			f.SetColWidth("Sheet1", wh_arr[0], wh_arr[0], 4)

			//遍历标题
			for i := 0; i < int(col_n); i++ {
				field_arr[i].flen = binary.LittleEndian.Uint16(bs)
				bs = bs[2:]
				field_arr[i].name = string(bs[:field_arr[i].flen])
				bs = bs[field_arr[i].flen:]
				field_arr[i].tp = bs[0]
				bs = bs[1:]

				wh := wh_arr[i+1] + strconv.Itoa(1)
				val := field_arr[i].name + "((" + strconv.Itoa(int(field_arr[i].tp))

				//fmt.Println(wh, val)

				f.SetCellValue("Sheet1", wh, val)

				f.SetColWidth("Sheet1", wh_arr[i+1], wh_arr[i+1], float64(field_arr[i].flen)+5)
			}

			//遍历行
			for i := 0; i < int(row_n); i++ {
				pkid := binary.LittleEndian.Uint32(bs)
				bs = bs[4:]
				f.SetCellValue("Sheet1", wh_arr[0]+strconv.Itoa(i+2), pkid)

				for j := 0; j < int(col_n); j++ {
					field := field_arr[j]

					wh := wh_arr[j+1] + strconv.Itoa(i+2)

					switch field.tp {
					case 1:
						//边长字符串
						flen := binary.LittleEndian.Uint16(bs)
						bs = bs[2:]
						str := string(bs[:flen])
						bs = bs[flen:]
						f.SetCellStr("Sheet1", wh, str)
					case 2:
						//bool
						val := binary.LittleEndian.Uint32(bs)
						bs = bs[4:]
						f.SetCellInt("Sheet1", wh, int(val))
					case 3:
						//int
						val := *(*int32)(unsafe.Pointer(&bs[0]))
						bs = bs[4:]
						f.SetCellInt("Sheet1", wh, int(val))
					case 4:
						//float 百分比
						val := *(*float32)(unsafe.Pointer(&bs[0]))
						bs = bs[4:]
						f.SetCellFloat("Sheet1", wh, float64(val), -1, 32)
					case 5:
						//float
						val := *(*float32)(unsafe.Pointer(&bs[0]))
						bs = bs[4:]

						f.SetCellFloat("Sheet1", wh, float64(val), -1, 32)
					default:
						logrus.Fatalf("unkwon type %d", field.tp)
					}
				}
			}

			f.SaveAs(dst_path)
			fmt.Printf("write to %s\n", dst_path)
			//fmt.Println(col_n, row_n, field_arr)
		}

	},
}

func init() {
	rootCmd.AddCommand(dnt2excelCmd)
}
