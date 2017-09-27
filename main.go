package main

import (
    "log"
)

func main() {
    log.SetFlags(log.Llongfile)

    var resume ResumeData

    err := ExtractJsonData(&resume)
    if err != nil {
        log.Fatal(err)
    }

    err = JsonToPdf(&resume)
    if err != nil {
        log.Fatal(err)
    }

    // resume.Print()
}

