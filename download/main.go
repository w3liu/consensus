package main

import "fmt"

type Course struct {
	url  string
	name string
}

func main() {
	courses := []Course{
		{url: "https://bf.topath.net.cn/zy/202002/5d14aa5a7c84f.mp4", name: "1.1.名师讲微课：认识亿以内数的计数单位以及数位顺序表"},
		{url: "https://bf.topath.net.cn/zy/202002/5d14a9427cc16.mp4", name: "1.2.名师讲微课：亿以内数的读法"},
		{url: "http://n.topath.cn/shipin/shangce_2016_qiu_shang/xiaoxue_shuxue/wei_ke_tang/R4S_zrwk/R4S-yiyineishudedaxiaobijiao1.mp4", name: "1.3.名师讲微课：亿以内数的大小比较"},
		{url: "https://bf.topath.net.cn/zy/202002/5d14aa2c9cc0b.mp4", name: "1.4.名师讲微课：亿以上数的改写及求近似数"},
		{url: "https://bf.topath.net.cn/zy/202002/5d16d349e32b2.mp4", name: "1.5.名师讲奥数：巧猜大数"},
		{url: "http://n.topath.cn/shipin/shangce_2016_qiu_shang/xiaoxue_shuxue/wei_ke_tang/R4S_zrwk/R4S-gongqing.mp4", name: "2.1.名师讲微课：公顷"},
		{url: "https://bf.topath.net.cn/zy/202002/5d16d332381db.mp4", name: "2.2.名师讲奥数：圆形面积问题"},
		{url: "https://bf.topath.net.cn/zy/202002/5d14a92ca1be7.mp4", name: "4.1.名师讲微课：三位数乘两位数"},
		{url: "http://n.topath.cn/shipin/shangce_2016_qiu_shang/xiaoxue_shuxue/wei_ke_tang/R4S_zrwk/R4S-7-yinshuzhongjianhuomoweiyou0dechengfa.mp4", name: "4.2.名师讲微课：因数中间或末尾有0的乘法"},
		{url: "https://bf.topath.net.cn/zy/202002/5d14a9e3ef479.mp4", name: "4.3.名师讲微课：积的变化规律"},
		{url: "https://bf.topath.net.cn/zy/202002/5d16d2f2d2e2d.mp4", name: "4.4.名师讲奥数：火车过桥问题"},
		{url: "http://n.topath.cn/shipin/shangce_2016_qiu_shang/xiaoxue_shuxue/wei_ke_tang/R4S_zrwk/R4S-yongsishewurufashishang.mp4", name: "6.1.名师讲微课：用“四舍五入”法试商"},
		{url: "http://n.topath.cn/shipin/shangce_2016_qiu_shang/xiaoxue_shuxue/wei_ke_tang/R4S_zrwk/R4S-shangshiliangweishudebisuanchufa.mp4", name: "6.2.名师讲微课：商是两位数的除法笔算"},
		{url: "https://bf.topath.net.cn/zy/202002/5d14a95f5033f.mp4", name: "6.3.名师讲微课：商随除数(或被除数)的变化而变化的规律"},
		{url: "https://bf.topath.net.cn/zy/202002/5d16d2c8e383f.mp4", name: "6.4.名师讲奥数：错中求解问题"},
	}
	for _, course := range courses {
		DownloadFileProgress(course.url, fmt.Sprintf("%s.mp4", course.name))
		fmt.Println(course.name, "download over")
	}
}
