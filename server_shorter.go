package gas

import (
	appd "appdynamics"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func ServerBTStarter(w http.ResponseWriter, r *http.Request, appdBTNameHDL string) appd.BtHandle {
	hdr := r.Header.Get(appd.APPD_CORRELATION_HEADER_NAME)
	fmt.Println(appdBTNameHDL, ":", hdr)         // Print the value of hdr
	btHandle := appd.StartBT(appdBTNameHDL, hdr) // Use the function name instead of the URL
	return btHandle
}

func ServerBTEnder(btHandle appd.BtHandle, r *http.Request, sqlStatement string) {

	// Set the snapshot URL if snapshotting is active
	appd.SetBTURL(btHandle, r.URL.Path)

	backendName := BackendName // A unique name for the backend
	backendType := BackendType // Specifies that this is a backend
	backendProperties := map[string]string{
		"DATABASE_TYPE": DatabaseType, // Specify the type of the database (e.g., MySQL, Oracle, SQL Server, etc.)
		"HOST":          Host,         // Hostname or IP address of the  server
		"PORT":          Port,         // Port number on which the server is running
		"DATABASE_NAME": DatabaseName, // Name of the database
	}
	resolveBackend := false // Set to false to use static backend configuration

	appd.AddBackend(backendName, backendType, backendProperties, resolveBackend)

	// Generate a new UUID
	my_ec_guid := uuid.New().String()

	// Convert btHandle to the correct type (appdynamics.BtHandle)
	businessTransactionHandle := appd.BtHandle(btHandle)
	ecHandle := appd.StartExitcall(businessTransactionHandle, backendName)

	// optionally store the handle in the global registry
	appd.StoreExitcall(ecHandle, my_ec_guid)

	// retrieve a stored handle from the global registry
	myEcHandle := appd.GetExitcall(my_ec_guid)

	// Add the SQL statement to the exit call details
	err := appd.SetExitcallDetails(ecHandle, sqlStatement)
	if err != nil {
		// Handle the error if needed
		fmt.Println("Error setting exit call details:", err)
	}

	appd.EndExitcall(myEcHandle)

	appd.EndBT(businessTransactionHandle)
}
