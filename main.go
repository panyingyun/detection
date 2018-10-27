package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
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
	_, err := exec_shell("raspistill  -t 1 -w 1024 -h 768 -br 60 -vf  -o a.jpg")
	fmt.Println("catch_orgjpg err = ", err)
}

func catch_newjpg() {
	_, err := exec_shell("raspistill  -t 1 -w 1024 -h 768 -br 60 -vf  -o c.jpg")
	fmt.Println("catch_newjpg err = ", err)
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
	fmt.Println("distance = ", distance)
	fmt.Println("err = ", err)
	return distance, err
}

//send post with pictures
func sendmessage(server string, username string, passwd string, team string, chname string, message string, newjpg string) error {
	//Login
	Client := model.NewAPIv4Client(server)
	mine, resp := Client.Login(username, passwd)
	fmt.Println("LOGIN: mine = ", mine)
	fmt.Println("LOGIN: resp = ", resp.StatusCode)
	//Get channel
	channel, chresp := Client.GetChannelByNameForTeamName(chname, team, "")
	fmt.Println("CH: channel = ", channel)
	fmt.Println("CH: chresp = ", chresp.StatusCode)
	//Upload file
	file, err := os.Open(newjpg)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, file)
	data := buf.Bytes()
	chid := channel.Id
	fmt.Println("UPLOAD: filename = ", newjpg)
	fileUploadResponse, uploadresp := Client.UploadFile(data, chid, newjpg)

	fmt.Println("UPLOAD: fileUploadResponse = ", fileUploadResponse)
	fmt.Println("UPLOAD: uploadresp = ", uploadresp.StatusCode)
	infos := fileUploadResponse.FileInfos
	fileid := infos[0].Id
	fmt.Println("UPLOAD: fileinfos = ", infos)
	fmt.Println("UPLOAD: fileid = ", fileid)

	//send a message with picture
	var post model.Post
	post.ChannelId = channel.Id
	post.Message = message
	var pictures model.StringArray
	pictures = append(pictures, fileid)
	post.FileIds = pictures
	pt, postresp := Client.CreatePost(&post)
	fmt.Println("CreatePost: post = ", post)
	fmt.Println("CreatePost: pt = ", pt)
	fmt.Println("CreatePost: postresp = ", postresp)
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

	catch_orgjpg()

	catch_newjpg()

	go func() {
		for {
			starttime := time.Now()

			catch_newjpg()

			distance, _ := comparejpg(orgjpg, newjpg)

			fmt.Printf("Distance between images: %v\n", distance)

			fmt.Printf("dtime = %v ", time.Now().Sub(starttime))

			if distance >= 15 {
				message := "@channel ^o^catch you^o^ 于时间：" + time.Now().Format("2006-01-02 15:04:05")
				err = sendmessage(server, username, passwd, team, chname, message, newjpg)
			}
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	fmt.Printf("signal received signal %v and showdown server", <-ch)
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
