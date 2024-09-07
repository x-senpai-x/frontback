package main
import (
    "bytes"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
)   
var (
    classHash  string
    jsonRpcUrl string
)
func init() {
    flag.StringVar(&classHash, "classHash", "", "The contract class hash")
    flag.StringVar(&jsonRpcUrl, "jsonRpcUrl", "", "The JSON-RPC URL")
}
func enableCors(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
    (*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
func GetABIHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request to fetch ABI")
    // Enable CORS for the response
    enableCors(&w)
    // Read the request body
    var requestData map[string]string
    err := json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil {
        log.Printf("Error decoding request body: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    // Log the decoded data
    log.Printf("Decoded request data: %v", requestData)
    classHash := requestData["classHash"]
    jsonRpcUrl := requestData["jsonRpcUrl"]
    log.Printf("classHash: %s, jsonRpcUrl: %s", classHash, jsonRpcUrl)
    // Call the ABI fetching function
    abi, err := GetStarknetABI(classHash, jsonRpcUrl)
    if err != nil {
        log.Printf("Error fetching ABI: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // Log ABI data fetched
    log.Println("Successfully fetched ABI")
    // Write the ABI JSON to the response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(abi)
}
func GetStarknetABI(classHash, jsonRpcUrl string) ([]byte, error) {
    // Log the command being run
    log.Printf("Running CLI command to fetch ABI with classHash: %s, jsonRpcUrl: %s", classHash, jsonRpcUrl)
    // Set up the command to run the Go file with the appropriate arguments
    cmd := exec.Command("go", "run", "cli/get/starknet.go", "--classHash", classHash, "--jsonRpcUrl", jsonRpcUrl, "--output", "abi.json")
    // Capture the output and error
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    // Run the command and log any errors
    err := cmd.Run()
    if err != nil {
        log.Printf("Error running CLI tool: %v, stderr: %s", err, stderr.String())
        return nil, fmt.Errorf("error running CLI tool: %v, stderr: %s", err, stderr.String())
    }
    // Log the success of command execution
    log.Println("CLI command executed successfully")
    // Read ABI JSON from file and log the file reading step
    abiJson, err := os.ReadFile("abi.json")
    if err != nil {
        log.Printf("Error reading ABI file: %v", err)
        return nil, fmt.Errorf("error reading ABI file: %v", err)
    }
    // Log the successful reading of ABI file
    log.Println("ABI file read successfully")
    return abiJson, nil
}
func GetBackfillHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request to backfill block data")
    enableCors(&w)

    var requestData map[string]interface{}
    err := json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil {
        log.Printf("Error decoding request body: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    fromBlock := uint64(requestData["fromBlock"].(float64))
    toBlock := uint64(requestData["toBlock"].(float64))
    rpcUrl := requestData["rpcUrl"].(string)
    outputFile := requestData["outputFile"].(string)
    transactionHashFlag := requestData["transactionHashFlag"].(bool)

    // Get the absolute path of the output file
    absOutputPath, err := filepath.Abs(outputFile)
    if err != nil {
        log.Printf("Error getting absolute path: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    err = RunBackfillCommand(fromBlock, toBlock, rpcUrl, absOutputPath, transactionHashFlag)
    if err != nil {
        log.Printf("Error running backfill: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    log.Printf("Backfill command: go run cli/backfill/starknet.go --from %d --to %d --rpc-url %s --output %s", fromBlock, toBlock, rpcUrl, outputFile)
    // Check if the file was actually created
    fileInfo, err := os.Stat(absOutputPath)
    if err != nil {
        if os.IsNotExist(err) {
            log.Printf("Output file was not created: %v", err)
            http.Error(w, fmt.Sprintf("Backfill command executed but output file was not created at %s", absOutputPath), http.StatusInternalServerError)
        } else {
            log.Printf("Error getting file info: %v", err)
            http.Error(w, fmt.Sprintf("Error verifying output file: %v", err), http.StatusInternalServerError)
        }
        return
    }

    log.Printf("Backfill successful. Output file created at: %s, Size: %d bytes", absOutputPath, fileInfo.Size())
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Backfill successful",
        "outputFile": absOutputPath,
        "fileSize": fileInfo.Size(),
    })
}

func RunBackfillCommand(fromBlock, toBlock uint64, rpcUrl, outputFile string, transactionHashFlag bool) error {
    args := []string{
        "run", "cli/backfill/starknet.go",
        "--from", fmt.Sprintf("%d", fromBlock),
        "--to", fmt.Sprintf("%d", toBlock),
        "--rpc-url", rpcUrl,
        "--output", outputFile,
    }
    
    if transactionHashFlag {
        args = append(args, "--transactionhash")
    }
    cmd := exec.Command("go", args...)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    err:=cmd.Run()
    if err != nil {
        log.Printf("Error running backfill command: %v", err)
        log.Printf("Stderr output: %s", stderr.String())
        return fmt.Errorf("command failed: %v, stderr: %s", err, stderr.String())
    }
    log.Println("Backfill command executed successfully")
    log.Printf("Command output: %s", out.String())

    return nil
}
/*
func GetBackfillHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request to backfill block data")
    enableCors(&w)
    var requestData map[string]interface{}
    err := json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil {
        log.Printf("Error decoding request body: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    fromBlock := uint64(requestData["fromBlock"].(float64))
    toBlock := uint64(requestData["toBlock"].(float64))
    rpcUrl := requestData["rpcUrl"].(string)
    outputFile := requestData["outputFile"].(string)
    transactionHashFlag := requestData["transactionHashFlag"].(bool)
    err = RunBackfillCommand(fromBlock, toBlock, rpcUrl, outputFile, transactionHashFlag)
    if err != nil {
        log.Printf("Error running backfill: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Backfill successful"))
}*/
func main() {
    flag.Parse()
    if classHash != "" && jsonRpcUrl != "" {
        // Directly handle the ABI fetching logic when flags are provided
        abi, err := GetStarknetABI(classHash, jsonRpcUrl)
        if err != nil {
            log.Fatalf("Error fetching ABI: %v", err)
        }
        fmt.Println(string(abi))
        return
    }
    // Handle the CORS preflight request for /api/abi endpoint
    http.HandleFunc("/api/abi", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "OPTIONS" {
            // Enable CORS for OPTIONS preflight request
            enableCors(&w)
            return
        }
        // Proceed to handle the actual POST request
        GetABIHandler(w, r)
    })
    http.HandleFunc("/api/backfill", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "OPTIONS" {
            enableCors(&w)
            return
        }
        GetBackfillHandler(w, r)
    })
    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
