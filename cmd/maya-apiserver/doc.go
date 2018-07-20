// Package main Maya API Server.
//
//Maya-apiserver helps with creation of CAS Volumes and provides API endpoints to manage those volumes. Maya-apiserver can also be considered as a template engine that can be easily extended to support any kind of CAS storage solutions. It takes as input a set of CAS templates that are converted into CAS K8s YAMLs based on user requests.
//
// Terms Of Service:
//
// https://github.com/openebs/openebs/blob/master/LICENSE
//
//     Schemes: http
//     Host: m-apiserverUrl
//     BasePath: /latest
//     Version: 0.5.3
//     License: Apache License  http://www.apache.org/licenses/
//
//
//     Consumes:
//     - application/json
//     - application/xml
//
//     Produces:
//     - application/json
//     - application/xml
//
//     Extensions:
//     x-meta-value: value
//     x-meta-array:
//       - value1
//       - value2
//     x-meta-array-obj:
//       - name: obj
//         value: field
// swagger:meta
package main