package main

import (
    "fmt"
    "github.com/jung-kurt/gofpdf"
    "log"
)

var OUTFILE = "pdf/example1.pdf"

// fonts to try: sans-serif, monospace

var (
    FontH1 float64
    FontH2 float64
    FontH3 float64
    FontH4 float64
    FontH5 float64
    FontP  float64

    pageLimit   float64 // overflow control
    pageUsed    float64 // overflow control

    x, y        float64 // global coords

    lineHt      float64 // cord intervals
    lineWd      float64

    fontSz      float64 // font config
    fontStr     string


    fill        bool // cell config
    borderStr   string
    alignStr    string

    page *PageInfo                    // page stats variables
    doc *gofpdf.Fpdf
    resume *ResumeData

    AboutH float64
    WorkH float64
    EduH float64
    SkillsH float64
    IntH float64
)

func initGlobals() {
    FontH1 = 36.0
    FontH2 = 30.0
    FontH3 = 24.0
    FontH4 = 18.0
    FontH5 = 14.0
    FontP  = 10.0

    pageLimit   = 0.95      // overflow control
    pageUsed    = 0.00      // overflow control

    lineHt      = FontP     // cord intervals
    lineWd      = FontP / 2

    fontSz      = FontP      // font config
    fontStr     = "robocon"


    fill        = false         // cell config
    borderStr   = ""
    alignStr    = ""
    doc.SetFont(fontStr, "", fontSz);
}


func calibrateGlobals() {
    FontH1 *= 0.9
    FontH2 *= 0.9
    FontH3 *= 0.9
    FontH4 *= 0.9
    FontH5 *= 0.9
    FontP  *= 0.9

    if (FontP <= 7.0) {
        log.Fatal("Could not calibrate\n")
    }

    lineHt = FontP
    lineWd = FontP / 2
    fontSz = FontP
    fontStr = "robocon"
    doc.SetFont(fontStr, "", fontSz);
}


func initDocSetup() *gofpdf.Fpdf {
    doc := gofpdf.New("P", "pt", "letter", "")

    // import fonts
    doc.SetFontLocation("font")
    doc.AddFont("libserif", "", "LiberationSerif-Regular.json")
    doc.AddFont("libserif", "B", "LiberationSerif-Bold.json")
    doc.AddFont("robocon", "", "RobotoCondensed-Regular.json")
    doc.AddFont("robocon", "B", "RobotoCondensed-Bold.json")

    doc.AddPage()
    doc.SetDrawColor(224, 224, 224)
    doc.SetLineWidth(1.0)

    return doc
}

func estimateHeight() float64 {
    // startY := page.mt
    curY := page.mt

    curY += FontH1
    curY += FontP * 2 // addr, email
    curY += FontP * 2 // break
    AboutH = curY
    curY += FontP * float64(len(doc.SplitLines([]byte(resume.Basics.Summary), page.w9 - page.mr)))
    curY += lineHt * 3 // break

    // work
    WorkH = curY
    fmt.Printf("Estimated curY: %f\n", curY)
    for _, w := range(resume.Work) {
        curY += FontH4
        curY += lineHt
        curY += lineHt

        l := doc.GetStringWidth(w.Summary)
        if l + page.w3 > page.w {
            curY += (lineHt * float64(len(doc.SplitLines([]byte(w.Summary), page.w9 - page.mr)))) + lineHt
        } else {
            curY += lineHt
        }
        curY += FontH4
    }
    curY -= FontH4
    curY += lineHt * 3


    EduH = curY
    fmt.Printf("Estimated curY: %f\n", curY)
    for _, e := range(resume.Education) {
        curY += FontH4
        curY += lineHt + lineHt

        numc := len(e.Courses)
        if numc > 2 {
            topy := curY;
            var i int
            for i = 0; i < numc / 2; i++ {
                curY += lineHt
            }
            curY = topy
            for i < numc {
                curY += lineHt
                i++
            }
        } else {
            curY += lineHt * float64(numc)
        }
        curY += lineHt + lineHt
    }
    curY -= lineHt + lineHt
    curY += lineHt * 3

    SkillsH = curY
    fmt.Printf("Estimated curY: %f\n", curY)
    curY += FontH4
    var max int = 0
    for _, s := range(resume.Skills) {
        l := len(s.Keywords)
        if l > max {
            max = l
        }
    }
    curY += lineHt * float64(max)
    curY += lineHt * 3

    fmt.Printf("Estimated curY: %f\n", curY)

    IntH = curY
    curY += FontH4
    max = 0
    for _, s := range(resume.Interests) {
        l := len(s.Keywords)
        if l > max {
            max = l
        }
    }
    curY += lineHt * float64(max)

    fmt.Printf("Estimated Height: %f\n", curY)

    return curY
}

