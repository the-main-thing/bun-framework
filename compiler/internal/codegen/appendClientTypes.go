package codegen

import (
	"fmt"
	"path/filepath"
	"strings"
)

const APPEND_CLIENT_TYPES_TEMPLATE = `
import type { 
`
