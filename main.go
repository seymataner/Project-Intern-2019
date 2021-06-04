package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strings"
)

type Data struct {
	Id        int
	JobId     string `json:"JobId"`
	JobName   string `json:"JobName"`
	JobType   string `json:"JobType"`
	JobParams string `json:"JobParams"`
	Date      string `json:"Date"`
}

type Content struct {
	Status  int
	Count   int
	Payload []map[string]map[string]string
}

func handler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "togepi:1@tcp(127.0.0.1:3306)/redisdb")
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("SELECT * FROM gearman_jobs")
	if err != nil {
		panic(err)
	}

	returnContent := Content{}
	returnContent.Status = 200
	count := 0
	for rows.Next() {
		var data Data
		err = rows.Scan(&data.Id, &data.JobId, &data.JobName, &data.JobType, &data.JobParams, &data.Date)
		innerMap := make(map[string]map[string]string)
		innnerMap, ok := innerMap[data.JobName]
		if !ok {
			innnerMap = make(map[string]string)
			innerMap[data.JobName] = innnerMap
		}
		innerMap[data.JobName]["jobId"] = data.JobId
		innerMap[data.JobName]["jobType"] = data.JobType
		innerMap[data.JobName]["JobParams"] = data.JobParams
		innerMap[data.JobName]["Date"] = data.Date
		returnContent.Payload = append(returnContent.Payload, innerMap)
		count++
	}
	returnContent.Count = count
	jsonContent, err := json.Marshal(returnContent)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonContent)
}

func Connect() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	return client
}

func main() {
	http.HandleFunc("/getjobs", handler)
	go func() {
		log.Fatal(http.ListenAndServe(":80", nil))
	}()

	db, err := sql.Open("mysql", "togepi:1@tcp(127.0.0.1:3306)/redisdb")
	if err != nil {
		panic(err)
	}

	fmt.Println("Succesfully Connected")
	defer db.Close()

	client := Connect()

	pubsub := client.PSubscribe("*")
	defer pubsub.Close()

	for {
		dt := Data{}
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			panic(err)
		}
		s := strings.Split(msg.Channel, ":")
		dt.JobType = s[1]
		t := strings.Split(msg.Payload, "-")
		dt.JobName = t[1]
		dt.JobId = t[2]
		result := client.HGetAll(msg.Payload)
		m := result.Val()
		dt.JobParams = m["data"]
		stmt, err := db.Prepare("INSERT INTO gearman_jobs(JobId,JobName,JobType,JobParams) VALUES(?,?,?,?)")
		if err != nil {
			panic(err)
		}

		stmt.Exec(dt.JobId, dt.JobName, dt.JobType, dt.JobParams)
	}
}
