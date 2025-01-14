package shape

import (
	"fmt"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/rod/lib/proto"

	"github.com/go-rod/stealth"
)

func init() {
	launcher.NewBrowser().MustGet()
}

func NewBrowser(proxy string) *rod.Browser {
	var browser *rod.Browser

	if proxy != "" {
		// Split the proxy string into its components.
		components := strings.Split(proxy, ":")
		ip := components[0]
		port := components[1]
		username := components[2]
		password := components[3]

		// Set the proxy server flag using the extracted components.
		l := launcher.New().Headless(false)
		l = l.Set(flags.ProxyServer, ip+":"+port)
		fmt.Println(l)
		// Launch the browser with the proxy server flag.
		controlURL, _ := l.Launch()
		browser = rod.New().ControlURL(controlURL).MustConnect().MustIncognito()
		// Handle the proxy server authentication.
		go browser.MustHandleAuth(username, password)()

		browser.MustIgnoreCertErrors(true)
	} else {
		// url, err := launcher.New().
		// 	Headless(false).
		// 	Devtools(false).
		// 	Launch()
		// if err != nil {
		// 	panic(err)
		// }
		url := "ws://127.0.0.1:9222/devtools/browser/96789d07-15ed-4a61-a5e8-194ab6bc4c64" //braveURL//
		browser = rod.New().ControlURL(url).MustConnect()
		//go browser.MustHandleAuth("OR750076101", "420wquh")()
	}

	return browser
}

func NewPage(browser *rod.Browser) *rod.Page {
	page := stealth.MustPage(browser)
	//page = page
	return page
}

func LoadSite(page *rod.Page, address string) {
	page.MustNavigate(address)
}

type ShapeHeaders struct {
	XDQ7Hy5L1a0 string `json:"X-DQ7Hy5L1-a0"`
	XDQ7Hy5L1a  string `json:"X-DQ7Hy5L1-a"`
	XDQ7Hy5L1b  string `json:"X-DQ7Hy5L1-b"`
	XDQ7Hy5L1c  string `json:"X-DQ7Hy5L1-c"`
	XDQ7Hy5L1d  string `json:"X-DQ7Hy5L1-d"`
	XDQ7Hy5L1f  string `json:"X-DQ7Hy5L1-f"`
	XDQ7Hy5L1z  string `json:"X-DQ7Hy5L1-z"`
}

type ShapeHarvester struct {
	Proxy          string
	Url            string
	ShapeUrl       string
	Identifier     string
	Method         string
	Body           string
	Headers        ShapeHeaders
	ReqHeaders     proto.NetworkHeaders
	ReqPayload     string
	Page           *rod.Page
	Browser        *rod.Browser
	BlockResources bool
}

func (harvester *ShapeHarvester) HarvestHeaders() {
	// harvester.Page.MustEval(`function test() {
	// 	try {
	// 			Object.defineProperty(navigator, 'userAgent', {
	// 					get: () => 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36',
	// 			});

	// 			Object.defineProperty(navigator, 'webdriver', {
	// 					get: () => false,
	// 			});

	// 			Object.defineProperty(navigator, 'plugins', {
	// 					get: () => [1,2,3,4,5],
	// 			});

	// 			Object.defineProperty(navigator, 'languages', {
	// 					get: () => ['en-US', 'en'],
	// 			});

	// 			Object.defineProperty(navigator, 'maxTouchPoints', {
	// 					get: () => 1
	// 			});

	// 			window.chrome = {
	// 					runtime: {},
	// 			};
	// 	} catch {}
	// };`)
	harvester.Page.MustEval(fmt.Sprintf(`function shape() {
		try {
			fetch("%s", {
				"method" : "%s",
				"referrerPolicy": "no-referrer-when-downgrade",
				"credentials": "include",
				"body": "%s",
				"headers": {
					"accept": "application/json",
					"accept-language": "en-US,en;q=0.9",
					"content-type": "application/json",
					"sec-ch-ua": "\"Google Chrome\";v=\"107\", \"Chromium\";v=\"107\", \"Not=A?Brand\";v=\"24\"",
					"sec-ch-ua-mobile": "?0",
					"sec-ch-ua-platform": "\"Windows\"",
					"sec-fetch-dest": "empty",
					"sec-fetch-mode": "cors",
					"sec-fetch-site": "same-site",
					"x-application-name": "web",
				},
			})
		} catch {}
	  }`, harvester.ShapeUrl, harvester.Method, harvester.Body))
}

func (harvester *ShapeHarvester) InitializeHijacking() {
	router := harvester.Page.HijackRequests()

	router.MustAdd("*", func(ctx *rod.Hijack) {
		if harvester.BlockResources {
			if ctx.Request.Method() == "GET" {
				if ctx.Request.Type() == proto.NetworkResourceTypeImage || ctx.Request.Type() == proto.NetworkResourceTypeStylesheet {
					ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
				}

				ctx.ContinueRequest(&proto.FetchContinueRequest{})
			}
		}

		if strings.Contains(ctx.Request.URL().Path, harvester.Identifier) {
			if ctx.Request.Method() == "OPTIONS" {
				ctx.ContinueRequest(&proto.FetchContinueRequest{})
			} else if ctx.Request.Method() == "POST" {
				//fmt.Println(ctx.Request.Headers())
				//fmt.Println(ctx.Request.Headers())
				ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
				harvester.ReqHeaders = ctx.Request.Headers()
				harvester.ReqPayload = ctx.Request.Body()
				// for key, value := range ctx.Request.Headers() {
				// 	switch key {
				// 	case "X-DQ7Hy5L1-a0":
				// 		harvester.Headers.XDQ7Hy5L1a0 = value.String()
				// 	case "X-DQ7Hy5L1-a":
				// 		harvester.Headers.XDQ7Hy5L1a = value.String()

				// 	case "X-DQ7Hy5L1-b":
				// 		harvester.Headers.XDQ7Hy5L1b = value.String()

				// 	case "X-DQ7Hy5L1-c":
				// 		harvester.Headers.XDQ7Hy5L1c = value.String()

				// 	case "X-DQ7Hy5L1-d":
				// 		harvester.Headers.XDQ7Hy5L1d = value.String()

				// 	case "X-DQ7Hy5L1-f":
				// 		harvester.Headers.XDQ7Hy5L1f = value.String()

				// 	case "X-DQ7Hy5L1-z":
				// 		harvester.Headers.XDQ7Hy5L1z = value.String()
				// 	}
				// }
			}
		}

	})

	go router.Run()
}

func (harvester *ShapeHarvester) InitializeHarvester(proxy string) {
	harvester.Browser = NewBrowser(proxy)
	harvester.Page = NewPage(harvester.Browser)

	harvester.Page.MustNavigate(harvester.Url).MustWaitLoad()

	harvester.InitializeHijacking()
	harvester.HarvestHeaders()
}
