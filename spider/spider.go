package spider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var cityId map[string]string

type Spider struct {
	clint     *http.Client
	userAgent string
	host      string
}

func NewSpider(userAgent, host string) *Spider {
	if userAgent == "" {
		userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.61 Safari/537.36"
	}
	s := &Spider{
		clint:     &http.Client{},
		userAgent: userAgent,
		host:      host,
	}

	return s
}

func (s *Spider) GetWeather(url, city string) (string, error) {
	cid, ok := cityId[city]
	if !ok {
		return "找不到该城市的天气", nil
	}
	url = url + cid
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("err:", err)
		return "", err
	}
	request.Header.Add("User-Agent", s.userAgent)
	request.Header.Add("Host", s.host)
	response, err := s.clint.Do(request)
	if err != nil {
		fmt.Println("err:", err)
		return "", err
	}
	fmt.Println(response)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("err:", err)
		return "", nil
	}
	result := WeatherResult{}
	json.Unmarshal([]byte(string(body)), &result)
	cityWeather := result.City + "近7日天气\n"
	for i := 0; i < 7; i++ {
		d := result.Data[i]
		cityWeather += d.Date + " " + d.Wea + " "
		cityWeather += "白天温度：" + d.Tem_day + " "
		cityWeather += "夜晚温度：" + d.Tem_night + "\n"
	}
	return cityWeather, nil
}

type WeatherResult struct {
	Cityid      string  `json:"cityid"`
	City        string  `json:"city"`
	Update_time string  `json:"update_time"`
	Data        []Wdata `json:"data"`
}

type Wdata struct {
	Date      string `json:"date"`
	Wea       string `json:"wea"`
	Wea_img   string `json:"wea_img"`
	Tem_day   string `json:"tem_day"`
	Tem_night string `json:"tem_night"`
	Win       string `json:"win"`
	Win_speed string `json:"win_speed"`
}

func init() {
	cityId = make(map[string]string)
	cityId["广州"] = "101280101"
	cityId["深圳"] = "101280601"
	cityId["珠海"] = "101280701"
	cityId["佛山"] = "101280800"
	cityId["肇庆"] = "101280901"
	cityId["湛江"] = "101281001"
	cityId["江门"] = "101281101"
	cityId["揭阳"] = "101281901"
	cityId["阳江"] = "101281801"
	cityId["潮州"] = "101281501"
	cityId["清远"] = "101281301"
}
