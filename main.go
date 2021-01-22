package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	bar "github.com/lycblank/goprogressbar"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	goQrCode "github.com/skip2/go-qrcode"
	"image"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"taobao_ok/serverJiang"
	"time"
)

func main() {
	///////////////////////////
	//chromdp依赖context上限传递参数
	ctx, _ := chromedp.NewExecAllocator(
		context.Background(),

		// 以默认配置的数组为基础，覆写headless参数
		// 当然也可以根据自己的需要进行修改，这个flag是浏览器的设置
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
		)...,
	)
	// 创建新的chromedp上下文对象，超时时间的设置不分先后
	// 注意第二个返回的参数是cancel()，只是我省略了
	ctx, _ = context.WithTimeout(ctx, 300*time.Second)
	ctx, _ = chromedp.NewContext(
		ctx,
		// 设置日志方法
		chromedp.WithLogf(log.Printf),
	)

	///////////////////////////

	// create chrome instance
	//ctx, cancel := chromedp.NewContext(
	//	context.Background(),
	//	chromedp.WithLogf(log.Printf),
	//)
	//defer cancel()

	///////////////////////////

	width, height := 1920, 1080

	if err := chromedp.Run(ctx, chromedp.EmulateViewport(int64(width), int64(height))); err != nil {
		log.Fatal(err)
	}

	// 加载cookies
	if err := chromedp.Run(ctx, loadCookies()); err != nil {
		log.Fatal(err)
	}

	//打开购物车页面
	if err := chromedp.Run(ctx, chromedp.Navigate("https://cart.taobao.com/cart.htm")); err != nil {
		log.Fatal(err)
	}

	var url string
	if err := chromedp.Run(ctx, chromedp.Evaluate(`window.location.href`, &url)); err != nil {
		log.Fatal(err)
	}

	if !strings.Contains(url, "https://cart.taobao.com/cart.htm") {
		log.Println("未登陆，进行登录操作")
		//打开购物车页面
		if err := chromedp.Run(ctx, login()); err != nil {
			log.Fatal(err)
		}
	}
	log.Println("已登录")

	for {
		fmt.Println("开始检查商品库存情况")
		go ss(ctx)
		time.Sleep(12 * time.Second)
		//执行任务
		if err := chromedp.Run(ctx, myTasks()); err != nil {
			log.Fatal(err)
		}

		rand.Seed(time.Now().UnixNano())
		//5-9之间的随机数
		num := rand.Intn(5) + 5
		fmt.Println("休息" + strconv.Itoa(num) + "分钟后查询")

		bar := bar.NewProgressBar(int64(num) * 60)
		for i := 1; i <= int(num*60); i++ {
			bar.Play(int64(i))
			time.Sleep(1 * time.Second)
		}
		bar.Finish()
	}

}
func ss(ctx context.Context) {
	defer chromedp.Stop()
	if err := chromedp.Run(ctx, chromedp.Navigate("https://detail.tmall.com/item.htm?id=624275284278")); err != nil {
		log.Fatal(err)
	}
	time.Sleep(10 * time.Second)
}

func login() chromedp.Tasks {
	return chromedp.Tasks{
		// 1. 点击登陆按钮
		// #login > div.corner-icon-view.view-type-qrcode > i
		chromedp.Click(`#login > div.corner-icon-view.view-type-qrcode > i`),

		// 2. 获取二维码
		// #login > div.login-content.nc-outer-box > div > div:nth-child(1) > div.qrcode-img > canvas
		getCode(),

		// 3.保存cookie
		saveCookies(),
	}

}

func myTasks() chromedp.Tasks {
	return chromedp.Tasks{
		checkGoods(),
	}

}

