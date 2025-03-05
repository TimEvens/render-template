/*
Copyright (c) 2022 Cisco Systems, Inc. and others.  All rights reserved.
*/
package render_tpl

import (
	"io/ioutil"
	"math/big"
	"net"
	"strconv"
)

import (
	log "github.com/sirupsen/logrus"
	"os"
	"sigs.k8s.io/yaml"
)

func ReadFile(name string) []byte {
	log.Info("Reading file ", name)
	data, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatalf("read file %v - err %v", name, err)
	}
	return data
}

func ReadYaml(name string, optional bool) map[string]interface{} {
	values := make(map[string]interface{})
	if _, err := os.Stat(name); os.IsNotExist(err) {
		if !optional {
			log.Fatalf("Missing file %s", name)
		}
		log.Warnf("Missing file %s", name)
		return values
	}

	yamlFile := ReadFile(name)
	err := yaml.Unmarshal(yamlFile, &values)
	if err != nil {
		log.Fatalf("Unmarshal yaml file %v: %v", name, err)
	}

	return values
}

// source: https://github.com/helm/helm/blob/master/pkg/chartutil/values.go
// coalesceTables merges a source map into a destination map.
//
// dest is considered authoritative.
func coalesceTables(dst map[string]interface{}, src map[string]interface{}, templateName string) map[string]interface{} {
	// Because dest has higher precedence than src, dest values override src
	// values.

	rv := make(map[string]interface{})
	for key, val := range src {
		dv, ok := dst[key]
		if !ok { // if not in dst, then copy from src
			rv[key] = val
			continue
		}
		if dv == nil { // if set to nil in dst, then ignore
			// When the YAML value is null, we skip the value's key.
			continue
		}

		srcTable, srcIsTable := val.(map[string]interface{})
		dstTable, dstIsTable := dv.(map[string]interface{})
		switch {
		case srcIsTable && dstIsTable: // both tables, we coalesce
			rv[key] = coalesceTables(dstTable, srcTable, templateName)
		case srcIsTable && !dstIsTable:
			log.Printf("Warning: Merging destination map for template '%s'. Overwriting table item '%s', with non table value: %v", templateName, key, dv)
			rv[key] = dv
		case !srcIsTable && dstIsTable:
			log.Printf("Warning: Merging destination map for template '%s'. The destination item '%s' is a table and ignoring the source '%s' as it has a non-table value of: %v", templateName, key, key, val)
			rv[key] = dv
		default: // neither are tables, simply take the dst value
			rv[key] = dv
		}
	}

	// do we have anything in dst that wasn't processed already that we need to copy across?
	for key, val := range dst {
		if val == nil {
			continue
		}
		_, ok := rv[key]
		if !ok {
			rv[key] = val
		}
	}

	return rv
}

// Convert integer IPv4 to dotted string
func ipv4_to_str(ipInt int64) string {

	b0 := strconv.FormatInt((ipInt>>24)&0xff, 10)
	b1 := strconv.FormatInt((ipInt>>16)&0xff, 10)
	b2 := strconv.FormatInt((ipInt>>8)&0xff, 10)
	b3 := strconv.FormatInt((ipInt & 0xff), 10)

	return b0 + "." + b1 + "." + b2 + "." + b3
}

// Convert net.IP to integer
func ipv4_to_int(IPv4Address net.IP) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(IPv4Address.To4())
	return IPv4Int.Int64()
}
