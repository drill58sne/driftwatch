// Package inventory provides types and functions for loading and querying
// a list of remote hosts (an "inventory") from a YAML file.
//
// Inventory files describe the servers driftwatch should connect to.
// Each host entry requires at minimum a name, address, and user.
// The port defaults to 22 if not specified.
//
// Example inventory file:
//
//	hosts:
//	  - name: web-01
//	    address: 192.168.1.10
//	    user: admin
//	    port: 22
//	    tags: [web, prod]
//
// Hosts can be filtered by tag using FilterByTag to allow targeted
// drift checks against logical groups of servers.
package inventory