func insertHeader() {
    // Name
    doc.SetFontSize(FontH1)
    doc.Write(0, resume.Basics.Name)
    doc.Ln(FontH1)

    // Address
    doc.SetFontSize(FontP)
    doc.Write(0, resume.Basics.Location.Address + " " + resume.Basics.Location.City + " " + resume.Basics.Location.PostalCode)

    // Phone
    doc.SetX(page.w6)
    doc.SetFont(fontStr, "B", FontP); doc.Write(0, "Phone: ")
    doc.SetFont(fontStr,  "", FontP); doc.Write(0, resume.Basics.Phone)
    doc.Ln(FontP)

    // Email
    doc.SetFont(fontStr, "B", FontP); doc.Write(0, "Email: ")
    doc.SetFont(fontStr,  "", FontP); doc.Write(0, resume.Basics.Email)

    // Website
    doc.SetX(page.w6)
    doc.SetFont(fontStr, "B", FontP); doc.Write(0, "Website: ")
    doc.SetFont(fontStr,  "", FontP); doc.Write(0, resume.Basics.Website)

    // doc.Ln(FontP)
    // doc.Ln(FontP)
}

func insertAbout() {
    // "About"
    doc.SetY(AboutH)
    doc.SetFontSize(FontH3)
    doc.Write(FontP * 2, "About")
    // Summary
    doc.SetFontSize(FontP)
    getPos()
    x = page.w3
    setPos()
    doc.MultiCell(page.w9 - page.mr, FontP, resume.Basics.Summary, borderStr, alignStr, fill)
    y += FontP * float64(len(doc.SplitLines([]byte(resume.Basics.Summary), page.w9 - page.mr)))
    setPos()
}

func insertWork() {
    topy := doc.GetY()

    doc.SetFontSize(FontH3)
    doc.Write(0, "Work"); doc.SetY(topy + FontH3)
    doc.Write(0, "Experience");

    doc.SetXY(page.w3, topy)
    getPos()

    for _, w := range(resume.Work) {
        setPos()
        doc.SetFontSize(FontH4) // add queue
        doc.Write(0, w.Company)
        y += FontH4
        x = page.w3

        setPos()
        doc.SetFontSize(FontP) // add quq
        doc.Write(0, w.Position)
        y += lineHt
        doc.SetXY(page.w3, y)
        doc.Write(0, w.StartDate + " " + w.EndDate) // add queue
        y += lineHt
        doc.SetXY(page.w3, y)


        doc.SetFontSize(FontP)
        l := doc.GetStringWidth(w.Summary)
        if l + page.w3 > page.w {
            doc.MultiCell(page.w9 - page.mr, FontP, w.Summary, borderStr, alignStr, fill)
            y += (lineHt * float64(len(doc.SplitLines([]byte(w.Summary), page.w9 - page.mr)))) + lineHt
        } else {
            y += lineHt
            doc.SetXY(page.w3, y)
            doc.Write(0, w.Summary)
        }
        y += FontH4
        x = page.w3
    }
    y -= FontH4
    setPos()
}

