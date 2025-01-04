package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestDownloadAndSendFileByRange(t *testing.T) {
    err := DebugDownloadRequest("https://rr5---sn-apn7en7l.googlevideo.com/videoplayback?expire=1721841110&ei=duGgZtqvAciJp-oPnZ2NgQQ&ip=41.211.141.37&id=o-ACJnvvbhezKb559ySHEGNkHvtaDkqcshDSyOHku42mNs&itag=18&source=youtube&requiressl=yes&xpc=EgVo2aDSNQ%3D%3D&bui=AXc671JZfTQ8ML1CEBRSJN5j_u9pt9U_lLeR7HWyqMFFGIuKCojG66kvGWgyIv8DZ9mr5tWv8vj038Bg&spc=NO7bASsvOgU1QaIvuEufJ6RVF1o2dpP80auOe-I8Oy0-upZbzhi7_JFCkgMpGIc&vprv=1&svpuc=1&mime=video%2Fmp4&ns=MRK6-A14aOvhuRwZbaj5EucQ&rqh=1&cnr=14&ratebypass=yes&dur=115.403&lmt=1719571555513089&c=WEB&sefc=1&txp=4538434&n=ftCC4_Kk1M83LA&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cxpc%2Cbui%2Cspc%2Cvprv%2Csvpuc%2Cmime%2Cns%2Crqh%2Ccnr%2Cratebypass%2Cdur%2Clmt&sig=AJfQdSswRgIhAJzypkOb-UTHmgOH52zyI2DHzWWj2Scl0rxP3aLb3cMMAiEAgE2Hamc3QgImISGoYCdcRvs-_iilNFB4V5b61So8GbQ%3D&title=hello%20world&rm=sn-b5cvouxaxjvh-ac8l7l,sn-apne76&rrc=79,104&fexp=24350516,24350517&req_id=fb64b5b74978a3ee&redirect_counter=2&cms_redirect=yes&cmsv=e&ipbypass=yes&mh=I3&mip=196.77.37.178&mm=29&mn=sn-apn7en7l&ms=rdu&mt=1721819118&mv=m&mvi=5&pl=21&lsparams=ipbypass,mh,mip,mm,mn,ms,mv,mvi,pl&lsig=AGtxev0wRQIgAjjjLpqmycJHTcMnuQxl7yyOjzIhoicIZYKpmdJpVr4CIQDC0rQiBqiO_BjhUIHMn5lmkNbS3RiwLrKjMPSTLFVGtA%3D%3D", 1024*1024, 1024*1024)
    if err != nil {
        t.Errorf("DebugDownloadRequest failed: %v", err)
    }
}

func DebugDownloadRequest(mediaUrl string, contentSize int, chunkSize int) error {
    // Initialize HTTP client
    client := &http.Client{}

    // Calculate number of chunks
    numChunks := (contentSize + chunkSize - 1) / chunkSize

    // Print initial debug information
    fmt.Printf("Debugging download request for URL: %s\n", mediaUrl)
    fmt.Printf("Content Size: %d, Chunk Size: %d, Number of Chunks: %d\n", contentSize, chunkSize, numChunks)

    // Create output file
    outputFile, err := os.Create("test.mp4")
    if err != nil {
        return fmt.Errorf("error creating output file: %v", err)
    }
    defer outputFile.Close()

    // Iterate through chunks
    for i := 0; i < numChunks; i++ {
        // Calculate byte range for this chunk
        startRange := i * chunkSize
        endRange := startRange + chunkSize - 1
        if endRange >= contentSize {
            endRange = contentSize - 1
        }

        // Create HTTP GET request
        req, err := http.NewRequest("GET", mediaUrl, nil)
        if err != nil {
            return fmt.Errorf("error creating request: %v", err)
        }

        // Set range header
        req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", startRange, endRange))

        // Print request details
        fmt.Printf("\n--- Request %d ---\n", i+1)
        fmt.Printf("URL: %s\n", req.URL)
        fmt.Printf("Method: %s\n", req.Method)
        fmt.Printf("Headers:\n")
        for key, values := range req.Header {
            for _, value := range values {
                fmt.Printf("  %s: %s\n", key, value)
            }
        }

        // Perform the HTTP request
        resp, err := client.Do(req)
        if err != nil {
            return fmt.Errorf("error performing request: %v", err)
        }
        defer resp.Body.Close()

        // Print response details
        fmt.Printf("\n--- Response %d ---\n", i+1)
        fmt.Printf("Status: %s\n", resp.Status)
        fmt.Printf("Status Code: %d\n", resp.StatusCode)
        fmt.Printf("Content-Length: %d\n", resp.ContentLength)
        fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))
        fmt.Printf("Headers:\n")
        for key, values := range resp.Header {
            for _, value := range values {
                fmt.Printf("  %s: %s\n", key, value)
            }
        }

        // Read the body and write to file
        n, err := io.Copy(outputFile, resp.Body)
        if err != nil {
            return fmt.Errorf("error writing to file: %v", err)
        }
        fmt.Printf("Wrote %d bytes to file\n", n)
    }

    fmt.Println("\nDebug request completed. Check test.mp4 for the downloaded content.")
    return nil
}