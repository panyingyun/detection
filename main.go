package main

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/corona10/goimagehash"
	"github.com/mattermost/mattermost-server/model"
	"github.com/panyingyun/detection/config"
	"github.com/urfave/cli"
)

func catch_orgjpg() {
	_, err := exec_shell("raspistill  -t 1 -w 800 -h 600 -br 60 -rot 180 -o a.jpg")
	if err != nil {
		fmt.Println("catch orgjpg err = ", err)
	}
}

func catch_newjpg() {
	_, err := exec_shell("raspistill  -t 1 -w 800 -h 600 -br 60 -rot 180  -o c.jpg")
	if err != nil {
		fmt.Println("catch newjpg err = ", err)
	}
}

//run extern shell
func exec_shell(s string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", s)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("err = ", err)
	}
	return out.String(), err
}

//compare two jpg and return distance
func comparejpg(orgjpg, newjpg string) (int, error) {
	orgfile, _ := os.Open(orgjpg)
	newfile, _ := os.Open(newjpg)
	defer orgfile.Close()
	defer newfile.Close()
	img1, _ := jpeg.Decode(orgfile)
	img2, _ := jpeg.Decode(newfile)
	hash1, _ := goimagehash.AverageHash(img1)
	hash2, _ := goimagehash.AverageHash(img2)
	distance, err := hash1.Distance(hash2)
	return distance, err
}

//send post with pictures
func sendmessage(server string, username string, passwd string, team string, chname string, message string, newjpg string) error {
	//Login
	Client := model.NewAPIv4Client(server)
	mine, resp := Client.Login(username, passwd)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("LOGIN: mine = ", mine)
		fmt.Println("LOGIN: resp = ", resp.StatusCode)
		return errors.New("Login fail,please check server or check user")
	}

	//Get channel
	channel, chresp := Client.GetChannelByNameForTeamName(chname, team, "")
	if chresp.StatusCode != http.StatusOK {
		fmt.Println("CH: channel = ", channel)
		fmt.Println("CH: chresp = ", chresp.StatusCode)
		return errors.New("GetChannel fail,please check server or check user")
	}

	//Upload file
	file, err := os.Open(newjpg)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, file)
	data := buf.Bytes()
	chid := channel.Id
	fileUploadResponse, uploadresp := Client.UploadFile(data, chid, newjpg)
	if uploadresp.StatusCode != http.StatusCreated {
		fmt.Println("UPLOAD: fileUploadResponse = ", fileUploadResponse)
		fmt.Println("UPLOAD: uploadresp = ", uploadresp.StatusCode)
		return errors.New("UploadFile fail,please check server or check user")
	}

	infos := fileUploadResponse.FileInfos
	fileid := infos[0].Id
	fmt.Println("UPLOAD: fileid = ", fileid)

	//send a message with picture
	var post model.Post
	post.ChannelId = channel.Id
	post.Message = message
	var pictures model.StringArray
	pictures = append(pictures, fileid)
	post.FileIds = pictures
	pt, postresp := Client.CreatePost(&post)
	if postresp.StatusCode != http.StatusCreated {
		fmt.Println("CreatePost: post = ", post)
		fmt.Println("CreatePost: pt = ", pt)
		fmt.Println("CreatePost: postresp = ", postresp)
		return errors.New("CreatePost fail,please check server or check user")
	}
	return nil
}

func run(c *cli.Context) error {
	conf, err := config.ReadConfig(c.String("conf"))
	if err != nil {
		fmt.Println("read from conf fail!", c.String("conf"))
		return err
	}
	fmt.Println("conf =  ", conf)

	server := conf.Server
	username := conf.Username
	passwd := conf.Passwd
	team := conf.Team
	chname := conf.Chname
	orgjpg := conf.Orgjpg
	newjpg := conf.Newjpg
	dist := conf.Distance

	fmt.Println("dist =  ", dist)

	catch_orgjpg()

	catch_newjpg()

	go func() {
		for {
			//starttime := time.Now()
			catch_newjpg()
			distance, _ := comparejpg(orgjpg, newjpg)
			//fmt.Printf("dtime = %v\n", time.Now().Sub(starttime))

			fmt.Printf("Distance between images: %v\n", distance)
			if distance >= dist {
				message := "@channel ^o^catch you^o^ 于时间：" + time.Now().Format("2006-01-02 15:04:05")
				err = sendmessage(server, username, passwd, team, chname, message, newjpg)
				if err != nil {
					fmt.Println("sendmessage err = ", err)
				} else {
					fmt.Println("sendmessage OK!")
				}
			} else if distance > 3 {
				catch_newjpg()
				fmt.Println("change org image OK!")
			}
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	fmt.Printf("signal received signal %v and showdown server\n", <-ch)
	return err
}

func main() {
	app := cli.NewApp()
	app.Name = "detection"
	app.Usage = "detection -f appserver.conf"
	app.Copyright = "panyingyun (at) gmail.com"
	app.Version = "0.0.1"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "conf,c",
			Usage:  "Set conf path here",
			Value:  "appserver.conf",
			EnvVar: "APP_CONF",
		},
	}
	app.Run(os.Args)
}
