package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/heroku/x/hmetrics/onload"
)

type Lock struct {
	Success string `json:"Success"`
	Guid    string `json:"Guid"`
	Brand   string `json:"Brand"`
	Locked  string `json:"Locked"`
}

func getRedisDb() *redis.Client {
	port := os.Getenv("REDISCLOUD_URL")
	opt, err := redis.ParseURL(port)
	if err != nil {
		panic(err)
	}
	return redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Password: opt.Password, // no password set
		DB:       opt.DB,       // use default DB
	})
	// return redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379", // use default Addr
	// 	Password: "",               // no password set
	// 	DB:       0,                // use default DB
	// })
}

// func keyExists(key string) *redis.IntCmd {
// 	rdb := getRedisDb()
// 	return rdb.Exists(key)
// }
func getKey(key string) string {
	rdb := getRedisDb()
	val, err := rdb.Get(key).Result()
	fmt.Println("get key val = ", val)
	fmt.Println("get key err = ", err)
	if err != nil {
		return ""
		// panic(err)
	}
	return val
}
func delKey(key string) int64 {
	rdb := getRedisDb()
	fmt.Println("    del key= ", key)
	val, err := rdb.Del(key).Result()
	fmt.Println("del key val = ", val)
	fmt.Println("del key err = ", err)
	if err != nil {
		return 0
		// panic(err)
	}
	return val
}
func setKey(key string, val string) string {
	rdb := getRedisDb()
	// fmt.Println("setKey key= ", key, "val ", val)
	val, err := rdb.Set(key, val, 0).Result()
	if err != nil {
		return ""
		// panic(err)
	}
	// fmt.Println("setKey err = ")
	// fmt.Println(err)
	// fmt.Println("setKey value = ")
	// fmt.Println(val)
	return val
}

// func getSetKey(key string, val string) *redis.StringCmd {
// 	rdb := getRedisDb()
// 	fmt.Println("getSetKey key= ", key, "val ", val)
// 	// val, err := rdb.Get(key).Result()
// 	getSetRes := rdb.GetSet("locks:b39bfc55-450c-4695-b63b-e15be698c377", "gcl")
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// fmt.Println("getKey err = ")
// 	// fmt.Println(err)
// 	// fmt.Println("getKey value = ")
// 	fmt.Println("getSetRes")
// 	fmt.Println(getSetRes)
// 	fmt.Println(getSetRes.Val())
// 	fmt.Println(getSetRes.Result())

// 	return getSetRes
// }
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	testData := Lock{Success: "1", Guid: "asdf-1234", Brand: "test"}

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, testData)
	})

	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379", // use default Addr
	// 	Password: "",               // no password set
	// 	DB:       0,                // use default DB
	// })
	// // pong, err := rdb.Ping().Result()
	// // fmt.Println(pong, err)
	// uniqKey := "test_cucumber_emails:gcl:2020-09-22:1600804232"
	// res := rdb.Exists(uniqKey)
	// fmt.Println("exists: ")
	// fmt.Println(res)

	// resGet := rdb.Get(uniqKey)
	// fmt.Println("resGet: ")
	// fmt.Println(resGet)
	// fmt.Println(resGet.Val())
	// fmt.Println(resGet.Result())

	// data2.Guid = resGet.Val()
	// brand := getKey("locks:b39bfc55-450c-4695-b63b-e15be698c377")
	// setKey("locks:b39bfc55-450c-4695-b63b-e15be698c377", "snnow")
	router.GET("/is_locked", func(c *gin.Context) {
		data := Lock{Success: "true", Guid: c.Query("guid"), Brand: c.Query("brand")}

		brand := getKey("locks:" + data.Guid)
		if brand == "" {
			data.Locked = "false"
		} else {
			data.Locked = "true"
		}
		c.JSON(http.StatusOK, data)
	})
	router.GET("/set_lock", func(c *gin.Context) {
		data := Lock{Success: "true", Guid: c.Query("guid"), Brand: c.Query("brand")}

		brand := getKey("locks:" + data.Guid)
		if brand == "" {
			brand = setKey("locks:"+data.Guid, data.Brand)
			if brand == "" {
				data.Locked = "false"
				data.Success = "false"
			} else {
				data.Locked = "true"
			}
		} else {
			data.Success = "false"
			data.Locked = "true"
		}

		c.JSON(http.StatusOK, data)
	})
	router.GET("/free_lock", func(c *gin.Context) {
		data := Lock{Success: "true", Guid: c.Query("guid"), Brand: c.Query("brand")}

		deletedCount := delKey("locks:" + data.Guid) // + ":" + data.Brand)
		if deletedCount == 0 {
			data.Success = "false"
			data.Locked = "false"
		} else {
			data.Locked = "true"
		}
		c.JSON(http.StatusOK, data)
	})

	// lastKey := "locks:b39bfc55-450c-4695-b63b-e15be698c377:gcl"
	// lastKeyExist := keyExists(lastKey)
	// fmt.Println("lastKey, last key exist ? ")
	// fmt.Println(lastKey)
	// fmt.Println(lastKeyExist)
	// fmt.Println(lastKeyExist.Val())
	// lastKey = "locks:e15be698c377:gcl"
	// lastKeyExist = keyExists(lastKey)
	// fmt.Println("lastKey, last key exist ? ")
	// fmt.Println(lastKey)
	// fmt.Println(lastKeyExist)
	// fmt.Println(lastKeyExist.Val())
	// fmt.Println("get key results..")
	// brand := getKey("locks:b39bfc55-450c-4695-b63b-e15be698c377")
	// fmt.Println("get key brand = ", brand)
	// getKey(lastKey)

	// val, err := rdb.Get("locks:b39bfc55-450c-4695-b63b-e15be698c377:gcl").Result()
	// fmt.Println("getKey err = ")
	// fmt.Println(err)
	// fmt.Println("getKey value = ")
	// fmt.Println(val)

	// setKey("locks:b39bfc55-450c-4695-b63b-e15be698c377", "snnow")
	// getSetKey("locks:b39bfc55-450c-4695-b63b-e15be698c377", "snnow")
	// if lastKeyExist {
	// 	fmt.Println("key exists!")
	// }

	// getSetRes := rdb.GetSet("locks:b39bfc55-450c-4695-b63b-e15be698c377", "gcl")
	// fmt.Println(getSetRes)
	// resGet2 := rdb.Get(lastKey)
	// fmt.Println("resGet2: ")
	// fmt.Println(resGet2)
	// fmt.Println(resGet2.Val())

	// router.GET("/test", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, homePage, nil)
	// })

	router.Run(":" + port)
}