func insertEdu() {
    getPos()

    doc.SetFontSize(FontH3)
    doc.Write(0, "Education")
    x = page.w3
    doc.SetXY(page.w3, y)
    for _, e := range(resume.Education) {
        x = page.w3
        setPos()
        doc.SetFontSize(FontH4)
        doc.Write(0, e.Institution)
        doc.SetFontSize(FontH5)
        doc.Write(1.8, "    " + e.Area)
        y += FontH4
        doc.SetXY(page.w3, y)

        doc.SetFontSize(FontP)
        doc.Write(0, e.StartDate + " - " + e.EndDate)
        y += lineHt + lineHt
        x = page.w3 + lineWd
        setPos()
        doc.SetFontSize(FontP)

        numc := len(e.Courses)
        if numc > 2 {
            topy := y;
            var i int
            for i = 0; i < numc / 2; i++ {
                setPos()
                doc.Write(0, " - " + e.Courses[i])
                y += lineHt
            }
            x = page.w6 + (page.w3 / 2) + lineWd
            y = topy
            for i < numc {
                setPos()
                doc.Write(0, " - " + e.Courses[i])
                y += lineHt
                i++
            }
        } else {
            x = page.w3 + lineWd
            for _, c := range(e.Courses) {
                setPos()
                doc.Write(0, " - " + c)
                y += lineHt
            }
        }
        y += lineHt + lineHt
        setPos()
    }
    y -= lineHt + lineHt
}

func insertGeneric(topicName string, g []GenericRecord) {
    getPos()
    strWidths := GetWidths(g)

    // draw loop
    var top float64
    var bottom float64

    doc.SetFontSize(FontH3)
    doc.Write(0, topicName)
    doc.SetXY(x, y)

    x = page.w3
    top = doc.GetY()
    for i, s := range(g) {
        y = top
        doc.SetXY(x, y)
        doc.SetFontSize(FontH4)
        doc.Write(0, s.Name)
        y += FontH4
        x += FontP
        doc.SetFontSize(FontP)
        for _, kw := range(s.Keywords) {
            doc.SetXY(x, y)
            doc.Write(0, " - " + kw)
            y += FontP
        }
        if y > bottom {
            bottom = y + FontP
        }
        x -= lineHt
        y -= lineHt
        x += strWidths[i]
    }
    y = bottom - lineHt
    x = page.w3
}

func insertLineDivide() {
    y += lineHt
    setPos()
    doc.Line(page.w3, y, page.w, y)
    y += lineHt + lineHt;
    x = page.ml
    setPos()
}

// funcs
func getPos() {
    x, y = doc.GetXY()
}

func setPos() {
    doc.SetXY(x, y)
}

/**
 * TODO: Fit all context on one page, or expand
 *       Format dates
 *       Font color/shades/family
 *       Refactor: datastructure to separate calculations from draws,
 *          group similar jobs together
 *          make more portable and reusable
 */
