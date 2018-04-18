/*Package global
 is responsible for holding constants or global variables needed by multiple other packages.
Helps resolve the issue of cyclical dependencies.
*/
package global

import (
	"log"
)

// Log represents the logging object to be used by Busted
var Log *log.Logger
