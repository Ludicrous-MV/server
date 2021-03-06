package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
	"syscall"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tsuru/config"
)

type LMVFile struct {
	Id        int64      `json:"-"`
	Size      int64      `json:"size"         binding:"required"`
	Name      string     `json:"name"         binding:"required"`
	Algorithm string     `json:"algorithm"    binding:"required"`
	Chunks    []LMVChunk `json:"chunks"       binding:"required"`
	Tar       bool       `json:"tar"`
	Token     string     `json:"token"`
}

type LMVChunk struct {
	Id        int64  `json:"-"`
	LMVFileId int64  `json:"-"`
	Hash      string `json:"hash"         binding:"required"`
	Size      int64  `json:"size"         binding:"required"`
	Index     int    `json:"index"        binding:"required"`
}

type Configuration struct {
	Web struct {
		Address string
	}
	System struct {
		Pid bool
	}
	Tokens struct {
		Pool   []byte
		Length int
	}
	Database struct {
		Type   string
		Source string
	}
}

var conf Configuration

func GenerateToken() string {

	return uniuri.NewLenChars(conf.Tokens.Length, conf.Tokens.Pool)

}

func processConfig() {

	foundConf := true

	if _, err := os.Stat("lmv-tracker.yml"); err == nil {
		config.ReadConfigFile("lmv-tracker.yml")
	} else {
		usr, err := user.Current()

		if err != nil {
			log.Fatal(err)
		}

		if _, err := os.Stat(usr.HomeDir + "/lmv-tracker.yml"); err == nil {
			config.ReadConfigFile(usr.HomeDir + "/lmv-tracker.yml")
		} else {
			if _, err := os.Stat("/etc/lmv-tracker.yml"); err == nil {
				config.ReadConfigFile("/etc/lmv-tracker.yml")
			} else {
				foundConf = false
			}
		}
	}

	if foundConf {
		address, _ := config.GetString("web:address")
		conf.Web.Address = address

		pid, _ := config.GetBool("system:pid")
		conf.System.Pid = pid

		token_pool, _ := config.GetString("tokens:pool")
		conf.Tokens.Pool = []byte(token_pool)

		token_length, _ := config.GetInt("tokens:length")
		conf.Tokens.Length = token_length

		database_type, _ := config.GetString("database:type")
		conf.Database.Type = database_type

		database_source, _ := config.GetString("database:source")
		conf.Database.Source = database_source
	} else {
		conf.Web.Address = ":8080"
		conf.System.Pid = false
		conf.Tokens.Pool = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
		conf.Tokens.Length = 10
		conf.Database.Type = "sqlite3"
		conf.Database.Source = "lmv-tracker.db"
	}

}

func main() {

	processConfig()

	if conf.System.Pid {
		ioutil.WriteFile("lmv-tracker.pid", []byte(strconv.Itoa(syscall.Getpid())), 0644)
	}

	db, err := gorm.Open(conf.Database.Type, conf.Database.Source)

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

		gc.JSON(200, lmv_file)

	})

	r.GET("/ping/", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.Run(conf.Web.Address)

}
