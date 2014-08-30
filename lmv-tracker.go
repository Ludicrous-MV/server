package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strconv"
	"syscall"

    "github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type LMVFile struct {
	Id        int64      `json:"-"`
	Size      int64      `bson:"size"         json:"size"         binding:"required"`
	Name      string     `bson:"name"         json:"name"         binding:"required"`
	Algorithm string     `bson:"algorithm"    json:"algorithm"    binding:"required"`
	Chunks    []LMVChunk `bson:"chunks"       json:"chunks"       binding:"required"`
	Tar       bool       `bson:"tar"          json:"tar"`
	Token     string     `bson:"token"        json:"token"`
}

type LMVChunk struct {
	Id        int64  `json:"-"`
	LMVFileId int64  `json:"-"`
	Hash      string `bson:"hash"         json:"hash"         binding:"required"`
	Size      int64  `bson:"size"         json:"size"         binding:"required"`
	Index     int    `bson:"index"        json:"index"        binding:"required"`
}

const (
	token_length = 10
)

func GenerateToken() string {

    pool := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    return uniuri.NewLenChars(token_length, pool)

}

func main() {

	pid := flag.Bool("pid", false, "Save the PID to lmv-server.pid")

	flag.Parse()

	if *pid {
		ioutil.WriteFile("lmv-tracker.pid", []byte(strconv.Itoa(syscall.Getpid())), 0644)
	}

	db, err := gorm.Open("sqlite3", "lmv-tracker.db")

	if err != nil {
		log.Fatal(err)
	}

	db.DB()
	db.DB().Ping()

	db.CreateTable(LMVFile{})
	db.CreateTable(LMVChunk{})

	r := gin.Default()

	r.GET("/files/", func(gc *gin.Context) {
		var lmvfiles []LMVFile
		var response []LMVFile

		db.Find(&lmvfiles)

		for _, file := range lmvfiles {

			var chunks []LMVChunk

			db.Where(&LMVChunk{LMVFileId: file.Id}).Find(&chunks)

			file.Chunks = chunks

			response = append(response, file)
		}

		gc.JSON(200, response)
	})

	r.GET("/files/:token/", func(gc *gin.Context) {

		token := gc.Params.ByName("token")
		var lmv_file LMVFile
		var response LMVFile

		db.Where(&LMVFile{Token: token}).First(&lmv_file)

		if lmv_file.Name != "" {

			var chunks []LMVChunk

			db.Where(&LMVChunk{LMVFileId: lmv_file.Id}).Find(&chunks)

			response = lmv_file
			response.Chunks = chunks

			gc.JSON(200, response)

		} else {

			gc.JSON(404, nil)

		}

	})

	r.POST("/files/", func(gc *gin.Context) {

		token := GenerateToken()
		var lmv_file LMVFile

		gc.Bind(&lmv_file)

		for {

			var testFile = LMVFile{}
			db.Where(&LMVFile{Token: token}).First(&testFile)

			if testFile.Name != "" {

				token = GenerateToken()

			} else {

				break

			}

		}

		lmv_file.Token = token

		db.Create(&lmv_file)

		gc.JSON(200, map[string]interface{}{"token": token})

	})

	r.GET("/ping/", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.Run(":8080")

}
