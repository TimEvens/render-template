/*
Copyright (c) 2022 Cisco Systems, Inc. and others.  All rights reserved.
*/
package render_tpl

import (
	"fmt"
	"github.com/alessio/shellescape"
	"log"
	"net"
	"regexp"
	"strings"
	"text/template"

	"sigs.k8s.io/yaml"
)

func GetCustomTplFuncMap() template.FuncMap {
	tplFuncMap := make(template.FuncMap)
	tplFuncMap["escEnvVarValue"] = escEnvVarValue
	tplFuncMap["escEnvVarName"] = escEnvVarName
	tplFuncMap["strslice"] = strslice
	tplFuncMap["strRegexReplace"] = strRegexReplace
	tplFuncMap["readFile"] = readFile
	tplFuncMap["toYaml"] = toYAML
	tplFuncMap["vasiLinkIps"] = vasiLinkIps
	return tplFuncMap
}

var envVarName = regexp.MustCompile(`(.*)\[(\d+)\]`)
var envVarValue = regexp.MustCompile(`\$\{[A-Za-z0-9_-]*\}`)

func escEnvVarName(str string) string {
	matches := envVarName.FindStringSubmatch(str)
	if matches != nil {
		str = fmt.Sprintf("%s_%s_", matches[1], matches[2])
	}
	return str
}

func escEnvVarValue(str string) string {
	return envVarValue.ReplaceAllStringFunc(shellescape.Quote(str), func(s string) string {
		return fmt.Sprintf("'\"%s\"'", s)
	})
}

//strslice 1 -1 should return a substring with out the first and last chars in the string
func strslice(start, end int, s string) string {
	if len(s) == 0 || s == "null" {
		return "<no value>"
	}
	if start < 0 {
		start = 0
	} else if start >= len(s) {
		start = len(s) - 1
	}

	if end <= 0 {
		end += len(s)
	} else if end >= len(s) {
		end = len(s)
	}
	if end <= start {
		return ""
	}
	return s[start:end]
}

func strRegexReplace(originalString, matchValue, replaceValue string) string {
	output := ""
	matcher, err := regexp.Compile(matchValue)
	if err != nil {
		log.Fatalf("Compile fail: %v\n", err)
	}
	if matcher.MatchString(originalString) {
		output = matcher.ReplaceAllString(originalString, replaceValue)
	} else {
		output = originalString + " " + replaceValue
	}
	return output
}

// readFile read file contents
func readFile(name string) string {
	return string(ReadFile(name))
}

// toYAML takes an interface, marshals it to yaml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func toYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}

func vasiLinkIps(base_network string, customer_id int) [2]string {
	var link_ips [2]string

	ip, ipnet, err := net.ParseCIDR(base_network)

	if err != nil {
		log.Fatal("Invalid base network: ", err)
	}

	net_ip := ip.Mask(ipnet.Mask)

	net_int := ipv4_to_int(net_ip)

	first_ip := int64(customer_id)*2 - 2

	link_ips[0] = ipv4_to_str(net_int + first_ip)
	link_ips[1] = ipv4_to_str(net_int + first_ip + 1)

	return link_ips
}
