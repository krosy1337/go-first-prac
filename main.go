package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
  var errorsCount uint8
  client := &http.Client{}

  for {
    resp, err := client.Get("http://srv.msk01.gigacorp.local/_stats")
    if err != nil || resp.StatusCode != http.StatusOK {
      errorsCount++
      if errorsCount >= 3 {
        fmt.Println("Unable to fetch server statistic.")
        break
      }
      time.Sleep(1 * time.Second)
      continue
    }

    body, err := io.ReadAll(resp.Body)
    _ = resp.Body.Close()
    if err != nil {
      errorsCount++
      if errorsCount >= 3 {
        fmt.Println("Unable to fetch server statistic.")
        break
      }
      time.Sleep(1 * time.Second)
      continue
    }

    bodyString := strings.TrimSpace(string(body))
    bodyParts := strings.Split(bodyString, ",")
    if len(bodyParts) != 7 {
      errorsCount++
      if errorsCount >= 3 {
        fmt.Println("Unable to fetch server statistic.")
        break
      }
      time.Sleep(1 * time.Second)
      continue
    }

    bodyFloats := make([]float64, len(bodyParts))
    for i, part := range bodyParts {
      f, err := strconv.ParseFloat(part, 64)
      if err != nil {
        errorsCount++
        if errorsCount >= 3 {
          fmt.Println("Unable to fetch server statistic.")
          break
        }
        continue
      }
      bodyFloats[i] = f
    }

    if bodyFloats[0] >= 30 {
      fmt.Printf("Load Average is too high: %d\n", int(bodyFloats[0]))
    }

    memUsage := bodyFloats[2] / bodyFloats[1] * 100
    if memUsage >= 80 {
      fmt.Printf("Memory usage too high: %d%%\n", int(memUsage))
    }

    diskUsage := bodyFloats[4] / bodyFloats[3] * 100
    if diskUsage >= 90 {
      freeMb := int((bodyFloats[3] - bodyFloats[4]) / (1024 * 1024))
      fmt.Printf("Free disk space is too low: %d Mb left\n", freeMb)
    }

    netUsage := bodyFloats[6] / bodyFloats[5] * 100
    if netUsage >= 90 {
      freeMbit := int(bodyFloats[5]-bodyFloats[6]) / (1000 * 1000)
      fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", freeMbit)
    }

    time.Sleep(time.Second)
  }
}
