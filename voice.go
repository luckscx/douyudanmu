package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	log "github.com/sirupsen/logrus"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func playWavFile(text string, content []byte) {
	buf := ioutil.NopCloser(bytes.NewReader(content))
	s, format, _ := wav.Decode(buf)

	// Init the Speaker with the SampleRate of the format and a buffer size of 1/10s
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// Now we Play our Streamer on the Speaker
	speaker.Play(s)

	sleepTime := time.Duration(len(text)) * 80 * time.Millisecond
	//log.Infof("Sleep %v", sleepTime)
	time.Sleep(sleepTime)
	// Channel, which will signal the end of the playback.
	//playing := make(chan struct{})

	//// Now we Play our Streamer on the Speaker
	//speaker.Play(beep.Seq(s, beep.Callback(func() {
	//log.Infof("finish play")
	//s.Close()
	//// Callback after the stream Ends
	//close(playing)
	//})))
	//<-playing
}

func randVoice() string {
	var voiceList = []string{"xiaoyun", "xiaogang", "xiaowei", "amei", "xiaoxue", "siqi", "sijia", "ruoxi", "xiaomeng"}

	voice := voiceList[rand.Intn(len(voiceList))]
	//log.Infof("pick %s", voice)
	return voice
}

func getAudioWav(text string) []byte {
	var textURLEncode = text
	textURLEncode = url.QueryEscape(textURLEncode)
	textURLEncode = strings.Replace(textURLEncode, "+", "%20", -1)
	textURLEncode = strings.Replace(textURLEncode, "*", "%2A", -1)
	textURLEncode = strings.Replace(textURLEncode, "%7E", "~", -1)

	url := "https://nls-gateway.cn-shanghai.aliyuncs.com/stream/v1/tts"
	url = url + "?appkey=" + aliyunAppKey
	url = url + "&token=" + aliyunToken
	url = url + "&text=" + textURLEncode
	url = url + "&format=wav"
	url = url + "&sample_rate=" + strconv.Itoa(16000)
	// voice 发音人，可选，默认是xiaoyun
	url = url + "&voice=" + randVoice()
	// volume 音量，范围是0~100，可选，默认50
	// url = url + "&volume=" + strconv.Itoa(50)
	// speech_rate 语速，范围是-500~500，可选，默认是0
	url = url + "&speech_rate=" + strconv.Itoa(100)
	// pitch_rate 语调，范围是-500~500，可选，默认是0
	// url = url + "&pitch_rate=" + strconv.Itoa(0)
	//fmt.Println(url)
	/**
	 * 发送HTTPS GET请求，处理服务端的响应
	 */
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("The GET request failed!")
		panic(err)
	}
	defer response.Body.Close()
	contentType := response.Header.Get("Content-Type")
	body, _ := ioutil.ReadAll(response.Body)
	if "audio/mpeg" == contentType {
	} else {
		statusCode := response.StatusCode
		fmt.Println("The HTTP statusCode: " + strconv.Itoa(statusCode))
		fmt.Println("The GET request failed: " + string(body))
	}
	return body
}

func trans(text string) {
	if c := utf8.RuneCountInString(text); c > 30 {
		log.Infof("skip long line %d", c)
		return
	}
	res := getAudioWav(text)
	playWavFile(text, res)
}