// 检查码数
func checkGoods() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		//#J_OrderList
		if err = chromedp.WaitVisible(`#J_DetailMeta > div.tm-clear > div.tb-property > div > div.tb-key > div > div > dl:nth-child(1) > dd > ul`, chromedp.ByID).Do(ctx); err != nil {
			return
		}
		var res string
		//#J_DetailMeta > div.tm-clear > div.tb-property > div > div.tb-key > div > div > dl:nth-child(1) > dd > ul
		if err = chromedp.Text(`#J_DetailMeta > div.tm-clear > div.tb-property > div > div.tb-key > div > div > dl:nth-child(1) > dd > ul`, &res, chromedp.ByID).Do(ctx); err != nil {
			return
		}
		log.Println(res)

		if strings.Contains(strings.TrimSpace(res), "175/88A") {
			//http://sc.ftqq.com/?c=code
			var s serverJiang.ServerJiang
			ss := make(map[string]string)
			ss["text"] = "商品到货通知"
			ss["desp"] = "优衣库到货了！！！！"
			s.Data = ss
			s.SCKey = "xxx"
			s.Do()
		}

		return
	}
}

func getCode() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		// 1. 用于存储图片的字节切片
		var code []byte

		// 2. 截图
		//#login > div.login-content.nc-outer-box > div > div:nth-child(1) > div.qrcode-img > canvas
		//#login > div.login-content.nc-outer-box > div > div:nth-child(1) > div.qrcode-img
		// 注意这里需要注明直接使用ID选择器来获取元素（chromedp.ByID）
		//#login
		if err = chromedp.Screenshot(`#login > div.login-content.nc-outer-box > div > div:nth-child(1) > div.qrcode-img > canvas`, &code, chromedp.ByID).Do(ctx); err != nil {
			return
		}

		// 3. 保存文件
		if err = ioutil.WriteFile("code.png", code, 0755); err != nil {
			return
		}
		// 4. 把二维码输出到标准输出流
		if err = printQRCode(code); err != nil {
			return err
		}

		return
	}
}

func printQRCode(code []byte) (err error) {
	// 1. 因为我们的字节流是图像，所以我们需要先解码字节流
	img, _, err := image.Decode(bytes.NewReader(code))
	if err != nil {
		return
	}

	// 2. 然后使用gozxing库解码图片获取二进制位图
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return
	}

	// 3. 用二进制位图解码获取gozxing的二维码对象
	res, err := qrcode.NewQRCodeReader().Decode(bmp, nil)
	if err != nil {
		return
	}

	// 4. 用结果来获取go-qrcode对象（注意这里我用了库的别名）
	qr, err := goQrCode.New(res.String(), goQrCode.High)
	if err != nil {
		return
	}

	// 5. 输出到标准输出流
	fmt.Println(qr.ToSmallString(true))
	return
}

// 保存Cookies
func saveCookies() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		// 等待二维码登陆
		//#J_OrderList
		if err = chromedp.WaitVisible(`#J_OrderList`, chromedp.ByID).Do(ctx); err != nil {
			return
		}

		// cookies的获取对应是在devTools的network面板中
		// 1. 获取cookies
		cookies, err := network.GetAllCookies().Do(ctx)
		if err != nil {
			return
		}

		// 2. 序列化
		cookiesData, err := network.GetAllCookiesReturns{Cookies: cookies}.MarshalJSON()
		if err != nil {
			return
		}

		// 3. 存储到临时文件
		if err = ioutil.WriteFile("cookies.tmp", cookiesData, 0755); err != nil {
			return
		}
		return
	}
}

// 加载Cookies
func loadCookies() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		// 如果cookies临时文件不存在则直接跳过
		if _, _err := os.Stat("cookies.tmp"); os.IsNotExist(_err) {
			return
		}

		// 如果存在则读取cookies的数据
		cookiesData, err := ioutil.ReadFile("cookies.tmp")
		if err != nil {
			return
		}

		// 反序列化
		cookiesParams := network.SetCookiesParams{}
		if err = cookiesParams.UnmarshalJSON(cookiesData); err != nil {
			return
		}

		// 设置cookies
		return network.SetCookies(cookiesParams.Cookies).Do(ctx)
	}
}
