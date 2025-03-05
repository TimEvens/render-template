/*
Copyright (c) 2022 Cisco Systems, Inc. and others.  All rights reserved.
*/
package render_tpl

/*
 * Global Configuration variables
 */
var (
	UseStdout    bool
	ValueFiles   []string
	TemplateFile string

	Values map[string]interface{}
)

/*
 * Constants that could be moved to dynamic/user defined configuration
 */
const ()
