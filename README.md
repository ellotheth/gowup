# gowup
### A Go client for the [Where's it Up API](https://api.wheresitup.com)

(Or more accurately, ellotheth learns to golang.)

#### Usage

```{.go}
package main

import (
    "github.com/ellotheth/gowup"
    "fmt"
)

func main() {
    api := gowup.WIU{
        Client: "<your WIU client ID>",
        Token:  "<your WIU client token>",
    }

    // get available source locations
    if locations, err := api.Locations(); err != nil {
        fmt.Println(err)
        return
    } else {
        fmt.Printf("%+v\n", locations)
    }

    // get a summary of your recent jobs
    if jobs, err := api.Jobs(); err != nil {
        fmt.Println(err)
        return
    } else {
        fmt.Printf("%+v\n", jobs)
    }

    // get details for one job
    if job, err := api.Job("<WIU job ID>"); err != nil {
        fmt.Println(err)
        return
    } else {
        fmt.Printf("%+v\n", job)
    }

    // submit a new job (in this case, pings from Denver to Google)
    id, err := api.Submit(&gowup.JobRequest{
        Url:       "https://google.com",
        Tests:     []string{"ping"},
        Locations: []string{"denver"},
    })
    if err != nil {
        fmt.Println(err)
        return
    } else {
        fmt.Printf("%+v\n", id)
    }
}
```
