package gas

import (
	
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func ServerBTStarter(w http.ResponseWriter, r *http.Request, appdBTNameHDL string) BtHandle {
	hdr := r.Header.Get(APPD_CORRELATION_HEADER_NAME)
	fmt.Println(appdBTNameHDL, ":", hdr)         // Print the value of hdr
	btHandle := StartBT(appdBTNameHDL, hdr) // Use the function name instead of the URL
	return btHandle
}

func ServerBTEnder(btHandle BtHandle, r *http.Request, sqlStatement string) {

	// Set the snapshot URL if snapshotting is active
	SetBTURL(btHandle, r.URL.Path)

	backendName := BackendName // A unique name for the backend, string type
	backendType := BackendType // Specifies that this is a backend, string type
	
	backendProperties := map[string]string{
		//backend properties  key-value varies based the type of your backend, see documentation for details. This is the  key-value for database
		"DATABASE_TYPE": DatabaseType, // Specify the type of the database (e.g., MySQL, Oracle, SQL Server, etc.), string type
		"HOST":          Host,         // Hostname or IP address of the  server, string type
		"PORT":          Port,         // Port number on which the server is running, string type
		"DATABASE_NAME": DatabaseName, // Name of the database, string type
	}
	resolveBackend := false // Set to false to use static backend configuration

	AddBackend(backendName, backendType, backendProperties, resolveBackend)

	// Generate a new UUID
	my_ec_guid := uuid.New().String()

	// Convert btHandle to the correct type (appdynamics.BtHandle)
	businessTransactionHandle := BtHandle(btHandle)
	ecHandle := StartExitcall(businessTransactionHandle, backendName)

	// optionally store the handle in the global registry
	StoreExitcall(ecHandle, my_ec_guid)

	// retrieve a stored handle from the global registry
	myEcHandle := GetExitcall(my_ec_guid)

	// Add the SQL statement to the exit call details
	err := SetExitcallDetails(ecHandle, sqlStatement)
	if err != nil {
		// Handle the error if needed
		fmt.Println("Error setting exit call details:", err)
	}

	EndExitcall(myEcHandle)

	EndBT(businessTransactionHandle)
}
