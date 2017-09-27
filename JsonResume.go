package main

import (
    "log"
    "fmt"
    "os"
    "encoding/json"
    "strings"
    "io"
)

func (r *ResumeData) Print() {
    fmt.Printf("Name:        %s\nLabel:       %s\nAddress:     %s\nPostalCode:  %s\nCity:        %s\nCountryCode: %s\nRegion:      %s\nPicture:     %s\nEmail:       %s\nPhone:       %s\nWebsite:     %s\n",
        r.Basics.Name,
        r.Basics.Label,
        r.Basics.Location.Address,
        r.Basics.Location.PostalCode,
        r.Basics.Location.City,
        r.Basics.Location.CountryCode,
        r.Basics.Location.Region,
        r.Basics.Picture,
        r.Basics.Email,
        r.Basics.Phone,
        r.Basics.Website)

    fmt.Printf("Summary: \n\t")
    if len(r.Basics.Summary) > 80 {
        rdr := strings.NewReader(r.Basics.Summary)
        var err error
        var b byte
        for i := 0; err != io.EOF; i++ {
            b, err = rdr.ReadByte()
            fmt.Printf("%c", b)
            if i + 10 >= 80 && b == ' ' {
                i = 0
                fmt.Printf("\n\t")
            }
        }
        fmt.Printf("\n")
    } else {
        fmt.Printf("%s\n", r.Basics.Summary)
    }

    for _, o := range(r.Basics.Profiles) {
        fmt.Printf("Network: %s\nUsername: %s\nUrl: %s\n",
        o.Network,
        o.Username,
        o.Url)
    }

    for _, o := range(r.Work) {
        fmt.Printf("Company: %s\nPosition: %s\nWebsite: %s\nStartDate: %s\nEndDate: %s\nSummary: %s\n",
            o.Company,
            o.Position,
            o.Website,
            o.StartDate,
            o.EndDate,
            o.Summary)

        if len(o.Highlights) > 0 {
            fmt.Printf("Highlights:\n")
            for i, oo := range(o.Highlights) {
                fmt.Printf("%d: %s\n", i, oo)
            }
        }
    }

    for _, s := range(r.Skills) {
        fmt.Printf("%s\n", s.Name)
        fmt.Printf("%s\n", s.Level)
        for _, kw := range(s.Keywords) {
            fmt.Printf("Keyword: %s\n", kw)
        }
    }

    log.Println(r.Basics.Name)

}

type ResumeData struct {
    Basics struct {
        Name    string `json:"name"`
        Label   string `json:"label"`
        Picture string `json:"picture"`
        Email   string `json:"email"`   // format "email"
        Phone   string `json:"phone"`
        Website string `json:"website"` // format "uri"
        Summary string `json:"summary"`

        Location struct {
            Address     string `json:"address"`
            PostalCode  string `json:"postalCode"`
            City        string `json:"city"`
            CountryCode string `json:"countryCode"`
            Region      string `json:"region"`
        } `json:"location"`

        Profiles []struct {
            Network  string `json:"network"`
            Username string `json:"username"`
            Url      string `json:"url"`
        } `json:"profiles"`

    } `json:"basics"`

    Work []struct {
        Company    string `json:"company"`
        Position   string `json:"position"`
        Website    string `json:"website"`   // format "uri"
        StartDate  string `json:"startDate"` // format "date"
        EndDate    string `json:"endDate"`   // format "date"
        Summary    string `json:"summary"`
        Highlights []string `json:"highlights"`
    } `json:"work"`

    Volunteer []struct {
        Organization string `json:"organization"`
        Position     string `json:"position"`
        Website      string `json:"website"`   // format "uri"
        StartDate    string `json:"startDate"` // format "date"
        EndDate      string `json:"endDate"`   // format "date"
        Summary      string `json:"summary"`
        Highlights   []string `json:"highlights"`
    } `json:"volunteer"`

    Education []struct {
        Institution string `json:"institution"`
        Area        string `json:"area"`
        StudyType   string `json:"studyType"`
        StartDate   string `json:"startDate"` // format "date"
        EndDate     string `json:"endDate"`   // format "date"
        Gpa         string `json:"gpa"`
        Courses     []string `json:"courses"`
    } `json:"education"`

    Awards []struct {
        Title   string `json:"title"`
        Date    string `json:"date"` // format "date"
        Awarder string `json:"awarder"`
        Summary string `json:"summary"`
    } `json:"awards"`

    Publications []struct {
        Name        string `json:"name"`
        Publisher   string `json:"publisher"`
        ReleaseDate string `json:"releaseDate"`
        Website     string `json:"website"`
        Summary     string `json:"summary"`
    } `json:"publications"`

    Skills []struct {
        Name     string `json:"name"`
        Level    string `json:"level"`
        Keywords []string `json:"keywords"`
    } `json:"skills"`

    Languages []struct {
        Language string `json:"language"`
        Fluency  string `json:"fluency"`
    } `json:"languages"`

    Interests []struct {
        Name     string `json:"name"`
        Keywords []string `json:"keywords"`
    } `json:"interests"`

    References []struct {
        Name      string `json:"name"`
        Reference string `json:"reference"`
    } `json:"references"`
}


func (r ResumeData)CharCount() int {
    total := 0
    work_total := 0
    edu_total := 0
    skills_total := 0
    refs_total := 0

    total += len(r.Basics.Name)
    total += len(r.Basics.Location.City)
    total += len(r.Basics.Phone)
    total += len(r.Basics.Email)

    // work
    for _, t := range(r.Work) {
        work_total += len(t.Company) + len(t.Position) + len(t.StartDate) + len(t.EndDate) + len(t.Summary)
        for _, hl := range(t.Highlights) {
            work_total += len(hl)
        }
    }

    // education
    for _, t := range(r.Education) {
        edu_total += len(t.Institution) + len(t.Area) + len(t.StudyType) + len(t.StartDate) + len(t.EndDate)
        for _, co := range(t.Courses) {
            edu_total += len(co)
        }
    }

    // skills
    for _, t := range(r.Skills) {
        skills_total += len(t.Name) + len(t.Level)
        for _, kw := range(t.Keywords) {
            skills_total += len(kw)
        }
    }
    // references
    for _, t := range(r.References) {
        refs_total += len(t.Name) + len(t.Reference)
    }

    total += work_total + edu_total + skills_total + refs_total

    return total
}


func ExtractJsonData(resume *ResumeData) error {
    file, err := os.Open("data/resume.json")
    if err != nil {
        return err
        // log.Fatal(err)
    }
    dec := json.NewDecoder(file)
    err = dec.Decode(resume)
    if err != nil {
        return err // log.Fatal(err)
    }

    err = file.Close()

    return err
}
