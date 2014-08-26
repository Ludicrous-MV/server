package main

import (
	"crypto/rand"
	"flag"
	"log"
	"io/ioutil"
	"strconv"
	"syscall"

    "github.com/gin-gonic/gin"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type LMVFile struct {
	Size      int64      `bson:"size"         json:"size"         binding:"required"`
	Name      string     `bson:"name"         json:"name"         binding:"required"`
	Algorithm string     `bson:"algorithm"    json:"algorithm"    binding:"required"`
	Chunks    []LMVChunk `bson:"chunks"       json:"chunks"       binding:"required"`
	Tar       bool       `bson:"tar"          json:"tar"`
	Token     string     `bson:"token"        json:"token"`
}

type LMVChunk struct {
	Hash  string `bson:"hash"         json:"hash"         binding:"required"`
	Size  int64  `bson:"size"         json:"size"         binding:"required"`
	Index int    `bson:"index"        json:"index"        binding:"required"`
}

const (
	token_length = 10
	mgo_host     = "localhost"
	mgo_db       = "Ludicrous-MV"
	mgo_col      = "Files"
)

func randstr(length int) string {

	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var bytes = make([]byte, length)

	rand.Read(bytes)

	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}

	return string(bytes)

}

func main() {

	pid := flag.Bool("pid", false, "Save the PID to lmv-server.pid")

	flag.Parse()

	if *pid {
		ioutil.WriteFile("lmv-tracker.pid", []byte(strconv.Itoa(syscall.Getpid())), 0644)
	}

	session, err := mgo.Dial(mgo_host)

	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB(mgo_db).C(mgo_col)

    r := gin.Default()

    r.GET("/files/", func(gc *gin.Context) {
        var lmv_files []LMVFile

        err := c.Find(bson.M{}).All(&lmv_files)

        if err != nil {
            log.Fatal(err)
        }

        gc.JSON(200, lmv_files)
    })

    r.GET("/files/:token/", func(gc *gin.Context) {

        token := gc.Params.ByName("token")

		n, err := c.Find(bson.M{"token": token}).Count()

		if err != nil {
			log.Fatal(err)
		}

		if n != 1 {
			gc.JSON(404, "")
		} else {
			var lmv_file LMVFile

			err = c.Find(bson.M{"token": token}).One(&lmv_file)

			if err != nil {
				log.Fatal(err)
			}

			gc.JSON(200, lmv_file)
		}

	})

    r.POST("/files/", func (gc *gin.Context) {

		token := randstr(token_length)
        var lmv_file LMVFile

        gc.Bind(&lmv_file)

		for {
			n, err := c.Find(bson.M{"token": token}).Count()

			if err != nil {
				log.Fatal(err)
			}

			if n > 0 {
				token = randstr(token_length)
			} else {
				break
			}
		}

		lmv_file.Token = token

		err := c.Insert(lmv_file)

		if err != nil {
			log.Fatal(err)
		}

		gc.JSON(200, map[string]interface{}{"token": token})

	})

    r.GET("/ping/", func(c *gin.Context) {
        c.String(200, "pong")
    })

    r.Run(":8080")

}
