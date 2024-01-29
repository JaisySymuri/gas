package gas

import (
	appd "appdynamics"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
)

func ClientBTStarter(w http.ResponseWriter, r *http.Request, appdBTNameHDL string) appd.BtHandle {
	btHandle := appd.StartBT(appdBTNameHDL, "") // Use the provided BT name
	return btHandle
}

func ClientBTEnder(btHandle appd.BtHandle, r *http.Request, endpoint string, sqlStatement string) {

	// Set the snapshot URL if snapshotting is active
	currentURL := GetCurrentURL(r)
	appd.SetBTURL(btHandle, currentURL)

	backendName := BackendName // A unique name for the backend
	backendType := BackendType // Specifies that this is a backend
	backendProperties := map[string]string{
		"HOST": Host,     // Hostname or IP address of the  server
		"PORT": Port,     // Port number on which the server is running
		"URL":  endpoint, // Name of the endpoint
	}
	resolveBackend := false // Set to false to use static backend configuration

	appd.AddBackend(backendName, backendType, backendProperties, resolveBackend)

	// Generate a new UUID
	my_ec_guid := uuid.New().String()

	// Convert btHandle to the correct type (appdynamics.BtHandle)
	businessTransactionHandle := appd.BtHandle(btHandle)

	ecHandle := appd.StartExitcall(businessTransactionHandle, backendName)
	hdr := appd.GetExitcallCorrelationHeader(ecHandle)

	req, err := CreateRequest(endpoint, "POST", nil)
	req.Header.Add(appd.APPD_CORRELATION_HEADER_NAME, hdr)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected response status code:", resp.StatusCode)
		return
	}

	// Read and process the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Process the response body
	if len(responseBody) > 30 {
		truncated := string(responseBody[:30])
		fmt.Println("Response from downstream service:", truncated+"...")
	} else {
		fmt.Println("Response from downstream service:", string(responseBody))
	}

	// Optionally store the handle in the global registry
	appd.StoreExitcall(ecHandle, my_ec_guid)

	// Retrieve a stored handle from the global registry
	myEcHandle := appd.GetExitcall(my_ec_guid)

	// Add the SQL statement to the exit call details
	err = appd.SetExitcallDetails(ecHandle, sqlStatement)
	if err != nil {
		// Handle the error if needed
		fmt.Println("Error setting exit call details:", err)
	}

	appd.EndExitcall(myEcHandle)
	appd.EndBT(businessTransactionHandle)
}