func JsonToPdf(resumeData *ResumeData) error {
    // topics := queue.New(5)

    // init
    doc = initDocSetup()
    page = NewPageInfo(doc)
    resume = resumeData

    initGlobals()

    for estimateHeight() > page.h {
        calibrateGlobals()
    }

    Break := func() {
        insertLineDivide()
        getPos()
    }

    Header := func() {
        insertHeader()

        // Dividing line
        // doc.Line(page.ml, y, page.w, y)
        doc.Line(page.ml, AboutH - lineHt, page.w, AboutH - lineHt)
        doc.Ln(FontP)
        getPos()
        pageUsed = page.PageUsed(y)


        insertAbout()
    }

    Skills := func() {
        g := make([]GenericRecord, len(resume.Skills))
        for i := 0; i < len(resume.Skills); i++ {
            g[i].Name = resume.Skills[i].Name
            g[i].Keywords = resume.Skills[i].Keywords
        }
        insertGeneric("Skills", g)
        setPos()
    }

    Interests := func() {
        g := make([]GenericRecord, len(resume.Interests))
        for i := 0; i < len(resume.Interests); i++ {
            g[i].Name = resume.Interests[i].Name
            g[i].Keywords = resume.Interests[i].Keywords
        }
        insertGeneric("Interests", g)
        setPos()
    }

    Header()
    procs := []func(){insertWork, insertEdu, Skills, Interests}

    var pid int = 0
    var numprocs int = len(procs)
    for pageUsed < pageLimit && pid < numprocs {
        Break()
        getPos()
        fmt.Println("CurY: ", y)
        procs[pid]()
        pageUsed = page.PageUsed(y)
        pid++
    }

    getPos()
    fmt.Println("CurY: ", y)
    used := page.PageUsed(y)
    fmt.Println("Percent page used: ", used)

    err := doc.OutputFileAndClose(OUTFILE)
    if err != nil {
        return err
    }
    fmt.Println("Successfully created file ", OUTFILE)

    return err
}

// Page Info
type PageInfo struct {
    aw, ah         float64 // absolute width, hight
    w, h           float64 // usable width, height (without margin space)
    mt, mb, ml, mr float64 // margins top, bottom, left, right
    w3, w6, w9     float64 // width of 1/4 of page, with of 1/2 of page, width of 3/4 of page
}

func NewPageInfo(doc *gofpdf.Fpdf) *PageInfo {
    page := new(PageInfo)
    page.ml, page.mt, page.mr, page.mb = doc.GetMargins()
    page.aw, page.ah = doc.GetPageSize()
    // page.w = page.aw - page.ml - page.mr
    // page.h = page.ah - page.mt - page.mb
    page.w = page.aw - page.mr
    page.h = page.ah - page.mb
    page.w3 = page.aw / 4
    page.w6 = page.aw / 2
    page.w9 = page.w3 * 3
    return page
}

func (p *PageInfo) PrintInfo() {
    fmt.Printf("\tPage Absolute w, h : %f, %f\n", p.aw, p.ah)
    fmt.Printf("\tPage w, h          : %f, %f\n", p.w, p.h)
    fmt.Printf("\tMargin top         : %f\n", p.mt)
    fmt.Printf("\tMargin bottom      : %f\n", p.mb)
    fmt.Printf("\tMargin left        : %f\n", p.ml)
    fmt.Printf("\tMargin right       : %f\n", p.mr)
    fmt.Printf("\tPageSize x0.25     : %f\n", p.w3)
    fmt.Printf("\tPageSize x0.50     : %f\n", p.w6)
    fmt.Printf("\tPageSize x0.75     : %f\n\n", p.w9)
}

func (p *PageInfo) PageUsed(y float64) float64 {
    return y / p.h
}

type GenericRecord struct {
    Name string
    Keywords []string
}

func GetWidths(g []GenericRecord) []float64 {
    strWidths := make([]float64, 10) //[10]float64{0.0}
    var sum float64 = 0.0
    num := len(g)
    for i, s := range(g) {
        doc.SetFontSize(FontH4)
        strWidths[i] = doc.GetStringWidth(s.Name)
        doc.SetFontSize(FontP)
        for _, kw := range(s.Keywords) {
            l := doc.GetStringWidth(kw)
            if l > strWidths[i] {
                strWidths[i] = l
            }
        }
        sum += strWidths[i]
    }

    rng := page.w9 - page.ml
    if sum >= rng {
        log.Fatal("Too many words!")
    }
    clr := rng - sum
    var a, b int
    for a = 0; a < num / 2; a++ {
        b = num - 2 - a
        sum -= strWidths[a] + strWidths[b]
        gap := (clr - sum) / 2
        strWidths[a] += gap
        if b > a {
            strWidths[b] += gap
        }
    }
    return strWidths
}
